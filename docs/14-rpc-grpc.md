# RPC, gRPC and Protobuf

> **NOTE**: Code examples for this section are stored in [`examples/13/`](../examples/13/).

Table of Contents:

- [RPC, gRPC and Protobuf](#rpc-grpc-and-protobuf)
  - [1. Remote Procedure Call (RPC)](#1-remote-procedure-call-rpc)
  - [2. gRPC and Protobuf](#2-grpc-and-protobuf)

This section is mainly taken from: <https://github.com/zalopay-oss/go-advanced/blob/master/ch3-rpc/ch3-01-rpc-go.md>

## 1. Remote Procedure Call (RPC)

A [remote procedure call (RPC)](https://en.wikipedia.org/wiki/Remote_procedure_call) is when a computer program causes a procedure (subroutine) to execute in a different address space (commonly on another computer on a shared network), which is written as if it were a normal (local) procedure call, without the programmer explicitly writing the details for the remote interaction.

![](https://media.geeksforgeeks.org/wp-content/uploads/operating-system-remote-call-procedure-working.png)

Write a simple example with [net/rpc](https://golang.org/pkg/net/rpc/) library.

```go
// rpcserver/main.go
package main

import (
 "log"
 "net"
 "net/rpc"
)

type HelloService struct{}

// Only methods that satisfy these criteria will be made available for remote access:
// - the method's type is exported.
// - the method is exported.
// - the method has two arguments, both exported (or builtin) types.
// - the method's second argument is a pointer.
// - the method has return type error.
// func (t *T) MethodName(argType T1, replyType *T2) error
func (p *HelloService) Hello(request string, reply *string) error {
 *reply = "Hello " + request
 return nil
}

func main() {
 rpc.RegisterName("HelloService", new(HelloService))
 listener, err := net.Listen("tcp", ":8081")
 if err != nil {
  log.Fatal("Listen TCP error:", err)
 }

 log.Println("Server is ready")

 for {
  conn, err := listener.Accept()
  if err != nil {
   log.Fatal("Accept error:", err)
  }

  go func() {
   log.Println("Accept new client:", conn.RemoteAddr())
   rpc.ServeConn(conn)
  }()
 }
}
```

```go
// rpcclient/main.go
package main

import (
 "log"
 "net/rpc"
)

func main() {
 client, err := rpc.Dial("tcp", "localhost:8081")
 if err != nil {
  log.Fatal("Dialing error:", err)
 }

 var reply string

 if err = client.Call("HelloService.Hello", "Kien", &reply); err != nil {
  log.Fatal(err)
 }

 log.Println(reply)
}
```

Run it:

```shell
# Server
$ go run rpcserver/main.go
2023/08/09 16:29:29 Server is ready
2023/08/09 16:29:30 Accept new client: 127.0.0.1:38728

# Client
$ go run rpcclient/main.go
2023/08/09 16:29:30 Hello Kien
```

## 2. gRPC and Protobuf

Protocol Buffers, also referred as **protobuf**, is Google's language-neutral, platform-neutral, extensible mechanism for serializing structured data. Protocol Buffers are smaller, faster, and simpler that provides high performance than other standards such as XML and JSON.

By using protocol buffers, you can define your structured data, then you generate source code for your choice of programming language using the protocol buffer compiler named **protoc**.

Install `protoc`:

```shell
# Ubuntu
$ sudo apt install -y protobuf-compiler
# install go plugin
$ go install github.com/golang/protobuf/protoc-gen-go@latest
```

Prepare `hello.proto`:

```protobuf
syntax = "proto3";
package main;
message String {
    string value = 1;
}
```

**gRPC** is a high performance, open-source remote procedure call (RPC) framework that can run anywhere. It enables client and server applications to communicate transparently, and makes it easier to build connected systems.

![](https://grpc.io/img/landing-2.svg)

The gRPC server implements the service interface and runs an RPC server to handle client calls to its service methods. On the client side, the client has a stub that provides the same methods as the server.

By default, gRPC uses Protocol Buffers as the Interface Definition Language (IDL) and as its underlying message interchange format.
