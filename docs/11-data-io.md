# Data IO in Go

Table of Contents:

- [Data IO in Go](#data-io-in-go)
  - [1. IO with readers and writers](#1-io-with-readers-and-writers)
  - [2. Formatted IO with fmts](#2-formatted-io-with-fmts)
  - [3. Buffered IO](#3-buffered-io)
  - [4. In-memory IO](#4-in-memory-io)

## 1. IO with readers and writers

Go models data input and output as a stream that flows from sources to targets. Data sources, such as files, network connections, or even some in-memory objects, can be modeled as streams of bytes from which data can be read or written to.

## 2. Formatted IO with fmts

The most common usage of the `fmt` package is for writting to standard output and reading from standard input.

```go
type metalloid struct {
    name string
    number int32
    weight float64
}

func main() {
    var metalloids = []metalloid{
        {"Boron", 5, 10.81},
         ...
         {"Polonium", 84, 209.0},
    }
     file, _ := os.Create("./metalloids.txt")
     defer file.Close()
     for _, m := range metalloids {
        fmt.Fprintf(
            file,
            "%-10s %-10d %-10.3f\n",
            m.name, m.number, m.weight,
        )
     }
}
```

## 3. Buffered IO

The `bufio` package offers several functions to do buffered writing of IO streams using an `io.Writer interface.

## 4. In-memory IO

In `bytes` package offers common primitives to achieve streaming IO on blocks of bytes stored in memory, represented by the `bytes.Buffer` byte. Since the `bytes.Buffer` type implements both `io.Reader` and `io.Writer` interfaces it is a great option to stream data into or out of memory using streaming IO primitives.
