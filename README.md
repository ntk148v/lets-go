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

    /* slicelengthcapacity.go */
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
