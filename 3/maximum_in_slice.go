package main

import "fmt"

func max(l []int) (max int) {
	max = l[0]
	for _, v := range l {
		if v > max {
			max = v
		}
	}

	return
}

func main() {
	l := []int{1, 2, 5, 7, 0, 12, 1}
	fmt.Println("Max: ", max(l))
}
