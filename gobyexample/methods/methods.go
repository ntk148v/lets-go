package main

import "fmt"

type rect struct {
	width, height int
}

// Methods can be defined for either pointer or value receiver types.
func (r *rect) area() int {
	return r.width * r.height
}

func (r rect) perim() int {
	return 2 * (r.width + r.height)
}

func main() {
	r := rect{width: 26, height: 9}
	fmt.Println("area: ", r.area())
	fmt.Println("perim: ", r.perim())
	// GO automatically handles conversion between values and pointers
	// for method call. You may want to use a pointer receiver type to
	// avoid copying on method calls or to allow the method to
	// mutate the receiving struct.
	rp := &r
	fmt.Println("area: ", rp.area())
	fmt.Print("perim: ", rp.perim())
}
