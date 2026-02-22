# Benchmarking

Source:

- <https://blog.logrocket.com/benchmarking-golang-improve-function-performance/#:~:text=A%20benchmark%20is%20a%20type,the%20code's%20overall%20performance%20level.>
- <https://dave.cheney.net/2013/06/30/how-to-write-benchmarks-in-go>

- Example:

```go
// main.go
func primeNumbers(max int) []int {
    var primes []int

    for i := 2; i < max; i++ {
        isPrime := true

        for j := 2; j <= int(math.Sqrt(float64(i))); j++ {
            if i%j == 0 {
                isPrime = false
                break
            }
        }

        if isPrime {
            primes = append(primes, i)
        }
    }

    return primes
}
```

- Create test file `_test.go`:
  - `func BenchmarkPrimeNumbers(*testing B)`, `testing.B` type managing the benchmark's timing.
  - `b.N` specifies the number of iterations; the value is not fixed, but dyanmically allocated, ensuring that the benchmark runs for at least one second by default.
- Running a benchmark:
  - `-6` suffix denotes the number of CPUs used to run the benchmark (`GOMAXPROCS`).
  - `509528`: the number of times the loop was executed.
  - `2702 ns/op`: the average amount of time each iteration took to complete, expressed in nanoseconds per operation.
  - Benchmark with various inputs.

```bash
$ go test -bench=.
goos: linux
goarch: amd64
cpu: Intel(R) Core(TM) i7-7700 CPU @ 3.60GHz
BenchmarkPrimeNumbers/input_size_100-6            509528              2702 ns/op
BenchmarkPrimeNumbers/input_size_1000-6            24826             49122 ns/op
BenchmarkPrimeNumbers/input_size_74382-6              72          15083821 ns/op
BenchmarkPrimeNumbers/input_size_382399-6              7         143468607 ns/op
PASS
ok      benchmarking 7.272s
# If there are any unit test functions present in the test files, when you run
# the benchmark, those will also be executed, causing the entire process to take longer
# or the benchmark to fail.
$ go test -bench=. -count 5 -run=^\#
# Adjust the minimum time
$ go test -bench=. -benchtime=10s
# Display memory allocation statistics
$ go test -bench=. -benchtime=10s -benchmem
goos: linux
goarch: amd64
cpu: Intel(R) Core(TM) i7-7700 CPU @ 3.60GHz
BenchmarkPrimeNumbers/input_size_100-6           4916781              2462 ns/op             504 B/op          6 allocs/op
BenchmarkPrimeNumbers/input_size_1000-6           222394             48224 ns/op            4088 B/op          9 allocs/op
BenchmarkPrimeNumbers/input_size_74382-6             764          15728152 ns/op          259320 B/op         18 allocs/op
BenchmarkPrimeNumbers/input_size_382399-6             84         150254140 ns/op         1160443 B/op         23 allocs/op
PASS
ok     benchmarking 61.914s
```
