# Stupid Gopher Tricks

[Andrew Gerrand](https://twitter.com/enneff)'s talk "Stupid Gopher Tricks" at the Golang UK conference.

## Type literals

Type declaration:

```golang
type Foo struct {
    i int
    s string
}
```

The latter part of a type declaration is the **type literal**:

```golang
struct {
    i int
    s string
}
```

An unnamed struct literal is often called an **anonymous struct**:

```golang
var t struct {
    i int
    s string
}
```

## Anonymous structs

### Template data

```golang
data := struct {
    Title               string
    Firstname, Lastname string
    Rank                int
}{
    "Dr", "Carl", "Sagan", 7,
}
if err := tmpl.Execute(os.Stdout, data); err != nil {
    log.Fatal(err)
}
```

### JSON

* Encode and decode JSON objects.

```golang
b, err := json.Marshal(struct {
    ID   int
    Name string
}{42, "The answer"})
if err != nil {
    log.Fatal(err)
}
fmt.Printf("%s\n", b)

var data struct {
    ID   int
    Name string
}
err := json.Unmarshal([]byte(`{"ID": 42, "Name": "The answer"}`), &data)
if err != nil {
    log.Fatal(err)
}
fmt.Println(data.ID, data.Name)
```

* Structs can be nested to describe more complex JSON objects

```golang
var data struct {
    ID int
    Person struct {
        Name string
        Job string
    }
}

const s = `{"ID":42,"Person":{"Name":"George Costanza","Job":"Architect"}}`
err := json.Unmarshal([]byte(s), &data)
if err != nil {
    log.Fatal(err)
}
fmt.Println(data.ID, data.Person.Name, data.Person.Job)
```

## Repeated literals and struct names

```golang
type Foo struct {
    i int
    s string
}

var s = []Foo{
    {6 * 9, "Question"},
    {42, "Answer"},
}

var m = map[int]Foo{
    7: {6 * 9, "Question"},
    3: {42, "Answer"},
}
```

* Combined with anonymous structs, this convenience shortens the code dramatically:

```golang
var s = []struct {
    i int
    s string
}{
    struct {
        i int
        s string
    }{6 * 9, "Question"},
    struct {
        i int
        s string
    }{42, "Answer"},
}

var t = []struct {
    i int
    s string
}{
    {6 * 9, "Question"},
    {42, "Answer"},
}
```
...
[WIP]
