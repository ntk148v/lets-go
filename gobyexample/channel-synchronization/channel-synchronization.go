/* Use channels to synchronize execution accross goroutines */
package main

import "fmt"
import "time"

func worker(done chan bool) {
	fmt.Print("Working...")
	time.Sleep(time.Second)
	fmt.Println("Done")
	done <- true
}

func main() {
	done := make(chan bool, 1)
	go worker(done)
	fmt.Println("Waiting...")
	<-done
}
