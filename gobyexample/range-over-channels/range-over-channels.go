package main

import "fmt"

func main() {
	queue := make(chan string, 2)
	queue <- "one"
	queue <- "two"
	close(queue)
	// It's possible to close a non-empty channel but
	// still have the remaining values be received
	for elem := range queue {
		fmt.Println(elem)
	}
}
