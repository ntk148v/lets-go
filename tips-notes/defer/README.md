# Defer

## 1. Interesting

What is the output of the following code?

```go
// https://go.dev/play/p/6lmhN0yFm9w
package main

import "fmt"

func main() {
	i := 0
	i++
	defer fmt.Println(i + 3*i)
	i++
}
```

You might think it's **8**, me either. But actually, the result is **4**. But why?

A defer statement defers execution of a function until the surrounding function returns. The deferred call's arguments are evaluated immediately, but the function call is not executed until the surrounding function returns. Parameters including method receivers, that is.

Therefore, to return **8**, use a closure.

```go
// https://go.dev/play/p/NixNkiVvnEE
package main

import "fmt"

func main() {
	i := 0
	i++
	defer func() { fmt.Println(i + 3*i) }()
	i++
}
```

## 2. Defer: From Basic to Traps

Source: <https://victoriametrics.com/blog/defer-in-go/index.html>

### 2.1. Defers are stacked

- When you use multiple `defer` statements in a function, they are executed in a 'stack' order, meaning the last deferred function is executed first.

```go
func main() {
  defer fmt.Println(1)
  defer fmt.Println(2)
  defer fmt.Println(3)
}

// Output:
// 3
// 2
// 1
```

- Every time you call a defer statement, you’re adding that function to the top of the current goroutine’s linked list, like this:

