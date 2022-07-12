# Goroutines size

Source:
- <https://tpaschalis.github.io/goroutines-size/>
- <https://dave.cheney.net/2013/06/02/why-is-a-goroutines-stack-infinite>
- <https://stackoverflow.com/questions/8509152/max-number-of-goroutines/8534711#8534711>


- Per [Go FAQ](https://golang.org/doc/faq#goroutines):

```
It is practical to create hundreds of thousands of goroutines in the same address space.
```

- To understand the max number of goroutines, note that per-goroutine cost is primarily the stack. Per FAQ again:

```
which we call goroutines, can be very cheap: they have little overhead beyond the memory for the stack, which is just a few kilobytes.
```

- A goroutine starts with a [2-kilobyte minimum stack size](https://github.com/golang/go/blob/go1.18.3/src/runtime/stack.go#L74) which grows and shrinks as needed without the risk of ever running out.
- Go runtime does also not allow goroutines to exceed [a maximum stack size](https://github.com/golang/go/blob/go1.18.3/src/runtime/stack.go#L1094); this maximum depends on the architecture and is [1 GB for 64-bit and 250MB for 32-bit systems](https://github.com/golang/go/blob/go1.18.3/src/runtime/proc.go#L155). If this limit is reached a call to `runtime.abort` will take place.

```go
package main

func foo(i int) int {
	if i < 1e8 {
		return foo(i + 1)
	}
	return -1
}

func main() {
	foo(0)
}
// runtime: goroutine stack exceeds 1000000000-byte limit
// runtime: sp=0xc0200e0398 stack=[0xc0200e0000, 0xc0400e0000]
// fatal error: stack overflow

// runtime stack:
// runtime.throw({0x4629fa?, 0x4b5880?})
//         /usr/local/go/src/runtime/panic.go:992 +0x71
// runtime.newstack()
//         /usr/local/go/src/runtime/stack.go:1101 +0x5cc
// runtime.morestack()
//         /usr/local/go/src/runtime/asm_amd64.s:547 +0x8b
```

- How many Goroutines you can run? two concerns:
  - Memory usage
  - Slower garbage collection

```go
package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"time"
)

var n = flag.Int("n", 1e5, "Number of goroutines to create")

var ch = make(chan byte)
var counter = 0

func f() {
	counter++
	<-ch // Block this goroutine
}

func main() {
	flag.Parse()
	if *n <= 0 {
		fmt.Fprintf(os.Stderr, "invalid number of goroutines")
		os.Exit(1)
	}

	// Limit the number of spare OS threads to just 1
	runtime.GOMAXPROCS(1)

	// Make a copy of MemStats
	var m0 runtime.MemStats
	runtime.ReadMemStats(&m0)

	t0 := time.Now().UnixNano()
	for i := 0; i < *n; i++ {
		go f()
	}
	runtime.Gosched()
	t1 := time.Now().UnixNano()
	runtime.GC()

	// Make a copy of MemStats
	var m1 runtime.MemStats
	runtime.ReadMemStats(&m1)

	if counter != *n {
		fmt.Fprintf(os.Stderr, "failed to begin execution of all goroutines")
		os.Exit(1)
	}

	fmt.Printf("Number of goroutines: %d\n", *n)
	fmt.Printf("Per goroutine:\n")
	fmt.Printf("  Memory: %.2f bytes\n", float64(m1.Sys-m0.Sys)/float64(*n))
	fmt.Printf("  Time:   %f µs\n", float64(t1-t0)/float64(*n)/1e3)
}
```