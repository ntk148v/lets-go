/* Create a simple stack which can hold a fixed number of ints */

package main

import "fmt"

type stack struct {
	i    int
	data [10]int
}

func (s *stack) push(k int) { // Must be *stack - a pointer to the stack. If not, the push function gets a copy of s
	s.data[s.i] = k
	s.i++
}

func (s *stack) pop() int {
	s.i--
	ret := s.data[s.i]
	s.data[s.i] = 0
	return ret
}

func main() {
	var s stack
	s.push(26)
	s.push(9)
	fmt.Printf("stack %v\n", s)
}
