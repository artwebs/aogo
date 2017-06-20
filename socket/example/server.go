package main

import (
	"log"

	"github.com/artwebs/aogo/socket"
)

type DemoMessage struct {
}

func (this DemoMessage) RecvMessage(c *socket.Client, mByte []byte) {
	log.Printf("%q", mByte)
	c.Send <- []byte("artwebs")
}

func main() {
	socket.Run(":8080", &DemoMessage{})
}
