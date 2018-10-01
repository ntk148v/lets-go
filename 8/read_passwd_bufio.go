package main

import (
	"bufio"
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
	r := bufio.NewReader(f)
	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()
	for {
		n, e := r.Read(buf)
		if e != nil {
			log.Fatalf(e.Error())
		}
		if n == 0 {
			break
		}
		w.Write(buf[0:n])
	}
}
