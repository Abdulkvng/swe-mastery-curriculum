# Phase 6 — Concurrency & OS

> Every interesting bug above the application layer is a concurrency bug. Every interesting bug below it is an OS bug. This phase teaches you the model — processes, threads, the kernel, the scheduler — and the patterns that handle concurrency safely. We compare Go's "share memory by communicating" with Rust's "no data races at compile time," and finish with a tour of embedded systems so the word stops being scary.

**Time:** 6–9 days.

**You'll know you're done when:** you can explain processes vs threads, write race-free Go and Rust, defend mutex vs channel choices, and explain why your laptop can run thousands of "concurrent" things on 8 CPU cores.

---

## Table of contents

1. [What does this even mean? — "Concurrency" vs "Parallelism"](#what)
2. [Module 6.1 — Processes vs threads](#processes-threads)
3. [Module 6.2 — The OS scheduler, virtual memory, syscalls](#os)
4. [Module 6.3 — Concurrency primitives](#primitives)
5. [Module 6.4 — Deadlock, livelock, race condition, starvation](#hazards)
6. [Module 6.5 — Go: goroutines and channels](#go-concurrency)
7. [Module 6.6 — Rust: ownership-based safety](#rust-concurrency)
8. [Module 6.7 — Async/await models compared](#async)
9. [Module 6.8 — What does "embedded" even mean?](#embedded)
10. [🛠️ Project: Concurrent web crawler (Go)](#crawler)
11. [Interview question bank](#interview-questions)
12. [What you should now know](#what-you-should-now-know)

---

<a name="what"></a>
## 🧠 What does this even mean? — "Concurrency" vs "Parallelism"

These two words get used interchangeably. They are NOT the same.

**Concurrency** = multiple things *making progress* over the same period. Could be on one CPU, fast-switching between tasks. Like a chef juggling 5 pans on 1 stove — at any instant they're stirring just one, but over the meal, all 5 progress.

**Parallelism** = multiple things *running simultaneously* on multiple CPUs. Two chefs, two stoves, two pans cooking at the literal same moment.

You can have:
- Concurrency without parallelism (single-threaded async JavaScript)
- Parallelism without concurrency (a SIMD instruction processing 8 numbers at once — same thing, in parallel)
- Both (Go program with 8 goroutines on 8 cores)

The key insight: **concurrency is a property of the program; parallelism is a property of the hardware.** Your program can be concurrent and run on 1 core (no parallelism but useful — you can wait on I/O without blocking). Or it can be parallel without being concurrent (a tight loop using AVX instructions). Most real programs want concurrency that *takes advantage of* parallelism.

> 🎯 **Interview answer:** "Concurrency is about *dealing with* many things at once. Parallelism is about *doing* many things at once. — Rob Pike."

---

<a name="processes-threads"></a>
## Module 6.1 — Processes vs threads

> 📖 **Definition — Process:** A running program with its own memory space (its own virtual address space), file descriptors, and identity (PID). Two processes cannot accidentally read each other's memory. The OS isolates them.

> 📖 **Definition — Thread:** A path of execution *within* a process. A process can have many threads. **Threads in the same process share the same memory.** That's their power and their danger.

### A picture

```
┌─────────────────── Process A (PID 1234) ──────────────────┐
│  Heap ──────────────────────────────                       │
│  ┌────────┐    ┌────────┐    ┌────────┐                    │
│  │Thread 1│    │Thread 2│    │Thread 3│                    │
│  │  stack │    │  stack │    │  stack │                    │
│  │  regs  │    │  regs  │    │  regs  │                    │
│  └────────┘    └────────┘    └────────┘                    │
│      ↘            ↓             ↙                          │
│       (all share the same heap, same fds, same code)       │
└────────────────────────────────────────────────────────────┘

Different process — totally separate world:

┌─────────────────── Process B (PID 5678) ──────────────────┐
│  Heap (different from A's!)                                │
│  ┌────────┐                                                │
│  │Thread 1│                                                │
│  └────────┘                                                │
└────────────────────────────────────────────────────────────┘
```

### Trade-offs

| | Process | Thread |
|---|---|---|
| Memory | Isolated (safer) | Shared (faster comms, more bugs) |
| Creation cost | Heavier (fork + setup) | Lighter |
| Communication | IPC (pipes, sockets, shared memory) | Just shared variables (with locks) |
| Crash blast radius | Just one process dies | All threads in the process die |
| Use when | Strong isolation needed (browser tabs, web servers) | Tight cooperation (parallel computation) |

### Modern reality: lightweight tasks

Real OS threads cost ~2MB stack each. You can't have 100,000 of them.

So languages built **lightweight tasks** that the runtime multiplexes onto a small number of OS threads:
- **Goroutines** (Go) — ~4KB initial stack, grow dynamically. 100k+ trivial.
- **Async tasks** (Rust, Python, JavaScript) — Future objects, no stack at all. Ridiculously cheap.
- **Virtual threads** (Java 21+) — same idea, finally.

When people say "thread," they could mean OS thread, language-level lightweight thread, or both. Always clarify.

---

<a name="os"></a>
## Module 6.2 — The OS scheduler, virtual memory, syscalls

### The OS scheduler

Your laptop has 8 cores. You're running ~500 processes (check Activity Monitor). They're not all running simultaneously — they're being **scheduled**.

The kernel scheduler decides, at any moment, which N (= core count) of the runnable threads gets to actually execute. Each one gets a *time slice* (typically a few ms), then the kernel preempts it and switches to another.

Linux uses CFS (Completely Fair Scheduler — until kernel 6.6, then EEVDF). It tracks how much CPU time each thread has used and gives the next slice to whoever is "behind." Priorities can adjust this (nice values, realtime classes).

The cost of switching threads — the **context switch** — is real: save registers, load registers, possibly flush TLB. On the order of 1–10 microseconds. Sounds small until you do millions a second; then your CPU is "spending all its time deciding what to spend time on."

### Virtual memory

Each process sees its own contiguous address space (e.g., 0x0 to 0x7fff_ffff_ffff on 64-bit). It's a **lie** — the OS maps virtual addresses to physical RAM (or disk via swap) on demand.

Why:
- **Isolation** — process A's address 0x1000 and B's 0x1000 are different RAM
- **Lazy allocation** — `malloc(1GB)` doesn't reserve 1GB until you actually touch each page
- **Swap** — pages can move to disk when memory's tight
- **Memory-mapped files** — you read a file by treating it as memory

The mapping is held in **page tables**, walked by the CPU's MMU on every memory access. The TLB caches recent translations. A TLB miss is ~10× slower than a hit. This is why "data locality" matters for performance.

### Syscalls

> 📖 **Definition — Syscall (system call):** A way for user code to request a privileged operation from the kernel (file I/O, network, fork, kill...). Costs hundreds of nanoseconds — orders of magnitude more than a function call.

```
write(fd, buf, n)
    ↓
user mode → kernel mode (a "trap" instruction)
    ↓
kernel does the work, returns
    ↓
back to user mode
```

That mode switch isn't free. It's why batching matters (one big `write` is faster than 1000 small ones), why `mmap` for big files beats `read` (no per-block syscall), and why kernel-bypass networking (DPDK, io_uring) exists for extreme throughput.

### Useful Linux/Mac commands

```bash
# How many threads is this process running?
ps -M -p <pid>            # Mac
ls /proc/<pid>/task/      # Linux

# What syscalls is it making? (Linux: strace; Mac: dtruss)
sudo dtruss -p <pid>      # Mac
strace -p <pid>           # Linux

# Memory layout
vmmap <pid>               # Mac
cat /proc/<pid>/maps      # Linux
```

---

<a name="primitives"></a>
## Module 6.3 — Concurrency primitives

The toolbox.

### Mutex (mutual exclusion)

> 📖 **Definition — Mutex:** A lock. At most one thread can hold it at a time. Others trying to lock it block until it's released.

```go
var (
    counter int
    mu      sync.Mutex
)

func incr() {
    mu.Lock()
    defer mu.Unlock()
    counter++
}
```

Always pair `Lock()` with `defer Unlock()` so panics don't leak the lock.

### RWMutex (read-write lock)

Many readers OR one writer. Useful when reads vastly outnumber writes (caches, config).

```go
var (
    cache = map[string]int{}
    mu    sync.RWMutex
)

func get(k string) (int, bool) {
    mu.RLock()
    defer mu.RUnlock()
    v, ok := cache[k]
    return v, ok
}

func set(k string, v int) {
    mu.Lock()
    defer mu.Unlock()
    cache[k] = v
}
```

### Semaphore

Like a mutex but with a count. "At most N concurrent holders." Useful for resource pools, rate limiters in code.

```go
sem := make(chan struct{}, 10)   // semaphore with capacity 10

func work() {
    sem <- struct{}{}   // acquire (blocks if 10 already held)
    defer func() { <-sem }()
    // do work
}
```

### Atomic operations

> 📖 **Definition — Atomic:** A read/write/update that the CPU guarantees no other thread can interleave with. No locking needed — single-instruction.

```go
import "sync/atomic"

var counter atomic.Int64
counter.Add(1)
v := counter.Load()
```

Faster than a mutex for simple counters/flags. Useless for anything that requires coordinating multiple variables.

### Condition variable

Wait until a condition is true. Often hidden behind higher-level primitives but worth knowing.

```go
var (
    mu    sync.Mutex
    cond  = sync.NewCond(&mu)
    ready bool
)

// Waiter
mu.Lock()
for !ready {
    cond.Wait()   // releases mu, blocks; reacquires mu on wake
}
mu.Unlock()

// Signaler
mu.Lock()
ready = true
cond.Broadcast()  // wake all waiters
mu.Unlock()
```

### Channel (Go's specialty)

A typed pipe between goroutines. Combines synchronization + data transfer.

```go
ch := make(chan int, 5)   // buffered channel, capacity 5
ch <- 42                  // send
v := <-ch                 // receive
close(ch)                 // close (further receives return zero, ok=false)
```

Channels are first-class — you pass them to functions, store them in structs. Different from threads/queues in most languages, where the queue is a separate concept.

### Once

Run a function exactly once across all goroutines. Singleton's best friend.

```go
var once sync.Once
once.Do(func() { setup() })   // setup runs once total
```

(You used this in Phase 2's connection pool.)

---

<a name="hazards"></a>
## Module 6.4 — Deadlock, livelock, race condition, starvation

The four classic concurrency hazards.

### Race condition

> 📖 **Definition — Race condition:** Two threads access shared data, at least one writes, and the result depends on timing.

```go
// Two goroutines run this concurrently
counter++   // NOT one operation. Reads counter, adds 1, writes back.
            // Both can read 5, both write 6. Lost update.
```

Detection: `go test -race` (Go), `thread sanitizer` (Rust, C++), `helgrind` (valgrind).

Cure: mutex, atomic, or don't share.

### Deadlock

> 📖 **Definition — Deadlock:** Two or more threads each waiting for a resource the other holds. Nobody makes progress.

Classic recipe (Coffman conditions, all four required):
1. Mutual exclusion (resources can't be shared)
2. Hold and wait (you hold one, you wait for another)
3. No preemption (you can't be forced to release)
4. Circular wait (A waits for B, B waits for A)

The fix: break one condition. The most practical fix: **always acquire locks in the same global order.**

```go
// BAD: deadlock-prone
go func() { a.Lock(); b.Lock(); ... }()
go func() { b.Lock(); a.Lock(); ... }()

// GOOD: always a then b
go func() { a.Lock(); b.Lock(); ... }()
go func() { a.Lock(); b.Lock(); ... }()
```

### Livelock

Threads keep responding to each other but make no progress. Like two people in a hallway, both stepping aside, then both stepping back, forever. Rarer than deadlock; often fixed by adding randomness.

### Starvation

A thread never gets the resource because higher-priority threads keep grabbing it. Fix: fair scheduling, FIFO queues, priority inheritance.

---

<a name="go-concurrency"></a>
## Module 6.5 — Go: goroutines and channels

Go's concurrency story is unusual and powerful.

### Goroutines

```go
go someFunc()              // start a goroutine
go func() { ... }()        // anonymous
```

Cost: ~4KB initial stack. Multiplexed onto OS threads by Go's runtime (the famous M:N scheduler — M goroutines on N OS threads). Fully preemptive since Go 1.14.

### Channels — patterns

```go
// Producer-consumer
ch := make(chan int, 100)
go func() {
    for i := 0; i < 1000; i++ {
        ch <- i
    }
    close(ch)
}()

for v := range ch {   // ranges until ch is closed AND drained
    fmt.Println(v)
}
```

### Fan-out / fan-in

```go
// Fan out: many workers consume from one channel
func runWorkers(jobs <-chan Job, results chan<- Result, n int) {
    var wg sync.WaitGroup
    for i := 0; i < n; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for j := range jobs {
                results <- process(j)
            }
        }()
    }
    go func() {
        wg.Wait()
        close(results)
    }()
}
```

### Select — multiplexing

```go
select {
case v := <-ch1:
    fmt.Println("got", v)
case ch2 <- 42:
    fmt.Println("sent")
case <-time.After(1 * time.Second):
    fmt.Println("timeout")
case <-ctx.Done():
    fmt.Println("canceled")
    return
}
```

`select` waits for any of the cases to become ready. If multiple, chooses randomly. The `default` case makes it non-blocking.

### Context — cancellation propagation

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

result, err := callDownstream(ctx)
```

The context flows down through every function. When canceled (timeout, client disconnect, parent cancellation), every layer can react. Critical for graceful shutdown and avoiding leaked goroutines.

### "Don't communicate by sharing memory; share memory by communicating"

Go's mantra. Translation: instead of mutex-protecting shared data, send it through a channel from one goroutine to another. Whoever holds the data right now is the only one allowed to mutate it.

It's not always the right answer (mutexes are fine for caches), but the bias toward channels prevents whole classes of bugs.

---

<a name="rust-concurrency"></a>
## Module 6.6 — Rust: ownership-based safety

Rust takes a radically different approach: **the compiler refuses to compile code with data races.** Not "warns" — refuses.

### Ownership rules

1. Each value has exactly one owner.
2. When the owner goes out of scope, the value is dropped.
3. You can have many *immutable* borrows OR exactly one *mutable* borrow at a time.

Combined with `Send` and `Sync` traits:
- `Send` = "this type can be moved to another thread"
- `Sync` = "this type can be shared (via &T) across threads"

`std::thread::spawn` requires its closure to be `Send`. So if you try to share a `Rc<RefCell<T>>` (not `Sync`) across threads, the compiler stops you.

```rust
use std::sync::{Arc, Mutex};
use std::thread;

fn main() {
    let counter = Arc::new(Mutex::new(0));   // Arc = atomic Rc, Sync
    let mut handles = vec![];

    for _ in 0..10 {
        let counter = Arc::clone(&counter);
        let h = thread::spawn(move || {
            let mut n = counter.lock().unwrap();
            *n += 1;
        });
        handles.push(h);
    }

    for h in handles { h.join().unwrap(); }
    println!("{}", *counter.lock().unwrap());   // 10
}
```

The `Mutex<T>` *contains* the data. You can only access `T` while holding the lock — the lock guard auto-releases when it goes out of scope. There's no way to forget to lock.

### Channels in Rust

```rust
use std::sync::mpsc;
use std::thread;

let (tx, rx) = mpsc::channel();

thread::spawn(move || {
    for i in 0..10 {
        tx.send(i).unwrap();
    }
});

while let Ok(v) = rx.recv() {
    println!("{}", v);
}
```

`mpsc` = multi-producer, single-consumer. (`tokio::sync::mpsc` for async; `crossbeam::channel` for SPMC/MPMC.)

### Async Rust

Rust's async is *zero-cost* — Futures are state machines with no allocation. But you need a runtime to execute them. The standard one is **tokio**.

```rust
#[tokio::main]
async fn main() {
    let task1 = tokio::spawn(async { fetch("a").await });
    let task2 = tokio::spawn(async { fetch("b").await });
    let (a, b) = tokio::join!(task1, task2);
    println!("{:?} {:?}", a, b);
}

async fn fetch(name: &str) -> String {
    tokio::time::sleep(std::time::Duration::from_millis(100)).await;
    format!("got {name}")
}
```

`tokio::spawn` is the analog of `go func()`. `await` is the suspension point.

---

<a name="async"></a>
## Module 6.7 — Async/await models compared

Different languages picked different concurrency styles. The *runtime model* matters more than syntax.

| Language | Model | Notes |
|---|---|---|
| Go | Goroutines + channels (M:N preemptive) | Simplest mental model. Stack-ful. |
| Rust | Stackless futures + `.await` + runtime (tokio) | Fastest. Hardest. |
| JavaScript | Single-threaded event loop, `async/await` | No real parallelism (without workers). I/O concurrency only. |
| Python | `asyncio`, GIL on CPU-bound code | Async for I/O; `multiprocessing` for CPU. |
| Java (21+) | Virtual threads, OS threads | Brand new; finally cheap threads. |
| C# | `Task<T>` + `async/await`, thread pool | Mature, polished. |

### Stackful (goroutines, virtual threads) vs stackless (Rust futures, JS Promises)

- **Stackful:** each task has its own stack. You can suspend in deeply nested function calls. Easy to reason about. Memory cost per task.
- **Stackless:** the compiler turns your `async fn` into a state machine. Suspension points are explicit (`await`). Tiny per-task cost. Function coloring problem ("can only call async from async").

### Function coloring — the async curse

In Rust/JS/Python: `async fn` and regular `fn` are different colors. You can't directly call an async function from a sync one. This propagates upward — once one function in your stack is async, everything above must be too.

Go and virtual-thread Java avoid this entirely: there's just "the function." It can block on I/O without making everything else async.

This is one reason Go is so loved for backend services: no coloring drama.

---

<a name="embedded"></a>
## Module 6.8 — What does "embedded" even mean?

> 📖 **Definition — Embedded system:** A computer that lives inside something that isn't typically called a "computer" — your car's engine controller, a microwave, a Nest thermostat, a pacemaker, the AirPods firmware.

Distinguishing features:

- **Constrained.** Maybe 32 KB of RAM. Maybe 256 KB of flash. No swap, no disk, sometimes no OS.
- **Real-time.** A pacemaker must respond within X ms. Not "usually fast" — *always* fast.
- **Direct hardware.** You toggle GPIO pins, read ADCs, write registers. No `printf` (or printf is your debugger).
- **Long-lived.** A car ECU is in the field for 15+ years. You can't easily push a fix.

### Why anyone would care for an Apple/Datadog interview

You probably won't be writing embedded code at Datadog. But:
- Apple has *huge* embedded teams (AirPods, Watch, sensors).
- Knowing the constraints helps you appreciate what desktop OSes do for you.
- "I understand the stack from MMU to React" is impressive.

### The embedded toolchain (Rust example)

Rust's `no_std` mode strips the standard library — no heap, no threads, no syscalls. You're left with `core` (the type system, basic algorithms) and bare-metal access.

```rust
#![no_std]
#![no_main]

use core::panic::PanicInfo;

#[no_mangle]
pub extern "C" fn _start() -> ! {
    // Toggle a GPIO pin to blink an LED.
    let gpio = 0x4002_0000 as *mut u32;
    loop {
        unsafe {
            *gpio = 1;
            for _ in 0..100_000 { core::hint::spin_loop(); }
            *gpio = 0;
            for _ in 0..100_000 { core::hint::spin_loop(); }
        }
    }
}

#[panic_handler]
fn panic(_info: &PanicInfo) -> ! {
    loop {}
}
```

You'd cross-compile this to ARM Cortex-M or RISC-V, flash to a board, and watch an LED blink.

The skeleton for QEMU-emulated blink lives in `projects/blink-qemu/`. Optional — only do this if you're curious.

### RTOS — Real-Time Operating Systems

For larger embedded projects, an RTOS (FreeRTOS, Zephyr, embedded-hal) provides primitive task scheduling, mutexes, queues, all in tens of KB. Halfway between bare-metal and Linux.

---

<a name="crawler"></a>
## 🛠️ Project: Concurrent web crawler (Go)

The classic — uses every concurrency primitive in one program.

**See `projects/web-crawler/`.**

### Spec

- Given a seed URL and `max_depth`, fetch and parse pages
- Extract links; queue new URLs (within depth)
- Limit total concurrent fetches (semaphore)
- Deduplicate URLs (sync.Map)
- Respect a context-driven timeout
- Print results to stdout
- Graceful shutdown on Ctrl-C

Concepts exercised:
- Goroutine pool
- Semaphore via channel
- Producer-consumer queues
- Worker fan-out, fan-in
- Context cancellation
- `sync.WaitGroup`, `sync.Map`

---

<a name="interview-questions"></a>
## 🎯 Interview question bank

1. **Concurrency vs parallelism — give a one-sentence definition and an example of each.**
2. **What's a thread? What's a process? When would you use one over the other?**
3. **What's a context switch and why is it expensive?**
4. **What's a race condition? How do you prevent one?**
5. **Walk through the four conditions for deadlock.**
6. **What's the difference between a mutex and a semaphore?**
7. **When is `sync.RWMutex` better than `sync.Mutex`? When is it worse?**
8. **Why is `i++` not atomic?**
9. **Implement a thread-safe counter without a mutex.** *(atomic.Int64.Add)*
10. **What's a goroutine and how is it different from an OS thread?**
11. **Why does Rust prevent data races at compile time?** *(Ownership: at most one mutable borrow OR many shared borrows.)*
12. **Implement a worker pool in Go.**
13. **Explain async/await. Why does Rust have "function coloring"?**
14. **How does virtual memory work and why does it matter?**
15. **What's an embedded system? How does Rust `no_std` differ from regular Rust?**

---

## ✅ What you should now know

- [ ] Concurrency vs parallelism
- [ ] Process vs thread, when each
- [ ] OS scheduler, virtual memory, syscall costs
- [ ] Mutex, RWMutex, semaphore, atomics, channels, condvar, sync.Once
- [ ] Race conditions, deadlock, livelock, starvation
- [ ] Goroutines, channels, select, context
- [ ] Rust's ownership-based concurrency safety
- [ ] Stackful vs stackless async, function coloring
- [ ] What "embedded" means and what `no_std` is
- [ ] You've built a concurrent web crawler

---

**Next:** [Phase 7 — Distributed Systems & System Design](../phase-07-distributed-systems/README.md)
