package main

import "fmt"

// zeroval will get a copy of ival distinct
// from the one in the calling function
func zeroval(ival int) {
	ival = 0
}

// zeroptr has an *int parameter, meaning that
// it takes an int pointer. The *iptr code in the
// function body hten dereferences the pointer from
// its memory address to the current value at the address.
// Assigning a value to a dereferenced pointer changes
// the value at the referenced address
func zeroptr(iptr *int) {
	*iptr = 0
}

func main() {
	i := 1
	fmt.Println("initial:", i)
	zeroval(i)
	fmt.Println("zeroval:", i)
	zeroptr(&i)
	fmt.Println("zeroptr:", i)
	fmt.Println("pointer:", &i)
}
