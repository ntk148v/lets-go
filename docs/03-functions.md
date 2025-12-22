# Functions

> **NOTE**: Code examples in this section are stored in [`examples/3/`](../examples/3/). Commands assume you run them from that directory.

Table of Contents:

- [Functions](#functions)
  - [1. Scope](#1-scope)
  - [2. Functions as values](#2-functions-as-values)
  - [3. Callbacks](#3-callbacks)
  - [4. Deferred Code](#4-deferred-code)
  - [5. Variadic Parameter](#5-variadic-parameter)
  - [6. Panic and recovering](#6-panic-and-recovering)

```go
// General Function
type mytype int
func (p mytype) funcname(q, int) (r, s int) { return 0,0 }
// p (optional) bind to a specific type called receiver (a function with a receiver is usually called a method)
// q - input parameter
// r,s - return parameters
```

- Functions can be declared in any order you wish.
- Go does not allow nested functions, but you can work around this with anonymous functions.

## 1. Scope

- Variables declared outside any functions are **global** in Go, those defined in functions are **local** to those functions.
- If names overlap - a local variable is decleard with the same name as a global one - the local variable hides the global one when the current function is executed.

## 2. Functions as values

- As with almost everything in Go, functions are also just values.

```go
import "fmt"

func main() {
    a := func() { // a is defined as an anonymous (nameless) function,
        fmt.Println("Hello")
    }
    a()
}
```

## 3. Callbacks

```go
func printit(x int) {
    fmt.Println("%v\n", x)
}

func callback(y int, f func(int)) {
    f(y)
}
```

## 4. Deferred Code

```go
/* Open a file & perform various writes & reads on it. */
func ReadWrite() bool {
    file.Open("file")
    // Do your thing
    if failureX {
        file.Close()
        return false
    }

    // Repeat a lot of code.
    if failureY {
        file.Close()
        return false
    }
    file.Close()
    return true
}

/* Same situation but using defer */
func ReadWrite() bool {
    file.Open("file")
    defer file.Close() // add file.Close() to the defer list
    // Do your thing
    if failureX {
        return false
    }

    if failureY {
        return false
    }
    return true
}
```

- Can put multiple functions on the "defer list".
- `Defer` functions are executed in _LIFO_ order.

```go
for i := 0; i < 5; i++ {
    defer fmt.Printf("%d ", i) // 4 3 2 1 0
}
```

- With `defer` you can even change return values, provided that you are using named result parameters & a function literal (`def func(x int) {/*....*/}(5)`).

```go
func f() (ret int)
    defer func() { // Initialized with zero
        ret++
    }()
    return 0 // This will not be the returned value, because of defer. Ths function f will return 1
}
```

## 5. Variadic Parameter

- Functions that take a variable number of parameters are known as variadic functions.

```go
func func1(arg... int) { // the variadic parameter is just a slice.
    for _, n := range arg {
        fmt.Printf("And the number is: %d\n", n)
    }
}
```

## 6. Panic and recovering

- Go does not have an exception mechanism: _you can not throw exception_. Instead it uses a _panic & recover mechanism_.
  - Panic: Built-in function that tstops the oridinary flow of control & begins panicking. When function F call `pacnic`, execution of `F` stops, any deferred functions in F are executed normally, & then F returns to its caller. To the caller, F then behaves like a call to panic. The process continues up the stack until all functions in the current goroutine have returned, at which point the program crashes. Panics can be initiated by invoking panic directly. They can also be caused by runtime errors, such as out-of-bounds array accesses.
  - Recover: Built-in function that regains control of a panicking goroutine. Recover is only useful inside deferred functions. During normal execution, a call to recover will return nil & have no other effect. If the current goroutine is panicking, a call to recover will capture the value given to panic & resume normal execution.

```go
/* defer_panic_recover.go */
package main

import "fmt"

func main() {
    f()
    fmt.Println("Returned normally from f.")
}

func f() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("Recovered in f", r)
        }
    }()
    fmt.Println("Calling g.")
    g(0)
    fmt.Println("Returned normally from g.")
}

func g(i int) {
    if i > 3 {
        fmt.Println("Panicking!")
        panic(fmt.Sprintf("%v", i))
    }
    defer fmt.Println("Defer in g", i)
    fmt.Println("Printing in g", i)
    g(i + 1)
}
/* Result */
// Calling g.
// Printing in g 0
// Printing in g 1
// Printing in g 2
// Printing in g 3
// Panicking!
// Defer in g 3
// Defer in g 2
// Defer in g 1
// Defer in g 0
// Recovered in f 4
// Returned normally from f.
```

- Still don't understanding how these works? Don't worry, I got you. Check [Go Defer Simplified with Praticial Visuals by Inanc Gunmus](https://blog.learngoprogramming.com/golang-defer-simplified-77d3b2b817ff).
- Other useful links about Defer:
  - [5 Gotchas of Defer in Go — Part I](https://blog.learngoprogramming.com/gotchas-of-defer-in-go-1-8d070894cb01)
  - [5 Gotchas of Defer in Go — Part II](https://blog.learngoprogramming.com/5-gotchas-of-defer-in-go-golang-part-ii-cc550f6ad9aa)
- Check out [use cases](../tips-notes/panic-and-recover-use-cases.md)
- To handle panics gracefully, check [handle tips](../tips-notes/handling-panics-gracefully.md).
