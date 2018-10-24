package main

import "fmt"

// intSeq returns another function, which we
// define anonymously in the body of intSeq.
// The returned function closes over the variable
//i to form a closure
func intSeq() func() int {
	i := 0
	return func() int {
		i++
		return i
	}
}

// call intSq, assigning the result (a func)
// to nextInt. This function value captures its
// own i value, which will be upated each time
// we call nextInt
func main() {
	nextInt := intSeq()
	fmt.Println(nextInt())
	fmt.Println(nextInt())
	fmt.Println(nextInt())

	newInts := intSeq()
	fmt.Println(newInts())
}
