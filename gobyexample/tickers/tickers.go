/* Timers are for when you want to do something once in
the future - tickers are for when you want to do something
repeatedly at regular intervals. */
package main

import "fmt"
import "time"

func main() {
	// Value - every 500ms
	ticker := time.NewTicker(900 * time.Millisecond)
	go func() {
		for t := range ticker.C {
			fmt.Println("Tick at", t)
		}
	}()

	time.Sleep(2600 * time.Millisecond)
	ticker.Stop()
	fmt.Println("Ticker stopped")
}
