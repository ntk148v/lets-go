# Handle panic gracefully

Source: <https://github.com/sourcegraph/conc#goal-2-handle-panics-gracefully>

A frequent problem with goroutines in long-running applications is handling panics. A goroutine spawned without a panic handler will crash the whole process on panic. This is usually undesirable.

We will use the following example, just a simple one with panic, goroutine and WaitGroup.

```go
package main

import (
    "errors"
    "sync"
)

func somethingThatPanic() {
    panic(errors.New("error in somethingThatPanic function")) // a demo-purpose panic
}

func main() {
    var wg sync.WaitGroup
    wg.Add(1)

    go func() {
        defer wg.Done()
        somethingThatPanic()
    }()

    wg.Wait()
}
```

Assume that you do add a panic handler to a goroutine, you can implement it with [defer, recover](https://go.dev/blog/defer-panic-and-recover). The question is 'what do you do with the panic once you catch it?'. There are some options:

1. Ignore it: This is a bad idea since panics usually mean there is actually something wrong and someone should fix it.

```go
package main

import (
    "errors"
    "fmt"
    "sync"
)

func somethingThatPanic() {
    panic(errors.New("error in somethingThatPanic function")) // a demo-purpose panic
}

func main() {
    var wg sync.WaitGroup
    wg.Add(1)

    go func() {
        defer func() {
            wg.Done()
            if v := recover(); v != nil {
                // do nothing
                fmt.Println("Avoid crashing the program")
            }
        }()
        somethingThatPanic()
    }()

    wg.Wait()
}
```

2. Log it: Logging panics isn't great either because then there is no indication to the spawner that something bad happened, and it might just continue on as normal even though your program is in a really bad state.

```go
package main

import (
    "errors"
    "fmt"
    "sync"
)

func somethingThatPanic() {
    panic(errors.New("error in somethingThatPanic function")) // a demo-purpose panic
}

func main() {
    var wg sync.WaitGroup
    wg.Add(1)

    go func() {
        defer func() {
            wg.Done()
            if v := recover(); v != nil {
                // print the error
                fmt.Println("Capture a panic", v)
                fmt.Println("Avoid crashing the program")
            }
        }()
        somethingThatPanic()
    }()

    wg.Wait()
}
```

3. Turn it into an error and return that to the goroutine spawner: This option sounds reasonable, but it requires the goroutine to have an owner that can actually receive the message that something went wrong.

4. Propagate the panic to the goroutine spawner: similar to (3). Let's go through an example:

```go
package main

import (
    "errors"
    "fmt"
    "runtime/debug"
)

type caughtPanicError struct {
    val   interface{}
    stack []byte
}

func (e caughtPanicError) Error() string {
    return fmt.Sprintf("panic: %q\n%s", e.val, string(e.stack))
}

func somethingThatPanic() {
    panic(errors.New("error in somethingThatPanic function")) // a demo-purpose panic
}

func main() {
    done := make(chan error)
    go func() {
        // Turn goroutine's panic to an error and return to the goroutine spawner
        // through channel
        defer func() {
            if v := recover(); v != nil {
                done <- caughtPanicError{
                    val:   v,
                    stack: debug.Stack(),
                }
            } else {
                done <- nil
            }
        }()
        somethingThatPanic()
    }()

    err := <-done
    // Propagate the panic to the goroutine spawner
    if err != nil {
        panic(err)
    }
}
```
