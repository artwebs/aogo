package socket

import (
	"fmt"
	"log"
	"net"
)

var hub *Hub

func init() {
	hub = newHub()
}

type SocketDelegate interface {
	RecvMessage(c *Client, mByte []byte)
}

func SendMessage(clientid, data string) {
	for client := range hub.clients {
		if client.Id == clientid {
			client.Send <- []byte(data)
		} else {
			log.Println(client.Id)
		}
	}
}

func Run(host string, delegate SocketDelegate) {
	go hub.run()
	fmt.Printf("Server is ready...\n")
	l, err := net.Listen("tcp", host)
	if err != nil {
		fmt.Printf("Failure to listen: %s\n", err.Error())
	}

	for {
		if c, err := l.Accept(); err == nil {
			go NewClient(hub, c, delegate)
		}
	}
}
