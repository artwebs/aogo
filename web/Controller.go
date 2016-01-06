package web

import (
	"html/template"
	"log"
	"net/http"
)

type Controller struct {
	Writer  http.ResponseWriter
	Request *http.Request
	Data    map[string]interface{}
}

func (this *Controller) Init(w http.ResponseWriter, r *http.Request) {
	this.Writer = w
	this.Request = r
	this.Data = make(map[string]interface{})
}

func (this *Controller) Display(tpl string) {
	t, err := template.ParseFiles(tpl)
	if err != nil {
		log.Println(err)
	}
	t.Execute(this.Writer, this.Data)
}

type ControllerInterface interface {
	Init(w http.ResponseWriter, r *http.Request)
}
