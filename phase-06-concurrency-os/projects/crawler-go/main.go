// Concurrent web crawler.
//
// Demonstrates:
//   - Worker pool with bounded parallelism (semaphore = buffered channel)
//   - Visited set with sync.Map (concurrent-safe map)
//   - context.Context for graceful cancellation (Ctrl-C)
//   - sync.WaitGroup for "wait for all workers to finish"
//
// Run: go run ./... -url=https://example.com -depth=2 -conc=10

package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os/signal"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

func main() {
	seed := flag.String("url", "https://example.com", "seed URL")
	maxDepth := flag.Int("depth", 2, "max crawl depth")
	concurrency := flag.Int("conc", 10, "max concurrent fetches")
	timeout := flag.Duration("timeout", 5*time.Second, "per-request timeout")
	flag.Parse()

	// Cancel on Ctrl-C.
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	c := newCrawler(*concurrency, *timeout)
	c.crawl(ctx, *seed, *maxDepth)
	fmt.Printf("\nDone. Visited %d URLs, %d errors\n", c.visited(), c.errors.Load())
}

type crawler struct {
	client    *http.Client
	sem       chan struct{}     // semaphore for max-concurrent
	seen      sync.Map           // url -> struct{}
	wg        sync.WaitGroup
	count     atomic.Int64
	errors    atomic.Int64
	allowHost string
}

func newCrawler(maxConc int, timeout time.Duration) *crawler {
	return &crawler{
		client: &http.Client{Timeout: timeout},
		sem:    make(chan struct{}, maxConc),
	}
}

func (c *crawler) visited() int64 { return c.count.Load() }

func (c *crawler) crawl(ctx context.Context, seed string, maxDepth int) {
	u, err := url.Parse(seed)
	if err != nil {
		log.Fatalf("invalid URL: %v", err)
	}
	c.allowHost = u.Host
	c.spawn(ctx, seed, 0, maxDepth)
	c.wg.Wait()
}

func (c *crawler) spawn(ctx context.Context, urlStr string, depth, maxDepth int) {
	if depth > maxDepth {
		return
	}
	if _, loaded := c.seen.LoadOrStore(urlStr, struct{}{}); loaded {
		return // already seen
	}

	c.wg.Add(1)
	go func() {
		defer c.wg.Done()

		// Acquire semaphore. Respect context cancellation while waiting.
		select {
		case c.sem <- struct{}{}:
		case <-ctx.Done():
			return
		}
		defer func() { <-c.sem }()

		links, err := c.fetch(ctx, urlStr)
		if err != nil {
			c.errors.Add(1)
			fmt.Printf("[err] %s: %v\n", urlStr, err)
			return
		}
		c.count.Add(1)
		fmt.Printf("[%d] %s (%d links)\n", depth, urlStr, len(links))

		for _, link := range links {
			c.spawn(ctx, link, depth+1, maxDepth)
		}
	}()
}

func (c *crawler) fetch(ctx context.Context, urlStr string) ([]string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", urlStr, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "kvng-crawler/0.1")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("status %d", resp.StatusCode)
	}

	// Read up to 1 MiB to avoid swallowing the heap on huge pages.
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return nil, err
	}

	return c.extractLinks(urlStr, string(body)), nil
}

// hrefRe is a quick-and-dirty href extractor. For real crawling, use a real
// HTML parser (golang.org/x/net/html). Regex over HTML is famous for being wrong;
// for a learning project it's fine.
var hrefRe = regexp.MustCompile(`href=["']([^"'#]+)["']`)

func (c *crawler) extractLinks(base, html string) []string {
	baseURL, err := url.Parse(base)
	if err != nil {
		return nil
	}
	matches := hrefRe.FindAllStringSubmatch(html, -1)
	out := make([]string, 0, len(matches))
	for _, m := range matches {
		raw := strings.TrimSpace(m[1])
		u, err := baseURL.Parse(raw) // resolves relative URLs
		if err != nil {
			continue
		}
		if u.Scheme != "http" && u.Scheme != "https" {
			continue
		}
		if u.Host != c.allowHost {
			continue // same-host only
		}
		u.Fragment = ""
		out = append(out, u.String())
	}
	return out
}
