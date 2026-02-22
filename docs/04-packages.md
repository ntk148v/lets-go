# Packages

> **NOTE**: Code examples in this section are stored in [`examples/4/`](../examples/4/). Commands assume you run them from that directory.

Table of Contents:

- [Packages](#packages)
  - [1. Identifiers](#1-identifiers)
  - [2. Documeting packages](#2-documeting-packages)
  - [3. Creating a package](#3-creating-a-package)
  - [4. Package initialization](#4-package-initialization)
  - [5. Program execution order](#5-program-execution-order)
  - [6. Installing 3rd party package](#6-installing-3rd-party-package)
  - [7. Testing packages](#7-testing-packages)
    - [7.1. Basic Testing](#71-basic-testing)
    - [7.2. Test Tables](#72-test-tables)
    - [7.3. HTTP Testing](#73-http-testing)
    - [7.4. Code Coverage](#74-code-coverage)
    - [7.5. Benchmarking](#75-benchmarking)
  - [8. Useful packages](#8-useful-packages)

A package is a collection of functions & data. The convention for package names is to use lowercase characters - the file does not have to match the package name.

```go
package even

func Even(i int) bool { // starts with capital -> exported
    return i%2 == 0
}

func odd(i int) bool { // start with lower-case -> private
    return i%2 == 1
}
```

Build the package:

```shell
mkdir $GOPATH/src/even
cp even.go $GOPATH/src/even
go build
go install
```

Now you can use the package in your program with `import "even"`.

## 1. Identifiers

- The Convention in Go is to use CamelCase rather than underscores to write multi-word names.
- The Convention in Go is that package names are lowercase, single word names.
- Override default package name: `import bar "bytes"`.
- Another convention is that the package name is the base name of its source directory; the package in `src/compress/gzip` is imported as `compress/gzip` but has name `gzip`, not `compress/gzip`.
- Avoid stuttering when naming things.
- The function to make new instance of `ring.Ring` package (package `container/ring`), would normally be called `NewRing`, but since `Ring` is the only type exported by the package, since the package is called `ring`, it's called just `New`. Clients of the package see that as `ring.New`.

## 2. Documeting packages

- Each package should have a \_package comment\*\*.
- When a package consists of multiple files the package comment should only appear in 1 file.
- A common convention (in really big packages) is to have a separate `doc.go` that only holds the package comment.

```go
/*
    The regexp package implements a simple library for
    regular expressions.

    The syntax of the regular expressions accepted is:

    regexp:
        concatenation { '|' concatenation }
*/
package regexp
```

- Each defined (and exported) function should have a same line of text documenting the behavior of the function.

## 3. Creating a package

- There are two types of packages: An **executable package** (main application since you will be running it) & an **utility package** (is not self executable, instead it enhances functionality of an executable package by providing utility functions & other import assets).
- Go exports a variable if a variable name starts with **Uppercase**. All other variables not starting with an uppercase letter is private to the package.
- For an executable package, a file with `main` function is entry file for execution.

## 4. Package initialization

- **Package scope - A scope is a region in code block where a defined variable is accessible**. A package scope is a region within a package where a declared variable is accessible from within a package (_across all files in the package_). This region is the top-most block of any file in the package. **You are not allowed to re-declare global variable with same name in the same package**
- **Variable initialization**.
- **Init function**: Like `main` function, `init` function is called by Go when a package is initialized. It does not take any arguments & doesn't return any value. `init`function is implicitly declared by Go. You can have multiple `init` functions in a file or a package. Order of the execution of `init` function in a file will be according to the order of their appearances.
- **Package alias**: Underscore is a special character in Go which act as `null` container.
- A main thing to remember is, **an imported package is initialized only once per package**. Hence if you have many import statements in a package, an imported package is going to be initialized only once in the lifetime of main package execution.

## 5. Program execution order

```shell
go run *.go
├── Main package is executed
├── All imported packages are initialized
|  ├── All imported packages are initialized (recursive definition)
|  ├── All global variables are initialized
|  └── init functions are called in lexical file name order
└── Main package is initialized
   ├── All global variables are initialized
   └── init functions are called in lexical file name order
```

## 6. Installing 3rd party package

Installing a 3rd party package is nothing but cloning the remote code into local `src/<package>` directory. Unfortunately, Go does not support package version or provide package manager but a proposal is waiting [**here**](https://github.com/golang/proposal/blob/master/design/24301-versioned-go.md).

## 7. Testing packages

- Writing test involves the `testing` package & the program `go test`.
- Source [Golang writing unit tests](https://blog.alexellis.io/golang-writing-unit-tests/)
- Unit testing in Go is just a opinionated as any other aspect of the language like formatting or naming.

### 7.1. Basic Testing

```go
package main

func Sum(x int, y int) int {
    return x + y
}

func main() {
    Sum(5, 5)
}

// Testcase
package main

import "testing"

func TestSum(t *testing.T) {
    total := Sum(5, 5)
    if total != 10 {
       t.Errorf("Sum was incorrect, got: %d, want: %d.", total, 10)
    }
}
```

### 7.2. Test Tables

Test tables is a set (slice array) of test input & output values:

```go
package main

import "testing"

func TestSum(t *testing.T) {
    tables := []struct {
        x int
        y int
        n int
    }{
        {1, 1, 2},
        {1, 2, 3},
        {2, 2, 4},
        {5, 2, 7},
    }

    for _, table := range tables {
        total := Sum(table.x, table.y)
        if total != table.n {
            t.Errorf("Sum of (%d+%d) was incorrect, got: %d, want: %d.", table.x, table.y, total, table.n)
        }
    }
}
```

Launching tests:

- Within the same directory as the test, _this picks up any files matching packagename_test.go_:

```shell
go test
```

- By fully-qualified package name:

```shell
go test github.com/username/package
```

### 7.3. HTTP Testing

- The `net/http/httptest` sub-package facilitates the testing automation of both HTTP server and client code.
- When writing HTTP server code, you will undoubtedly run into the need to test your code in a robust and repeatable manner, without having to set up some fragile code harness to simulate end-to-end testing. Type `httptest.ResponseRecorder` is designed specifically to provide unit testing capabilities for exercising the HTTP handler methods by inspecting state changes to the http.ResponseWriter in the tested function.
- Creating test code for an HTTP client is more involved, since you actually need a server running for proper testing. Luckily, package `httptest` provides type `httptest.Server` to programmatically create servers to test client requests and send back mock responses to the client.

### 7.4. Code Coverage

```shell
$ go test -cover
PASS
coverage: 50.0% of statements
ok      github.com/alexellis/golangbasics1    0.009s
# Generate a HTML coverage report.
$ go test -cover -coverprofile=c.out
$ go tool cover -html=c.out -o coverage.html
```

### 7.5. Benchmarking

Code benchmark: The purpose of benchmarking is to measure a code's performance. The go test command-line tool comes with support for the automated generation and measurement of benchmark metrics. Similar to unit tests, the test tool uses benchmark functions to specify what portion of the code to measure.

Running the benchmark:

```shell
$> go test -bench=.
 PASS
 BenchmarkVectorAdd-2 2000000 761 ns/op
 BenchmarkVectorSub-2 2000000 788 ns/op
 BenchmarkVectorScale-2 5000000 269 ns/op
 BenchmarkVectorMag-2 5000000 243 ns/op
 BenchmarkVectorUnit-2 3000000 507 ns/op
 BenchmarkVectorDotProd-2 3000000 549 ns/op
 BenchmarkVectorAngle-2 2000000 659 ns/op
 ok github.com/vladimirvivien/learning-go/ch12/vector 14.123s
```

Skipping test functions:

```shell
> go test -bench=. -run=NONE -v
 PASS
 BenchmarkVectorAdd-2 2000000 791 ns/op
 BenchmarkVectorSub-2 2000000 777 ns/op
 ...
 BenchmarkVectorAngle-2 2000000 653 ns/op
 ok github.com/vladimirvivien/learning-go/ch12/vector 14.069s
```

Comparative benchmarks: to compare the performance of different algorithms that implement similar functionalities. Exercising the algorithms using performance benchmarks will indicate which of the implementations may be more compute and memory efficient.

Isolating dependencies: The Key factor that defines a unit test is isolation from runtime dependencies or collaborators. Check out [Dependency Injection](../tips-notes/dependency-injection.md).

## 8. Useful packages

- **fmt**: Package `fmt` implements formatted I/O with functions analogous to C's`printf` & `scanf`. The format verbs are derived from C's but are simpler. Some verbs (%-sequences) that can be used:
  - %v, the value in a default format, when printing structs, the plus flag (%+v) adds fields names.
  - %#v, a Go-syntax representation of the value.
  - %T, a Go-sytanx representation of the type of the value.

- **io**: The package provides basic interfaces to I/O primitives. Its primary job is to wrap existing implementation of such primitives, such as those in package `os`, into shared public interfaces that abstract the functionality, plus some other related primitives.

- **bufio**: This package implements buffered I/O. It wraps an io.Reader or io.Writer object, creating another object (Reader or Writer) that also implements the interface but provides buffering & some help for textual I/O.

- **sort**: The sort package provides primitives for sorting arrays & user-defined collections.

- **strconv**: The strconv package implements conversions to & from string representations of basic data types.

- **os**: The os package provides a platform-independent interface to operating system functionality. The design is Unix-like.

- **sync**: The package sync provides basic synchronization primitives such as mutual exclusion locks.

- **flag**: The flag package implements command-line flag parsing.

- **encoding/json**: The encoding/json package implements encoding & decoding of JSON objects as defined in RFC 4627.

- **html/template**: Data-driven templates for generating textual output such as HTML.

- **net/http**: The net/http package implements parsing of HTTP requests, replies, & URLs & provides an extensible HTTP server & a basic HTTP client.

- **unsafe**: The unsafe package contains operations that step around the type safety of Go programs. Normally you don't need this package, but it is worth mentioning that unsafe Go programs are possible.

- **reflect**: The reflect package implements run-time reflection, allowing a program to manipulate objects with arbitrary types. The typical use is to take a value with static type interface{} & extract its dynamic type information by calling TypeOf, which returns an object with interface type Type.

- **os/exec**: The os/exec package runs external commands.
