package main

import "fmt"

func main() {

	// make(map[key-type]val-type)
	m := make(map[string]int)
	m["k1"] = 9
	m["k2"] = 26
	fmt.Println("map:", m)

	v1 := m["k1"]
	fmt.Println("v1:", v1)
	fmt.Println("len:", len(m))

	delete(m, "k2")
	fmt.Println("map:", m)

	_, prs := m["k2"]
	fmt.Println("prs:", prs)
	_, prs = m["k1"]
	fmt.Println("prs:", prs)

	n := map[string]int{"foo": 1, "bar": 2}
	fmt.Println("map:", n)
}