![](https://victoriametrics.com/blog/defer-in-go/defer-chain.png)

- It does not execute all the defer in the linked list of goroutine, it's only run the defer in the returned function, because our defer linked list could contain many defers from many different functions.
  - So, only the deferred functions in the current function (or current stack frame) are executed.

```go
func B() {
  defer fmt.Println(1)
  defer fmt.Println(2)
  A()
}

func A() {
  defer fmt.Println(3)
  defer fmt.Println(4)
}

// 4
// 3
// 2
// 1
```

![](https://victoriametrics.com/blog/defer-in-go/defer-frame.png)

### 2.2. Defer, Panic and Recover

- Panic is a way to stop the execution of the current goroutine, unwind the stack, and execute the deferred functions in the current goroutine, causing our application to crash.
- To handle unexpected errors and prevent the application from crashing, you can use the `recover` function within a deferred function to regain control of a panicking goroutine.

```go
func main() {
  defer func() {
    if r := recover(); r != nil {
      fmt.Println("Recovered:", r)
    }
  }()

  panic("This is a panic")
}

// Output:
// Recovered: This is a panic
```

- Common mistakes:

  - Use `recover` directly as a deferred function:

  ```go
  func main() {
  	// This code still panics, the recover function is meant to catch a panic
   	// but it has to be called within a deferred function to work properly.
    // our call to `recover` is actually the `runtime.gorecover`, and it
    // checks that the recover call is happening in the right context, specifically
    // from the correct deferred function that was active when the panic occurred.
  	defer recover()
   	panic("This is a panic")
  }
  ```

  - Our call to `recover` is actually the `runtime.gorecover`, and it checks that the recover call is happening in the right context, specifically from the correct deferred function that was active when the panic occurred.

  ```go
  // It won't work too, because recover isn't called directly from a deferred function
  // but from a nested function
  func myRecover() {
    if r := recover(); r != nil {
      fmt.Println("Recovered:", r)
    }
  }

  func main() {
      defer func() {
      myRecover()
      // ...
      }()

      panic("This is a panic")
  }
  ```

  - Try to catch a panic from different goroutine:

  ```go
  func myRecover() {
    if r := recover(); r != nil {
      fmt.Println("Recovered:", r)
    }
  }

  func main() {
    defer func() {
      myRecover()
      // ...
    }()

    panic("This is a panic")
  }
  ```

  ### 2.3. Defer arguments, including receiver are immediately evaluated

- The following code output is 10, not 20. That's because when you use the defer statement, it grabs the values right then. This is called "capture by value". So, the value of `a` that gets sent to `pushAnalytic` is set to 10 when the defer is scheduled, even though `a` changes later.

  ```go
  package main

  import "fmt"

  func pushAnalytic(a int) {
  fmt.Println(a)
  }

  func main() {
  a := 10
  defer pushAnalytic(a)

  a = 20
  }

  // Output:
  // 10
  ```

- How to fix it?

  - Use a closure. This means wrapping the deferred function call inside another function. That way, you capture the variable by reference, not by value like before.

  ```go
  package main

  import "fmt"

  func pushAnalytic(a int) {
  	fmt.Println(a)
  }

  func main() {
    a := 10
    defer func() {
    	pushAnalytic(a)
    }()

    a = 20
  }

  // Output:
  // 20
  ```

  - Pass the memory address of the variable instead of its value.

  ```go
  package main

  import "fmt"

  func pushAnalytic(a *int) {
  	fmt.Println(*a)
  }

  func main() {
    a := 10
    defer pushAnalytic(&a)

    a = 20
  }

  // Output: 20
  ```

- The trap:
  - The output is 10, just like before.

```go
package main

import "fmt"

type Data struct {
	a int
}

func (d Data) pushAnalytic() {
	fmt.Println(d.a)
}

func main() {
	d := Data{a: 10}
	defer d.pushAnalytic()

	d.a = 20
}

// Output:
// 10
```

- This happens because the defer argument also evaluates its receiver immediately, capturing the value of `d` at that moment. Under the hood, the receiver is like an argument, so the defer statement works like this:

```go
defer Data.pushAnalytic(d) // defer d.pushAnalytic()
```

- So, the same rule applies: the arguments of the deferred function are evaluated right away.
- How to fix:

  - Use a closure.
  - Use a pointer but you need to change a bit.

  ```go
  d := &Data{}
  defer Data.PushAnalytic(*d)

  func (d *Data) pushAnalytic() {
    fmt.Println(d.a)
  }
  ```

### 2.4. Defer with error handling

- It is a good illustration point to show how defer works, but it’s also a bad example of how to use defer.

```go
func doSomething() error {
  f, err := os.Open("phuong-secrets.txt")
  if err != nil {
    return err
  }
  defer f.Close()

  // ...
}
```

- The problem is that if we use `defer f.Close()`, we miss the chance to handle the error gracefully because the `Close` method returns an error, but we miss it.

```go
func doSomething() (err error) {
  f, err := os.Open("phuong-secrets.txt")
  if err != nil {
    return err
  }
  defer func() {
    err = errors.Join(err, f.Close())
  }()

  // ...
}
```

### 2.5. Defer types: Heap-allocated, Stack-allocated and Open-coded defer

- When we call `defer`, we’re creating a structure called a defer object `_defer`, which holds all the necessary information about the deferred call. This object gets pushed into the goroutine’s defer chain.
- The difference between heap-allocated and stack-allocated types is where the defer object is allocated. Below Go 1.13, we only had heap-allocated defer.

![](https://victoriametrics.com/blog/defer-in-go/defer-3-types.png)

- Currently, in Go 1.22, if you use defer in a loop, it will be heap-allocated.

```go
func main() {
  for i := 0; i < unpredictableNumber; i++ {
    defer fmt.Println(i) // Heap-allocated defer
  }
}
```

- The heap allocation here is necessary because the number of defer objects can change at runtime. So, the heap ensures that the program can handle any number of defers, no matter how many or where they appear in the function, without bloating the stack.
- If the defer statement within the if block is invoked only once and not in a loop or another dynamic context, it benefits from the optimization introduced in Go 1.13, meaning the `defer` object will be stack-allocated.

```go
func testDefer(a int) {
	if a == unpredictableNumber {
		defer println("Defer in if") // stack-allocated defer
	}
	if a == unpredictableNumber+1 {
		defer println("Defer in if") // stack-allocated defer
	}

  for range a {
    defer println("Defer in for") // heap-allocated defer
  }
}
```

- With this optimization, according to the [Open-coded defers proposal](https://github.com/golang/proposal/blob/master/design/34481-opencoded-defers.md), in the cmd/go binary, this optimization applies to 363 out of 370 static defer sites. As a result, these sites see a 30% performance improvement compared to the previous approach where defer objects were heap-allocated. If it’s that good, why do we need something called ‘open-coded defer’?
- What if we just put the defer at the end of the function? The performance of a direct call is much better than the other two. As of Go 1.13, most defer operations take about 35ns (down from about 50ns in Go 1.12). In contrast, a direct call takes about 6ns.
- Go will inline our defer call directly at the end of the function and also before every return statement in the assembly code, but there are some restrictions for this type to be applied.
- If a function has at least one heap-allocated defer, any defer in the function will NOT be inlined or open-coded. That means, to optimize the above function, we should remove or move the heap-allocated defer elsewhere.

```go
func testDefer(a int) {
	if a == unpredictableNumber {
		defer println("Defer in if") // open-coded defer
	}
	if a == unpredictableNumber+1 {
		defer println("Defer in if") // open-coded defer
	}
}
```

- Another thing to keep in mind is that the product of the number of defers in the function and the number of return statements needs to be 15 or less to fit into this category.
