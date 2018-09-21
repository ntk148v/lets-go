package main

import "fmt"

func main() {
	var a = 26.9
	var p = &a

	fmt.Println("Value stored in variable a = ", a)
	fmt.Println("Address of variable a = ", &a)
	fmt.Println("Value stored in variable p = ", p)
}
