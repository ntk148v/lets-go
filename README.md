# Let's Go!

## Learn

* [Official Golang Documetation](https://golang.org/doc/#learning)
* [Community-driven initiatives](https://github.com/golang/go/wiki/Learn)

These above links contain a lot of sources, so finding documentation isn't a big deal.. Just choose one and let's start. In my case, I follow [Learning Go - Miek Gieben](https://miek.nl/go/).

> **NOTE**: Every examples in this documentation are stored in directories named by section. I assume that every commands in section X will be executed **in X directory**, so I don't write a full path to Go script file.

## 1. Introduction

> Is Go an object-oriented language? Yes and no!
>                   [FAQ - Golang documentation](https://golang.org/doc/faq#Is_Go_an_object-oriented_language)

* What is Golang? Golang Programming language is an open source programming language that makes it easy to build simple, reliable and efficient software.
* Go is a compliaed statically typed language that feels like a dynamically typed, interpreted language.

## 2. Basic

### 2.1. Say Hello World in Golang!

* Get started with Go in the classic way: printing "Hello World" (Ken Thompson and Dennies Ritchie started this when they presented the C language in the 1970s #til)

```Go
/* hello_world.go */
package main

import "fmt" // Implements formatted I/O

/* Say Hello-World */
func main() {
    fmt.Printf("Hello World")
}
```

### 2.2. Compiling and Running Code

* To build [helloworld.go](./2/hello_world.go), just type:

```
$ go build helloworld.go # Return an executable called helloworld 
```

* Run a previous step result

```
$ ./helloworld
```

* Want to ombine these two steps? Ok, Golang got you.

```
$ go run helloworld.go
```

### 2.4. Variables, Types and Keywords

* Go is different from most other language in that type of a variable is specified *after* the variable name: ~~int a~~ a int.

```Go
/* When you declare a variable it is assigned the "natural" null value for the type */
var a int // a has a value of 0
var s string // s is assigned the zero string, which is ""
a = 26
s = "hello"

/* Declaring and assigning in Go is a two step process, but they may be combined */
a := 26 // In this case the variable type is deduced from the value. A value of 26 indicates an int for example.
b := "hello" // The type should be string

/* Multiple var declarations may also be grouped (import and const also allow this) */
var (
    a int
    b string
)

/* Multiple variables of the same type ca also be declared on a single line */
var a, b int
a, b := 26, 9

/* A special name for a variable is _, any value assigned to it is discarded. */
_, b := 26, 9
```

* Boolean Types: `bool`
* Numerical Types:
    * Go has most of the well-know types such as `int` - it has the appropriate length for your machine (32-bit machine - 32 bits, 64-bit machine - 64 bits)
    * The full list for (signed and unsigned) integers is `int8`, `int16`, `int32`, `int64` and `byte` (an alias for `uint8`), `uint8`, `uint16`, `uint32`, `uint64`.
    * For floating point values there is `float32`, `float64`, ~~float~~.

    ```Go
    /* numerical_types.go */
    package main

    func main() {
        var a int
        var b int32
        b = a + a // Give an error: cannot use a + a (type int) as type int32 in assignment.
        b = b + 5
    }
    ```

* Constants: Constants are created at compile time, and can only be numbers, strings, or booleans. You can use `iota` to enumerate values.

```Go
const (
    a = iota // First use of iota will yield 0. Whenever iota is used again on a new line its value is incremented with 1, so b has a vaue of 1.
    b
)
```

* Strings:
    * Strings in Go are a sequence of UTF-8 characters enclosed in double quotes. If you use the single quote you mean one character (encoded in UTF-8) - which is *not* a `string` in Go. Note that! In Python (my favourite programming language), I can use both of them for string assignment.
    * String in Go are immutable. To change one character in string, you have to create a new one.

    ```Go
    s1 := "Hello"
    c := []rune(s) // Convert s1 to an array of runes
    c[0] := 'M'
    s2 := string(c) // Create a new string s2 with the alteration
    fmt.Printf("%s\n", s2)
    ```

* Rune: `Rune` is an alias for `int32`, (use when you're iterating over characters in a string).
* Complex Numbers: `complex128` (64 bit real and imaginary parts) or `complex32`.
* Errors: Go has a builtin type specially for errors, called `error.var e`.

### 2.5. Operators and Built-in Functions

* Go supports the normal set of numerical operators.

```
Precedence	Operator(s)
Highest   	* / % << >> & &^
            `+ -
            == != < <= > >=
            <-
            &&
Lowest    	||
```

* `&` bitwise and, `|` bitwise or, `^` bitwise xor, `&^` bit clear respectively.

### 2.6. Go Keywords

```
break     default      func    interface   select
case      defer        go      map         struct
chan      else         goto    package     switch
const     fallthrough  if      range       type
continue  for          import  return      var
```

* `var`, `const`, `package`, `import` are used in the previous sections.
* `func` is used to declare functions and methods.
* `return` is used to return from functions.
* `go` is used for concurrency.
* `select` is used to choose from different types of communication.
* `interface`.
* `struct` is used for abstract data types.
* `type`.

### 2.7. Control Structures

* If-Else

```Go
if x > 0 {
    return y
} esle {
    return x
}

if err := MagicFunction(); err != nil {
    return err
}

// do something
```

* Goto: With `goto` you jump to a label which must be defined within the current function.

```Go
/* goto_test */
/* Create a loop */
func gototestfunc() {
    i := 0
Here:
    fmt.Println()
    i++
    goto Here
}
```

* For: `for` loop has three forms, only one of which has semicolons:

```Go
for init; condition; post { } // aloop using the syntax borrowed from C
for condition { } // a while loop
for { } // a endless loop

sum := 0
for i := 0; i < 10; i++ {
    sum = sum + i
}
```

* Break and continue

```Go
for i := 0; i < 10; i++ {
    if i > 5 {
        break
    }
    fmt.Println(i)
}

/* With loops within loop you can specify a label after `break` to identify which loop to stop */
J: for j := 0; j < 5; j++ {
    for i := 0; i < 10; i++ {
        if i > 5 {
            break J
        }
        fmt.Println(i)
    }
}
```

* Range:
    * `range` can be used for loops. It can loop over slices, arrays, strings, maps and channels.
    * `range` is an iterator that, when called, returns the next key-value pair from the "thing" it loops over.

    ```Go
    list := []string{"a", "b", "c", "d", "e", "f"}
    for k, v := range list {
        // do some fancy thing with k and v
    }
    ```

* Switch:
    * The case are evaluated top to bottom until a match is found, and if the `switch` has no expression it switches on `true`.
    * It's therefore possible - and idomatic - to write an `if-else-if-else` chain as a `switch`.

    ```Go
    /* Convert hexadecimal character to an int value */
    switch { // switch without condition = switch true
        case '0' <= c && c <= '9':
            return c - '0'
        case 'a' <= c && c <= 'f':
            return c - 'a' + 10
        case 'A' <= c && c <= 'F':
            return c - 'A' + 10
    }
    return 0

    /* Automatic fall through */
    switch i {
        case 0: fallthrough
        case 1:
            f()
        default:
            g()
    }
    ```

### 2.8. Built-in functions

```
closei 	 new    	panic   	complex
delete 	 make   	recover 	real
len    	 append 	print   	imag
cap    	 copy   	println 	
```

    * close: is used in channel communication. It closes a channel (obviously XD)
    * delete: is used for deleting entries in maps.
    * len and cap: are used on a number of different types, `len` is used to return the lengths of strrings, maps, slices and arrays.
    * new: is used for allocating memory for user defined data types.
    * copy, append: `copy` is for copying slices. And `append` is for concatenating slices.
    * panic, recover: are used for an exception mechanism.
    * complex, real, imag: all deal with complex numbers.

### 2.9. Arrays, Slices and Maps

* Brief: list -> arrays, slices. dict -> map
* Arrays:
    * An array is defined by `[n]<type>`.

    ```Go
    var arr [10]int // The size is part of the type, fixed size
    arr[0] = 42
    arr[1] = 13
    fmt.Printf("The first element is %s\n", arr[0])

    // Initialize an array to something other than zero, using composuite literal
    a := [3]int{1, 2, 3}
    a := [...]int{1, 2, 3}
    ```

    * Array are **value types**: Assigning one array to another copies all the elements. In particular, if you pass an array to a function it will receive a copy of the array, not a pointer to it.

* Slices:
    * Similar to an array, but it can grow when new elements are added.
    * A slice is a pointer to an (underlaying) array, slices are **reference types**.

    ```Go
    // Init array primes
    primes := [6]int{2, 3, 5, 7, 11, 13}

    // Init slice s
    var s []int = primes[1:4]

    fmt.Println(s) // Return [3, 5, 7]

    /* slice_length_capacity.go */
    package main

    import "fmt"

    func main() {
        s := []int{2, 3, 5, 7, 11, 13}
        printSlice(s)

        // Slice the slice to give it zero length.
        s = s[:0]
        printSlice(s)

        // Extend its length.
        s = s[:4]
        printSlice(s)

        // Drop its first two values.
        s = s[2:]
        printSlice(s)
    }

    func printSlice(s []int) {
        fmt.Printf("len=%d cap=%d %v\n", len(s), cap(s), s)
    }
    ```

    * `len` is the number of elements it contains.
    * `cap` is the number of elements in the underlyting array, counting from the 1st element in the slice.

    ```Go
    var array [m]int
    slice := array[:n]
    // len(slice) == n
    // cap(slice) == m
    // len(array) == cap(array) == m
    ```

    * To extend a slice, there are a couple of built-in functions that make life easier: `append` and `copy`.

    ```Go
    s0 := []int{0, 0}
    s1 := append(s0, 2) // same type as s0 - int.
    // If the original slice isn't big enough to fit the added values,
    // append will allocate a new slice that is big enough. So the slice
    // returned by append may refer to a different underlaying array than
    // the original slices does.
    s2 := append(s1, 3, 5, 7)
    s3 := append(s2, s0...) // []int{0, 0, 2, 3, 5, 7, 0, 0} - three dots used after s0 is needed make it clear explicit that you're appending another slice, instead of a single value

    var a = [...]int{0, 1, 2, 3, 4, 5, 6, 7}
    var s = make([]int, 6)
    // copy function copies slice elements from source to a destination
    // returns the number of elements it copied
    n1 := copy(s, a[0:]) // n1 = 6; s := []int{0, 1, 2, 3, 4, 5}
    n2 := copy(s, s[2:]) // n2 = 4; s := []int{2, 3, 4, 5, 4, 5}
    ```

* Maps:
    * Python has its dictionaries. In go we have the `map` type.

    ```Go
    monthdays := map[string]int{
        "Jan": 31, "Feb": 28, "Mar": 31,
        "Apr": 30, "May": 31, "Jun": 30,
        "Jul": 31, "Aug": 31, "Sep": 30,
        "Oct": 31, "Nov": 30, "Dec": 31, // A trailing comma is required
    }

    value, key := monthdays["Jan"]
    ```

    * Use `make` when only declaring a map. A map is **reference type**.

## 3. Functions

```Go
// General Function
type mytype int
func (p mytype) funcname(q, int) (r, s int) { return 0,0 }
// p (optional) bind to a specific type called receiver (a function with a receiver is usually called a method)
// q - input parameter
// r,s - return parameters
```

* Functions can be declared in any order you wish.
* Go does not allow nested functions, but you can work around this with anonymous functions.

### 3.1. Scope

* Variables declared outside any functions are **global** in Go, those defined in functions are **local** to those functions.
* If names overlap - a local variable is decleard with the same name as a global one - the local variable hides the global one when the current function is executed.

### 3.2. Functions as values

* As with almost everything in Go, functions are also just values.

```Go
import "fmt"

func main() {
    a := func() { // a is defined as an anonymous (nameless) function,
        fmt.Println("Hello")
    }
    a()
}
```

### 3.3. Callbacks

```Go
func printit(x int) {
    fmt.Println("%v\n", x)
}

func callback(y int, f func(int)) {
    f(y)
}
```

### 3.4. Deferred Code

```Go
/* Open a file and perform various writes and reads on it. */
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

* Can put multiple functions on the "defer list".
* `Defer` functions are executed in *LIFO* order.

```Go
for i := 0; i < 5; i++ {
    defer fmt.Printf("%d ", i) // 4 3 2 1 0
}
```

* With `defer` you can even change return values, provided that you are using named result parameters and a function literal (`def func(x int) {/*....*/}(5)`).

```Go
func f() (ret int)
    defer func() { // Initialized with zero
        ret++
    }()
    return 0 // This will not be the returned value, because of defer. Ths function f will return 1
}
```

### 3.5. Variadic Parameter

* Functions that take a variable number of parameters are known as variadic functions.

```Go
func func1(arg... int) { // the variadic parameter is just a slice.
    for _, n := range arg {
        fmt.Printf("And the number is: %d\n", n)
    }
}
```

### 3.6. Panic and recovering

* Go does not have an exception mechanism: *you can not throw exception*. Instead it uses a *panic and recover mechanism*.
    * Panic: Built-in function that tstops the oridinary flow of control and begins panicking. When function F call `pacnic`, execution of `F` stops, any deferred functions in F are executed normally, and then F returns to its caller. To the caller, F then behaves like a call to panic. The process continues up the stack until all functions in the current goroutine have returned, at which point the program crashes. Panics can be initiated by invoking panic directly. They can also be caused by runtime errors, such as out-of-bounds array accesses.
    * Recover: Built-in function that regains control of a panicking goroutine. Recover is only useful inside deferred functions. During normal execution, a call to recover will return nil and have no other effect. If the current goroutine is panicking, a call to recover will capture the value given to panic and resume normal execution.


```Go
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

* Still don't understanding how these works? Don't worry, I got you. Check [Go Defer Simplified with Praticial Visuals by Inanc Gunmus](https://blog.learngoprogramming.com/golang-defer-simplified-77d3b2b817ff).
* Other useful links about Defer:
    * [5 Gotchas of Defer in Go — Part I](https://blog.learngoprogramming.com/gotchas-of-defer-in-go-1-8d070894cb01)
    * [5 Gotchas of Defer in Go — Part II](https://blog.learngoprogramming.com/5-gotchas-of-defer-in-go-golang-part-ii-cc550f6ad9aa)

## 4. Packages

* A package is a collection of functions and data.
* The convention for package names is to use lowercase characters - the file does not have to match the package name.

```Go
package even

func Even(i int) bool { // starts with capital -> exported
    return i%2 == 0
}

func odd(i int) bool { // start with lower-case -> private
    return i%2 == 1
}
```

* Build the package

```
$ mkdir $GOPATH/src/even
$ cp even.go $GOPATH/src/even
$ go build
$ go install
```

* Now you can use the package in your program with `import "even"`.

### 4.1. Identifiers

* The Convention in Go is to use CamelCase rather than underscores to write multi-word names.
* The Convention in Go is that package names are lowercase, single word names.
* Override default package name: `import bar "bytes"`.
* Another convention is that the package name is the base name of its source directory; the package in `src/compress/gzip` is imported as `compress/gzip` but has name `gzip`, not `compress/gzip`.
* Avoid stuttering when naming things.
* The function to make new instance of `ring.Ring` package (package `container/ring`), would normally be called `NewRing`, but since `Ring` is the only type exported by the package, since the package is called `ring`, it's called just `New`. Clients of the package see that as `ring.New`.

### 4.2. Documeting packages

* Each package should have a *package comment**.*
* When a package consists of multiple files the package comment should only appear in 1 file.
* A common convention (in really big packages) is to have a separate `doc.go` that only holds the package comment.

```Go
/*
    The regexp package implements a simple library for
    regular expressions.

    The syntax of the regular expressions accepted is:

    regexp:
        concatenation { '|' concatenation }
*/
package regexp
```

* Each defined (and exported) function should have a samll line of text documenting the behavior of the function.

### 4.3. Testing packages

* Writing test involves the `testing` package and the program `go test`.
* Fill this section later...

### 4.4. Useful packages

* **fmt**: Package `fmt` implements formatted I/O with functions analogous to C's`printf` and `scanf`. The format verbs are derived from C's but are simpler. Some verbs (%-sequences) that can be used:
    * %v, the value in a default format, when printing structs, the plus flag (%+v) adds fields names.
    * %#v, a Go-syntax representation of the value.
    * %T, a Go-sytanx representation of the type of the value.

* **io**: The package provides basic interfaces to I/O primitives. Its primary job is to wrap existing implementation of such primitives, such as those in package `os`, into shared public interfaces that abstract the functionality, plus some other related primitives.

* **bufio**: This package implements buffered I/O. It wraps an io.Reader or io.Writer object, creating another object (Reader or Writer) that also implements the interface but provides buffering and some help for textual I/O.

* **sort**: The sort package provides primitives for sorting arrays and user-defined collections.

* **strconv**: The strconv package implements conversions to and from string representations of basic data types.

* **os**: The os package provides a platform-independent interface to operating system functionality. The design is Unix-like.

* **sync**: The package sync provides basic synchronization primitives such as mutual exclusion locks.

* **flag**: The flag package implements command-line flag parsing.

* **encoding/json**: The encoding/json package implements encoding and decoding of JSON objects as defined in RFC 4627.

* **html/template**: Data-driven templates for generating textual output such as HTML.

* **net/http**: The net/http package implements parsing of HTTP requests, replies, and URLs and provides an extensible HTTP server and a basic HTTP client.

* **unsafe**: The unsafe package contains operations that step around the type safety of Go programs. Normally you don’t need this package, but it is worth mentioning that unsafe Go programs are possible.

* **reflect**: The reflect package implements run-time reflection, allowing a program to manipulate objects with arbitrary types. The typical use is to take a value with static type interface{} and extract its dynamic type information by calling TypeOf, which returns an object with interface type Type.

* **os/exec**: The os/exec package runs external commands.

## 5. Beyond the basics
