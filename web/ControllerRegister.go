package web

import (
	"net/http"
	"reflect"
)

type ControllerRegistor struct {
	routes map[string]*Handler
}

func NewControllerRegistor() *ControllerRegistor {
	return &ControllerRegistor{routes: make(map[string]*Handler)}
}

func (this *ControllerRegistor) ServeHTTP(w http.ResponseWriter, r *http.Request) {

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
