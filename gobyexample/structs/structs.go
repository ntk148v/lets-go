package main

import "fmt"

type person struct {
	name string
	age  int
}

func main() {
	fmt.Println(person{"Kien", 25})
	fmt.Println(person{name: "Another Kien", age: 25})
	fmt.Println(person{name: "Another another Kien"})
	// An & prefix yields a pointer to the struct
	fmt.Println(&person{name: "Another another another Kien", age: 26})
	s := person{name: "Not Kien anymore", age: 27}
	fmt.Println(s.name)
	// You can also use dots with struct pointers - the pointers
	// automatically deferenced.
	sp := &s
	fmt.Println(s.age)
	sp.age = 28
	fmt.Println(sp.age)
	fmt.Println(s.age)
}
