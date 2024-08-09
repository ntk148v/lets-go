package main

import "fmt"

func pushAnalytic(a *int) {
	fmt.Println(*a)
}

func main() {
	a := 10
	defer pushAnalytic(&a)

	a = 20
}
