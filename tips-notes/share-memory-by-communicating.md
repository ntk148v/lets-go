# Share Memory By Communicating

Traditional threading models (commonly used when writing Java, C++ and Python programs for example) require the programmer to communicate between threads using shared memory. Typically, shared data structes are protected by locks, and threads will contend over those locks to access the data. In some cases, this is made easier by the use of thread-safe data structes such as Python's Queue.

Go's concurrency primitives - goroutines and channels - provide an elegant and distinct means of structuring concurrent software. Instead of explicitly using locks to mediate access to shared data, go encourages the use of channels to pass references to data between goroutines.

_Do not communicate by sharing memory; instead, share memory by communicating._

In a traditional threading environment, one might structure its data like so:

```go
type Resource struct {
    url        string
    polling    bool
    lastPolled int64
}

type Resources struct {
    data []*Resource
    lock *sync.Mutex
}

func Poller(res *Resources) {
    for {
        // get the least recently-polled Resource
        // and mark it as being polled
        res.lock.Lock()
        var r *Resource
        for _, v := range res.data {
            if v.polling {
                continue
            }
            if r == nil || v.lastPolled < r.lastPolled {
                r = v
            }
        }
        if r != nil {
            r.polling = true
        }
        res.lock.Unlock()
        if r == nil {
            continue
        }

        // poll the URL

        // update the Resource's polling and lastPolled
        res.lock.Lock()
        r.polling = false
        r.lastPolled = time.Nanoseconds()
        res.lock.Unlock()
    }
}
```

--> Convert to use Go idiom. Poller is a function that receives Resources to be polled from an input channel, and sends them to an output channel when they're done.

```go
type Resource string

func Poller(in, out chan *Resource) {
    for r := range in {
        // poll the URL

        // send the processed Resource to out
        out <- r
    }
}
```
