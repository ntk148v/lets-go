# Getting Started with Go

Table of Contents:

- [Getting Started with Go](#getting-started-with-go)
  - [1. What is Go?](#1-what-is-go)
  - [2. Installing Go](#2-installing-go)
    - [2.1. Linux](#21-linux)
    - [2.2. macOS](#22-macos)
    - [2.3. Windows](#23-windows)
    - [2.4. Verify Installation](#24-verify-installation)
  - [3. Your First Go Program](#3-your-first-go-program)
  - [4. Setting Up Your Workspace](#4-setting-up-your-workspace)
  - [5. Essential Go Commands](#5-essential-go-commands)
  - [6. Next Steps](#6-next-steps)
  - [7. Resources](#7-resources)

## 1. What is Go?

Go (also known as Golang) is an open-source programming language developed by Google. Key characteristics:

- **Statically typed** - Types are checked at compile time
- **Compiled** - Compiles directly to native machine code (no JVM)
- **Concurrent** - Built-in primitives for concurrent programming (goroutines & channels)
- **Simple** - Clean syntax with minimal keywords
- **Fast** - Quick compilation and efficient execution

## 2. Installing Go

### 2.1. Linux

```bash
# Ubuntu/Debian
sudo apt update
sudo apt install golang-go

# Or download from https://go.dev/dl/
wget https://go.dev/dl/go1.25.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.25.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
```

### 2.2. macOS

```bash
# Using Homebrew
brew install go

# Or download from https://go.dev/dl/
```

### 2.3. Windows

Download the installer from [go.dev/dl](https://go.dev/dl/) and follow the installation wizard.

### 2.4. Verify Installation

```bash
go version
# Output: go version go1.25 linux/amd64
```

## 3. Your First Go Program

Create a file named `hello.go`:

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
```

Run it:

```bash
go run hello.go
# Output: Hello, World!
```

Or build and run separately:

```bash
go build hello.go  # Creates executable
./hello           # Run the executable
```

## 4. Setting Up Your Workspace

Go uses modules for dependency management (since Go 1.11+):

```bash
# Create a new project
mkdir myproject
cd myproject

# Initialize a new module
go mod init github.com/username/myproject

# Your project structure:
# myproject/
# ├── go.mod
# └── main.go
```

## 5. Essential Go Commands

| Command            | Description                           |
| ------------------ | ------------------------------------- |
| `go run <file.go>` | Compile and run a Go program          |
| `go build`         | Compile packages and dependencies     |
| `go test`          | Run tests                             |
| `go mod init`      | Initialize a new module               |
| `go mod tidy`      | Add missing and remove unused modules |
| `go get <package>` | Download and install packages         |
| `go fmt`           | Format source code                    |
| `go vet`           | Report likely mistakes in packages    |
| `go doc`           | Show documentation for a package      |

> **NOTE**: Go 1.25 introduces `go doc -http` which launches a local documentation server.

## 6. Next Steps

- [Basics](02-basics.md) - Variables, types, control structures
- [Functions](03-functions.md) - Functions, methods, and closures
- [Packages](04-packages.md) - Creating and using packages

## 7. Resources

- [Official Go Documentation](https://golang.org/doc/)
- [A Tour of Go](https://tour.golang.org/)
- [Go by Example](https://gobyexample.com/)
- [Effective Go](https://golang.org/doc/effective_go.html)
