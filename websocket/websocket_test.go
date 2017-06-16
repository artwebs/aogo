package ws

import (
	"log"
	"testing"
)

type DemoMessage struct {
}

func (this DemoMessage) RecvMessage(client *Client, messageType int, messageByte []byte) {
	log.Println(string(messageType))
	client.Send <- []byte("artwebs")
}
func TestWebsocket(t *testing.T) {

	AddRouter("/test", &DemoMessage{})
	// server := httptest.NewServer(http.DefaultServeMux)
	// defer server.Close()
	// e := httpexpect.New(t, server.URL)
	//
	// e.GET("/test").
	// 	Expect().
	// 	Status(http.StatusOK).JSON().Array().Empty()
}
