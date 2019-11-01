# Debugging performance issues

https://rakyll.org/custom-profiles/

https://software.intel.com/en-us/blogs/2014/05/10/debugging-performance-issues-in-go-programs

## Builtin profiles provided by the runtime/pprof

* **profile**: CPU profile determines where a program spends its time while actively consuming CPU cycles.
* **heap**: Heap profile reports the currently live allocations; used to monitor current memory usage or check for memory leaks.
* **threadcreate**: Thread creation profile reports the sections of the program that lead the creation of new OS threads.
* **goroutine**: Goroutine profile report the stack traces of all current goroutines.
* **block**: Block profile show where goroutines block waiting on synchronization primitives (including time channels) (disabled by default).
* **mutex**: Mutex profile reports the lock contentions. When you think your CPU is not fully utilized due to a mutex contention, use this profile (disabled by default).

...
// WIP WIP WIP
