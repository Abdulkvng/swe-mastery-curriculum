// Concurrent web crawler.
//
// Uses every concurrency primitive in one place:
//   - Goroutine pool (worker pattern)
//   - Semaphore via buffered channel (limit total in-flight fetches)
//   - sync.Map for visited-URL deduplication
//   - sync.WaitGroup to know when all work is done
//   - context.Context for graceful cancellation (Ctrl-C, timeout)
//   - Producer/consumer via a buffered channel for the URL queue
//
// Run:
//   go run . -url https://example.com -depth 2 -workers 8

package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"golang.org/x/net/html"
)

type crawlJob struct {
	url   string
	depth int
}

type crawler struct {
	client    *http.Client
	maxDepth  int
	sem       chan struct{}     // semaphore — concurrency limit
	visited   sync.Map          // map[string]struct{}
	jobs      chan crawlJob     // queue
	wg        sync.WaitGroup    // tracks outstanding jobs
	pages     atomic.Uint64
	bytes     atomic.Uint64
	startTime time.Time
}

func newCrawler(maxDepth, maxConcurrent int) *crawler {
	return &crawler{
		client:    &http.Client{Timeout: 5 * time.Second},
		maxDepth:  maxDepth,
		sem:       make(chan struct{}, maxConcurrent),
		jobs:      make(chan crawlJob, 1024),
		startTime: time.Now(),
	}
}

// enqueue puts a URL on the queue if we haven't seen it.
// Caller must have already incremented wg.
func (c *crawler) enqueue(j crawlJob) {
	if _, dup := c.visited.LoadOrStore(j.url, struct{}{}); dup {
		c.wg.Done()
		return
	}
	if j.depth > c.maxDepth {
		c.wg.Done()
		return
	}
	// Non-blocking send — if jobs is full, just drop (we're saturated).
	select {
	case c.jobs <- j:
	default:
		log.Printf("queue full, dropping %s", j.url)
		c.wg.Done()
	}
}

// run starts N workers and waits for completion.
func (c *crawler) run(ctx context.Context, seed string, workers int) {
	for i := 0; i < workers; i++ {
		go c.worker(ctx)
	}

	c.wg.Add(1)
	c.enqueue(crawlJob{url: seed, depth: 0})

	// Wait for all jobs to finish.
	done := make(chan struct{})
	go func() {
		c.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-ctx.Done():
		log.Printf("canceled: %v", ctx.Err())
	}

	close(c.jobs)
}

func (c *crawler) worker(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case j, ok := <-c.jobs:
			if !ok {
				return
			}
			c.process(ctx, j)
			c.wg.Done()
		}
	}
}

func (c *crawler) process(ctx context.Context, j crawlJob) {
	// Acquire concurrency slot.
	select {
	case c.sem <- struct{}{}:
	case <-ctx.Done():
		return
	}
	defer func() { <-c.sem }()

	links, n, err := c.fetchAndExtract(ctx, j.url)
	if err != nil {
		log.Printf("err %s: %v", j.url, err)
		return
	}
	c.pages.Add(1)
	c.bytes.Add(uint64(n))
	fmt.Printf("[depth=%d] %s (%d bytes, %d links)\n", j.depth, j.url, n, len(links))

	if j.depth >= c.maxDepth {
		return
	}
	base, _ := url.Parse(j.url)
	for _, link := range links {
		abs, err := base.Parse(link)
		if err != nil {
			continue
		}
		if abs.Scheme != "http" && abs.Scheme != "https" {
			continue
		}
		// Stay on the same host to avoid crawling the entire web.
		if abs.Host != base.Host {
			continue
		}
		c.wg.Add(1)
		c.enqueue(crawlJob{url: abs.String(), depth: j.depth + 1})
	}
}

func (c *crawler) fetchAndExtract(ctx context.Context, u string) ([]string, int, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", u, nil)
	if err != nil {
		return nil, 0, err
	}
	req.Header.Set("User-Agent", "kvng-crawler/1.0")
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return nil, 0, fmt.Errorf("status %d", resp.StatusCode)
	}
	body, err := io.ReadAll(io.LimitReader(resp.Body, 5<<20)) // 5MB cap
	if err != nil {
		return nil, 0, err
	}
	return extractLinks(string(body)), len(body), nil
}

func extractLinks(htmlStr string) []string {
	var out []string
	tok := html.NewTokenizer(strings.NewReader(htmlStr))
	for {
		tt := tok.Next()
		if tt == html.ErrorToken {
			return out
		}
		if tt == html.StartTagToken || tt == html.SelfClosingTagToken {
			t := tok.Token()
			if t.Data != "a" {
				continue
			}
			for _, a := range t.Attr {
				if a.Key == "href" && a.Val != "" {
					out = append(out, a.Val)
				}
			}
		}
	}
}

func main() {
	urlFlag := flag.String("url", "", "seed URL")
	depth := flag.Int("depth", 2, "max depth")
	workers := flag.Int("workers", 8, "worker count")
	concurrency := flag.Int("concurrency", 16, "max concurrent fetches")
	timeout := flag.Duration("timeout", 60*time.Second, "overall timeout")
	flag.Parse()

	if *urlFlag == "" {
		log.Fatal("must pass -url")
	}

	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	// Cancel on Ctrl-C.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigCh
		log.Println("interrupt — shutting down")
		cancel()
	}()

	c := newCrawler(*depth, *concurrency)
	c.run(ctx, *urlFlag, *workers)

	dur := time.Since(c.startTime)
	fmt.Printf("\n=== Done ===\n")
	fmt.Printf("Pages: %d\n", c.pages.Load())
	fmt.Printf("Bytes: %d\n", c.bytes.Load())
	fmt.Printf("Time:  %v\n", dur)
	fmt.Printf("Rate:  %.1f pages/sec\n", float64(c.pages.Load())/dur.Seconds())
}
