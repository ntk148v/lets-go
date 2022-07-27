package main

import "fmt"

func main() {
	ptr := new(int)
	*ptr = 100
	fmt.Printf("Ptr = %#x, Ptr value = %d\n", ptr, *ptr)
}
