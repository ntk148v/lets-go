package main

import "fmt"

func main() {
	// By default channels are unbuffered, meaning
	// that they will only accept sends (chan <-) if
	// there is a corresponding receive (<-chan) ready
	// to receive the sent value.
	// Buffered channels accept a limited number
	// of values without a corresponding receiver
	// for those values.
	messages := make(chan string, 2)
	messages <- "buffered"
	messages <- "channel"
	fmt.Println(<-messages)
	fmt.Print(<-messages)
}
