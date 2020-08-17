# Debugging performance issues

https://jvns.ca/blog/2017/09/24/profiling-go-with-pprof/
https://rakyll.org/custom-profiles/
https://software.intel.com/en-us/blogs/2014/05/10/debugging-performance-issues-in-go-programs

## builtin profiles provided by the runtime/pprof

A Profile is a collection of stack traces showing the call sequences that led to instances of a particular event, such as allocation. Packages can create and maintain their own profiles; the most common use is for tracking resources that must be explicitly closed, such as files or network connections.

- **profile**: CPU profile determines where a program spends its time while actively consuming CPU cycles.
- **heap**: Heap profile reports the currently live allocations; used to monitor current memory usage or check for memory leaks.
- **threadcreate**: Thread creation profile reports the sections of the program that lead the creation of new OS threads.
- **goroutine**: Goroutine profile report the stack traces of all current goroutines.
- **block**: Block profile show where goroutines block waiting on synchronization primitives (including time channels) (disabled by default).
- **mutex**: Mutex profile reports the lock contentions. When you think your CPU is not fully utilized due to a mutex contention, use this profile (disabled by default).

## pprof basics

- Setup a webserver for getting Go profiles:

  ```go
  import _ "net/http/pprof"
  ```

- Go to `localhost:$PORT/debug/pprof` to get a list of available profiles.
- Use `go tool pprof` to analyze profile.

  ```bash
  go tool pprof localhost:$PORT/debug/pprof/$PROFILE_TYPE
  # Generate a svg
  go tool pprof -svg localhost:$PORT/debug/pprof/$PROFILE_TYPE > $PROFILE_TYPE.svg
  ```

**NOTE**: The trace endpoint (`/debug/pprof/trace?seconds=5`), unlike all the rest, outputs a file that is **not** a pprof, it's a **trace** and you can view it using `go tool trace`.

## heap profile

- `alloc_space` vs `inuse_space`:

```
-inuse_space      Display in-use memory size
-inuse_objects    Display in-use object counts
-alloc_space      Display allocated memory size
-alloc_objects    Display allocated object counts
```

// WIP WIP...
