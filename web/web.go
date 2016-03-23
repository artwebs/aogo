package web

import (
	aolog "github.com/artwebs/aogo/log"
	"github.com/artwebs/aogo/utils"
	"log"
	"net/http"
	"path"
	"reflect"
	"strings"
	"time"
)

var (
	register     *ControllerRegistor
	exceptMethod = []string{"Init", "Display"}
)

type Handler struct {
	controller ControllerInterface
	method     string
}

func init() {
	register = NewControllerRegistor()
	aolog.NewLogger(1000)
	aolog.SetLogger("console", "")
	aolog.SetLevel(aolog.LevelDebug)
}

func Run() {
	log.Println(register.routes)
	conn := &http.Server{Addr: ":8080", Handler: register, ReadTimeout: 5 * time.Second}
	http.Handle("/css/", http.FileServer(http.Dir("")))
	http.Handle("/js/", http.FileServer(http.Dir("")))
	http.Handle("/images/", http.FileServer(http.Dir("")))
	err := conn.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}

func Router(pattern string, c ControllerInterface, method string) {
	register.routes[pattern] = &Handler{controller: c, method: method}
}

func AutoRouter(prefix string, c ControllerInterface) {
	reflectVal := reflect.ValueOf(c)
	rt := reflectVal.Type()
	ct := reflect.Indirect(reflectVal).Type()
	controllerName := strings.TrimSuffix(ct.Name(), "Controller")
	for i := 0; i < rt.NumMethod(); i++ {
		if !util.InSlice(rt.Method(i).Name, exceptMethod) {
			pattern := path.Join(prefix, strings.ToLower(controllerName), strings.ToLower(rt.Method(i).Name))
			Router(pattern, c, rt.Method(i).Name)
		}

	}
}
