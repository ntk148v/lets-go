# Pipelines and cancellation

Source: <https://go.dev/blog/pipelines>

- Go’s concurrency primitives make it easy to construct streaming data pipelines that make efficient use of I/O and multiple CPUs.
- A series of stages connected by channels, where each stage is a gourp of goroutines running in the same function. In each stage, the goroutines:
  - receive values from upstream via inbound channels
  - perform some function on that data, usually producing new values
  - send values downstream via outbound channels
- Example: A pipeline with 3 stages

```go
// First stage - `gen`: a function that converts a list of integers to a channel that emits the integers in the list.
func gen(nums ...int) <-chan int {
    out := make(chan int)
    go func() {
        for _, n := range nums {
            out <- n
        }
        close(out)
    }()
    return out
}

// Second stage - `sq`: receives integers from a channel and returns a channel that emits the square of each received integer.

func sq(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        for n := range in {
            out <- n * n
        }
        close(out)
    }()
    return out
}

// Third stage - `main`: sets up the pipeline and runs the final stage; receives values from the second stage and prints each one, until the channel is closed.o
func main() {
    // Set up the pipeline.
    c := gen(2, 3)
    out := sq(c)

    // Consume the output.
    fmt.Println(<-out) // 4
    fmt.Println(<-out) // 9
}
```

## 1. Fan-out, fan-in

- Fan-out: multiple functions can read from the same channel until that channel is closed.
- Fan-in: A function can read from multiple inputs and proceed until all are closed by multiplexing the input channel onto a single channel that's closed when all the inputs are closed.
- Run two instance of `sq`, each reading from the same input channel.

```go
func main() {
    in := gen(2, 3)

    // Distribute the sq work across two goroutines that both read from in.
    c1 := sq(in)
    c2 := sq(in)

    // Consume the merged output from c1 and c2.
    for n := range merge(c1, c2) {
        fmt.Println(n) // 4 then 9, or 9 then 4
    }
}

// merge function converts a list of channels to a single channel by starting
// a goroutine for each inbound channel that copies tha values to the sole
// outbound channel
func merge(cs ...<-chan int) <-chan int {
    var wg sync.WaitGroup
    out := make(chan int)

    // Start an output goroutine for each input channel in cs.  output
    // copies values from c to out until c is closed, then calls wg.Done.
    output := func(c <-chan int) {
        for n := range c {
            out <- n
        }
        wg.Done()
    }
    wg.Add(len(cs))
    for _, c := range cs {
        go output(c)
    }

    // Start a goroutine to close out once all the output goroutines are
    // done.  This must start after the wg.Add call.
    go func() {
        wg.Wait()
        close(out)
    }()
    return out
}
```

## 2. Stopping short

- There is a pattern to our pipeline functions:
  - stages close their outbouind channels when all the send opreations are done
  - stages keep receiving values from inbound channels until those channels are closed
- But in real pipelines, stages don't always receive all the inbound values - sometimes, the receiver may only need a subset of values to make progress.
- In our example pipeline, if a stage fails to consume all the inbound values, the goroutines attempting to send those values will block indefinitely -> a resource leak.

```go
    // Consume the first value from the output.
    out := merge(c1, c2)
    fmt.Println(<-out) // 4 or 9
    return
    // Since we didn't receive the second value from out,
    // one of the output goroutines is hung attempting to send it.
}
```

- Arrange for the upstream stages to exit even when the downstream stages fail to receive all the inbound -> change the outbound channels to have a buffer. A buffer can hold a fixed number of values; send operations complete immediately if there's room in the buffer.

```go
c := make(chan int, 2) // buffer size 2
c <- 1  // succeeds immediately
c <- 2  // succeeds immediately
c <- 3  // blocks until another goroutine does <-c and receives 1
```

- When the number of values to be sent is known -> OK

```go
func gen(nums ...int) <-chan int {
    out := make(chan int, len(nums))
    for _, n := range nums {
        out <- n
    }
    close(out)
    return out
}

func merge(cs ...<-chan int) <-chan int {
    var wg sync.WaitGroup
    out := make(chan int, 1) // enough space for the unread inputs
    // ... the rest is unchanged ...
```

- If we pass an additional value to gen, or if the downstream stage reads any fewer values, we will again have blocked goroutines -> NOT OK
- We need to provide a way for downstream stages to indicate to the senders that they will stop accepting input.

## 3. Explicit cancellation

- When `main` decides to exit without receiving all the values from `out`, it must tell the goroutines in the upstream stages to abandon the values.

```go
func main() {
    in := gen(2, 3)

    // Distribute the sq work across two goroutines that both read from in.
    c1 := sq(in)
    c2 := sq(in)

    // Consume the first value from output.
    done := make(chan struct{}, 2)
    out := merge(done, c1, c2)
    fmt.Println(<-out) // 4 or 9

    // Tell the remaining senders we're leaving.
    done <- struct{}{}
    done <- struct{}{}
}
```

- Replace send operation with a `select` statement.

  - Problem: Each downstream receiver needs to know the number of potentially blocked upstream senders and arrange to signal those senders on early return.

  ```go
  func merge(done <-chan struct{}, cs ...<-chan int) <-chan int {
      var wg sync.WaitGroup
      out := make(chan int)

      // Start an output goroutine for each input channel in cs.  output
      // copies values from c to out until c is closed or it receives a value
      // from done, then output calls wg.Done.
      output := func(c <-chan int) {
          for n := range c {
              select {
              case out <- n:
              case <-done:
              }
          }
          wg.Done()
      }
      // ... the rest is unchanged ...
  ```

  - Simply close the `done` channel. This close is effectively a broadcasts signal to the senders.

  ```go
  func main() {
    // Set up a done channel that's shared by the whole pipeline,
    // and close that channel when this pipeline exits, as a signal
    // for all the goroutines we started to exit.
    done := make(chan struct{})
    defer close(done)

    in := gen(done, 2, 3)

    // Distribute the sq work across two goroutines that both read from in.
    c1 := sq(done, in)
    c2 := sq(done, in)

    // Consume the first value from output.
    out := merge(done, c1, c2)
    fmt.Println(<-out) // 4 or 9

    // done will be closed by the deferred call.
  }
  ```

  - Ensure `wg.Done` is called on all return paths.

  ```go
  func merge(done <-chan struct{}, cs ...<-chan int) <-chan int {
    var wg sync.WaitGroup
    out := make(chan int)

    // Start an output goroutine for each input channel in cs.  output
    // copies values from c to out until c or done is closed, then calls
    // wg.Done.
    output := func(c <-chan int) {
        defer wg.Done()
        for n := range c {
            select {
            case out <- n:
            case <-done:
                return
            }
        }
    }
    // ... the rest is unchanged ...
  ```

  - `sq` can return as soon as done is closed.

  ```go
  func sq(done <-chan struct{}, in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for n := range in {
            select {
            case out <- n * n:
            case <-done:
                return
            }
        }
    }()
    return out
  }
  ```
