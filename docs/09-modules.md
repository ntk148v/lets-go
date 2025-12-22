# Modules

Table of Contents:

- [Modules](#modules)
  - [1. Concepts](#1-concepts)
  - [2. Quickstart](#2-quickstart)
  - [3. Go Proxy](#3-go-proxy)
  - [4. Workspaces](#4-workspaces)
  - [5. Organizing a Go module](#5-organizing-a-go-module)
    - [5.1. Basic package](#51-basic-package)
    - [5.2. Basic command](#52-basic-command)
    - [5.3. Package with internal packages](#53-package-with-internal-packages)
    - [5.4. Multiple packages](#54-multiple-packages)
    - [5.5. Server project](#55-server-project)

Go 1.11 includes preliminary support for [modules](https://go.dev/doc/go1.11#modules), Go's [new dependency management system](https://blog.golang.org/versioning-proposal) that makes dependency version information explicit and easier to manage.

## 1. Concepts

- **Modules**: a collection of related Go packages that are versioned together as a single unit.
- Summarizing the relationship between repositories, modules, & packages:
  - A repository contains one or more Go modules.
  - Each module contains one or more Go packages.
  - Each package consists of one or more Go source files in a single directory.
- **go.mod**: A module is defined by a tree of Go source files with a `go.mod` file in the tree's root directory. Module source code may be located outside of GOPATH. There are four directives: `module`, `require`, `replace`, `exclude`.
- **Version selection**: If you add a new import to your source code that is not yet covered by a `require`in `go.mod`, most go commands like 'go build' & 'go test' will automatically look up the proper module & add the _highest_ version of that new direct dependency to your module's `go.mod` as a `require` directive.
- **Semantic Import versioning**: The result of following both the import compatibility rule & semver is called _Semantic Import Versioning_, where the major version is included in the import path.
- As a result of Semantic Import Versioning, code opting in to Go modules **must comply with these rules**:
  - Follow [semver](https://semver.org/) (with tags such as `v1.2.3`).
  - If the module is version v2 or higher, the major version of the module _must_ be included as a `/vN` at the end of the module paths.
  - If the module is version v0 or v1, do _not_ include the major version in either the module path or the import path.
- Starting in Go 1.13, module mode will be the default for all development.
- Check [Go Modules wiki](https://github.com/golang/go/wiki/Modules) for updated information.

## 2. Quickstart

Go Module's hello-world: init module and add dependencies.

```shell
# Create a directory outside of your GOPATH:
$ mkdir -p /tmp/scratchpad/hello
$ cd /tmp/scratchpad/hello
# Initialize a new module:
$ go mod init github.com/you/hello

go: creating new go.mod: module github.com/you/hello
# Write your code
$ cat <<EOF > hello.go
package main

import (
    "fmt"
    "rsc.io/quote"
)

func main() {
    fmt.Println(quote.Hello())
}
EOF

# Introduce `go mod tidy`
$ go mod tidy
go: finding module for package rsc.io/quote
go: downloading rsc.io/quote v1.5.2
go: found rsc.io/quote in rsc.io/quote v1.5.2

$ cat go.mod
module github.com/you/hello
require rsc.io/quote v1.5.2

# List the current module and all its dependencies
$ go list -m all

# Build & run
$ go build
$ ./hello
Hello, world.
```

Upgrade your dependencies:

```shell
$ go get golang.org/x/text
go: downloading golang.org/x/text v0.7.0
go: upgraded golang.org/x/text v0.0.0-20170915032832-14c0d48ead0c => v0.7.0
```

Remove unused dependencies: Just run `go mod tidy` and Go does the rest.

Vendor dependencies with `go mod vendor`:

- `go mod vendor` constructs a directory named `vendor` in the main module's root directory.
- You can check in `vendor` to your Version Control System.

```shell
$ go mod vendor
$ tree -L 3
├── go.mod
├── go.sum
├── hello.go
└── vendor
    ├── golang.org
    │   └── x
    ├── modules.txt
    └── rsc.io
```

## 3. Go Proxy

- Started from Go 1.13, `go` tool defaults to downloading modules from the public Go module mirror: <https://proxy.golang.org>.
- Change the default Go proxy:

```shell
export GOPROXY=https://goproxy.io,direct
```

- For private code, configure `GOPRIVATE`:

```shell
go env -w GOPRIVATE=*.corp.com,github.com/secret/repo
```

## 4. Workspaces

Go 1.18 introduces workspaces. Workspaces allow you to create projects of several modules that share a common list of dependencies through a `go.work` file.

```go
go 1.18

use ./my/first/thing
use ./my/second/thing

replace example.com/bad/thing v1.4.5 => example.com/good/thing v1.4.5
```

Example:

```shell
$ mkdir workspace
$ cd workspace
$ mkdir hello && cd hello
$ go mod init example.com/hello
# ... create hello.go

# Create the workspace
$ cd ../
$ go work init ./hello
$ go work use ./example

$ cat go.work
go 1.20

use (
    ./example
    ./hello
)
```

## 5. Organizing a Go module

Source: <https://go.dev/doc/modules/layout>

### 5.1. Basic package

```shell
project-root-directory/
  go.mod
  modname.go
  modname_test.go
```

### 5.2. Basic command

```shell
project-root-directory/
  go.mod
  auth.go
  client.go
  main.go
```

### 5.3. Package with internal packages

```shell
project-root-directory/
  internal/
    auth/
      auth.go
    hash/
      hash.go
  go.mod
  modname.go
```

### 5.4. Multiple packages

```shell
project-root-directory/
  go.mod
  modname.go
  auth/
    auth.go
  hash/
    hash.go
  internal/
    trace/
      trace.go
```

### 5.5. Server project

```shell
project-root-directory/
  go.mod
  internal/
    auth/
    metrics/
    model/
  cmd/
    api-server/
      main.go
    metrics-analyzer/
      main.go
```
