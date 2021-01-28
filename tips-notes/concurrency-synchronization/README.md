# Concurrency Synchronization

- [Concurrency Synchronization](#concurrency-synchronization)
  - [1. What are Concurrency synchronizations?](#1-what-are-concurrency-synchronizations)
  - [2. Syncronization techniques](#2-syncronization-techniques)
  - [3. Channel](#3-channel)
    - [3.1. Channel types and values](#31-channel-types-and-values)
    - [3.2. Channel operations](#32-channel-operations)
    - [3.3. Notes](#33-notes)
    - [3.3. Use cases](#33-use-cases)

## 1. What are Concurrency synchronizations?

Concurrency synchronizations means how to control concurrent computations:

- to avoid data races between them.
- to avoid them consuming CPU resources when they have nothing to do.

## 2. Syncronization techniques

- Channel.
- Gracefully Close channels.
- `sync` standard package.
- Atomic operations provided in the `sync/atomic` standard package.

## 3. Channel

Source:

- https://go101.org/article/channel.html
- https://go101.org/article/channel-use-cases.html

```
Don't (let computations) communicate by sharing memory, (let them) share memory by communicating (through channels) - Rob Pike -
```

- Channels make goroutines share memory by communicating.
- A FIFO queue.

### 3.1. Channel types and values

- Each channel type has an element type.
- Channel types can be bi-directional or single-directional:
  - `chan T` denotes a bidirectional channel type (send + receive).
  - `chan<- T` denotes a send-only channel type.
  - `<-chan T` denotes a receive-only channel type.
- Each channel value has a capacity. A channel value with a zero capacity is called unbuffered channel and a channel value with a non-zero capacity is called buffered channel.

### 3.2. Channel operations

- 5 operations:

  - Close the channel:

  ```go
  close(ch)
  ```

  - Send a value:

  ```go
  ch <- v
  ```

  - Receive a value:

  ```go
  <-ch
  ```

  - Query tue value buffer capacity of the channel:

  ```go
  // return zero, if ch is nil channel.
  cap(ch)
  ```

  - Query the current number of values in the value buffer (or the length) of the channel:

  ```go
  // return zero, if ch is nil channel.
  len(ch)
  ```

- Most basic operations in Go are not sysnchronized.
- The behaviors for all kind of operations applying on nil, closed and not-closed non-nil channels.

| Operation          | A nil channel  | A closed channel | A not-closed non-nil channel |
| ------------------ | -------------- | ---------------- | ---------------------------- |
| Close              | panic          | panic            | succeed to close             |
| Send value to      | block for ever | panic            | block or succeed to send     |
| Receive value from | block for ever | never block      | block or succeed to receive  |

### 3.3. Notes

- Channel element value are transferred by copy.
- A channel cannot be garbage collected. If a goroutine is blocked and stays in either the sending or the receiving goroutine queue of a chnnel, then goroutine also cannot be be garbage collected, even if the channel is referenced only this goroutine. In fact, a goroutine can only be garbage collected when it has already exited.
- `for-range` on channels: the loop will try to iteratively receive the values sent to a channel, until the channel is closed and its value buffer queue becomes blank.

```go
// if channel is nil, the loop will bolcok there forever
for v := range channel {
    // use v
}
```

- `select-case`:

```go
package main

import "fmt"

func main() {
    var c chan struct{} // nil
    select {
    case <-c:             // blocking operation
    case c <- struct{}{}: // blocking operation
    default:
        fmt.Println("Go here.")
    }
}
```

### 3.3. Use cases

- Use channels as Futures/Promises.

  - Return receive-only channels as results.

    ```go
    package main

    import (
        "time"
        "math/rand"
        "fmt"
    )

    func longTimeRequest() <-chan int32 {
        r := make(chan int32)

        go func() {
            // Simulate a workload.
            time.Sleep(time.Second * 3)
            r <- rand.Int31n(100)
        }()

        return r
    }

    func sumSquares(a, b int32) int32 {
        return a*a + b*b
    }

    func main() {
        rand.Seed(time.Now().UnixNano())

        a, b := longTimeRequest(), longTimeRequest()
        // 3 seconds only
        fmt.Println(sumSquares(<-a, <-b))
    }
    ```

  - Pass send-only channels as arguments.

  ```go
  package main

  import (
      "time"
      "math/rand"
      "fmt"
  )

  func longTimeRequest(r chan<- int32)  {
      // Simulate a workload.
      time.Sleep(time.Second * 3)
      r <- rand.Int31n(100)
  }

  func sumSquares(a, b int32) int32 {
      return a*a + b*b
  }

  func main() {
      rand.Seed(time.Now().UnixNano())

      ra, rb := make(chan int32), make(chan int32)
      go longTimeRequest(ra)
      go longTimeRequest(rb)

      fmt.Println(sumSquares(<-ra, <-rb))
  }
  ```

  - The first response wins.

  ```go
   package main

   import (
       "fmt"
       "time"
       "math/rand"
   )

   func source(c chan<- int32) {
       ra, rb := rand.Int31(), rand.Intn(3) + 1
       // Sleep 1s/2s/3s.
       time.Sleep(time.Duration(rb) * time.Second)
       c <- ra
   }

   func main() {
       rand.Seed(time.Now().UnixNano())

       startTime := time.Now()
       // c must be a buffered channel.
       c := make(chan int32, 5)
       for i := 0; i < cap(c); i++ {
           go source(c)
       }
       // Only the first response will be used.
       rnd := <- c
       fmt.Println(time.Since(startTime))
       fmt.Println(rnd)
   }
  ```

- Use channels for notifications.

// WIP
