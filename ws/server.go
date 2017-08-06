package ws

import (
	"log"
	"net/http"
)

var hub *Hub

func init() {
	hub = NewHub()
}

type WebSocketDelegate interface {
	RecvMessage(c *Client, mType int, mByte []byte)
}

func AddRouter(pattern string, delegate WebSocketDelegate) {
	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		ServeWs(hub, w, r, delegate)
	})
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

func Run(addr string, handler http.Handler) {
	go hub.Run()
	err := http.ListenAndServe(addr, handler)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
