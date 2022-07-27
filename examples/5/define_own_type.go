package main

import "fmt"

type NameAge struct {
	name string // both non exported fiedls
	age  int
}

func main() {
	a := new(NameAge)
	a.name = "Pete"
	a.age = 42
	fmt.Printf("%v\n", a)
}
