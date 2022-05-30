# Panic and Recover Uses cases

Source: <https://go101.org/article/panic-and-recover-use-cases.html>

- [Panic and Recover Uses cases](#panic-and-recover-uses-cases)
  - [1. Avoid Panics crashing programs](#1-avoid-panics-crashing-programs)
  - [2. Automatically restart a crashed goroutine](#2-automatically-restart-a-crashed-goroutine)
  - [3. Use `panic/recover` calls to simulate long jump statements](#3-use-panicrecover-calls-to-simulate-long-jump-statements)
  - [4. Use `panic/recover` calls to reduce errors checks](#4-use-panicrecover-calls-to-reduce-errors-checks)

## 1. Avoid Panics crashing programs

- Used commonly in concurrent programs, especially client-server programs.
- Example:

```go
package main

import "errors"
import "log"
import "net"

func main() {
    listener, err := net.Listen("tcp", ":12345")
    if err != nil {
        log.Fatalln(err)
    }
    for {
        conn, err := listener.Accept()
        if err != nil {
            log.Println(err)
        }
        // Handle each client connection
        // in a new goroutine.
        go ClientHandler(conn)
    }
}

func ClientHandler(c net.Conn) {
    defer func() {
        if v := recover(); v != nil {
            log.Println("capture a panic:", v)
            log.Println("avoid crashing the program")
        }
        c.Close()
    }()
    panic(errors.New("just a demo.")) // a demo-purpose panic
}
// Start the server
// Telnet localhost 12345 -> Server will not crash down for the panics created in each client handler goroutine.
```

## 2. Automatically restart a crashed goroutine

- When a panic is detected in a goroutine, we can create a new goroutine for it.

```go
package main

import "log"
import "time"

func shouldNotExit() {
    for {
        // Simulate a workload.
        time.Sleep(time.Second)

        // Simulate an unexpected panic.
        if time.Now().UnixNano() & 0x3 == 0 {
            panic("unexpected situation")
        }
    }
}

func NeverExit(name string, f func()) {
    defer func() {
        if v := recover(); v != nil {
            // A panic is detected.
            log.Println(name, "is crashed. Restart it now.")
            go NeverExit(name, f) // restart
        }
    }()
    f()
}

func main() {
    log.SetFlags(0)
    go NeverExit("job#A", shouldNotExit)
    go NeverExit("job#B", shouldNotExit)
    select{} // block here for ever
}
```

## 3. Use `panic/recover` calls to simulate long jump statements

- We can use `panic/recover` as a way to simulate crossing-function long jump statements and crossing-function returns
- Not recommended! This way does harm for both code readability and execution effiency.

```go
package main

import "fmt"

func main() {
    n := func () (result int)  {
        defer func() {
            if v := recover(); v != nil {
                if n, ok := v.(int); ok {
                    result = n
                }
            }
        }()

        func () {
            func () {
                func () {
                    // ...
                    panic(123) // panic on succeeded
                }()
                // ...
            }()
        }()
        // ...
        return 0
    }()
    fmt.Println(n) // 123
}
```

## 4. Use `panic/recover` calls to reduce errors checks

- Code will be less verbose, but it's not recommended.

```go
func doSomething() (err error) {
    defer func() {
        err = recover()
    }()

    doStep1()
    doStep2()
    doStep3()
    doStep4()
    doStep5()

    return
}

// In reality, the prototypes of the doStepN functions
// might be different. For each of them,
// * panic with nil for success and no needs to continue.
// * panic with error for failure and no needs to continue.
// * not panic for continuing.
func doStepN() {
    ...
    if err != nil {
        panic(err)
    }
    ...
    if done {
        panic(nil)
    }
}

func doSomething() (err error) {
    shouldContinue, err := doStep1()
    if !shouldContinue {
        return err
    }
    shouldContinue, err = doStep2()
    if !shouldContinue {
        return err
    }
    shouldContinue, err = doStep3()
    if !shouldContinue {
        return err
    }
    shouldContinue, err = doStep4()
    if !shouldContinue {
        return err
    }
    shouldContinue, err = doStep5()
    if !shouldContinue {
        return err
    }

    return
}

// If err is not nil, then shouldContinue must be true.
// If shouldContinue is true, err might be nil or non-nil.
func doStepN() (shouldContinue bool, err error) {
    ...
    if err != nil {
        return false, err
    }
    ...
    if done {
        return false, nil
    }
    return true, nil
}
```
