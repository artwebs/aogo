package main

import (
	"encoding/json"
	"github.com/artwebs/aogo/log"
	"github.com/artwebs/aogo/web"
	"io/ioutil"
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
	web.Router("/index", &IndexController{}, "Index")
	web.HandleFunc("/test", HelloServer)
	web.Run()

	// http.Handle("/images/", http.FileServer(http.Dir("views")))
	// http.ListenAndServe(":8080", nil)
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
	if val, ok := routers[router]; ok {
		log.InfoTag(this, val)
	}
	log.InfoTag(this, this.UrlKey)
	this.WriteString("Index")
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
