package main

import (
	"fmt"
	"strings"
)

const PathSeparator = "/"

func Join(parts ...string) string {
	return strings.Join(parts, PathSeparator)
}

func main() {
	s := Join("a", "b", "c")
	fmt.Println(s)
}
