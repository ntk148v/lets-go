/* Rate limiting is an important mechanism for
controlling resource utilization and maintaining
quality of service. */
package main

import "fmt"
import "time"

func main() {
	requests := make(chan int, 5)
	for i := 1; i <= 5; i++ {
		requests <- i
	}
	close(requests)

	limiter := time.Tick(200 * time.Millisecond)
	for req := range requests {
		<-limiter
		fmt.Println("request", req, time.Now())
	}

	fmt.Println("------------------------------")

	// Allow short bursts of requests in our
	// rate limiting scheme while presereving
	// the overall rate limit.
	burstyLimiter := make(chan time.Time, 3)

	// Fill up the channel to represent allowed
	// bursting.
	for i := 0; i < 3; i++ {
		burstyLimiter <- time.Now()
	}

	// Every 200 milliseconds we'll try to
	// add a new value to burstyLimiter, up
	// to its limit of 3.
	go func() {
		for t := range time.Tick(200 * time.Millisecond) {
			burstyLimiter <- t
		}
	}()

	burstyRequests := make(chan int, 5)
	for i := 1; i <= 5; i++ {
		burstyRequests <- i
	}
	close(burstyRequests)
	for req := range burstyRequests {
		// server the first 3 immediately because
		// of the burstable rate limiting
		// then serve the remaining 2 with ~200ms delays each
		<-burstyLimiter
		fmt.Println("request", req, time.Now())
	}
}
