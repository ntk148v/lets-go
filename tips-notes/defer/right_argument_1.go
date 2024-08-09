package main

import "fmt"

func pushAnalytic(a int) {
	fmt.Println(a)
}

func main() {
	a := 10
	defer func() {
		pushAnalytic(a)
	}()

	a = 20
}

// Output:
// 20
