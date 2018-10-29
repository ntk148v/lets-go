package main

import "fmt"

func main() {
	messages := make(chan string)
	signals := make(chan bool)

	select {
	case msg := <-messages:
		fmt.Println("received messages", msg)
	default:
		fmt.Println("no message received")
	}

	msg := "hi"
	select {
	// msg cannot be sent to the message channel,
	// because the channel has no buffer and there
	// is no receiver.
	case messages <- msg:
		fmt.Println("sent message", msg)
	default:
		fmt.Println("no message sent")
	}

	// multi-way non-blocking select
	select {
	case msg := <-messages:
		fmt.Println("received message", msg)
	case sig := <-signals:
		fmt.Println("received signal", sig)
	default:
		fmt.Println("no activity")
	}
}
