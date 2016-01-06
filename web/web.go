package web

import (
	"log"
	"net/http"

	"time"
)

var (
	register *ControllerRegister
)

type Handler struct {
	controller ControllerInterface
	method     string
}

func init() {
	register = NewControllerRegister()
}

func Run() {
	conn := &http.Server{Addr: ":8080", Handler: register, ReadTimeout: 5 * time.Second}
	http.Handle("/css/", http.FileServer(http.Dir("")))
	http.Handle("/js/", http.FileServer(http.Dir("")))
	http.Handle("/images/", http.FileServer(http.Dir("")))
	err := conn.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func Router(pattern string, c ControllerInterface, method string) {
	register.routes[pattern] = &Handler{controller: c, method: method}
}
