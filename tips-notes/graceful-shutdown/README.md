# Graceful shutdown

Source: <https://victoriametrics.com/blog/go-graceful-shutdown/>

Table of contents:

- [Graceful shutdown](#graceful-shutdown)
  - [1. Catching the signal](#1-catching-the-signal)
  - [2. Timeout awareness](#2-timeout-awareness)
  - [3. Stop accepting new requests](#3-stop-accepting-new-requests)
  - [4. Handle pending requests](#4-handle-pending-requests)
    - [4.1. Use context middleware to inject cancellation logic](#41-use-context-middleware-to-inject-cancellation-logic)
    - [4.2. Use BaseContext to provide a global context to all connections](#42-use-basecontext-to-provide-a-global-context-to-all-connections)
  - [5. Release critical resources](#5-release-critical-resources)
  - [Summary](#summary)

> [!IMPORTANT]
> Graceful shutdown generally satisfies three minimum conditions:
>
> 1. Close the entrypoint by stopping new requests or messages from sources like HTTP, pub/sub systems, etc. However, keep outgoing connections to third-party services like databases or caches active.
> 2. Wait for all ongoing request to finish. If a request takes too long, respond with a graceful error.
> 3. Release critical resources such as database conections, file locks, or network listeners. Do any final cleanup.
>
> The core principle of graceful shutdown is the same across all systems: **Stop accepting new requests or messages, and give existing operations time to finish within a defined grace period.**

## 1. Catching the signal

- Our applications must be know it's time to exit and begin the shutdown process, and we tell them by using **signal**.
  - [Signals](https://cis.temple.edu/~ingargio/cis307/readings/signals.html) represent a very limited form of interprocess communication.
  - They notify a process that something has happened and it should take action. When a signal is sent, the operating system interrupts the normal flow of the process to deliver the notification.
  - There are three actions that can take place when a signal is delivered to a process:
    - it can be ignored; or
    - the process can be terminated (with or without core dumping); or
    - a handler function can be called. This function receives as its only argument the number identifying the signal it is handling.
- For graceful shutdown, only 3 termination signals are typically important:

```text
SIGNAL      ID   DEFAULT  DESCRIPTION
======================================================================
SIGHUP      1    Termin.  Hang up on controlling terminal (mostly for
                          reload configs)
SIGINT      2    Termin.  Interrupt. Generated when we enter CNRTL-C
                          and it is delivered to all processes/threads
                          associated to the current terminal. If
                          generated with kill, it is delivered to only
                          one process/thread.
SIGTERM    15    Termin.  Software termination signal
```

```go
func main() {
    // signal.Notify tells the Go runtime to deliver specified signals
    // to a channel instead of using the default behavior.
    signalChan := make(chan os.Signal, 1)
    signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

    // Setup work here

    <-signalChan

    fmt.Println("Received termination signal, shutting down...")
}
```

![](https://victoriametrics.com/blog/go-graceful-shutdown/signal-setup-before-init.webp)

- When you press `Ctrl+C` more than once, it doesn't automatically kill the app.
  - The first `Ctrl+C` sends a `SIGINT` to the foreground process.
  - Pressing it again usually sends another `SIGNINT`, not `SIGKILL`.
  - Most shells do not escalate the signal automatically. If you want to force a stop, you must send `SIGKILL` by using `kill -9`.
  - If you want to the second `Ctrl+C` to terminate the app forcefully, you can stop the app from listening to further signals by using `signal.Stop` right after the first signal is received:

```go
func main() {
    signalChan := make(chan os.Signal, 1)
    signal.Notify(signalChan, syscall.SIGINT)

    <-signalChan

    signal.Stop(signalChan)
    select {}
}
```

- Or you can just use `signal.NotifyContext` (Go >= 1.16):

```go
ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
defer stop()

// Setup tasks here

<-ctx.Done()
stop()
```

## 2. Timeout awareness

- It's important to knoow how long your application to shut down after receiving a termination signal.
- Your shutdown must complete within this time, including processing any remaining requests and releasing resources.
- Kubernetes default grace period is 30 seconds, you can use it as starting point and observe to find the best number for your application.

## 3. Stop accepting new requests

- When using `net/http`, you can handle graceful shutdown by calling the `http.Server.Shutdown` method. This method stops the server from accepting new connections and waits for all active requests to complete before shutting down idle connections.
- To avoid connection errors during this short window, the correct strategy is [to fail the readiness probe first](https://github.com/kubernetes/kubernetes/blob/95860cff1c418ea6f5494e4a6168e7acd1c390ec/test/images/agnhost/netexec/netexec.go#L357). This tells the orchestrator that your pod should no longer receive traffic:

```go
var isShuttingDown atomic.Bool

func readinessHandler(w http.ResponseWriter, r *http.Request) {
    if isShuttingDown.Load() {
        w.WriteHeader(http.StatusServiceUnavailable)
        w.Write([]byte("shutting down"))
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("ok"))
}
```

## 4. Handle pending requests

- The server is shutting down gracefully, we need to choose a timeout based on your shutdown budget:

```go
ctx, cancelFn := context.WithTimeout(context.Background(), timeout)
err := server.Shutdown(ctx)
```

- A common issue is that handlers are not automatically aware when the server is shutting down. So, how can we notify our handlers that the server is shutting down? The answer is by using context. There are two main ways to do this:

### 4.1. Use context middleware to inject cancellation logic

```go
func WithGracefulShutdown(next http.Handler, cancelCh <-chan struct{}) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ctx, cancel := WithCancellation(r.Context(), cancelCh)
        defer cancel()

        r = r.WithContext(ctx)
        next.ServeHTTP(w, r)
    })
}
```

### 4.2. Use BaseContext to provide a global context to all connections

- It allows you to create a global context with cancellation that applies to the entire server, and you can cancel it to signal all active requests that the server is shutting down.

```go
ongoingCtx, cancelFn := context.WithCancel(context.Background())
server := &http.Server{
    Addr: ":8080",
    Handler: yourHandler,
    BaseContext: func(l net.Listener) context.Context {
        return ongoingCtx
    },
}

// After attempting graceful shutdown:
cancelFn()
time.Sleep(5 * time.Second) // optional delay to allow context propagation
```

![](https://victoriametrics.com/blog/go-graceful-shutdown/shutdown-context-propagation-timeline.webp)

> [!NOTE]
> All of this work around graceful shutdown won’t help if your functions do not respect context cancellation. Try to avoid using `context.Background()`, `time.Sleep()`, or any other function that ignores context.

## 5. Release critical resources

- A common mistake is releasing critical resources as soon as the termination signal is received. At that point, your handlers and in-flight requests may still be using those resources. _You should delay the resource cleanup until the shutdown timeout has passed or all requests are done_.
- In many cases, simple letting the process exit is enough. The operating system will automatically reclaim resources. However, there are important cases where explicit cleanup is still necessary during shutdown:
  - **Database connections** should be closed properly. If any transactions are still open, they need to be committed or rolled back. Without a proper shutdown, the database has to rely on connection timeouts.
  - **Message queues and brokers** often require a clean shutdown. This may involve flushing messages, committing offsets, or signaling to the broker that the client is exiting. Without this, there can be rebalancing issues or message loss.
  - **External services** may not detect the disconnect immediately. Closing connections manually allows those systems to clean up faster than waiting for TCP timeouts.
- A good rule is to shutdown components in the reverse order of how they were initialized. This respects dependencies between components -> `defer`:

```go
db := connectDB()
defer db.Close()

cache := connectCache()
defer cache.Close()
```

## Summary

Complete example of a graceful shutdown mechanism:

```go
const (
	_shutdownPeriod      = 15 * time.Second
	_shutdownHardPeriod  = 3 * time.Second
	_readinessDrainDelay = 5 * time.Second
)

var isShuttingDown atomic.Bool

func main() {
	// Setup signal context
	rootCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Readiness endpoint
	http.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		if isShuttingDown.Load() {
			http.Error(w, "Shutting down", http.StatusServiceUnavailable)
			return
		}
		fmt.Fprintln(w, "OK")
	})

	// Sample business logic
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-time.After(2 * time.Second):
			fmt.Fprintln(w, "Hello, world!")
		case <-r.Context().Done():
			http.Error(w, "Request cancelled.", http.StatusRequestTimeout)
		}
	})

	// Ensure in-flight requests aren't cancelled immediately on SIGTERM
	ongoingCtx, stopOngoingGracefully := context.WithCancel(context.Background())
	server := &http.Server{
		Addr: ":8080",
		BaseContext: func(_ net.Listener) context.Context {
			return ongoingCtx
		},
	}

	go func() {
		log.Println("Server starting on :8080.")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	// Wait for signal
	<-rootCtx.Done()
	stop()
	isShuttingDown.Store(true)
	log.Println("Received shutdown signal, shutting down.")

	// Give time for readiness check to propagate
	time.Sleep(_readinessDrainDelay)
	log.Println("Readiness check propagated, now waiting for ongoing requests to finish.")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), _shutdownPeriod)
	defer cancel()
	err := server.Shutdown(shutdownCtx)
	stopOngoingGracefully()
	if err != nil {
		log.Println("Failed to wait for ongoing requests to finish, waiting for forced cancellation.")
		time.Sleep(_shutdownHardPeriod)
	}

	log.Println("Server shut down gracefully.")
}
```
