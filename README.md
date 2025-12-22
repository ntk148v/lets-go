<p align="center">
    <img src="./logo.png" width="20%" height="20%">
</p>
<h1 align="center">Let's Go</h1>

<p align="center">
    <a href="https://github.com/ntk148v/lets-go/blob/master/LICENSE">
        <img alt="GitHub license" src="https://img.shields.io/github/license/ntk148v/lets-go?style=for-the-badge">
    </a>
    <a href="https://github.com/ntk148v/lets-go/stargazers"> <img alt="GitHub stars" src="https://img.shields.io/github/stars/ntk148v/lets-go?style=for-the-badge"></a>
</p>

A comprehensive Go learning guide. From basics to advanced topics, covering Go 1.18+ features including generics, and the latest Go 1.24/1.25 updates.

> [!important]
> Checkout the previous one-page version [here](./README_onepage.md)

Table of Contents:
- [1. Quick Start](#1-quick-start)
- [2. Documentation](#2-documentation)
- [3. Repository Structure](#3-repository-structure)
- [4. Advanced Topics (tips-notes)](#4-advanced-topics-tips-notes)
- [5. Learning Path](#5-learning-path)
  - [5.1. For Beginners](#51-for-beginners)
  - [5.2. For Intermediate Developers](#52-for-intermediate-developers)
  - [5.3. For Advanced Users](#53-for-advanced-users)
- [6. Resources](#6-resources)
- [7. Contributing](#7-contributing)
- [8. License](#8-license)

## 1. Quick Start

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, Gopher!")
}
```

```bash
go run hello.go
```

## 2. Documentation

| Topic                                         | Description                                                | Level        |
| --------------------------------------------- | ---------------------------------------------------------- | ------------ |
| [Getting Started](docs/01-getting-started.md) | Installation, first program, essential commands            | Beginner     |
| [Basics](docs/02-basics.md)                   | Variables, types, control structures, arrays, slices, maps | Beginner     |
| [Functions](docs/03-functions.md)             | Functions, closures, defer, panic/recover                  | Beginner     |
| [Packages](docs/04-packages.md)               | Creating packages, testing, useful packages                | Beginner     |
| [Pointers](docs/05-pointers.md)               | Pointers, allocation, conversions                          | Intermediate |
| [Interfaces](docs/06-interfaces.md)           | Duck typing, methods, reflection                           | Intermediate |
| [Concurrency](docs/07-concurrency.md)         | Goroutines, channels, patterns                             | Intermediate |
| [Communication](docs/08-communication.md)     | IO, networking, command line                               | Intermediate |
| [Modules](docs/09-modules.md)                 | Go modules, workspaces, dependency management              | Intermediate |
| [Testing](docs/10-testing.md)                 | Unit tests, benchmarks, mocking                            | Intermediate |
| [Data IO](docs/11-data-io.md)                 | Readers, writers, formatted IO                             | Intermediate |
| [Encoding](docs/12-encoding.md)               | JSON marshaling, encoding                                  | Intermediate |
| [Web Programming](docs/13-web-programming.md) | HTTP servers, forms, middleware, websockets                | Intermediate |
| [RPC and gRPC](docs/14-rpc-grpc.md)           | Remote procedures, Protocol Buffers                        | Advanced     |
| [New Packages](docs/15-new-packages.md)       | unique package and more                                    | Advanced     |

## 3. Repository Structure

```shell
golang/
├── README.md           # This file (documentation hub)
├── README_onepage.md   # Complete reference documentation (one-page version)
├── docs/               # Focused topic guides
│   ├── 01-getting-started.md
│   ├── 02-basics.md
│   ├── 03-functions.md
│   ├── ...
│   └── 15-new-packages.md
├── examples/           # Code examples by section (matches README_onepage numbering)
├── gobyexample/        # Examples by topic
└── tips-notes/         # Advanced tips and patterns
```

## 4. Advanced Topics (tips-notes)

Deep-dive guides for specific topics:

| Topic                                                                                                        | Description                                      |
| ------------------------------------------------------------------------------------------------------------ | ------------------------------------------------ |
| [Concurrency Deep Dive](tips-notes/go-concurrency.md)                                                        | Goroutines vs threads, patterns, best practices  |
| [Error Handling](tips-notes/handling-errors-gracefully.md)                                                   | Error wrapping, custom errors, Go 1.13+ features |
| [Context Usage](tips-notes/context.md)                                                                       | Context best practices, propagation              |
| [Dependency Injection](tips-notes/dependency-injection.md)                                                   | Testing and mocking patterns                     |
| [Defer Patterns](tips-notes/defer/)                                                                          | Defer mechanics and gotchas                      |
| [Pipelines and Cancellation](tips-notes/pipelines-cancellation.md)                                           | Pipeline patterns with context                   |
| [Performance](tips-notes/debugging-performance-issues.md)                                                    | Debugging and optimization                       |
| [Build for Multiple OS](tips-notes/build-go-applications-for-different-operating-systems-and-architectures/) | Cross-compilation                                |
| [Shrink Binaries](tips-notes/shrink-binary.md)                                                               | Reducing binary size                             |
| [Container CPU Throttling](tips-notes/cpu-throttling-in-containerized-go-apps.md)                            | GOMAXPROCS in containers                         |

## 5. Learning Path

### 5.1. For Beginners

1. [Getting Started](docs/01-getting-started.md) - Install Go and write your first program
2. [Basics](docs/02-basics.md) - Learn variables, types, and control flow
3. [Functions](docs/03-functions.md) - Understand functions and closures
4. [Packages](docs/04-packages.md) - Organize code into packages

### 5.2. For Intermediate Developers

1. [Concurrency](docs/07-concurrency.md) - Master goroutines and channels
2. [Testing](docs/10-testing.md) - Write effective tests
3. [Web Programming](docs/13-web-programming.md) - Build web applications
4. [Error Handling](tips-notes/handling-errors-gracefully.md) - Handle errors properly

### 5.3. For Advanced Users

1. [Concurrency Patterns](tips-notes/go-concurrency.md) - Advanced patterns
2. [Performance Tuning](tips-notes/debugging-performance-issues.md) - Optimize your code
3. [gRPC and Protobuf](docs/14-rpc-grpc.md) - Build RPC services
4. [New Packages](docs/15-new-packages.md) - Latest additions to standard library

## 6. Resources

- [Official Go Documentation](https://golang.org/doc/)
- [A Tour of Go](https://tour.golang.org/)
- [Go by Example](https://gobyexample.com/)
- [Effective Go](https://golang.org/doc/effective_go.html)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Uber Go Style Guide](https://github.com/uber-go/guide)
- [100 Go Mistakes](https://100go.co/)

## 7. Contributing

Contributions are welcome! Feel free to:

- Fix typos or errors
- Add new examples
- Improve explanations
- Add coverage for new Go features

## 8. License

[Apache-2.0](LICENSE)
