package main

import "fmt"

type Duck interface {
	Quack()
}

type Donald struct {
}

func (d Donald) Quack() {
	fmt.Println("quack quack!")
}

type Daisy struct {
}

func (d Daisy) Quack() {
	fmt.Println("-quack -quack")
}

func sayQuack(duck Duck) {
	duck.Quack()
}

type Dog struct {
}

func (d Dog) Bark() {
	fmt.Println("go go")
}

func main() {
	donald := Donald{}
	sayQuack(donald) // quack
	daisy := Daisy{}
	sayQuack(daisy) // --quack
	dog := Dog()
	sayQuack(dog) // compile error - cannot use dog (type Dog) as type Duck
}
