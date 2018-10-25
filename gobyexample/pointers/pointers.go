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

var ap *int

func main() {
	i := 1
	fmt.Println("initial:", i)
	zeroval(i)
	fmt.Println("zeroval:", i)
	zeroptr(&i)
	fmt.Println("zeroptr:", i)
	fmt.Println("pointer:", &i)

	fmt.Println("---------------------------------------")

	a := 1 // define int
	b := 2 // define int

	ap = &a
	fmt.Println("set ap to address of a (&a)")
	//   ap address: 0x2101f1018
	//   ap value  : 1
	fmt.Println("ap address:", ap)
	fmt.Println("ap value:  ", *ap)

	*ap = 3
	fmt.Println("change the value at address &a to 3")
	//   ap address: 0x2101f1018
	//   ap value  : 3
	fmt.Println("ap address:", ap)
	fmt.Println("ap value:  ", *ap)

	a = 4
	fmt.Println("change the value of a to 4")
	//   ap address: 0x2101f1018
	//   ap value  : 4
	fmt.Println("ap address:", ap)
	fmt.Println("ap value:  ", *ap)

	ap = &b
	fmt.Println("set ap to the address of b (&b)")
	//   ap address: 0x2101f1020
	//   ap value  : 2
	fmt.Println("ap address:", ap)
	fmt.Println("ap value:  ", *ap)

	ap2 := ap
	fmt.Println("set ap2 to the address in ap")
	//   ap  address: 0x2101f1020
	//   ap  value  : 2
	//   ap2 address: 0x2101f1020
	//   ap2 value  : 2
	fmt.Println("ap address: ", ap)
	fmt.Println("ap value:   ", *ap)
	fmt.Println("ap2 address:", ap2)
	fmt.Println("ap2 value:  ", *ap2)

	*ap = 5
	fmt.Println("change the value at the address &b to 5")
	//   ap  address: 0x2101f1020
	//   ap  value  : 5
	//   ap2 address: 0x2101f1020
	//   ap2 value  : 5
	// If this was a reference ap & ap2 would now
	// have different values
	fmt.Println("ap address: ", ap)
	fmt.Println("ap value:   ", *ap)
	fmt.Println("ap2 address:", ap2)
	fmt.Println("ap2 value:  ", *ap2)

	ap = &a
	fmt.Println("change ap to address of a (&a)")
	//   ap  address: 0x2101f1018
	//   ap  value  : 4
	//   ap2 address: 0x2101f1020
	//   ap2 value  : 5
	// Since we've changed the address of ap, it now
	// has a different value then ap2
	fmt.Println("ap address: ", ap)
	fmt.Println("ap value:   ", *ap)
	fmt.Println("ap2 address:", ap2)
	fmt.Println("ap2 value:  ", *ap2)
}
