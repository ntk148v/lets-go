# Pointers, Allocation and Conversions

> **NOTE**: Code examples for this section are stored in [`examples/5/`](../examples/5/).

Table of Contents:

- [Pointers, Allocation and Conversions](#pointers-allocation-and-conversions)
  - [1. Pointers](#1-pointers)
  - [2. Allocation and Constructor](#2-allocation-and-constructor)
    - [2.1. new vs make](#21-new-vs-make)
    - [2.2. Constructors and Composite Literals](#22-constructors-and-composite-literals)
    - [2.3. Defining Your Own Types](#23-defining-your-own-types)
    - [2.4. Methods](#24-methods)
  - [3. Conversions](#3-conversions)

## 1. Pointers

- Go has pointers but not pointer arthmetic, so they act more like references than pointers that you may know from C.
  - A pointer is a variable which stores the address of another variable. A pointer is thus the location at which a value is stored. Not every value has an address but every variable does.
  - A reference is a variable which refers to another value.
  - There is no pointer arithmetic. You cannot write in Go. That is, you cannot alter the address p points to unless you assign another address to it.

  ```go
  var p *int
  p++
  ```

  - Want to understand more? Check this [article](http://spf13.com/post/go-pointers-vs-references/) or [another `this`](https://dave.cheney.net/2017/04/26/understand-go-pointers-in-less-than-800-words-or-your-money-back).

- Pointers are useful. Remember that when you call a function in Go, the variables are _pass-by-value_. So, for efficiency & the possibility to modify a passed value in functions we have pointers.

- Pointer type (\* type) & address-of (&) operators \*: If a variable is declared `var x int`, the expression `&x` ("address of x") yields a pointer to an integer variable (a value of type `* int`). If this value is called `p`, we say "`p` points to to `x`", or equivalently "`p` contains the address of `x`". The variable to which `p` points is written `*p`. The expression `*p` yields the value of that variable, an `int`, but since `*p` denotes a variable, it may also appear on the left-hand side of an assignment, in which case the assignment updates the variable. [Reference here](https://notes.shichao.io/gopl/ch2/#pointers)

  ```go
  x := 1
  p := &x          // p, of type *int, points to x
  fmt.Println(*p)  // "1"
  *p = 2           // equivalent to x = 2
  fmt.Println(x)   // "2"
  ```

- All newly declared variables are assigned their zero value & pointers are no different. A newly declared pointer, or just a pointer that points to nothing, has a nil-value.

```go
var p *int // declare a pointer
fmt.Printf("%v", p)

var i int
p = &i // Make p point to i

fmt.Printf("%v", p) // Print somthing like 0x7ff96b81c000a
```

- Pointer illustration:

![](https://media.geeksforgeeks.org/wp-content/uploads/20190710183146/Pointer-To-Pointer.jpg)

```go
// Go program to illustrate the
// concept of the Pointer to Pointer
package main

import "fmt"

// Main Function
func main() {

        // taking a variable
        // of integer type
    var V int = 100

    // taking a pointer
    // of integer type
    var pt1 *int = &V

    // taking pointer to
    // pointer to pt1
    // storing the address
    // of pt1 into pt2
    var pt2 **int = &pt1

    fmt.Println("The Value of Variable V is = ", V)
    fmt.Println("Address of variable V is = ", &V)

    fmt.Println("The Value of pt1 is = ", pt1)
    fmt.Println("Address of pt1 is = ", &pt1)

    fmt.Println("The value of pt2 is = ", pt2)

    // Dereferencing the
    // pointer to pointer
    fmt.Println("Value at the address of pt2 is or *pt2 = ", *pt2)

    // double pointer will give the value of variable V
    fmt.Println("*(Value at the address of pt2 is) or **pt2 = ", **pt2)
}

// The Value of Variable V is =  100
// Address of variable V is =  0x414020
// The Value of pt1 is =  0x414020
// Address of pt1 is =  0x40c128
// The value of pt2 is =  0x40c128
// Value at the address of pt2 is or *pt2 =  0x414020
// *(Value at the address of pt2 is) or **pt2 =  100
```

- Check [pointer vs reference](../tips-notes/pointer-vs-references.md).

> **NOTE (Go 1.24)**: Go 1.24 introduces **Weak Pointers** in the new `weak` package. Weak pointers are references that don't prevent garbage collection:
>
> ```go
> import "weak"
> obj := new(int)
> wp := weak.Make(obj)
> if val := wp.Value(); val != nil {
>     fmt.Println(*val)
> }
> ```

## 2. Allocation and Constructor

- Go also has garbage collection.
- To allocate memory Go has 2 primitives, `new` & `make`.

### 2.1. new vs make

- **new** allocates; **make** initializes.
  - _new(T)_ returns \*T pointing to a zerod T.
  - _make(T)_ returns an initialized T.
  - _make_ is only used for slices, maps, channels.

### 2.2. Constructors and Composite Literals

```go
// A lot of boiler plate
func NewFile(fd int, name string) *File {
    if fd < 0 {
        return nil
    }
    f := new(File)
    f.fd = fd
    f.name = name
    f.dirinfo = nil
    f.nepipe = 0
    return f
}
// Using a composite literal
func NewFile(fd int, name string) *File {
    if fd < 0 {
        return nil
    }
    f := File{fd, name, nil. 0}
    return &f // Return the address of a local variable. The storage associated with the variable survives after the function returns.
    // return &File{fd, name, nil, 0}
    // return &File{fd: fd, name: name}
}
```

- As a limiting case, if a composite literal contains no fields at all, it creates a zero value for the type. The expression `new(File)` & `&File{}` are equivalent.

- Composite literal can also be created for arrays, slices, & maps, with the field labels being indices or map keys as appropriate.

```go
ar := [...]string{Enone: "no error", Einval: "invalid argument"}
sl := []string{Enone: "no error", Einval: "invalid argument"}
ma := map[int]string {Enone: "no error", Einval: "invalid argument"}
```

### 2.3. Defining Your Own Types

```go
/* defining_own_type.go */
package main

import "fmt"

type NameAge struct {
    name string // both non exported fiedls
    age  int
}

func main() {
    a := new(NameAge)
    a.name = "Kien"
    a.age = 25
    fmt.Printf("%v\n", a) // &{Kien, 25}
}
```

- More on structure fields

```go
struct {
    x, y int
    A *[]int
    F func()
}
```

### 2.4. Methods

- Create a function that takes the type as an argument:

```go
func doSomething1(n1 *NameAge, n2 int) {/* */}
// method call
var n *NameAge
n.doSomething1(2)
```

- Create a function that works on the type:

```go
func (n1 *NameAge) doSomething2(n2 int) {/* */}
```

> **NOTE**: If x is addressable & &x's method set contains m, x.m() is shorthand for (&x).m().

- Suppose we have:

```go
// A mutex is a data type with two methods, Lock & Unlock
type Mutex struct {/* Mutex fields */}
func (m *Mutex) Lock() {/* Lock impl */}
func (m *Mutext) Unlock {/* Unlock impl */}

// NewMutex is equal to Mutex, but it does not have any of the methods of Mutex.
type NewMutex Mutex
// PrintableMutex hash inherited the method set from Mutex, contains the methods
// Lock & Unlock bound to its anonymous field Mutex
type PrintableMutex struct{Mutex}
```

## 3. Conversions

```
FROM    b []byte    i []int    r []rune    s string    f float32    i int
TO
[]byte  .                                  []byte(s)
[]int               .                      []int(s)
[]rune                                     []rune(s)
string  string(b)   string(i)  string(r)   .
float32                                                .            float32(i)
int                                                    int(f)       .
```

- float64 works the same as float32
- From a `string` to a slice of bytes or runes

```go
mystring := "hello this is string"
byteslice := []byte(mystring)
runeslice := []rune(string)
```

- From a slice of bytes or runes to a string

```go
b := []byte{'h', 'e', 'l', 'l', 'o'} // Composite literal
s := string(b)
i := []rune(26, 9, 1994)
r := string(i)
```

- For numeric values:
  - Convert to an integer with a specific (bit) length: `uint8(int)`.
  - From floating point to an integer value: `int(float32)`. This discards the fraction part from the floating point value.
  - And other way around `float32(int)`

- User defined types & conversions

```go
type foo struct { int } // Anonymous struct field
type bar foo            // bar is an alias for foo

var b bar = bar{1} // Declare `b` to be a `bar`
var f foo = b      // Assign `b` to `f` --> Cannot use b (type bar) as type foo in assignment
var f foo = foo(b) // OK!
```
