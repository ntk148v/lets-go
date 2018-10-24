package main

import "fmt"

func plus(a int, b int) int {
	return a + b
}

// When you have multiple consecutive parameters
// the same type, you may omit the type name
// for the like-typed parameters
// up to the final parameter that declares the type
func plusPlus(a, b, c int) int {
	return a + b + c
}

func main() {
	res := plus(26, 9)
	fmt.Println("26+9 =", res)

	res = plusPlus(26, 9, 1994)
	fmt.Println("26+9+1994 =", res)
}
