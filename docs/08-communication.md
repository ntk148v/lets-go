# Communication

> **NOTE**: Code examples for this section are stored in [`examples/8/`](../examples/8/).

Table of Contents:

- [Communication](#communication)
  - [1. io.Reader](#1-ioreader)
  - [2. Command line arguments](#2-command-line-arguments)
  - [3. Executing commands](#3-executing-commands)
  - [4. Networking](#4-networking)

Building blocks in Go for communcating with the outside world (files, directories, networking & executing other programs).

Central to Go's I/O are the interfaces `io.Reader` & `io.Writer`.

## 1. io.Reader

- `io.Reader` is an important interface in the language Go. A lot (if not all) functions that need to read from something take an `io.Reader` as input.
- The writing side `io.Writer` has the `Write` method.
- If you think of new type in your program or package & you make it fulfill the `io.Reader` or `io.Writer` interface, _the whole standard Go library can be used_ on that type.

## 2. Command line arguments

- Arguments from the command line are available inside your program via the string slide `os.Args`.
- The `flag` package has a more sophisticated interface, & also provided a way to parse flags.

## 3. Executing commands

- The `os/exec` package has functions to run external commands, & is the premier way to execute commands from within a Go program.

```go
import "os/exec"

cmd := exec.Command("/bin/ls", "-l")
// Just run without doing anything with the returned data
err := cmd.Run()
// Capturing the standard output
buf, err := cmd.Output() // buf is byte slice
```

## 4. Networking

- All network related types & functions can be found in the package `net`.
- One of the most important functions in there is `Dial`. When you `Dial` into a remote system the function returns a `Conn` interface type, which can be used to send & receive information. The function `Dial` neatly abstracts away the network family & transport.

```go
conn, e := Dial("tcp", "192.0.32.10:80")
conn, e := Dial("udp", "192.0.32.10:80")
conn, e := Dial("tcp", "[2620:0:2d0:200::10]:80")
```
