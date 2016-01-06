package web

import (
	"net/http"
	"reflect"
)

type ControllerRegister struct {
	routes map[string]*Handler
}

func NewControllerRegister() *ControllerRegister {
	return &ControllerRegister{routes: make(map[string]*Handler)}
}

func (this *ControllerRegister) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if handler, ok := this.routes[r.URL.String()]; ok {
		handler.controller.Init(w, r)
		reflectVal := reflect.ValueOf(handler.controller)
		if val := reflectVal.MethodByName(handler.method); val.IsValid() {
			val.Call([]reflect.Value{})
		} else {
			panic("'' method doesn't exist in the controller " + handler.method)
		}
	}
}
