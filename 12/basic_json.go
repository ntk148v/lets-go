package main

import (
	"encoding/json"
	"fmt"
)

type Person struct {
	Name string
	Age  int
}

func main() {
	bob := &Person{
		Name: "Bob",
		Age:  20,
	}
	bobRaw, _ := json.Marshal(bob)
	fmt.Println(string(bobRaw))

	aliceRaw := []byte(`{"name": "Alice", "age": 23}`)
	var alice Person

	if err := json.Unmarshal(aliceRaw, &alice); err != nil {
		panic(err)
	}
	fmt.Println(alice)
}
