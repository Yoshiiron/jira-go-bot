package main

import (
	"github.com/nats-io/nats.go"
	"log"
)

func main() {
	nc, err := nats.Connect("nats://192.168.1.121:4222")
	if err != nil {
		log.Fatal(err)
	}
	nc.Publish("foo", []byte("Hello"))
}
