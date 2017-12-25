package web

import (
	"log"
	"net/http"
	"path"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/artwebs/aogo/logger"
	"github.com/artwebs/aogo/utils"
)

var (
	register           *ControllerRegistor
	exceptMethod       = []string{"Init", "WillDid", "Release"}
	LoggerLevelDefault int
)

type Handler struct {
	controller ControllerInterface
	method     string
}

func init() {
	InitAppConfig()
	InitSession()
	register = NewControllerRegistor()
	logger.SetLevel(logger.LoggerDebug)
}

func Run() {
	for _, item := range register.namespaces {
		HandleFile(item+"/css", ViewsPath)
		HandleFile(item+"/js", ViewsPath)
		HandleFile(item+"/images", ViewsPath)
	}
	register.tree.PrintTree("")
	conn := &http.Server{Addr: HttpAddress + ":" + strconv.Itoa(HttpPort), Handler: register, ReadTimeout: 5 * time.Second, WriteTimeout: 5 * time.Second}
	logger.Info("server " + HttpAddress + ":" + strconv.Itoa(HttpPort) + " started")
	if Debug == 1 {
		LoggerLevelDefault = logger.LoggerDebug
	} else {
		LoggerLevelDefault = logger.LoggerError
	}
	logger.SetLevel(LoggerLevelDefault)
	err := conn.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}

func Router(pattern string, c ControllerInterface, method string) {
	register.tree.AddRouter(pattern, &Handler{controller: c, method: method})
}

func Handle(pattern string, handler http.Handler) {
	register.tree.AddRouter(pattern, handler)
}

func HandleFile(pattern, root string) {
	register.tree.AddRouter(pattern, http.FileServer(http.Dir(root)))
}

func HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	register.tree.AddRouter(pattern, HandlerFunc(handler))
}

func AutoRouter(prefix string, c ControllerInterface) {
	reflectVal := reflect.ValueOf(c)
	rt := reflectVal.Type()
	ct := reflect.Indirect(reflectVal).Type()
	controllerName := strings.TrimSuffix(ct.Name(), "Controller")
	for i := 0; i < rt.NumMethod(); i++ {
		if !utils.InSlice(rt.Method(i).Name, exceptMethod) {
			pattern := path.Join(prefix, strings.ToLower(controllerName), strings.ToLower(rt.Method(i).Name))
			Router(pattern, c, rt.Method(i).Name)
		}

	}
}

// The HandlerFunc type is an adapter to allow the use of
// ordinary functions as HTTP handlers. If f is a function
// with the appropriate signature, HandlerFunc(f) is a
// Handler that calls f.
type HandlerFunc func(http.ResponseWriter, *http.Request)

// ServeHTTP calls f(w, r).
func (f HandlerFunc) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	f(w, r)
}
