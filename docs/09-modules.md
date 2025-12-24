# Modules (Golang version >=1.11)

Table of Contents:

- [Modules (Golang version \>=1.11)](#modules-golang-version-111)
  - [1. Concepts](#1-concepts)
  - [2. Quickstart](#2-quickstart)
  - [3. Go Proxy](#3-go-proxy)
  - [4. Workspaces](#4-workspaces)
  - [5. Organizing a Go module](#5-organizing-a-go-module)
    - [5.1. Basic package](#51-basic-package)
    - [5.2. Basic command](#52-basic-command)
    - [5.3. Package or command with supporting packages](#53-package-or-command-with-supporting-packages)
    - [5.4. Multiple packages](#54-multiple-packages)
    - [5.5. Multiple commands](#55-multiple-commands)
    - [5.6. Packages and commands in the same repository](#56-packages-and-commands-in-the-same-repository)
    - [5.7. Server project](#57-server-project)

Go 1.11 includes preliminary support for [modules](https://go.dev/doc/go1.11#modules), Go's [new dependency management system](https://blog.golang.org/versioning-proposal) that makes dependency version information explicit and easier to manage.

## 1. Concepts

- **Modules**: a collection of related Go packages that are versioned together as a single unit.
- Summarizing the relationship between repositories, modules, & packages:
  - A repository contains one or more Go modules.
  - Each module contains one or more Go packages.
  - Each package consists of one or more Go source files in a single directory.
- **go.mod**: A module is defined by a tree of Go source files with a `go.mod` file in the tree's root directory. Module source code may be located outside of GOPATH. There are four directives: `module`, `require`, `replace`, `exclude`.
- **Version selection**: If you add a new import to your source code that is not yet covered by a `require`in `go.mod`, most go commands like 'go build' & 'go test' will automatically look up the proper module & add the _highest_ version of that new direct dependency to your module's `go.mod` as a `require` directive. For example, if your new import corresponds to dependency M whose latest tagged release version is `v1.2.3`, your module's `go.mod` will end up with `require M v1.2.3`, which indicates module M is a dependency with allowed version >= v1.2.3 (and < v2, given v2 is considered incompatible with v1).
- **Semantic Import versioning**: The result of following both the import compatibility rule & semver is called _Semantic Import Versioning_, where the major version is included in the import path — this ensures the import path changes any time the major version increments due to a break in compatibility.
- As a result of Semantic Import Versioning, code opting in to Go modules **must comply with these rules**:
  - Follow [semver](https://semver.org/) (with tags such as `v1.2.3`).
  - If the module is version v2 or higher, the major version of the module _must_ be included as a `/vN` at the end of the module paths used in `go.mod` files (e.g., `module github.com/my/mod/v2`, `require github.com/my/mod/v2 v2.0.0`) & in the package import path (e.g., `import "github.com/my/mod/v2/mypkg"`).
  - If the module is version v0 or v1, do _not_ include the major version in either the module path or the import path.
- As of Go 1.11, the go command enables the use of modules when the current directory or any parent directory has a `go.mod`, provides the directory is _outside_ `$GOPATH/src`. (Inside `$GOPATH/src`, for compatibility, the go command still runs in the old `GOPATH` mode, even if a `go.mod` is found)
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
# Tidy makes sure go.mod matches the source code in the module.
# It adds any missing modules necessary to build the current module's
# packages and dependencies, and it removes unused modules that
# don't provide any relevant packages. It also adds any missing entries
# to go.sum and removes any unnecessary ones
$ go mod tidy                                                                                                                                              t/s/hello ﳑ
go: finding module for package rsc.io/quote
go: downloading rsc.io/quote v1.5.2
go: found rsc.io/quote in rsc.io/quote v1.5.2
go: downloading rsc.io/sampler v1.3.0
go: downloading golang.org/x/text v0.0.0-20170915032832-14c0d48ead0c

$ cat go.mod

module github.com/you/hello

require rsc.io/quote v1.5.2

# Add a new dependency often brings in other indirect dependencies too
# List the current module and all its dependencies
$ go list -m all
github.com/you/hello
golang.org/x/text v0.0.0-20170915032832-14c0d48ead0c
rsc.io/quote v1.5.2
rsc.io/sampler v1.3.0

# In addition to go.mod, there is a go.sum file containing the expected
# cryptographic hashes of the content of specific module versions
$ cat go.sum                                                                                                                                                    t/s/hello ﳑ
golang.org/x/text v0.0.0-20170915032832-14c0d48ead0c h1:qgOY6WgZOaTkIIMiVjBQcw93ERBE4m30iBm00nkL0i8=
golang.org/x/text v0.0.0-20170915032832-14c0d48ead0c/go.mod h1:NqM8EUOU14njkJ3fqMW+pc6Ldnwhi/IjpwHt7yyuwOQ=
rsc.io/quote v1.5.2 h1:w5fcysjrx7yqtD/aO+QwRjYZOKnaM9Uh2b40tElTs3Y=
rsc.io/quote v1.5.2/go.mod h1:LzX7hefJvL54yjefDEDHNONDjII0t9xZLPXsUe+TKr0=
rsc.io/sampler v1.3.0 h1:7uVkIFmeBqHfdjD+gZwtXXI+RODJ2Wc4O7MPEh/QiW4=
rsc.io/sampler v1.3.0/go.mod h1:T1hPZKmBbMNahiBKFy5HrXp6adAjACjK9JXDnKaTXpA=

# Build & run
$ go build
$ ./hello

Hello, world.
```

Upgrade your dependencies:

```shell
# From the output of go list -m all, we're using an untagged version of golang.org/x/text
# Let's upgrade to the latest tagged version
$ go get golang.org/x/text                                                                                                                                      t/s/hello ﳑ
go: downloading golang.org/x/text v0.7.0
go: upgraded golang.org/x/text v0.0.0-20170915032832-14c0d48ead0c => v0.7.0

$ cat go.mod                                                                                                                                                    t/s/hello ﳑ
module github.com/you/hello

go 1.19

require rsc.io/quote v1.5.2

require (
        golang.org/x/text v0.7.0 // indirect
        rsc.io/sampler v1.3.0 // indirect
)

$ go list -m all                                                                                                                                                t/s/hello ﳑ
github.com/you/hello
golang.org/x/mod v0.6.0-dev.0.20220419223038-86c51ed26bb4
golang.org/x/sys v0.0.0-20220722155257-8c9f86f7a55f
golang.org/x/text v0.7.0
golang.org/x/tools v0.1.12
rsc.io/quote v1.5.2
rsc.io/sampler v1.3.0
```

Remove unused dependencies: If you want to remove any dependencies, just simple run `go mod tidy` and Go does the rest.

Go modules takes care of versioning, but it doesn't necessarily take care of modules disappearing off the Internet or the Internet not being available. If a module is not available, the code cannot be built. Go Proxy will mitigate disappearing modules to some extent by mirroring modules, but it may not do it for all modules for all time. That's why `go` tool provides `go mod vendor` command.

- `go mod vendor` command constructs a directory named `vendor` in the main module's root directory that contains copies of all packages needed to support builds and tests of packages in the main modules.
- `go mod vendor` also creates the file `vendor/modules.txt` that contains a list of vendored packages and the module versions they were copied from.
- You can check in `vendor` to your Version Control System, then copy this around.

```shell
$ go mod vendor

# Main module's directory structure
$ tree -L 3
├── go.mod
├── go.sum
├── hello.go
└── vendor
    ├── golang.org
    │   └── x
    ├── modules.txt
    └── rsc.io
        ├── quote
        └── sampler
```

## 3. Go Proxy

- Started from Go 1.13, `go` tool defaults to downloading modules from the public Go module mirror: <https://proxy.golang.org> and also defaults to validating downloaded modules (regardless of source) against the public Go checksum database at <https://sum.golang.org>.
- If you want to change the default Go proxy, you can use the following command:

```shell
export GOPROXY=https://goproxy.io,direct
```

- The `go` command defaults to downloading modules from the public Go module mirror, therefore if you have private code, you most likely should configure the `GOPRIVATE` setting (such as `go env -w GOPRIVATE=*.corp.com,github.com/secret/repo`), or the more fine-grained variants `GONOPROXY` or `GONOSUMDB` that support less frequent use cases. See the [documentation](https://go.dev/ref/mod#private-module-privacy) for more details.

## 4. Workspaces

Go introduces the concept of workspaces in 1.18. Workspaces allows you to create projects of several modules that share a common list of dependencies through a new file called `go.work`. The dependencies in this file can span multiple modules and anything declared in the `go.work` file will override dependencies in the module's `go.mod`.

- A workspace is a collection of modules on disk that are used as the main modules when running minimal version selection (MVS).
- A workspace can be declared in a `go.work` file that specifies relative paths to the module directories of each the modules in the workspace. When no `go.work` file exists, the workspace consists of the single module containing the current directory.
  - Lexical elements in `go.work` files are defined in exactly the same way as for `go.mod` files.
  - Check out [here](https://go.dev/ref/mod#workspaces).

```go
go 1.18

use ./my/first/thing
use ./my/second/thing

// or
// use (
//     ./my/first/thing
//     ./my/second/thing
// )

replace example.com/bad/thing v1.4.5 => example.com/good/thing v1.4.5
```

Example:

```shell
$ mkdir workspace
$ cd workspace
# Create hello module
$ mkdir hello
$ cd hello
$ go mod init eaxmple.com/hello
go: creating new go.mod: module example.com/hello
$ cat <<EOF > hello.go
package main

import (
    "fmt"

    "golang.org/x/example/stringutil"
)

func main() {
    fmt.Println(stringutil.Reverse("Hello"))
}
EOF

$ go mod tidy
go: finding module for package golang.org/x/example/stringutil
go: found golang.org/x/example/stringutil in golang.org/x/example v0.0.0-20220412213650-2e68773dfca0

$ go run example.com/hello
olleH

# Create the workspace
$ cd ../
$ go work init ./hello
$ tree
.
├── go.work
└── hello
    ├── go.mod
    ├── go.sum
    └── hello.go

1 directory, 4 files

# Go command includes all the modules in the workspace as main modules. This allow us to refer to a package in the module
# even outside the module.
$ go run example.com/hello
olleH

# Download and modify the golang.org/x/example module
$ git clone https://go.googlesource.com/example
Cloning into 'example'...
remote: Total 165 (delta 27), reused 165 (delta 27)
Receiving objects: 100% (165/165), 434.18 KiB | 1022.00 KiB/s, done.
Resolving deltas: 100% (27/27), done.
# Add module to the workspace
$ go work use ./example
$ tree -L 1
.
├── example
├── go.work
└── hello

2 directories, 1 file

$ cat go.work
go 1.20

use (
    ./example
    ./hello
)

$ cd example/stringutil
# Create a new file
$ cat <<EOF > toupper.go
package stringutil

import "unicode"

// ToUpper uppercases all the runes in its argument string.
func ToUpper(s string) string {
    r := []rune(s)
    for i := range r {
        r[i] = unicode.ToUpper(r[i])
    }
    return string(r)
}
EOF

# Modify hello program
$ cd ../../hello
$ cat <<EOF > hello.go
package main

import (
    "fmt"

    "golang.org/x/example/stringutil"
)

func main() {
    fmt.Println(stringutil.ToUpper("Hello"))
}
EOF

$ cd ..
# Go command finds the example.com/hello module specified in the command line
# in the hello directory specified by the go.work file, and similiarly
# resolves the golang.org/x/example import using the go.work file.
$ go run example.com/hello
HELLO
```

## 5. Organizing a Go module

Source: <https://go.dev/doc/modules/layout>

Go projects can include packages, command-line programs or a combination of the two. This guide is organized by project type.

> **NOTE**: throughout this document, file/package names are entirely arbitrary

### 5.1. Basic package

A basic Go package has all its code in the project's root directory. The project consists of a single module, which consists of a single package. The package name matches the last path component of the module name. For a very simple package requiring a single Go file, the project structure is:

```shell
project-root-directory/
  go.mod
  modname.go
  modname_test.go
  auth.go
  auth_test.go
  hash.go
  hash_test.go
```

The code in `modname.go` declares the package with:

```go
package modname

// ... package code here
```

### 5.2. Basic command

A basic executable program (or command-line tool) is structured according to its complexity and code size. The simplest program can consist of a single Go file where func main is defined. Larger programs can have their code split across multiple files, all declaring package main:

```shell
project-root-directory/
  go.mod
  auth.go
  auth_test.go
  client.go
  main.go
```

Here the `main.go` file contains `func main`, but this is just a convention. The "main" file can also be called `modname.go` (for an appropriate value of modname) or anything else.

### 5.3. Package or command with supporting packages

Larger packages or commands may benefit from splitting off some functionality into supporting packages. Initially, it's recommended placing such packages into a directory named `internal`; [this prevents](https://pkg.go.dev/cmd/go#hdr-Internal_Directories) other modules from depending on packages we don't necessarily want to expose and support for external uses. Since other projects cannot import code from our `internal` directory, we're free to refactor its API and generally move things around without breaking external users. The project structure for a package is thus:

```shell
project-root-directory/
  internal/
    auth/
      auth.go
      auth_test.go
    hash/
      hash.go
      hash_test.go
  go.mod
  modname.go
  modname_test.go
```

### 5.4. Multiple packages

A module can consist of multiple importable packages; each package has its own directory, and can be structured hierarchically. Here's a sample project structure:

```shell
project-root-directory/
  go.mod
  modname.go
  modname_test.go
  auth/
    auth.go
    auth_test.go
    token/
      token.go
      token_test.go
  hash/
    hash.go
  internal/
    trace/
      trace.go
```

As a reminder, we assume that the module line in go.mod says:

```go
module github.com/someuser/modname
```

### 5.5. Multiple commands

Multiple programs in the same repository will typically have separate directories:

```shell
project-root-directory/
  go.mod
  internal/
    ... shared internal packages
  prog1/
    main.go
  prog2/
    main.go
```

In each directory, the program's Go files declare package `main`. A top-level `internal` directory can contain shared packages used by all commands in the repository.

A common convention is placing all commands in a repository into a `cmd` directory; while this isn't strictly necessary in a repository that consists only of commands, it's very useful in a mixed repository that has both commands and importable packages, as we will discuss next.

### 5.6. Packages and commands in the same repository

Sometimes a repository will provide both importable packages and installable commands with related functionality. Here's a sample project structure for such a repository:

```shell
project-root-directory/
  go.mod
  modname.go
  modname_test.go
  auth/
    auth.go
    auth_test.go
  internal/
    ... internal packages
  cmd/
    prog1/
      main.go
    prog2/
      main.go
```

### 5.7. Server project

Go is a common language choice for implementing servers. There is a very large variance in the structure of such projects, given the many aspects of server development: protocols (REST? gRPC?), deployments, front-end files, containerization, scripts and so on. We will focus our guidance here on the parts of the project written in Go.

Server projects typically won't have packages for export, since a server is usually a self-contained binary (or a group of binaries). Therefore, it's recommended to keep the Go packages implementing the server's logic in the `internal` directory. Moreover, since the project is likely to have many other directories with non-Go files, it's a good idea to keep all Go commands together in a cmd directory:

```shell
project-root-directory/
  go.mod
  internal/
    auth/
      ...
    metrics/
      ...
    model/
      ...
  cmd/
    api-server/
      main.go
    metrics-analyzer/
      main.go
    ...
  ... the project's other directories with non-Go code
```

In case the server repository grows packages that become useful for sharing with other projects, it's best to split these off to separate modules.
