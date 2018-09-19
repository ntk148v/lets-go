/* A mapF function is a function that takes a function and a list.
The function is applied to each member in the list and a new list
containing these calculated values is returned. */
package main

import "fmt"

func mapF(f func(int) int, l []int) []int {
	j := make([]int, len(l))
	for k, v := range l {
		j[k] = f(v)
	}

	return j
}

func main() {
	m := []int{2, 3, 4}
	f := func(i int) int {
		return i * i
	}

	fmt.Println("Map function results:", (mapF(f, m)))
}
