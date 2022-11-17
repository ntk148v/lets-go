# The empty struct

Source: https://dave.cheney.net/2014/03/25/the-empty-struct

```go
// struct type that has no fields
type Q struct{}
var q struct{}
```

- Width:
  - Width describes the number of bytes of storage an instance of a type occupies.
  - Width is a property of a type, always a multiple of 8 bits.
  - Can discover the width of any value using `unsafe.Sizeof`.
- An empty struct's width = 0.

```go
var s struct{}
fmt.Println(unsafe.Sizeof(s)) // prints 0

var x [1000000000]struct{}
fmt.Println(unsafe.Sizeof(x)) // prints 0

// Slices of struct{}s consume only the space for their slice header
// type sliceHeader struct {
//     ptr unsafe.Pointer
//     len int
//     cap int
// }
// var s sliceHeader
// fmt.Println(unsafe.Sizeof(s)) // print 24 -> 8 (ptr) + 8 (len) + 8 (cap)
var x = make([]struct{}, 1000000000)
fmt.Println(unsafe.Sizeof(x)) // prints 24 in the playground

// The address of empty struct
var a, b struct{}
var c = &a // addressable
fmt.Println(&a == &b) // true
```

- struct{} as a method receiver.

```go
type S struct{}

func (s *S) addr() { fmt.Printf("%p\n", s) }

func main() {
        var a, b S
        a.addr() // 0x585448 - address of all zero sized values
        b.addr() // 0x585448
}
```

- `chan struct{}` construct used for signaling between go routines.
- What uses a type with an empty struct has in Go?
  - It's the smallest building block in Go. Its size is literally 0 bytes.
  - If it has zero size, you may create a slice of 1000's empty structures and this slice will be very tiny. Because really Go stores only a number of them in the slice but not them itself. The same story with channels.
  - All pointers to it always point to the same special place in memory.
  - Very useful in channels when you notify about some event but you don't need to pass any information about it, only a fact. Best solution is to pass an empty structure because it will only increment a counter in the channel but not assign memory, copy elements and so on. Sometime people use Boolean values for this purpose, but it's much worse.
  - Zero size container for methods. You may want have a mock for testing interfaces. Often you don't need data on it, just methods with predefined input and output.
  - Go has no `Set` object. Bit can be easily realized as a `map[keyType]struct{}`. This way map keeps only keys and no values.
