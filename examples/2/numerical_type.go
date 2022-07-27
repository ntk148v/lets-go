package main

func main() {
    var a int
    var b int32
    b = a + a // Give an error: cannot use a + a (type int) as type int32 in assignment.
    b = b + 5
}
