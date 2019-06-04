# Golang WaitGroups tutorial

[Source](https://tutorialedge.net/golang/go-waitgroup-tutorial/)

When you starting writing applications in Go that leverage `goroutines`, you start hitting scenarios where you need to block the execution of certain parts of your code base, until these `goroutines` have successfully executed.

## The Solution - WaitGroups

WaitGroups essentially allow us to tackle this problem by blocking until any goroutines within that `WaitGroup` have successfully executed.

We first call `.Add(1)` (before execute your `goroutines`) on our `WaitGroup` to set the number of `goroutines` we want to wait for, and subsequently, we call `.Done()` within any `goroutines` to signal the end of its execution.
