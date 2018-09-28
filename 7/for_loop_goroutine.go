package main

import "fmt"

func shower(c chan int, quit chan bool) {
	for {
		select {
		case j := <-c:
			fmt.Printf("%d\n", j)
		case <-quit:
			break
		}
	}
}

func main() {
	ch := make(chan int)
	quit := make(chan bool)
	go shower(ch, quit)
	for i := 0; i < 10; i++ {
		ch <- i
	}
	quit <- false // or true, does not matter
}
