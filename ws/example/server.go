package main

import (
	"log"

	"github.com/artwebs/aogo/websocket"
)

type DemoMessage struct {
}

func (this DemoMessage) RecvMessage(c *ws.Client, mType int, mByte []byte) {
	log.Printf("%q", mByte)
	c.Send <- []byte("artwebs")
}

func main() {
	ws.AddRouter("/test", &DemoMessage{})
	ws.Run(":8080", nil)
}
