# Web Programming

[Go Web Example](https://gowebexamples.com)

Table of Contents:

- [Web Programming](#web-programming)
  - [1. HTTP Server](#1-http-server)
  - [2. Templating](#2-templating)
  - [3. Requests and Forms](#3-requests-and-forms)
  - [4. Assets and Files](#4-assets-and-files)
  - [5. Middleware (Basic)](#5-middleware-basic)
  - [6. Middleware (Advanced)](#6-middleware-advanced)
  - [7. Session](#7-session)
  - [8. Websockets](#8-websockets)

## 1. HTTP Server

A basic HTTP server has a few key jobs to take care of:

- _Process dynamic request_: Process incoming requests from users who browse the website, log into their accounts or post images.
- _Serve static assets_: Serve JavaScript, CSS & images to browsers to create a dynamic experience for the user.
- _Accept connections_: The HTTP Server must listen on a specific port to be able to accept connections from the Internet.

```go
package main

import (
    "fmt"
    "net/http"
)

func main() {
    // Process dynamic request
    http.HandleFunc("/", func (w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Welcome to my website!")
    })

    // Serving static assets
    fs := http.FileServer(http.Dir("static/"))
    http.Handle("/static/", http.StripPrefix("/static/", fs))

    // Accept connections
    http.ListenAndServe(":80", nil)
}
```

## 2. Templating

Go's `html/template` package provides data-driven templates for generating safe HTML output.

```go
package main

import (
    "html/template"
    "net/http"
)

type PageData struct {
    Title string
    Items []string
}

func main() {
    tmpl := template.Must(template.New("index").Parse(`
<!DOCTYPE html>
<html>
<head><title>{{.Title}}</title></head>
<body>
    <h1>{{.Title}}</h1>
    <ul>
    {{range .Items}}
        <li>{{.}}</li>
    {{end}}
    </ul>
</body>
</html>
`))

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        data := PageData{
            Title: "My Go Website",
            Items: []string{"Learn Go", "Build APIs", "Deploy"},
        }
        tmpl.Execute(w, data)
    })

    http.ListenAndServe(":8080", nil)
}
```

- Key template actions: `{{.}}` (current value), `{{.Field}}` (struct field), `{{range}}` (iteration), `{{if}}` (conditionals)
- Use `template.Must()` to panic on template parse errors at startup
- For external template files: `template.ParseFiles("templates/*.html")`

## 3. Requests and Forms

Handling HTTP requests and form data:

```go
package main

import (
    "fmt"
    "net/http"
)

func main() {
    http.HandleFunc("/contact", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case "GET":
            // Display the form
            fmt.Fprint(w, `
                <form method="POST" action="/contact">
                    <input type="text" name="name" placeholder="Name">
                    <input type="email" name="email" placeholder="Email">
                    <textarea name="message" placeholder="Message"></textarea>
                    <button type="submit">Send</button>
                </form>
            `)
        case "POST":
            // Parse form data
            if err := r.ParseForm(); err != nil {
                http.Error(w, "Error parsing form", http.StatusBadRequest)
                return
            }

            name := r.FormValue("name")
            email := r.FormValue("email")
            message := r.FormValue("message")

            fmt.Fprintf(w, "Received: %s (%s): %s", name, email, message)
        }
    })

    http.ListenAndServe(":8080", nil)
}
```

- `r.ParseForm()` - parse URL query params and POST form data
- `r.FormValue("key")` - get form value by key
- `r.URL.Query().Get("key")` - get URL query parameters
- For file uploads, use `r.ParseMultipartForm()` and `r.FormFile()`

## 4. Assets and Files

Serving static files and assets:

```go
package main

import "net/http"

func main() {
    // Serve files from ./static directory at /static/ URL path
    fs := http.FileServer(http.Dir("./static"))
    http.Handle("/static/", http.StripPrefix("/static/", fs))

    // Serve a single file
    http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "./static/favicon.ico")
    })

    // Main handler
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "./static/index.html")
    })

    http.ListenAndServe(":8080", nil)
}
```

- `http.FileServer()` - create a handler for serving files from a directory
- `http.StripPrefix()` - remove URL prefix before passing to file server
- `http.ServeFile()` - serve a specific file
- Security: be careful with path traversal attacks; `http.FileServer` handles this automatically

## 5. Middleware (Basic)

A simple logging middleware.

```go
package main

import (
    "fmt"
    "log"
    "net/http"
)

func logging(f http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        log.Println(r.URL.Path)
        f(w, r)
    }
}

func foo(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "foo")
}

func bar(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "bar")
}

func main() {
    http.HandleFunc("/foo", logging(foo))
    http.HandleFunc("/bar", logging(bar))

    http.ListenAndServe(":8080", nil)
}
```

## 6. Middleware (Advanced)

A middleware in itself simple takes a `http.HandleFunc` as one of its parameters, wraps it & returns a new `http.HandlerFunc` for the server to call.

```go
package main

import (
    "fmt"
    "log"
    "net/http"
    "time"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

// Logging logs all requests with its path & the time it took to process
func Logging() Middleware {
    return func(f http.HandlerFunc) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()
            defer func() { log.Println(r.URL.Path, time.Since(start)) }()
            f(w, r)
        }
    }
}

// Method ensures that url can only be requested with a specific method
func Method(m string) Middleware {
    return func(f http.HandlerFunc) http.HandlerFunc {
        return func(w http.ResponseWriter, r *http.Request) {
            if r.Method != m {
                http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
                return
            }
            f(w, r)
        }
    }
}

// Chain applies middlewares to a http.HandlerFunc
func Chain(f http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
    for _, m := range middlewares {
        f = m(f)
    }
    return f
}

func Hello(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "hello world")
}

func main() {
    http.HandleFunc("/", Chain(Hello, Method("GET"), Logging()))
    http.ListenAndServe(":8080", nil)
}
```

## 7. Session

Go's standard library doesn't include session management, but you can use packages like `gorilla/sessions`:

```go
package main

import (
    "fmt"
    "net/http"

    "github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("super-secret-key"))

func main() {
    http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
        session, _ := store.Get(r, "session-name")
        session.Values["authenticated"] = true
        session.Values["user"] = "john"
        session.Save(r, w)
        fmt.Fprintln(w, "Logged in!")
    })

    http.HandleFunc("/dashboard", func(w http.ResponseWriter, r *http.Request) {
        session, _ := store.Get(r, "session-name")
        auth, ok := session.Values["authenticated"].(bool)
        if !ok || !auth {
            http.Error(w, "Forbidden", http.StatusForbidden)
            return
        }
        user := session.Values["user"].(string)
        fmt.Fprintf(w, "Welcome, %s!", user)
    })

    http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
        session, _ := store.Get(r, "session-name")
        session.Values["authenticated"] = false
        session.Options.MaxAge = -1
        session.Save(r, w)
        fmt.Fprintln(w, "Logged out!")
    })

    http.ListenAndServe(":8080", nil)
}
```

- Install: `go get github.com/gorilla/sessions`
- Use secure, random keys in production
- Consider server-side session stores (Redis, PostgreSQL) for scalability

## 8. Websockets

Websockets enable real-time bidirectional communication. The `gorilla/websocket` package is commonly used:

```go
package main

import (
    "fmt"
    "log"
    "net/http"

    "github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
    CheckOrigin:     func(r *http.Request) bool { return true },
}

func main() {
    http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
        conn, err := upgrader.Upgrade(w, r, nil)
        if err != nil {
            log.Println("Upgrade error:", err)
            return
        }
        defer conn.Close()

        for {
            messageType, message, err := conn.ReadMessage()
            if err != nil {
                log.Println("Read error:", err)
                break
            }

            log.Printf("Received: %s", message)

            reply := fmt.Sprintf("Echo: %s", message)
            if err := conn.WriteMessage(messageType, []byte(reply)); err != nil {
                log.Println("Write error:", err)
                break
            }
        }
    })

    log.Println("Server started on :8080")
    http.ListenAndServe(":8080", nil)
}
```

- Install: `go get github.com/gorilla/websocket`
- Use secure `CheckOrigin` validation in production
- Handle connection lifecycle (ping/pong, timeouts)
