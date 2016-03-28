package web

import (
	"fmt"
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
		// reg := regexp.MustCompile(`^` + key + `\/*`)
		// rsReg := reg.FindStringIndex(url)
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
			fmt.Println(urlarr[len(keyarr):])
			handler.controller.Init(w, r, urlarr[len(keyarr):])
			reflectVal := reflect.ValueOf(handler.controller)
			if val := reflectVal.MethodByName(handler.method); val.IsValid() {
				val.Call([]reflect.Value{})
			} else {
				panic("'' method doesn't exist in the controller " + handler.method)
			}
			return
		}
	next:
		{
			continue
		}

	}

	// if handler, ok := this.routes[r.URL.String()]; ok {
	// 	handler.controller.Init(w, r)
	// 	reflectVal := reflect.ValueOf(handler.controller)
	// 	if val := reflectVal.MethodByName(handler.method); val.IsValid() {
	// 		val.Call([]reflect.Value{})
	// 	} else {
	// 		panic("'' method doesn't exist in the controller " + handler.method)
	// 	}
	// }
}
