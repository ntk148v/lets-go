package main

import (
	"log"
	"os"
)

func main() {
	buf := make([]byte, 1024)
	f, e := os.Open("/etc/passwd")
	if e != nil {
		log.Fatalf(e.Error())
	}
	defer f.Close()
	for {
		n, e := f.Read(buf)
		if e != nil {
			log.Fatalf(e.Error())
		}
		if n == 0 {
			break
		}
		os.Stdout.Write(buf[:n])
	}
}
