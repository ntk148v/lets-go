# Testing in Go

> **NOTE**: Code examples in this section are stored in the respective test files. Run tests with `go test` from the directory.

Go has excellent built-in testing support. This guide covers everything from basic unit tests to advanced patterns.

Table of Contents:

- [Testing in Go](#testing-in-go)
  - [1. Basic Testing](#1-basic-testing)
    - [Test File Naming](#test-file-naming)
    - [Running Tests](#running-tests)
  - [2. Table-Driven Tests](#2-table-driven-tests)
  - [3. Subtests](#3-subtests)
  - [4. Testing Concurrent Code](#4-testing-concurrent-code)
    - [Using Goroutines in Tests](#using-goroutines-in-tests)
    - [Race Detection](#race-detection)
  - [5. HTTP Testing](#5-http-testing)
    - [Testing HTTP Clients](#testing-http-clients)
  - [6. Benchmarking](#6-benchmarking)
    - [Benchmark Example Output](#benchmark-example-output)
  - [7. Code Coverage](#7-code-coverage)
  - [8. Mocking and Dependency Injection](#8-mocking-and-dependency-injection)
    - [Interface-Based Mocking](#interface-based-mocking)
  - [9. testing/synctest Package *(Go 1.25+)*](#9-testingsynctest-package-go-125)
    - [The Problem](#the-problem)
    - [The Solution](#the-solution)
    - [Use Cases](#use-cases)
  - [Further Reading](#further-reading)

## 1. Basic Testing

### Test File Naming

- Test files must end with `_test.go`
- Test functions must start with `Test`
- Test functions take `*testing.T` as parameter

```go
// math.go
package math

func Sum(a, b int) int {
    return a + b
}
```

```go
// math_test.go
package math

import "testing"

func TestSum(t *testing.T) {
    result := Sum(2, 3)
    if result != 5 {
        t.Errorf("Sum(2, 3) = %d; want 5", result)
    }
}
```

### Running Tests

```bash
# Run all tests in current directory
go test

# Run with verbose output
go test -v

# Run specific test
go test -run TestSum

# Run tests in all subdirectories
go test ./...
```

## 2. Table-Driven Tests

Table-driven tests make it easy to add new test cases:

```go
func TestSum(t *testing.T) {
    tests := []struct {
        name     string
        a, b     int
        expected int
    }{
        {"positive numbers", 2, 3, 5},
        {"negative numbers", -1, -2, -3},
        {"zero", 0, 0, 0},
        {"mixed", -1, 5, 4},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := Sum(tt.a, tt.b)
            if result != tt.expected {
                t.Errorf("Sum(%d, %d) = %d; want %d",
                    tt.a, tt.b, result, tt.expected)
            }
        })
    }
}
```

## 3. Subtests

Subtests allow hierarchical test organization:

```go
func TestMath(t *testing.T) {
    t.Run("Sum", func(t *testing.T) {
        t.Run("positive", func(t *testing.T) {
            if Sum(1, 2) != 3 {
                t.Error("failed")
            }
        })
        t.Run("negative", func(t *testing.T) {
            if Sum(-1, -2) != -3 {
                t.Error("failed")
            }
        })
    })

    t.Run("Multiply", func(t *testing.T) {
        // ...
    })
}
```

Run specific subtest:

```bash
go test -run TestMath/Sum/positive
```

## 4. Testing Concurrent Code

### Using Goroutines in Tests

```go
func TestConcurrent(t *testing.T) {
    var wg sync.WaitGroup
    errCh := make(chan error, 10)

    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            if err := doWork(id); err != nil {
                errCh <- err
            }
        }(i)
    }

    wg.Wait()
    close(errCh)

    for err := range errCh {
        t.Error(err)
    }
}
```

### Race Detection

```bash
go test -race ./...
```

## 5. HTTP Testing

Use `httptest` package for testing HTTP handlers:

```go
import (
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestHandler(t *testing.T) {
    // Create request
    req := httptest.NewRequest("GET", "/api/users", nil)

    // Create response recorder
    w := httptest.NewRecorder()

    // Call handler
    handler(w, req)

    // Check response
    resp := w.Result()
    if resp.StatusCode != http.StatusOK {
        t.Errorf("status = %d; want %d", resp.StatusCode, http.StatusOK)
    }
}
```

### Testing HTTP Clients

```go
func TestAPIClient(t *testing.T) {
    // Create test server
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"id": 1, "name": "test"}`))
    }))
    defer server.Close()

    // Use test server URL
    client := NewAPIClient(server.URL)
    result, err := client.GetUser(1)

    if err != nil {
        t.Fatal(err)
    }
    if result.Name != "test" {
        t.Errorf("name = %s; want test", result.Name)
    }
}
```

## 6. Benchmarking

Benchmark functions start with `Benchmark`:

```go
func BenchmarkSum(b *testing.B) {
    for i := 0; i < b.N; i++ {
        Sum(1, 2)
    }
}
```

Running benchmarks:

```bash
# Run benchmarks
go test -bench=.

# With memory allocation stats
go test -bench=. -benchmem

# Run specific benchmark
go test -bench=BenchmarkSum

# Run N times for accuracy
go test -bench=. -count=5
```

### Benchmark Example Output

```
BenchmarkSum-8     1000000000     0.29 ns/op     0 B/op     0 allocs/op
```

## 7. Code Coverage

```bash
# Generate coverage profile
go test -cover

# Detailed coverage
go test -coverprofile=coverage.out
go tool cover -html=coverage.out

# Coverage by function
go tool cover -func=coverage.out
```

## 8. Mocking and Dependency Injection

### Interface-Based Mocking

```go
// Define interface
type UserRepository interface {
    GetUser(id int) (*User, error)
}

// Real implementation
type PostgresUserRepo struct{}

func (r *PostgresUserRepo) GetUser(id int) (*User, error) {
    // Database call
}

// Mock implementation for tests
type MockUserRepo struct {
    users map[int]*User
}

func (m *MockUserRepo) GetUser(id int) (*User, error) {
    if user, ok := m.users[id]; ok {
        return user, nil
    }
    return nil, errors.New("not found")
}

// Service using interface
type UserService struct {
    repo UserRepository
}

func (s *UserService) GetUserName(id int) (string, error) {
    user, err := s.repo.GetUser(id)
    if err != nil {
        return "", err
    }
    return user.Name, nil
}

// Test with mock
func TestGetUserName(t *testing.T) {
    mock := &MockUserRepo{
        users: map[int]*User{
            1: {ID: 1, Name: "Alice"},
        },
    }

    service := &UserService{repo: mock}

    name, err := service.GetUserName(1)
    if err != nil {
        t.Fatal(err)
    }
    if name != "Alice" {
        t.Errorf("name = %s; want Alice", name)
    }
}
```

## 9. testing/synctest Package *(Go 1.25+)*

The `testing/synctest` package provides robust testing for concurrent code by virtualizing time.

### The Problem

Testing time-dependent code is challenging:

- Tests that use `time.Sleep` are slow
- Timing-based tests can be flaky

### The Solution

`synctest.Run` virtualizes time within its scope:

```go
import (
    "testing"
    "testing/synctest"
    "time"
)

func TestTimeout(t *testing.T) {
    synctest.Run(func() {
        start := time.Now()

        // This would normally take 1 hour!
        time.Sleep(1 * time.Hour)

        elapsed := time.Since(start)
        if elapsed < 1*time.Hour {
            t.Error("expected at least 1 hour to pass")
        }
    })
    // Test completes immediately, not in 1 hour!
}
```

### Use Cases

- Testing timeouts without waiting
- Testing ticker/timer behavior
- Testing retry logic with exponential backoff
- Testing cache expiration

## Further Reading

- [Writing Unit Tests](../tips-notes/writing-unit-tests.md)
- [Dependency Injection](../tips-notes/dependency-injection.md)
- [Benchmarking](../tips-notes/benchmarking/)
