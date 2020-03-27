# Building Go applications for different operating systems and architectures

Source: https://www.digitalocean.com/community/tutorials/building-go-applications-for-different-operating-systems-and-architectures

It is important to consider the operating system and underlying processor architecture that you would like to compile your binary for.

Go solves this problem by building support for multiple platforms directly into the `go build` tool as well as the rest of the Go toolchain:
* Environment variables.
* Build tags.

## Possible platforms for `GOOS` and `GOARCH`

To find the list of possible platforms:

```bash
$ go tool dist list
```

* Output with Format: `<operating system>/<architecture>`

```
aix/ppc64
android/386
android/amd64
android/arm
android/arm64
darwin/386
darwin/amd64
darwin/arm
darwin/arm64
dragonfly/amd64
freebsd/386
freebsd/amd64
freebsd/arm
freebsd/arm64
illumos/amd64
js/wasm
linux/386
linux/amd64
linux/arm
linux/arm64
linux/mips
linux/mips64
linux/mips64le
linux/mipsle
linux/ppc64
linux/ppc64le
linux/riscv64
linux/s390x
netbsd/386
netbsd/amd64
netbsd/arm
netbsd/arm64
openbsd/386
openbsd/amd64
openbsd/arm
openbsd/arm64
plan9/386
plan9/amd64
plan9/arm
solaris/amd64
windows/386
windows/amd64
windows/arm
```

## Sample

Check [the sample source code](./origin)

## Using `GOOS` or `GOARCH` build tags

Check [the sample source code](./build-tags)

## Using your local `GOOS` and `GOARCH` environment variables

The `go build` command behaves in a similar manner to the `go env` command. You can set either the `GOOS` or `GOARCH` environment variables to build for a different platform using go build.

Build a `windows` binary by settings the `GOOS` environment variable to `windows` when running the `go build` command:

```bash
$ GOOS=windows go build
```

## Using `GOOS` and `GOARCH` filename suffixes

You could specify the OS and architecture by changing the filename to `filename_GOOS_GOARCH.go`.

Check [the sample source code]()./filename-suffixes
