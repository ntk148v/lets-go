# Shrink Golang binary

Source: <https://github.com/xaionaro/documentation/blob/master/golang/reduce-binary-size.md>

Table of contents:

- [Shrink Golang binary](#shrink-golang-binary)
  - [0. Prepare](#0-prepare)
  - [1. Go build flags](#1-go-build-flags)
  - [2. Compress the binary](#2-compress-the-binary)
  - [3. Illegal: It works but ...](#3-illegal-it-works-but-)
  - [4. Result](#4-result)

## 0. Prepare

```shell
$ go version
go version go1.20.2 linux/amd64

$ mkdir /tmp/helloworld

$ cd /tmp/helloworld

$ cat > hello.go <<EOF
package main

import "fmt"

func main() {
   fmt.Println("hello-world")
}
EOF
```

- Default build:

```shell
$ go build -o default ./hello.go; stat -c %s default
1845890

$ export CGO_ENABLED=0

$ go build -o default ./hello.go; stat -c %s default
1845882
```

## 1. Go build flags

- Strip the debugging information with [-s and -w linker flags](https://golang.org/cmd/link/):

```shell
$ go build -ldflags="-s -w" -o stripped hello.go; stat -c %s stripped
1228800


```

- Disable function inlining:

```shell
$ go build -ldflags="-w -s" -gcflags=all=-l -o disinlining hello.go; stat -c %s disinlining
1208320

# It's not helpful on hello-world example, but it's helpful on big projects. It could save ~10%
```

- Disable bound checks:

```shell
$ go build -a -gcflags=all="-l -B" -ldflags="-w -s" -o disbound hello.go; stat -c %s disbound
1179648
```

## 2. Compress the binary

- You can compress the binary using [UPX](http://upx.sourceforge.net/).
- But, just keep in mind:
  - The binary will be much slower
  - It will consume more RAM
  - It will be almost useless if you already store your binary in a compressed state (for example in `initrd`, compressed by `xz`).
- Let's go:

```shell
$ upx disbound -o disbound-upx
                       Ultimate Packer for eXecutables
                          Copyright (C) 1996 - 2020
UPX 3.96        Markus Oberhumer, Laszlo Molnar & John Reiser   Jan 23rd 2020

        File size         Ratio      Format      Name
   --------------------   ------   -----------   -----------
   1179648 ->    479572   40.65%   linux/amd64   disbound-upx

Packed 1 file.
479572

$ upx --best --ultra-brute disbound -o disbound-upx-brute; stat -c %s disbound-upx-brute
                       Ultimate Packer for eXecutables
                          Copyright (C) 1996 - 2020
UPX 3.96        Markus Oberhumer, Laszlo Molnar & John Reiser   Jan 23rd 2020

        File size         Ratio      Format      Name
   --------------------   ------   -----------   -----------
   1179648 ->    386044   32.73%   linux/amd64   disbound-upx-brute

Packed 1 file.
386044
```

## 3. Illegal: It works but ...

- 32bits instead of 64bits: It has obvious limitations:
  - 32bit address space
  - 32bit integers
  - Less registers (less performance in some cases)
  - 32bit syscalls (for example there's no `kexec_file_load`)

```shell
$ go env GOARCH
amd64

$ GOARCH=386 go build -a -gcflags=all="-l -B" -ldflags="-w -s" -o disbound-32 hello.go; stat -c %s disbound-32
1122304
```

- You can find more in the Source link but I don't think that we should follow it,...

## 4. Result

```shell
$ ls -la
total 9.0M
drwxrwxr-x  2 kiennt kiennt 4.0K Apr 24 16:47 ./
drwxrwxrwt 25 kiennt kiennt  20K Apr 24 16:53 ../
-rwxrwxr-x  1 kiennt kiennt 1.8M Apr 24 15:48 default*
-rwxrwxr-x  1 kiennt kiennt 1.2M Apr 24 15:52 disbound*
-rwxrwxr-x  1 kiennt kiennt 1.1M Apr 24 16:47 disbound-32*
-rwxrwxr-x  1 kiennt kiennt 469K Apr 24 15:52 disbound-upx*
-rwxrwxr-x  1 kiennt kiennt 377K Apr 24 15:52 disbound-upx-brute*
-rwxrwxr-x  1 kiennt kiennt 1.2M Apr 24 15:51 disinlining*
-rwxrwxr-x  1 kiennt kiennt 1.8M Apr 24 15:47 hello*
-rw-rw-r--  1 kiennt kiennt   74 Apr 24 15:46 hello.go
-rwxrwxr-x  1 kiennt kiennt 1.2M Apr 24 15:50 stripped*
```
