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
	routes     map[string]interface{}
	namespaces []string
}

func NewControllerRegistor() *ControllerRegistor {
	return &ControllerRegistor{routes: make(map[string]interface{}), namespaces: []string{""}}
}

func (this *ControllerRegistor) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	url := r.URL.String()
	log.InfoTag(this, url)
	urlarr := strings.Split(strings.Split(url, "?")[0], "/")
	for key, handler := range this.routes {
		keyarr := strings.Split(key, "/")
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
			this.doController(keyarr, urlarr, handler, w, r)
			return
		}
	next:
		{
			continue
		}

	}

	if handler, ok := this.routes["/"]; ok {
		keyarr := []string{}
		this.doController(keyarr, urlarr, handler, w, r)
		log.InfoTag(this, "to /")
		return
	} else {
		if url == "/" {
			url = "/index"
			http.Redirect(w, r, "/index", http.StatusFound)
			return
		}
	}

	log.ErrorTag(this, url+" do not find")

}

func (this *ControllerRegistor) doController(keyarr, urlarr []string, h interface{}, w http.ResponseWriter, r *http.Request) {
	switch handler := h.(type) {
	case *Handler:
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
		break
	case http.Handler:
		handler.ServeHTTP(w, r)
		break
	default:
		log.ErrorTag(this, h, " do not find")
	}

}
