# Common cocurrent programming mistakes

- [Common cocurrent programming mistakes](#common-cocurrent-programming-mistakes)
  - [1. No synchronizations when synchronizations are needed](#1-no-synchronizations-when-synchronizations-are-needed)
  - [2. Use `time.Sleep` calls to do synchronizations](#2-use-timesleep-calls-to-do-synchronizations)
  - [3. Leaves Goroutines hanging](#3-leaves-goroutines-hanging)
  - [4. Copy values of the types in the `sync` standard package](#4-copy-values-of-the-types-in-the-sync-standard-package)
  - [5. Call the `sync.WaitGroup.Add method` at wrong places](#5-call-the-syncwaitgroupadd-method-at-wrong-places)
  - [6. Use Channels as futures/promises improperly](#6-use-channels-as-futurespromises-improperly)
  - [7. Close channels not from the last active sender goroutine](#7-close-channels-not-from-the-last-active-sender-goroutine)
  - [8. Do 64-bit atomic operations on values which are not guaranteed to be 8-byte aligned](#8-do-64-bit-atomic-operations-on-values-which-are-not-guaranteed-to-be-8-byte-aligned)
  - [9. Not pay attention to too many resources are consumed by calls to the `time.After` function](#9-not-pay-attention-to-too-many-resources-are-consumed-by-calls-to-the-timeafter-function)
  - [10. Use `time.Timer` values incorrectly](#10-use-timetimer-values-incorrectly)

## 1. No synchronizations when synchronizations are needed

- The program may run well on one computer, but may panic on another one, or it runs well when it is compiled by one computer, but panics when another computer is used.

```go
package main

import (
    "time"
    "runtime"
)

// the code lines might be not executed by their appearance order
func main() {
    var a []int // nil
    var b bool  // false

    // a new goroutine
    go func () {
        a = make([]int, 3)
        b = true // write b
    }()

    // compilers and CPUs may make optimizations by reordering
    // instructions in the new goroutine, so the assignment of `b`
    // may happen before the assignment of `a` at run time.
    // -> slice `a` is still nill.
    for !b { // read b
        time.Sleep(time.Second)
        runtime.Gosched()
    }
    a[0], a[1], a[2] = 0, 1, 2 // might panic
}
```

- We should use channel or the synchronization techniques provided in the `sync` to ensure the memory orders.

```go
package main

func main() {
    var a []int = nil
    c := make(chan struct{})

    go func () {
        a = make([]int, 3)
        c <- struct{}{}
    }()

    <-c
    // The next line will not panic for sure.
    a[0], a[1], a[2] = 0, 1, 2
}
```

## 2. Use `time.Sleep` calls to do synchronizations

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    var x = 123

    // Go runtime doesn't guarantee the write of `x`
    // happens before the read of `x` for sure.
    go func() {
        x = 789 // write x
    }()

    time.Sleep(time.Second)
    fmt.Println(x) // read x
}
```

## 3. Leaves Goroutines hanging

- Reasons:
  - A goroutine tries to receive a value from a channel which no more other goroutines will send values to.
  - A goroutine tries to send a value to nil channel or to a channel which no more other goroutines will receive values from.
  - A goroutine is dead locked by itself.
  - A group of goroutines are dead locked by each other.
  - A goroutine is blocked when executing a select code block without default branch, and all the channel operations following the case keywords in the select code block keep blocking forever.
- If the capacity of the channel which is used a future is not large enough, some slower response goroutines will hang when trying to send a result to the future channel.

```go
func request() int {
    c := make(chan int)
    for i := 0; i < 5; i++ {
        i := i
        go func() {
            c <- i // 4 goroutines will hang here.
        }()
    }
    return <-c
}
```

## 4. Copy values of the types in the `sync` standard package

- In practice, values of the types (except the Locker interface values) in the sync standard package should never be copied. We should only copy pointers of such values.

```go
import "sync"

type Counter struct {
    sync.Mutex
    n int64
}

// This method is okay.
func (c *Counter) Increase(d int64) (r int64) {
    c.Lock()
    c.n += d
    r = c.n
    c.Unlock()
    return
}

// The method is bad. When it is called,
// the Counter receiver value will be copied.
// As a field of the receiver value, the respective Mutex field of the Counter receiver
// value will also be copied. The copy is ot synchronized, so the copied Mutex value
// might be corrupted.
func (c Counter) Value() (r int64) {
    c.Lock()
    r = c.n
    c.Unlock()
    return
}
```

- Use `go vet` command provided in Go toolchain will report potentional bad value copies.

## 5. Call the `sync.WaitGroup.Add method` at wrong places

- To make the uses of `WaitGroup` value meaningful, when the counter of a `WaitGroup` value is zero, the next call to the `Add` method of the `WaitGroup` value must happen before the next call to the `Wait` method of the `WaitGroup` value.

```go
package main

import (
    "fmt"
    "sync"
    "sync/atomic"
)

func main() {
    var wg sync.WaitGroup
    var x int32 = 0
    // None of the Add method calls are guaranteed to ahppen
    // before the Wait method call.
    for i := 0; i < 100; i++ {
        go func() {
            wg.Add(1)
            atomic.AddInt32(&x, 1)
            wg.Done()
        }()
    }
    // Move the Add method calls out of the new goroutines
    // for i := 0; i < 100; i++ {
    //     wg.Add(1)
    //     go func() {
    //         atomic.AddInt32(&x, 1)
    //         wg.Done()
    //     }()
    // }

    fmt.Println("Wait ...")
    wg.Wait()
    fmt.Println(atomic.LoadInt32(&x))
}
```

## 6. Use Channels as futures/promises improperly

```go
// The generations of the two arguments are processed sequentially
// instead of concurrently.
doSomethingWithFutureArguments(<-fa(), <-fb())

// process them concurrently
ca, cb := fa(), fb()
doSomethingWithFutureArguments(<-ca, <-cb)
```

## 7. Close channels not from the last active sender goroutine

## 8. Do 64-bit atomic operations on values which are not guaranteed to be 8-byte aligned

## 9. Not pay attention to too many resources are consumed by calls to the `time.After` function

```go
import (
    "fmt"
    "time"
)

// The function will return if a message
// arrival interval is larger than one minute.
func longRunning(messages <-chan string) {
    for {
        select {
        case <-time.After(time.Minute):
            return
        case msg := <-messages:
            fmt.Println(msg)
        }
    }
}

// The right one, avoid too many Timer values being created
func longRunning(messages <-chan string) {
    timer := time.NewTimer(time.Minute)
    defer timer.Stop()

    for {
        select {
        case <-timer.C: // expires (timeout)
            return
        case msg := <-messages:
            fmt.Println(msg)

            // This "if" block is important.
            if !timer.Stop() {
                <-timer.C
            }
        }

        // Reset to reuse.
        timer.Reset(time.Minute)
    }
}
```

## 10. Use `time.Timer` values incorrectly
