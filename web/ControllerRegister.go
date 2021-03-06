package web

import (
	// "fmt"
	"net/http"
	"reflect"
	// "regexp"
	"strings"

	"github.com/artwebs/aogo/log"
)

type ControllerRegistor struct {
	tree       *RouterTree
	namespaces []string
}

func NewControllerRegistor() *ControllerRegistor {
	return &ControllerRegistor{tree: &RouterTree{}, namespaces: []string{""}}
}

func (this *ControllerRegistor) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	url := r.URL.String()
	log.InfoTag(this, url)
	if url == "/favicon.ico" {
		http.ServeFile(w, r, "./favicon.ico")
		return
	}

	urlarr := strings.Split(strings.Trim(strings.Split(url, "?")[0], "/"), "/")
	data, handler := this.tree.FindRouter(strings.Split(url, "?")[0])
	if handler != nil {
		this.doController(data, urlarr, handler, w, r)
		return
	}

	log.ErrorTag(this, url+" do not find")

}

func (this *ControllerRegistor) doController(data, urlarr []string, h interface{}, w http.ResponseWriter, r *http.Request) {
	switch handler := h.(type) {
	case *Handler:
		reflectVal := reflect.ValueOf(handler.controller)
		handler.controller.Init(w, r, handler.controller, handler.method, data)
		if handler.controller.WillDid() {
			handler.controller.SetUrl(urlarr[:len(urlarr)-len(data)])
			if val := reflectVal.MethodByName(handler.method); val.IsValid() {
				val.Call([]reflect.Value{})
			} else {
				panic("'' method doesn't exist in the controller " + handler.method)
			}
		}
		handler.controller.Release()
		break
	case http.Handler:
		handler.ServeHTTP(w, r)
		break
	default:
		log.ErrorTag(this, h, " do not find")
	}

}
