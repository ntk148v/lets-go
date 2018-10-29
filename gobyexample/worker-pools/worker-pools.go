package main

import "fmt"
import "time"

// calculate execution time
func elapsed(what string) func() {
	start := time.Now()
	return func() {
		fmt.Printf("%s took %v\n", what, time.Since(start))
	}
}

// Run several concurrent instances
func worker(id int, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		fmt.Println("worker", id, "started job", j)
		time.Sleep(time.Second)
		fmt.Println("worker", id, "finished job", j)
		results <- j * 2
	}
}

func main() {
	defer elapsed("workers")()
	jobs := make(chan int, 100)
	results := make(chan int, 100)

	// 3 workers
	for w := 1; w <= 3; w++ {
		go worker(w, jobs, results)
	}

	// 5 jobs
	for j := 1; j <= 5; j++ {
		jobs <- j
	}

	// close jobs channel to indicate that's all
	// the work we have
	close(jobs)

	for a := 1; a <= 5; a++ {
		<-results
	}
}
