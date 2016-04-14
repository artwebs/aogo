package web

import (
	aolog "github.com/artwebs/aogo/log"
	"github.com/artwebs/aogo/utils"
	"log"
	"net/http"
	"path"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var (
	register     *ControllerRegistor
	exceptMethod = []string{"Init", "WillDid", "Redirect", "WriteString", "Display", "WriteJson", "Release", "SetSession", "GetSession", "FlushSession", "SaveToFile", "SetUrl"}
)

type Handler struct {
	controller ControllerInterface
	method     string
}

func init() {
	InitAppConfig()
	InitSession()
	register = NewControllerRegistor()
	aolog.NewLogger(1000)
	aolog.SetLogger("console", "")
	aolog.SetLevel(aolog.LevelDebug)
}

func Run() {
	aolog.Info(register.routes)
	conn := &http.Server{Addr: HttpAddress + ":" + strconv.Itoa(HttpPort), Handler: register, ReadTimeout: 5 * time.Second}
	aolog.Info("server " + HttpAddress + ":" + strconv.Itoa(HttpPort) + " started")
	for _,item := range register.namespaces{
		http.Handle(item+"/css/", http.FileServer(http.Dir(ViewsPath)))
		http.Handle(item+"/js/", http.FileServer(http.Dir(ViewsPath)))
		http.Handle(item+"/images/", http.FileServer(http.Dir(ViewsPath)))
	}
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
		if !utils.InSlice(rt.Method(i).Name, exceptMethod) {
			pattern := path.Join(prefix, strings.ToLower(controllerName), strings.ToLower(rt.Method(i).Name))
			Router(pattern, c, rt.Method(i).Name)
		}

	}
}
