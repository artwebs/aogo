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

func (this *Controller) Init(w http.ResponseWriter, r *http.Request, data []string) {
	this.Writer = w
	this.Request = r
	this.Data = make(map[string]interface{})
	if len(data)%2 == 0 {
		index := 0
		for {
			if index >= len(data) {
				break
			}
			this.Data[data[index]] = data[index+1]
			index += 2
		}
	}

	for k, v := range r.Form {
		if len(v) > 0 {
			this.Data[k] = v[0]
		} else {
			this.Data[k] = v
		}
	}

	log.Println(this.Data)
}

func (this *Controller) Display(tpl string) {
	t, err := template.ParseFiles(tpl)
	if err != nil {
		log.Println(err)
	}
	t.Execute(this.Writer, this.Data)
}

type ControllerInterface interface {
	Init(w http.ResponseWriter, r *http.Request, data []string)
}
