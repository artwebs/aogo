package web

import (
	// "fmt"
	"net/http"
	"reflect"
	"runtime/debug"
	"time"
	// "regexp"
	"strings"

	"github.com/artwebs/aogo/logger"
	"github.com/artwebs/aogo/utils"
)

type ControllerRegistor struct {
	tree       *RouterTree
	namespaces []string
}

func NewControllerRegistor() *ControllerRegistor {
	return &ControllerRegistor{tree: &RouterTree{}, namespaces: []string{""}}
}

func (this *ControllerRegistor) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(err)
			debug.PrintStack()
		}
	}()
	stime := time.Now()

	url := r.URL.String()
	// ip, port, _ := utils.HttpClientIP(r)

	if url == "/favicon.ico" {
		http.ServeFile(w, r, "./favicon.ico")
		return
	}
	urlarr := strings.Split(strings.Trim(strings.Split(url, "?")[0], "/"), "/")
	data, handler := this.tree.FindRouter(strings.Split(url, "?")[0])
	if handler != nil {
		this.doController(data, urlarr, handler, w, r)
		etime := time.Now()
		_, port, _ := utils.HttpClientIP(r)
		logger.InfoTag(this, "[", etime.Sub(stime), "]", r.Header.Get("X-Real-IP"), "[", port, "]", url)
		return
	}
	etime := time.Now()
	_, port, _ := utils.HttpClientIP(r)
	logger.InfoTag(this, "[", etime.Sub(stime), "]", r.Header.Get("X-Real-IP"), "[", port, "]", url)
	logger.ErrorTag(this, url+" do not find")

}

func (this *ControllerRegistor) doController(data, urlarr []string, h interface{}, w http.ResponseWriter, r *http.Request) {
	defer func() {
		err := r.Body.Close()
		if err != nil {
			logger.ErrorTag(this, "关闭http response 失败", err)
		}
	}()
	switch handler := h.(type) {
	case *Handler:
		reflectVal := reflect.ValueOf(handler.controller)
		ctx := &Context{}
		ctx.Init(w, r, handler.controller, handler.method, data)
		handler.controller.Init(ctx)
		if handler.controller.WillDid(ctx) {
			ctx.SetUrl(urlarr[:len(urlarr)-len(data)])
			if val := reflectVal.MethodByName(handler.method); val.IsValid() {
				val.Call([]reflect.Value{reflect.ValueOf(ctx)})
			} else {
				panic("'' method doesn't exist in the controller " + handler.method)
			}
		}
		handler.controller.Release(ctx)
		break
	case http.Handler:
		handler.ServeHTTP(w, r)
		break
	default:
		logger.ErrorTag(this, h, " do not find")
	}

}
