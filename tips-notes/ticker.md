# Time Ticker and Sleep

Tickers are for when you want to do something repeatedly at regular intervals.

But what if the `something` task takes time longer than interval?

Let's test it:

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	done := make(chan bool)
	go func() {
		time.Sleep(10 * time.Second)
		done <- true
	}()
	for {
		select {
		case <-done:
			fmt.Println("Done!")
			return
		case t := <-ticker.C:
			fmt.Println("Current time: ", t)
			time.Sleep(time.Second * 10)
			fmt.Println("After a tick 2s", time.Now())
		}
	}
}
```

The result is interesting, the ticker has to wait until the interrupt Sleep is done but the returned tick doesn't include the Sleep time.

```
Current time:  2020-07-10 09:03:37.70405479 +0700 +07 m=+1.000134765
After a tick 2s 2020-07-10 09:03:39.704664613 +0700 +07 m=+3.000744605
Current time:  2020-07-10 09:03:38.704132675 +0700 +07 m=+2.000212636
After a tick 2s 2020-07-10 09:03:41.7048811 +0700 +07 m=+5.000961093
Current time:  2020-07-10 09:03:40.704144127 +0700 +07 m=+4.000224096
After a tick 2s 2020-07-10 09:03:43.705124327 +0700 +07 m=+7.001204329
Current time:  2020-07-10 09:03:42.704140894 +0700 +07 m=+6.000220865
After a tick 2s 2020-07-10 09:03:45.705391382 +0700 +07 m=+9.001471372
Current time:  2020-07-10 09:03:44.704143883 +0700 +07 m=+8.000223851
After a tick 2s 2020-07-10 09:03:47.705639864 +0700 +07 m=+11.001719857
Current time:  2020-07-10 09:03:46.704047296 +0700 +07 m=+10.000127260
After a tick 2s 2020-07-10 09:03:49.705904585 +0700 +07 m=+13.001984508
Done!
```
