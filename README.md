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
/* helloworld.go */
package main

import "fmt" // Implements formatted I/O

/* Say Hello-World */
func main() {
    fmt.Printf("Hello World")
}
```

### 2.2. Compiling and Running Code

* To build [helloworld.go](./2/helloworld.go), just type:

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
    /* numericaltypes.go */
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
/* gototest */
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
