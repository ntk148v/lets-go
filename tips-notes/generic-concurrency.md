# Generic Concurrency

Source: <https://sergey.kamardin.org/articles/generic-concurrency-in-go/>

Table of contents:
- [Generic Concurrency](#generic-concurrency)
	- [TL;DR](#tldr)
	- [1. Introduction](#1-introduction)
		- [1.1. Pre-Generics era](#11-pre-generics-era)
		- [1.2. Generics era](#12-generics-era)
	- [2. Concurrent mapping](#2-concurrent-mapping)
		- [2.1. Context cancellation](#21-context-cancellation)
		- [2.2. Limit concurrency](#22-limit-concurrency)
		- [2.3. Reusing goroutines](#23-reusing-goroutines)
	- [3. Generalisation of `transform()`](#3-generalisation-of-transform)

## TL;DR

Generics and goroutines (and iterators in the future) are great tools we can leverage to have reusable general purpose concurrent processing in our programs.

## 1. Introduction

### 1.1. Pre-Generics era

Let's define the first integer numbers mapping:

```go
func transform([]int, func(int) int) []int

// For example
func transform(xs []int, f func(int) int) []int {
    ret := make([]int, len(xs))
    for i, x := range xs {
        ret[i] = f(x)
    }
    return ret
}

// Example use of such function
// Output: [1, 4, 9]
transform([]int{1, 2, 3}, func(n int) int {
    return n * n
})
```

Now lets assume we want to map integers to strings. That's easy, define `transform()` *just slightly* different:

```go
func transform([]int, func(int) string) []string
// Output: ["1", "2", "3"]
transform([]int{1, 2, 3}, strconv.Itoa)
```

What about reporting whether a number is odd or even? Just another *tiny correction*:

```go
func transform([]int, func(int) bool) []bool
```

### 1.2. Generics era

Thanks to the generics, we now have an ability to *parametrize* functions and types with *type parameters* and define `tranform()` this way:

```go
func transform[A, B any]([]A, func(A) B) []B

# great
transform([]int{1, 2, 3}, square)       // [1, 4, 9]
transform([]int{1, 2, 3}, strconv.Itoa) // ["1", "2", "3"]
transform([]int{1, 2, 3}, isEven)       // [false, true, false]
```

## 2. Concurrent mapping

Getting back to the `tranform()` function. Let's assume that all the calls to `f()` can be done concurrently without breaking our (or anyone else's) program.

```go
func transform[A, B any](as []A, f func(A) B) []B {
	bs := make([]B, len(as))

	var wg sync.WaitGroup
	for i := 0; i < len(as); i++ {
		wg.Add(1)
    // we stasrt a goroutine per each element of the input and call `f(elem)`
		go func(i int) {
			defer wg.Done()
			bs[i] = f(as[i])
		}(i)
	}
	wg.Wait()

	return bs
}
```

### 2.1. Context cancellation

In real world many or even most of the concurrent tasks, especially the i/o related, would be controlled by `context.Context` instance. Since there is a context, there could be timeout or cancellation.

```go
func transform[A, B any](
	ctx context.Context,
	as []A,
	f func(context.Context, A) (B, error),
) (
	[]B,
	error,
) {
	bs := make([]B, len(as))
  // store errors potentially returned by `f()`.
	es := make([]error, len(as))

	subctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	for i := 0; i < len(as); i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			bs[i], es[i] = f(subctx, as[i])
      // If any goroutine's `f()` fails, we cancel the entire `transform()` context
			if es[i] != nil {
				cancel()
			}
		}(i)
	}
	wg.Wait()

	err := errors.Join(es...)
	if err != nil {
		return nil, err
	}
	return bs, nil
}
```

### 2.2. Limit concurrency

In reality, we cannot assume too much about `f()` implicitly. Users of `transform()` might want to limit the number of concurrent calls to `f()`. For example, `f()` can map a url to the result of an http request. Without any limits we can overwhelm the server or get banned ourselves.

At this point we need to switch from using `sync.WaitGroup` to a semaphore `chan`, as we want to control the (maximum) number of simultaneously running goroutines as well as to handle the context cancellation, both by using `select`.

```go
func transform[A, B any](
	ctx context.Context,
	parallelism int,
	as []A,
	f func(context.Context, A) (B, error),
) (
	[]B,
	error,
) {
	bs := make([]B, len(as))
	es := make([]error, len(as))

	// FIXME: if the given context is already cancelled, no worker will be
	// started but the transform() call will return bs, nil.
	subctx, cancel := context.WithCancel(ctx)
	defer cancel()

	sem := make(chan struct{}, parallelism)
sched:
	for i := 0; i < len(as); i++ {
		// We are checking the sub-context cancellation here, in addition to
		// the user-provided context, to handle cases where f() returns an
		// error, which leads to the termination of transform.
		if subctx.Err() != nil {
			break
		}
		select {
		case <-subctx.Done():
			break sched

		case sem <- struct{}{}:
			// Being able to send a tick into the channel means we can start a
			// new worker goroutine. This could be either due to the completion
			// of a previous goroutine or because the number of started worker
			// goroutines is less than the given parallism value.
		}
		go func(i int) {
			defer func() {
				// Signal that the element has been processed and the worker
				// goroutine has completed.
				<-sem
			}()
			bs[i], es[i] = f(subctx, as[i])
			if es[i] != nil {
				cancel()
			}
		}(i)
	}
	// Since each goroutine reads off one tick from the semaphore before exit,
	// filling the channel with artificial ticks makes us sure that all started
	// goroutines completed their execution.
	//
	// FIXME: for the high values of parallelism this loop becomes slow.
	for i := 0; i < cap(sem); i++ {
		// NOTE: we do not check the user-provided context here because we want
		// to return from this function only when all the started worker
		// goroutines have completed. This is to avoid surprising users with
		// some of the f() function calls still running in the background after
		// transform() returns.
		//
		// This implies f() should respect context cancellation and return as
		// soon as its context gets cancelled.
		sem <- struct{}{}
	}

	err := errors.Join(es...)
	if err != nil {
		return nil, err
	}
	return bs, nil
}
```

### 2.3. Reusing goroutines

In the previous iteration we were starting a goroutine per each task, but no more `parallelism` goroutines at a time. Hmm, users might want to have a custom execution context per each goroutine. For example, suppose we have `N` tasks with maximum `P` running concurrently (and `P` can be significantly less than `N`). If each task requires some form of resource preparation, such as a large memory allocation, a database session, or maybe a single-threaded Cgo "coroutine", it would seem logical to prepare only `P` resources and reuse them among workers through context.

```go
func transform[A, B any](
	ctx context.Context,
	prepare func(context.Context) (context.Context, context.CancelFunc),
	parallelism int,
	as []A,
	f func(context.Context, A) (B, error),
) (
	[]B,
	error,
) {
	bs := make([]B, len(as))
	es := make([]error, len(as))

	// FIXME: if the given context is already cancelled, no worker will be
	// started but the transform() call will return bs, nil.
	subctx, cancel := context.WithCancel(ctx)
	defer cancel()

	sem := make(chan struct{}, parallelism)
	// Start up P goroutines, and distribute tasks across them using non-buffered channel wrk.
	// The channel is non-buffered because we want to have an immediate runtime "feedback" to know if there are any idle workers at the momment
	// or if we should consider starting a new one.
	wrk := make(chan int)
sched:
	for i := 0; i < len(as); i++ {
		// We are checking the sub-context cancellation here, in addition to
		// the user-provided context, to handle cases where f() returns an
		// error, which leads to the termination of transform.
		if subctx.Err() != nil {
			break
		}
		select {
		case <-subctx.Done():
			break sched

		case wrk <- i:
			// There is an idle worker goroutine that is ready to process the
			// next element.
			continue

		case sem <- struct{}{}:
			// Being able to send a tick into the channel means we can start a
			// new worker goroutine. This could be either due to the completion
			// of a previous goroutine or because the number of started worker
			// goroutines is less than the given parallism value.
		}
		go func(i int) {
			defer func() {
				// Signal that the element has been processed and the worker
				// goroutine has completed.
				<-sem
			}()

			// Capture the subctx from the dispatch loop. This prevents
			// overriding it if the given prepare() function is not nil.
			subctx := subctx
			if prepare != nil {
				var cancel context.CancelFunc
				subctx, cancel = prepare(subctx)
				defer cancel()
			}
			for {
				bs[i], es[i] = f(subctx, as[i])
				if es[i] != nil {
					cancel()
					return
				}
				var ok bool
				i, ok = <-wrk
				if !ok {
					// Work channel has been closed, which means we will not
					// get any new tasks for this worker and can return.
					break
				}
			}
		}(i)
	}
	// Since each goroutine reads off one tick from the semaphore before exit,
	// filling the channel with artificial ticks makes us sure that all started
	// goroutines completed their execution.
	//
	// FIXME: for the high values of parallelism this loop becomes slow.
	for i := 0; i < cap(sem); i++ {
		// NOTE: we do not check the user-provided context here because we want
		// to return from this function only when all the started worker
		// goroutines have completed. This is to avoid surprising users with
		// some of the f() function calls still running in the background after
		// transform() returns.
		//
		// This implies f() should respect context cancellation and return as
		// soon as its context gets cancelled.
		sem <- struct{}{}
	}

	err := errors.Join(es...)
	if err != nil {
		return nil, err
	}
	return bs, nil
}
```

> As in the previous section, this might be done inside `f()` too, for example, by using `sync.Pool`.

## 3. Generalisation of `transform()`

So far our focus has been on mapping slices, which in many cases is enough. However, what if we want to map `map` types, or maybe `chan` even? Can we map anything that we can range over? And as in for loops, do we always need to map values really?

> TBD
