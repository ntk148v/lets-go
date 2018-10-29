package main

import "fmt"
import "time"

func main() {
	// Timer represent a single event in the future
	timer1 := time.NewTimer(2 * time.Second)
	// blocks on the timer;s channel C until it sends
	// a value indicating that the timer expired
	// If you just wanted to wait, you could have used
	// time.Sleep.
	<-timer1.C
	fmt.Println("Timer 1 expired")

	// Cancel the timer before it expires
	timer2 := time.NewTimer(2 * time.Second)
	go func() {
		<-timer2.C
		fmt.Println("Timer 2 expired")
	}()

	stop2 := timer2.Stop()
	if stop2 {
		fmt.Println("Timer 2 stopped")
	}
}
