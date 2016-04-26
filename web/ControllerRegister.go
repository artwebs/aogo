package web

import (
	// "fmt"
	"net/http"
	"reflect"
	// "regexp"
	"github.com/artwebs/aogo/log"
	"strings"
)

type ControllerRegistor struct {
	routes     map[string]*Handler
	namespaces []string
}

func NewControllerRegistor() *ControllerRegistor {
	return &ControllerRegistor{routes: make(map[string]*Handler), namespaces: []string{""}}
}

func (this *ControllerRegistor) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	url := r.URL.String()
	if url == "/" {
		url = "/index"
		http.Redirect(w, r, "/index", http.StatusFound)
		return
	}
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
			if handler.controller.WillDid() {
				handler.controller.SetUrl(keyarr)
				if val := reflectVal.MethodByName(handler.method); val.IsValid() {
					val.Call([]reflect.Value{})
				} else {
					panic("'' method doesn't exist in the controller " + handler.method)
				}
			}
			handler.controller.Release()
			return
		}
	next:
		{
			continue
		}

	}
	log.ErrorTag(this, url+" do not find")

}
