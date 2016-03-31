package web

import (
	// "fmt"
	"net/http"
	"reflect"
	// "regexp"
	"strings"
)

type ControllerRegistor struct {
	routes map[string]*Handler
}

func NewControllerRegistor() *ControllerRegistor {
	return &ControllerRegistor{routes: make(map[string]*Handler)}
}

func (this *ControllerRegistor) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	url := r.URL.String()
	for key, handler := range this.routes {
		keyarr := strings.Split(key, "/")
		urlarr := strings.Split(strings.Split(url, "?")[0], "/")
		if len(keyarr) > len(urlarr) {
			continue
		}

		for index := range keyarr {
			if keyarr[index] != urlarr[index] {
				goto next
			}
		}
		goto success
	success:
		{
			reflectVal := reflect.ValueOf(handler.controller)
			data := urlarr[len(keyarr):]
			handler.controller.Init(w, r, handler.controller, handler.method, data)
			handler.controller.SetUrl(keyarr)
			if val := reflectVal.MethodByName(handler.method); val.IsValid() {
				val.Call([]reflect.Value{})
			} else {
				panic("'' method doesn't exist in the controller " + handler.method)
			}
			handler.controller.Release()
			return
		}
	next:
		{
			continue
		}

	}

}
