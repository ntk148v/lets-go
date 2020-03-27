package main

import (
	"fmt"
	"strings"
)

func Join(parts ...string) string {
	return strings.Join(parts, PathSeparator)
}

func main() {
	s := Join("a", "b", "c")
	fmt.Println(s)
}
