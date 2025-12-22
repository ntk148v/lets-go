# Interfaces

> **NOTE**: Code examples for this section are stored in [`examples/6/`](../examples/6/).

Table of Contents:

- [Interfaces](#interfaces)
  - [1. Which is what?](#1-which-is-what)
  - [2. Empty interface](#2-empty-interface)
  - [3. Methods](#3-methods)
  - [4. Listing interfaces in interfaces](#4-listing-interfaces-in-interfaces)
  - [5. Introspection and reflection](#5-introspection-and-reflection)

Every type has an interface, which is the _set of methods defined_ for that type.

```go
/* a struct type S with 1 field, 2 methods */
type S struct { i int }
func (p *S) Get() int { return p.i }
func (p *S) Put(v int) { p.i = v }

/* an interface type */
type I interface {
    Get() int
    Put() int
}
/* S is a valid implementation for interface I */
```

S is a valid implementation for ineterface I. A Go program can use this fact via yet another meaning of interface, which is an interface value:

```go
func f(p I) {
    fmt.Println(p.Get())
    p.Put(1)
}
var s S
/* Because S implements I, we can call the
function f passing in a pointer to a value
of type S */
/* The reason we need to take the address of s,
rather than a value of type S, is because
we defined the methods on s to operae on pointers */
f(&s)
```

The fact that you do not need to declare whether or not a type implements an interface means that Go implements a form of [duck typing](https://en.wikipedia.org/wiki/Duck_typing). This is not pure duck typing, because when possible the Go complier will statically check whether the type implements the inerface. However, Go does have a purely dynamic aspect, in that you can convert from one interface to another. In the general case, that conversion is checked at run time. If the conversion is invalid - if the type of the value stored in the existing interface value does not satisfy the interface to which it is being converted - the program will fail with a run time error.

_Duck typing - If it looks like a duck, & it quacks like a duck, then it is a duck_. It means if it has a set of methods that match an interface, then you can use it wherever that interface is needed without explicitly defining that your types implement that interface.

```go
package main

import "fmt"

type Duck interface {
    Quack()
}

type Donald struct {
}

func (d Donald) Quack() {
    fmt.Println("quack quack!")
}

type Daisy struct {
}

func (d Daisy) Quack() {
    fmt.Println("-quack -quack")
}

func sayQuack(duck Duck) {
    duck.Quack()
}

type Dog struct {
}

func (d Dog) Bark() {
    fmt.Println("go go")
}

func main() {
    donald := Donald{}
    sayQuack(donald) // quack
    daisy := Daisy{}
    sayQuack(daisy) // --quack
    dog := Dog()
    sayQuack(dog) // compile error - cannot use dog (type Dog) as type Duck
}
```

Go's interfaces let you use `duck typing` like you would in a purely dynamic language like Python but still have the compiler catch obvious mistakes like passing an `int` where an object with a `Read` method was expected, or like calling the `Read` method with the wrong number of arguments.

## 1. Which is what?

Let's define another type R that also implements the interface I.

```go
type R struct { i int }
func (p * R) Get() int { return p.i }
func (p *R) Put(v int) { p.i = v }

func f(p I) {
    switch t := p.(type) {
        case *S:
        case *R:
        default:
    }
}
```

## 2. Empty interface

Create a generic function which has an empty interface as its argument

```go
func g(something interface{}) int {
    return something.(I).Get()
}
```

The `.(I)` is a type assertion which converts `something` to an interface of type I. If we have the type we can invoke the `Get()` function.

```go
s = new(S)
fmt.Println(g(s))
```

## 3. Methods

- Methods are functions that have a receiver.
- You can definen methods on any type (except on non-local types, this includes built-in types: the type `int` can not have methods).
- Methods on interface types
  - An interface defines a set of methods. A method contains the actual code.
  - An interface is the definition & the methods are the implementation.
- By convention, one-method interfaces are named by the method name plus the -er suffix: Reader, Writer, Formatter,...
- Pointer & Non-pointer method receivers.

  ```go
  func (s *MyStruct) pointerMethod() {} // method on pointer
  func (s MyStruct) valueMethod() {} // method on value
  ```

  - When defining a method on a type, the receiver behaves exactly as if it were an argument to the method. Whether to define the receiver as a value or as a pointer is the same question, then, as whether a function argument should be a value or a pointer.
  - 1st: Does the method need to modify the receiver? If it _does_, the receiver must be a _pointer_ (Slices & maps act as references, so their story is a little more subtle, but for instance to change the length of a slice in a method the receiver must still be a pointer). Otherwise, it should be _value_.

  ```go
  package main

  import "fmt"

  type Mutatable struct {
      a int
      b int
  }

  func (m Mutatable) StayTheSame() {
      m.a = 5
      m.b = 7
  }

  func (m *Mutatable) Mutate() {
      m.a = 5
      m.b = 7
  }

  func main() {
      m := &Mutatable{0, 0}
      fmt.Println(m)
      m.StayTheSame()
      fmt.Println(m)
      m.Mutate()
      fmt.Println(m)
  }
  ```

  - 2nd: efficiency. If the receiver is large, a big `struct` for instance, it will be much cheaper to use a pointer receiver.
  - 3rd: consistency. If some of the methods of the type must have pointer receivers, the rest should too, so the method set is consistent regardless of how the type is used.
  - For types such as basic types, slices & small `struct`, a value receiver is very cheap so unless the semantics of the methods requires a pointer, a value receiver is effient & clear.

## 4. Listing interfaces in interfaces

## 5. Introspection and reflection

Type switching: A type switching is like a regular switch statement, but the cases in a type switch specify types (not values) which are compared against the type of the value held by the given interface value.

```go
package main

import "fmt"

func do(i interface{}) {
    switch v := i.(type) {
    case int:
        fmt.Printf("Twice %v is %v\n", v, v*2)
    case string:
        fmt.Printf("%q is %v bytes long\n", v, len(v))
    default:
        fmt.Printf("I don't know about type %T!\n", v)
    }
}

func main() {
    do(21)
    do("hello")
    do(true)
}

// Twice 21 is 42
// "hello" is 5 bytes long
// I don't know about type bool!
```

More examples at [a8m/reflect-examples](https://github.com/a8m/reflect-examples).
