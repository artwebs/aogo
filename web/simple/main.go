package main

import (
	"encoding/json"
	"io/ioutil"

	"github.com/artwebs/aogo/log"
	"github.com/artwebs/aogo/web"
	// "net/http"
	"io"
	"net/http"
	"os"
	"strings"
)

var (
	sessions map[string]*Session
	routers  map[string]*Router
)

func main() {
	sessions = make(map[string]*Session)
	routers = make(map[string]*Router)
	sstr := readSimple("session.json")
	if sstr != "" {
		json.Unmarshal([]byte(sstr), &sessions)
	}
	log.Info("session", sessions)
	sstr = readSimple("router.json")
	if sstr != "" {
		json.Unmarshal([]byte(sstr), &routers)
	}
	log.Info("router", routers)
	web.Router("/", &IndexController{}, "Index")
	web.HandleFunc("/test", HelloServer)
	web.Run()

}

func HelloServer(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "hello, world!\n")
}

type IndexController struct {
	web.Controller
}

func (this *IndexController) Index() {
	router := strings.Join(this.UrlKey, "/")
	log.InfoTag(this, router)
	model := &DefaultModel{}
	web.D(model)
	if val, ok := routers[router]; ok {
		if val.Tpl != "" {
			for key, value := range val.Data {
				this.Data[key] = model.Aws(value, this.Form)
			}
			log.InfoTag(this, this.Data)
			this.Display(val.Tpl)
		} else {
			this.WriteJson(this.Data)
		}

	} else {
		this.WriteString(router + " do not find!")
	}

}

type DefaultModel struct {
	web.Model
}

func (this *DefaultModel) Aws(name string, args map[string]interface{}) map[string]interface{} {
	data, _ := json.Marshal(args)
	var notused string
	this.Drv.Conn()
	err := this.Drv.Db().QueryRow("SELECT aws($1,$2)", name, string(data)).Scan(&notused)
	defer this.Drv.Close()
	if err != nil {
		log.ErrorTag(this, err)
		return make(map[string]interface{})
	}
	result := make(map[string]interface{})
	json.Unmarshal([]byte(notused), &result)
	return result
}

// {
// 	"session":{
// 		"session1":{"name":"user","fail":"/index/error"}
// 		},
// 	"router":{
// 		"/index":{
// 			"tpl":"/index.html",
// 			"data":{"data1":"test1","data2":"test2"},
// 			"session":"session1"
// 			},
// 		"/index/test:session1":{
// 			"tpl":"/test.html",
// 			"data":{"data1":"test1","data2":"test2"}
// 		}
// 	}
// }

type Session struct {
	Name, Fail string
}

type Router struct {
	Tpl, Session string
	Data         map[string]string
}

func readSimple(path string) string {
	fi, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer fi.Close()
	fd, err := ioutil.ReadAll(fi)
	// fmt.Println(string(fd))
	return string(fd)
}
