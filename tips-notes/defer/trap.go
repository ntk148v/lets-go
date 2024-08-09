package main

import "fmt"

type Data struct {
	a int
}

func (d Data) pushAnalytic() {
	fmt.Println(d.a)
}

func main() {
	d := Data{a: 10}
	defer d.pushAnalytic()

	d.a = 20
}

// Output:
// 10
