# Pointer vs References

[Source](https://spf13.com/post/go-pointers-vs-references/)

* C, C++: pointers
* C++, Java, Python, Perl & PHP: references
* On the surface both references and pointers are very similar, both are used to have one variable provide  access to another.

## What is the difference?

A pointer is a variable which stores the address of another variable.

A reference is a variable which refers to another variable.

```C++
int i = 3; // define a variable
int *ptr = &i; // define a pointer to that variable's memory address
int &ref = i; // define a reference to the first variable

// change the value of i to 13
*ptr = 13;
ref = 13;
```

What happens if I try to access the ptr directly without dereferencing first?
-> Pointer can be reassigned while references cannot. In other words a pointer can be assigned to a different address.

```go
package main

import "fmt"

var ap *int

func main() {
	a := 1 // define int
	b := 2 // define int

	ap = &a
	fmt.Println(ap)
	fmt.Println(*ap)
	// set ap to address of a (&a)
	//   ap address: 0x2101f1018
	//   ap value  : 1

	*ap = 3
	fmt.Println(ap)
	fmt.Println(*ap)
	// change the value at address &a to 3
	//   ap address: 0x2101f1018
	//   ap value  : 3

	a = 4
	fmt.Println(ap)
	fmt.Println(*ap)
	// change the value of a to 4
	//   ap address: 0x2101f1018
	//   ap value  : 4

	ap = &b
	fmt.Println(ap)
	fmt.Println(*ap)
	// set ap to the address of b (&b)
	//   ap address: 0x2101f1020
	//   ap value  : 2

	ap2 := ap
	fmt.Println(ap2)
	fmt.Println(*ap2)
	// set ap2 to the address in ap
	//   ap  address: 0x2101f1020
	//   ap  value  : 2
	//   ap2 address: 0x2101f1020
	//   ap2 value  : 2

	*ap = 5
	fmt.Println(ap)
	fmt.Println(*ap)
	// change the value at the address &b to 5
	//   ap  address: 0x2101f1020
	//   ap  value  : 5
	//   ap2 address: 0x2101f1020
	//   ap2 value  : 5
	// If this was a reference ap & ap2 would now
	// have different values

	ap = &a
	fmt.Println(ap)
	fmt.Println(*ap)
	// change ap to address of a (&a)
	//   ap  address: 0x2101f1018
	//   ap  value  : 4
	//   ap2 address: 0x2101f1020
	//   ap2 value  : 5
	// Since we've changed the address of ap, it now
	// has a different value then ap2
}
```
