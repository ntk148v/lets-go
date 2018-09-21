package main

import "fmt"

func main() {
	var a = 26.9
	var p = &a

	fmt.Println("a (before) = ", a)
	fmt.Println("p = ", p)
	fmt.Println("*p = ", *p)
	// Change the value stored in the pointed variable through the pointer
	*p = 269
	fmt.Println("a (after) = ", a)
}
