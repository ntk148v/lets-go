# Go Enums and iota

Source: https://blog.learngoprogramming.com/golang-const-type-enums-iota-bc4befd096d3

[Enums](https://en.wikipedia.org/wiki/Enumerated_type) are useful parts of many languages. Unfortunately, enums in Go aren't as useful due to Go's implementattion. The biggest drawback is that they aren’t strictly typed, thus you have to manually validate them.

For example, suppose that you want to create `an enum of the weekdays`.

## The enums without iota

```go
package main

import (
    "fmt"
)

// Declare a new type named Weekday which will unify our enum values
// It has an underlying type of unsigned integer (uint).
type Weekday int

// Declare typed constants each with type of Weekday
const (
   Sunday    Weekday = 0
   Monday    Weekday = 1
   Tuesday   Weekday = 2
   Wednesday Weekday = 3
   Thursday  Weekday = 4
   Friday    Weekday = 5
   Saturday  Weekday = 6
)

// String returns the name of the day
func (day Weekday) String() string {
    names := [...]string{"Sunday", "Monday", "Tuesday", "Wednesday",
        "Thursday", "Friday", "Saturday"}

    if day < Sunday || day > Saturday {
      return "Unknown"
    }

    return names[day]
}

// Weekend return true for a weekend day
func (day Weekday) Weekend() bool {
    switch day {
    case Sunday, Saturday:
        return true
    default:
        return false
    }
}

func main() {
    // Which day it is? Sunday
    fmt.Printf("Which day it is? %s\n", Sunday)

    // Is Saturday a weekend day? true
    fmt.Printf("Is Saturday a weekend day? %t\n", Saturday.Weekend())

    // Is Monday a weekend day? false
    fmt.Printf("Is Monday a weekend day? %t\n", Monday.Weekend())
}
```

## The enums with iota

```go
package main

import (
    "fmt"
)

// Declare a new type named Weekday which will unify our enum values
// It has an underlying type of unsigned integer (uint).
type Weekday int

// Declare typed constants each with type of Weekday
const (
   Sunday    Weekday = iota
   Monday
   Tuesday
   Wednesday
   Thursday
   Friday
   Saturday
)

// String returns the name of the day
func (day Weekday) String() string {
    names := [...]string{"Sunday", "Monday", "Tuesday", "Wednesday",
        "Thursday", "Friday", "Saturday"}

    if day < Sunday || day > Saturday {
      return "Unknown"
    }

    return names[day]
}

// Weekend return true for a weekend day
func (day Weekday) Weekend() bool {
    switch day {
    case Sunday, Saturday:
        return true
    default:
        return false
    }
}

func main() {
    // Which day it is? Sunday
    fmt.Printf("Which day it is? %s\n", Sunday)

    // Is Saturday a weekend day? true
    fmt.Printf("Is Saturday a weekend day? %t\n", Saturday.Weekend())

    // Is Monday a weekend day? false
    fmt.Printf("Is Monday a weekend day? %t\n", Monday.Weekend())
}
```

### Iota basic

`iota` in short:

- A numeric universal counter starting at 0
- Used only with constant declarations.

iota increases by 1 after each line except empty and comment lines.

Don't use iota for a list of predefined values like [FTP server status codes](https://en.wikipedia.org/wiki/List_of_FTP_server_return_codes).

### Iota expressions & rules

- Reset iota

```go
// iota reset: it will be 0.
const (
    Zero = iota  // Zero = 0
    One          // One = 1
)
// iota reset: will be 0 again
const (
    Two = iota   // Two = 0
)
// iota: reset
const Three = iota // Three = 0
```

- Skip some values

```go
type Timezone int

const (
    // iota: 0, EST: -5
    EST Timezone = -(5 + iota)
    // _ is the blank identifier
    // iota: 1
    _
    // iota: 2, MST: -7
MST
)
```

- One a comment or an empty line, iota will not increase

- Use iota in the middle

```go
const (
    One = 1
    Two = 2
    // Three = 2 + 1 => 3
    Three = iota + 1
)
```

- Multiple iotas in a single line

```go
const (
    // Active = 0, Moving = 0,
    // Running = 0
    Active, Moving, Running = iota, iota, iota
    // Passive = 1, Stopped = 1,
    // Stale = 1
    Passive, Stopped, Stale
)
```

- Repeat and cancel expressions

```go
const (
    // iota: 0, One: 1 (type: int64)
    One  int64  = iota + 1
    // iota: 1, Two: 2 (type: int64)
    // Two will be declared as if:
    // Two int64 = iota + 1
    Two
    // iota: 2, Four: 4 (type: int32)
    Four int32  = iota + 2
    // iota: 3, Five: 5 (type: int32)
    // Five will be declared as if:
    // Five int32 = iota + 2
    Five
    // (type: int)
    Six = 6
    // (type: int)
    // Seven will be declared as if:
    // Seven = 6
    Seven
)
```

- Produce alphabets

```go
const (
    // string will convert the
    // expression into string.
    //
    // or, it'll assign character
    // codes.
    a = string(iota + 'a') // a
    b                      // b
    c                      // c
    d                      // d
    e                      // e
)
```

- Beware the zero-value

```go
type Activity int
const (
    Sleeping = iota
    Walking
    Running
)
func main() {
    var activity Activity
    // activity initialized to
    // its zero-value of int
    // which is Sleeping
// iota + 1 trick
const (
    Sleeping = iota + 1
    Walking
    Running
)
func main() {
    var activity Activity
    // activity will be zero,
    // so it's not initialized
    activity = Sleeping
    // now you know that it's been
    // initialized
}
```
