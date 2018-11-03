/* The primary mechanism for managing state in Go
is communication over channels.
*/
package main

import "fmt"
import "time"
import "sync/atomic"

func main() {
	var ops uint64
	for i := 0; i < 50; i++ {
		// TO simulate concurrent updtes, start
		// 50 goroutines that each increment the
		// counter about once a millisecond
		go func() {
			for {
				atomic.AddUint64(&ops, 1)
				time.Sleep(time.Millisecond)
			}
		}()
	}

	time.Sleep(time.Second)

	// In order to safely use the counter while
	// it's still being updated by other goroutines
	// extract a copy of the current value.
	opsFinal := atomic.LoadUint64(&ops)
	fmt.Println("ops:", opsFinal)
}
