package socket

import (
	"fmt"
	"log"
	"net"
	"time"
)

type Client struct {
	hub  *Hub
	conn net.Conn

	Send     chan []byte
	Id       string
	Login    string
	Index    int
	Delegate SocketDelegate
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 30 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

func NewClient(hub *Hub, conn net.Conn, delegate SocketDelegate) {
	client := &Client{hub: hub, conn: conn, Send: make(chan []byte, 1024), Delegate: delegate}
	client.hub.register <- client
	go client.writePump()

	client.readPump()
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	// c.conn.SetReadDeadline(time.Now().Add(pongWait))
	data := make([]byte, maxMessageSize)
	for {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		n, err := c.conn.Read(data)
		if err != nil {
			fmt.Printf("read message from lotus failed")
			return
		}
		if c.Delegate != nil {
			c.Delegate.RecvMessage(c, data[0:n])
		}

	}

}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.Close()
				return
			}
			c.conn.Write(message)

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if _, err := c.conn.Write([]byte{}); err != nil {
				log.Println("ticker error:", err)
				return
			}
		}
	}
}
