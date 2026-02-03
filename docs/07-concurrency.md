# Concurrency in Go

> **NOTE**: Related code examples can be found in [`examples/7/`](../examples/7/) which contains goroutine examples.

Firstly, don't mess between [parallelism & concurrency](https://go.dev/blog/waza-talk). Go's concurrency model, based on goroutines and channels, is one of its most powerful features. This guide covers everything from basics to advanced patterns.

Table of Contents:

- [Concurrency in Go](#concurrency-in-go)
  - [1. Goroutines vs OS Threads](#1-goroutines-vs-os-threads)
    - [1.1. OS Threads](#11-os-threads)
    - [1.2. Goroutines](#12-goroutines)
  - [2. Goroutines](#2-goroutines)
    - [2.1. Basic Goroutine](#21-basic-goroutine)
    - [2.2. Passing Parameters](#22-passing-parameters)
    - [2.3. Goroutine with Proper Cleanup](#23-goroutine-with-proper-cleanup)
  - [3. Channels](#3-channels)
    - [3.1. Channel Types](#31-channel-types)
    - [3.2. Buffered vs Unbuffered](#32-buffered-vs-unbuffered)
    - [3.3. Channel Example](#33-channel-example)
    - [3.4. Closing Channels](#34-closing-channels)
  - [4. Select Statement](#4-select-statement)
    - [4.1. Non-blocking Channel Operations](#41-non-blocking-channel-operations)
  - [5. Synchronization Primitives](#5-synchronization-primitives)
    - [5.1. WaitGroup](#51-waitgroup)
    - [5.2. Mutex](#52-mutex)
    - [5.3. Atomic Operations](#53-atomic-operations)
  - [6. Context](#6-context)
    - [6.1. Creating Contexts](#61-creating-contexts)
    - [6.2. Using Context in Functions](#62-using-context-in-functions)
    - [6.3. Context Best Practices](#63-context-best-practices)
  - [7. Concurrency Patterns](#7-concurrency-patterns)
    - [7.1. Worker Pool](#71-worker-pool)
    - [7.2. Fan-Out Fan-In](#72-fan-out-fan-in)
    - [7.3. Rate Limiter](#73-rate-limiter)
  - [8. Best Practices](#8-best-practices)
    - [8.1. Always Use Context for Cancellation](#81-always-use-context-for-cancellation)
    - [8.2. Sender Closes Channels](#82-sender-closes-channels)
    - [8.3. Handle Errors from Goroutines](#83-handle-errors-from-goroutines)
  - [9. Common Pitfalls](#9-common-pitfalls)
    - [9.1. Goroutine Leak](#91-goroutine-leak)
    - [9.2. Race Condition](#92-race-condition)
    - [9.3. Mutex Copying](#93-mutex-copying)
  - [10. Container-aware GOMAXPROCS _(Go 1.25+)_](#10-container-aware-gomaxprocs-go-125)
    - [10.1. The Problem (Before Go 1.25)](#101-the-problem-before-go-125)
    - [10.2. The Solution (Go 1.25+)](#102-the-solution-go-125)
    - [10.3. Opting Out](#103-opting-out)
  - [Further Reading](#further-reading)

## 1. Goroutines vs OS Threads

### 1.1. OS Threads

- A thread is a sequence of instructions that can be executed independently by a processor
- Modern processors can execute multiple threads (multi-threading)
- Threads share memory and don't need new virtual memory space

**Why threads can be slow:**

- Large stack size (>= 1MB)
- Need to save/restore many registers on context switch
- Setup/teardown requires OS calls

### 1.2. Goroutines

- Goroutines exist only in Go's virtual space, not in the OS
- A goroutine is created with only **2KB of stack** (vs 1MB+ for threads)
- Goroutines are **cooperatively scheduled** by the Go runtime
- Very few registers need to be saved/restored on switch

```go
// Start a goroutine
go func() {
    fmt.Println("Hello from goroutine!")
}()
```

## 2. Goroutines

### 2.1. Basic Goroutine

```go
go func() {
    fmt.Println("Hello from goroutine!")
}()
```

### 2.2. Passing Parameters

Always pass parameters to avoid closure issues:

```go
// ✅ Good: Pass parameters explicitly
go func(msg string, id int) {
    fmt.Printf("Message %d: %s\n", id, msg)
}("Hello", 1)

// ❌ Bad: Closure captures variable
for i := 0; i < 5; i++ {
    go func() {
        fmt.Println(i) // All goroutines may print "5"!
    }()
}
```

### 2.3. Goroutine with Proper Cleanup

```go
func main() {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel() // Important: ensures cleanup

    go func(ctx context.Context) {
        for {
            select {
            case <-ctx.Done():
                fmt.Println("Cleanup and exit")
                return
            default:
                // Do work
                time.Sleep(time.Second)
            }
        }
    }(ctx)

    time.Sleep(5 * time.Second)
}
```

### 3.4. Go scheduler

The Go scheduler is the behind-the-scenes magical machine that powers Go programs. It efficiently runs goroutines, and also coordinates network IO and memory management.

[Kavya’s talk](https://youtu.be/YHRO5WQGh0k) will explore the inner workings of the scheduler machinery. She will delve into the M:N multiplexing of goroutines on system threads, and the mechanisms to schedule, unschedule, and rebalance goroutines. Kavya will also touch upon how the scheduler supports the netpoller and the memory management systems for goroutine stack resizing and heap garbage collection. Finally, she will evaluate the effectiveness and performance of the scheduler.

## 3. Channels

Channels are typed conduits for communication between goroutines. You can check out the [Kavya Joshi's talk - Understanding channels](https://youtu.be/KBZlN0izeiY) for the inner working of channels and channel operations, including how they're supported by the runtime scheduler and memory management systems.

### 3.1. Channel Types

```go
// Unbuffered - blocks until receiver is ready
ch := make(chan string)

// Buffered - blocks only when buffer is full
buffered := make(chan string, 5)

// Read-only (receive-only)
func readOnly(ch <-chan string) {}

// Write-only (send-only)
func writeOnly(ch chan<- string) {}
```

### 3.2. Buffered vs Unbuffered

**Buffered Channels:**

- Buffered channel has capacity.
- When a goroutine attempts to send a resource to a buffered and the channel is full, the channel will lock the goroutine and make it wait until a buffer becomes available.
- When a goroutine attempts to receive from a buffered channel and the buffered channel is empty, the channel will lock the goroutine and make it wait until a resource has been sent.

![Buffered Channel](https://www.ardanlabs.com/images/goinggo/Screen+Shot+2014-02-17+at+8.38.15+AM.png)

- Buffered channel is used to perform asynchronous communication.
- A receive will block only if there's no value in the channel to receive.
- A send will block only if there's no available buffer to place the value being sent.

**Unbuffered Channels:**

- Unbuffered channel has no capacity and therefore require both goroutines to be ready to make any exchange.
- When a goroutine attempts to send a resource to an unbuffered channel and there is no goroutine waiting to receive the resource, the channel will lock the sending goroutine and make it wait.
- When a goroutine attempts to receive from an unbuffered channel, and there is no goroutine waiting to send a resource, the channel will lock the receiving goroutine and make it wait.

![Unbuffered Channel](https://www.ardanlabs.com/images/goinggo/Screen+Shot+2014-02-16+at+10.10.54+AM.png)

- Unbuffered channel is used to perform synchronous communication between goroutines. Unbuffered channel provides a guarantee that an exchange between 2 goroutines is performed at the instant the send and receive take place.
- Synchronization is fundamental in the interaction between the send and receive on the channel.

```go
unbuffered := make(chan int)      // Unbuffered channel of integer type
buffered := make(chan int, 10)    // Buffered channel of integer type
```

You can check [ArdanLabs blog post for more detail](https://www.ardanlabs.com/blog/2014/02/the-nature-of-channels-in-go.html).

### 3.3. Channel Example

```go
func main() {
    ch := make(chan int)

    go func() {
        ch <- 42  // Send value
    }()

    value := <-ch  // Receive value
    fmt.Println(value)  // 42
}
```

### 3.4. Closing Channels

```go
// Sender closes the channel
close(ch)

// Check if channel is closed
value, ok := <-ch
if !ok {
    fmt.Println("Channel closed")
}

// Range over channel until closed
for value := range ch {
    fmt.Println(value)
}
```

## 4. Select Statement

Select lets a goroutine wait on multiple channel operations.

```go
select {
case msg := <-ch1:
    fmt.Println("Received from ch1:", msg)
case msg := <-ch2:
    fmt.Println("Received from ch2:", msg)
case <-time.After(time.Second):
    fmt.Println("Timeout!")
default:
    fmt.Println("No message available")
}
```

### 4.1. Non-blocking Channel Operations

```go
select {
case msg := <-ch:
    fmt.Println(msg)
default:
    fmt.Println("No message, moving on")
}
```

## 5. Synchronization Primitives

### 5.1. WaitGroup

```go
var wg sync.WaitGroup

for i := 0; i < 3; i++ {
    wg.Add(1)  // Must be called before goroutine starts
    go func(id int) {
        defer wg.Done()  // Called when goroutine finishes
        fmt.Printf("Worker %d\n", id)
    }(i)
}

wg.Wait()  // Blocks until counter is zero
```

### 5.2. Mutex

```go
type Counter struct {
    mu    sync.RWMutex
    count map[string]int
}

func (c *Counter) Increment(key string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.count[key]++
}

func (c *Counter) Get(key string) int {
    c.mu.RLock()         // Read lock allows multiple readers
    defer c.mu.RUnlock()
    return c.count[key]
}
```

### 5.3. Atomic Operations

`sync/atomic` package provides low-level, primitive operations that perform common tasks like adding, comparing-and-swapping, or loading values in a thread-safe manner without explicit locking. These operations are typically implemented using special CPU instructions that guarantee atomicity, meaning they complete in a single, indivisible step, even in a multi-core environment. This makes them highly efficient for specific use cases.

But why Atomic operations?

Consider a scenario where multiple goroutines need to increment a shared counter.

```go
package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

func main() {
	counter := 0
	numGoroutines := 1000
	var wg sync.WaitGroup

	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				counter++ // Data race!
			}
		}()
	}

	wg.Wait()
	fmt.Println("Final Counter (potential race):", counter)
}
```

Running this code multiple times might yield different results, this is because the `counter++` is not atomic; it involves three steps: read, increment, and write. A context switch can happen between these steps, leading to lost updates. One way to fix this is using `sync.Mutex`, but acquiring and releasing a mutex for every tiny increment can introduce unnecessary overhead. For simple arithmetic operations or value swapping, `sync/atomic` provides a more performant alternative.

`Add*` functions atomically a delta to a value and return the new value.

```go
package main

import (
	"fmt"
	"sync"
	"sync/atomic" // Import the atomic package
)

func main() {
	var counter int64 // Use int64 for atomic operations
	numGoroutines := 1000
	var wg sync.WaitGroup

	wg.Add(numGoroutines)
	for i := 0; i < numGoroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < 1000; j++ {
				atomic.AddInt64(&counter, 1) // Atomically add 1
			}
		}()
	}

	wg.Wait()
	fmt.Println("Final Counter (with atomic):", counter) // Should be 1,000,000
}
```

`Load*` functions atomically load (read) the value stored at an address. It's crucial to always use atomic loads when reading a value that might be written to atomically by another goroutine. This ensures you get the most up-to-date, consistent value.

```go
import "sync/atomic"

type Counter struct {
    count atomic.Int64
}

func (c *Counter) Increment() {
    c.count.Add(1)  // Atomic, no lock needed
}

func (c *Counter) Value() int64 {
    return c.count.Load()
}
```

When to use `sync/atomic` vs. `sync.Mutex`

Use sync/atomic when:

- You need to perform simple read, write, add, swap, or compare-and-swap operations on primitive integer types or pointers.
- Performance is critical, and the operations are truly atomic at the hardware level.
- You want to avoid the overhead of mutex locking/unlocking.
- You're implementing lock-free data structures (though this is advanced and error-prone).

Use sync.Mutex (or sync.RWMutex) when:

- The shared state is a complex data structure (e.g., maps, slices, structs) where multiple fields might be modified in a related way, or where operations involve multiple discrete reads/writes that need to be grouped as a single logical unit.
- The operations are more complex than simple arithmetic or assignments (e.g., appending to a slice, deleting from a map).
- You need to protect an entire critical section, not just a single value.
- Simplicity and correctness are prioritized over micro-optimizations. Mutexes are generally easier to reason about and less error-prone for complex scenarios.

## 6. Context

Context provides cancellation signals and request-scoped values.

### 6.1. Creating Contexts

```go
// Background context (root)
ctx := context.Background()

// With cancellation
ctx, cancel := context.WithCancel(parent)
defer cancel()

// With timeout
ctx, cancel := context.WithTimeout(parent, 5*time.Second)
defer cancel()

// With deadline
ctx, cancel := context.WithDeadline(parent, time.Now().Add(time.Hour))
defer cancel()

// With value (use sparingly!)
ctx := context.WithValue(parent, "key", "value")
```

### 6.2. Using Context in Functions

```go
func doWork(ctx context.Context) error {
    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
            // Do work
        }
    }
}
```

### 6.3. Context Best Practices

- Context should be the first parameter of a function
- Don't store context in a struct
- Use context.Value sparingly - it should inform, not control
- Always call cancel() to release resources

## 7. Concurrency Patterns

### 7.1. Worker Pool

```go
func WorkerPool(jobs []Job, numWorkers int) []Result {
    jobsChan := make(chan Job, len(jobs))
    results := make(chan Result, len(jobs))
    var wg sync.WaitGroup

    // Start workers
    for i := 0; i < numWorkers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for job := range jobsChan {
                results <- process(job)
            }
        }()
    }

    // Send jobs
    for _, job := range jobs {
        jobsChan <- job
    }
    close(jobsChan)

    // Wait and close results
    go func() {
        wg.Wait()
        close(results)
    }()

    // Collect results
    var output []Result
    for r := range results {
        output = append(output, r)
    }
    return output
}
```

### 7.2. Fan-Out Fan-In

```go
// Fan-Out: Distribute work to multiple goroutines
func fanOut(input <-chan int, n int) []<-chan int {
    outputs := make([]<-chan int, n)
    for i := 0; i < n; i++ {
        outputs[i] = worker(input)
    }
    return outputs
}

// Fan-In: Merge multiple channels into one
func fanIn(channels ...<-chan int) <-chan int {
    merged := make(chan int)
    var wg sync.WaitGroup

    for _, ch := range channels {
        wg.Add(1)
        go func(c <-chan int) {
            defer wg.Done()
            for v := range c {
                merged <- v
            }
        }(ch)
    }

    go func() {
        wg.Wait()
        close(merged)
    }()

    return merged
}
```

### 7.3. Rate Limiter

```go
type RateLimiter struct {
    tokens   chan struct{}
    interval time.Duration
}

func NewRateLimiter(rate int, interval time.Duration) *RateLimiter {
    rl := &RateLimiter{
        tokens:   make(chan struct{}, rate),
        interval: interval,
    }

    // Fill initial tokens
    for i := 0; i < rate; i++ {
        rl.tokens <- struct{}{}
    }

    // Refill periodically
    go func() {
        ticker := time.NewTicker(interval)
        for range ticker.C {
            select {
            case rl.tokens <- struct{}{}:
            default:
            }
        }
    }()

    return rl
}

func (rl *RateLimiter) Wait() {
    <-rl.tokens
}
```

## 8. Best Practices

### 8.1. Always Use Context for Cancellation

```go
// ❌ Bad
func longTask() {
    for {
        // No way to stop!
    }
}

// ✅ Good
func longTask(ctx context.Context) error {
    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
            // Work
        }
    }
}
```

### 8.2. Sender Closes Channels

```go
// ❌ Bad: Receiver closes
go func() {
    <-ch
    close(ch)  // Wrong!
}()

// ✅ Good: Sender closes
go func() {
    defer close(ch)
    ch <- data
}()
```

### 8.3. Handle Errors from Goroutines

```go
// ✅ Use error channels
errCh := make(chan error, 1)
go func() {
    if err := work(); err != nil {
        errCh <- err
    }
}()

// Or use errgroup
g, ctx := errgroup.WithContext(ctx)
g.Go(func() error {
    return work()
})
if err := g.Wait(); err != nil {
    log.Fatal(err)
}
```

## 9. Common Pitfalls

### 9.1. Goroutine Leak

```go
// ❌ Goroutine stuck forever
ch := make(chan int)
go func() {
    val := <-ch  // Never receives!
}()
// ch never gets a value

// ✅ Use context for cleanup
go func(ctx context.Context) {
    select {
    case val := <-ch:
        process(val)
    case <-ctx.Done():
        return
    }
}(ctx)
```

### 9.2. Race Condition

```go
// ❌ Data race
counter := 0
for i := 0; i < 1000; i++ {
    go func() {
        counter++  // Race!
    }()
}

// ✅ Use atomic or mutex
var counter atomic.Int64
for i := 0; i < 1000; i++ {
    go func() {
        counter.Add(1)
    }()
}
```

### 9.3. Mutex Copying

```go
// ❌ Mutex copied (value receiver)
func (c Config) Get() string {
    c.Lock()  // Operates on copy!
    defer c.Unlock()
    return c.data
}

// ✅ Use pointer receiver
func (c *Config) Get() string {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.data
}
```

## 10. Container-aware GOMAXPROCS _(Go 1.25+)_

Go 1.25 introduces container-aware CPU scheduling. On Linux, the runtime now considers the CPU bandwidth limit of the containing cgroup.

### 10.1. The Problem (Before Go 1.25)

In a container with CPU limits (e.g., Kubernetes pod with 2 CPU limit on an 8-core host):

- `GOMAXPROCS` would default to 8
- This causes excessive goroutine scheduling and CPU throttling

### 10.2. The Solution (Go 1.25+)

```go
// GOMAXPROCS now automatically respects container CPU limits
// In a 2-CPU limited container on an 8-core host:
// - Before: GOMAXPROCS = 8
// - After:  GOMAXPROCS = 2
```

### 10.3. Opting Out

If you need the old behavior:

```bash
GODEBUG=containercpucount=0 ./myapp
```

This change significantly improves performance in containerized environments like Kubernetes by preventing CPU throttling.

## Further Reading

- [Go Concurrency Patterns](https://go.dev/blog/pipelines)
- [Share Memory By Communicating](../tips-notes/share-memory-by-communicating.md)
- [Pipelines and Cancellation](../tips-notes/pipelines-cancellation.md)
- [Generic Concurrency](../tips-notes/generic-concurrency.md)
- [CPU Throttling in Containerized Go Apps](../tips-notes/cpu-throttling-in-containerized-go-apps.md)
