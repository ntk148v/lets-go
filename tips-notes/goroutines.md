# Goroutines vs OS Threads

Source:

- https://codeburst.io/why-goroutines-are-not-lightweight-threads-7c460c1f155f
- https://medium.com/@riteeksrivastava/a-complete-journey-with-goroutines-8472630c7f5c

## Threads

- A thread is just a sequence of instructions that can be executed independently by a processor.
- Modern processors can executed multiple threads at once (multi-threading) and also switch between threads to achieve parallelism.
- Threads share memory and don't need to create a new virtual memory space when they are created and thus don't require a MMU context switch.
- Communication between threads is simpler as they have a shared memory while processes requires various modes of IPC (Inter-Process Communications) like semaphores, messages queues, pipes etc.
- Things make threads slow:
  - Threads consume a lot of memory due to their large stack size (>= 1MB).
  - Threads need to restore a lot registers some of which include AVX, SSE, PC, SP (???) which hurts the application performance.
  - Threads setup and teardown requires call to OS for resources (such as memory) which is slow.
- Threads are scheduled _preemptively_: If a process is running for more than a scheduler time slice, it would prempt the process and schedule execution of another runnable process on the same CPU), the scheduler needs to save/restore all register.

## Goroutines

- The idea of Goroutines was inspired by [Coroutines](https://en.wikipedia.org/wiki/Coroutine).
- Goroutines exists only in the virtual space of Go runtime and not in the OS. Hence, a Go runtime scheduler is needed which manages their lifecycle.
- On startup, Go runtime starts a number of goroutines for GC, scheduler and user code (3 structs: G struct, M struct and Sched struct). An OS thread is created to handle these goroutines. These threads can be at most equal to GOMAXPROCS.

![](https://miro.medium.com/max/933/1*ntxTfMNaxclAE7AJgBuAtw.png)

- _Goroutines are multiplexed onto multiple OS threads so if one should block, such as while waiting for I/O, others continue to run. Their design hides many of the complexities of thread creation and management._
- A goroutine is created with initial only 2KB of stack size.
- Goroutines are scheduled _cooperatively_, they do not directly talk to the OS kernel. When a Goroutine switch occurs, very few registers like program counter and stack pointer need to be saved/restored.
