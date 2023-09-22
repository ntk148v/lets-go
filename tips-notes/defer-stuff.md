# Defer

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
