package websocket

import (
	"log"
	"net/http"
)

var hub *Hub

func init() {
	hub = newHub()
}

type Delegate interface {
	RecvMessage(messageType int, messageByte []byte)
}

func AddRouter(pattern string, delegate Delegate) {
	http.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r, delegate)
	})
}

func SendMessage(clientid, data string) {
	for client := range hub.clients {
		if client.id == clientid {
			client.send <- []byte(data)
			return
		} else {
			log.Println(client.id)
		}
	}
}

func Run(addr string, handler http.Handler) {
	err := http.ListenAndServe(addr, handler)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
