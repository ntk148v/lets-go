# Encoding & Decoding

> **NOTE**: Code examples for this section are stored in [`examples/12/`](../examples/12/).

Table of Contents:

- [Encoding \& Decoding](#encoding--decoding)
  - [1. JSON](#1-json)
    - [1.1. Marshal and Unmarshal Structured data](#11-marshal-and-unmarshal-structured-data)
    - [1.2. JSON struct tags - custom field names](#12-json-struct-tags---custom-field-names)
    - [1.3. Decoding JSON to Maps - Unstructured Data](#13-decoding-json-to-maps---unstructured-data)
    - [1.4. Ignore empty fields with omitempty](#14-ignore-empty-fields-with-omitempty)

Go's standard library comes packed with some great encoding and decoding packages covering a wide array of encoding schemes. Everything from CSV, XML, JSON, and even gob - a Go specific encoding format - is covered, and all of these packages are incredibly easy to get started with.

## 1. JSON

Go offers built-in support for JSON encoding and decoding, including to and from built-in and custom data types.

There are two types of data:

- Structured data
- Unstructured data

### 1.1. Marshal and Unmarshal Structured data

`json` package will assign values only to fields found in the JSON; other fields will just keep their [Go zero values](https://golang.org/ref/spec#The_zero_value).

```go
package main

import (
    "encoding/json"
    "fmt"
)

type Measurement struct {
    Height int
    Weight int
}

type Person struct {
    Name string
    Age  int
    Measurement Measurement // Nested object
}

func main() {
    bob := &Person{
        Name: "Bob",
        Age:  20,
    }
    bobRaw, _ := json.Marshal(bob)
    fmt.Println(string(bobRaw))

    // Raw data without Measurement field
    aliceRaw := []byte(`{"name": "Alice", "age": 23}`)
    var alice Person

    if err := json.Unmarshal(aliceRaw, &alice); err != nil {
        panic(err)
    }
    fmt.Printf("%+v\n", alice)
}
// {"Name":"Bob","Age":20,"Measurement":{"Height":190,"Weight":75}}
// {Name:Alice Age:23 Measurement:{Height:0 Weight:0}}
```

### 1.2. JSON struct tags - custom field names

Sometimes we want a different attribute name than the one provided in your JSON data.

```go
package main

import (
        "encoding/json"
        "fmt"
)

type Measurement struct {
        Height int `json:"height"`
        Weight int `json:"weight"`
}

type Person struct {
        Name        string      `json:"who"`
        Age         int         `json:"how old"`
        Measurement Measurement `json:"mm"`
}

func main() {
        bob := &Person{
                Name: "Bob",
                Age:  20,
        }
        bobRaw, _ := json.Marshal(bob)
        fmt.Println(string(bobRaw))

        // Raw data without Measurement field
        aliceRaw := []byte(`{"who": "Alice", "how old": 23, "mm": {"height": 150, "weight": 40}}`)
        var alice Person

        if err := json.Unmarshal(aliceRaw, &alice); err != nil {
                panic(err)
        }
        fmt.Printf("%+v", alice)
}
// {"who":"Bob","how old":20,"mm":{"height":0,"weight":0}}
// {Name:Alice Age:23 Measurement:{Height:150 Weight:40}}
```

### 1.3. Decoding JSON to Maps - Unstructured Data

```go
package main

import (
    "encoding/json"
    "fmt"
)

func main() {
    // Raw data without Measurement field
    aliceRaw := []byte(`{"name": "Alice", "age": 23, "measurement": {"height": 150, "weight": 40}}`)
    var alice map[string]interface{}

    if err := json.Unmarshal(aliceRaw, &alice); err != nil {
        panic(err)
    }
    // the object stored in the "mesurement" key is also stored
    // as a map[string]interface{} type, and its type is asserted
    // the interface{} type
    measurement := alice["measurement"].(map[string]interface{})
    fmt.Printf("%+v\n", alice)
    fmt.Printf("%+v\n", measurement)
}

// map[age:23 measurement:map[height:150 weight:40] name:Alice]
// map[height:150 weight:40]

```

### 1.4. Ignore empty fields with omitempty

In some cases, we would want to ignore a field in our JSON output, if its value is empty. We can use the `omitempty` property.

```go
package main

import (
    "encoding/json"
    "fmt"
)

type Measurement struct {
    Height int `json:"height"`
    Weight int `json:"weight"`
    }

type Person struct {
    Name        string      `json:"name"`
    Age         int         `json:"age,omitempty"`
    Measurement Measurement `json:"measurement"`
}

func main() {
    bob := &Person{
        Name: "Bob",
        Measurement: Measurement{
            Height: 190,
            Weight: 75,
        },
    }
    bobRaw, _ := json.Marshal(bob)
    fmt.Println(string(bobRaw))
}

// Age field is ignored
// {"name":"Bob","measurement":{"height":190,"weight":75}}
```

About the Advanced Encoding and Decoding techniques, you can check [this blog](https://blog.gopheracademy.com/advent-2016/advanced-encoding-decoding/).

> **NOTE**: There are a lot more helpful things in [tips-notes](../tips-notes/). You may want to check it out.
