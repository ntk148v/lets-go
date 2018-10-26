package main

import "fmt"

func ping(pings chan<- string, msg string) {
	fmt.Println("Ping!")
	pings <- msg
}

func pong(pings <-chan string, pongs chan<- string) {
	fmt.Println("Pong!")
	pongs <- <-pings
}

func main() {
	pings := make(chan string, 1)
	pongs := make(chan string, 1)
	ping(pings, "passed message")
	pong(pings, pongs)
	fmt.Println(<-pongs)
}
