package main

import (
	"log"
	"net/rpc"
)

func main() {
	client, err := rpc.Dial("tcp", "localhost:8081")
	if err != nil {
		log.Fatal("Dialing error:", err)
	}

	var reply string

	if err = client.Call("HelloService.Hello", "Kien", &reply); err != nil {
		log.Fatal(err)
	}

	log.Println(reply)
}
