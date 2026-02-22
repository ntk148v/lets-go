# CPU Throttling for containerized Go applications explained

Source:

- <https://kanishk.io/posts/cpu-throttling-in-containerized-go-apps/>
- <https://github.com/golang/go/issues/33803>

## 1. What could happen when a service experiences CPU throttling?

- Increase latencies in your application.
- Under extremely high traffic your application will stop responding to liveness checks, leading to Kubernetes restarting your pods, which in turn would result in a further pile up of connections on a load balancer which results in overloading the load balancer where it stops serving traffic to other services in the network too.

## 2. The example

```yaml
spec:
  containers:
    - name: app
      image: images.my-company.example/app:v4
      resources:
        requests:
          cpu: "100m"
        limits:
          cpu: "1000m"
```

- 100m translates to 0.1 (or 10%) of a CPU core. 1000m translates to 1 full CPU core.
- Requests define the amount of resources guaranteed to your service.
- Limits define the maximum amount of resources your service can utilize.

## 3. How are CPU limits enforced by Kubernetes?

- Kubernetes uses something known as the Linux Completely Fair Scheduler (CFS) to monitor, ensure, and enforce resource usage and limits.
- The job of CFS is fairly simple:
  - Ensure every container receives CPU time equal to what’s set under `resources.request.cpu`.
  - Prevent any container from using more CPU time than what’s set under `resources.limits.cpu`
- The quota is calculated as:

```text
Quota per period    = CPU Limit * Period

CPU limit is 1000m  = 1 core
Period              = 100ms (default)

Thus, Quota in a single period = 1 * 100 = 100ms
```

## 4. Go concurrency primer

- Go works on the principle of maximizing the concurrency of the application by leveraging parallelism on multi-core hardware.
  - G - Goroutines.
  - P - Processors. This is a Go construct, used to schedule a goroutine on to the machine threads. The Go scheduler is responsible for moving goroutines on and off these Processors. This is the PROC in GOMAXPROCS.
  - M - Machine threads. Also known as OS threads or kernel threads as they operate in kernel space. These are expensive so the GO scheduler is designed to reuse them as much as possible.
    - Go runtime will create Ms = the value of GOMAXPROCS.

![](https://kanishk.io/images/go_concurrency_model.png)

- By default, it’s set to the value returned by runtime.NumCPU(). You can also set it with: `runtime.GOMAXPROCS(n)` and `GOMAXPROCS` environment variable.

## 5. Why would a Go program experience throttling?

- When the application is running inside a container, the Go runtime isn't aware of the fact that the container is allotted just 1 core (1000m = 1 core). It sees all the available cores on the node.
- `runtime.NumCPU()` returns the total number of logical cores on the NODE and not the cores allotted to the CONTAINER! This is because the file inside Linux from which Go reads CPU information is not modified by Linux when a container is created.

![](https://kanishk.io/images/go_svc_throttled.png)

- If Go schedules some goroutines to carry out operations on both cores, they may start using the CPUs, and within the first 50ms of the period, CFS will realize that the service has exhausted its quota of 100ms, and throttle the ongoing CPU usage, i.e. prevent the service from using CPU for the rest of the 50ms remaining in the period.
- This is important because now your service can only resume its work after 50ms when the next period of 100ms starts. This is key to understanding how CPU throttling impacts your application latencies.

## 7. Will setting my requests == limit prevent throttling?

- Setting requests == limits ensures Kubernetes provides a “Guaranteed” QoS for your container.

## 8. A note on limit < 1 and GOMAXPROCS

- Even if you do set GOMAXPROCS but define a value < 1 for CPU limit, you may still end up being throttled if the GO service uses a significant amount of CPU.
- Consider that you’re tuning your GOMAXPROCS=1 (it can’t go any lower than this) and your service has a CPU limit of 100m. Let's quickly calculate the maximum CPU time we'll be allowed to use in a single 100ms period:

```text
CPU Limit = 100m = 0.1 core
CPU Period = 100ms

Quota per period = 0.1 * 100 = 10ms
```

- If your service attempts to use more than 10ms of CPU time, even though it’s using a single core (because GOMAXPROCS=1), it will be throttled if the operation it’s carrying out requires more CPU time than your quota. Assume one of these operations takes 30ms, your CPU usage would then look something like this:

![](https://kanishk.io/images/fractional_cpu_svc_throttled.png)

## 9. I use CPU limits, how do I set GOMAXPROCS appropriately?

- [automaxproces](https://github.com/uber-go/automaxprocs). This is a one-line import into the cmd of your application which sets your GOMAXPROCS at startup time by reading the cgroup quota information inside the container.

```go
import _ "go.uber.org/automaxprocs"

func main() {
  // Your application logic here.
}
```

- Via resourceFieldRef:

```yaml
env:
  - name: GOMAXPROCS
    valueFrom:
      resourceFieldRef:
        resource: limits.cpu
```

- Wait for Go runtime to be CFS aware: [This](https://github.com/golang/go/issues/33803) is a long-standing open issue with the maintainers of Go.
