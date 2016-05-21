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
	reload()
	web.Router("/", &IndexController{}, "Index")
	web.HandleFunc("/reload", GOReload)
	web.Run()

}

func reload() {
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

}

func GOReload(w http.ResponseWriter, req *http.Request) {
	reload()
	io.WriteString(w, "reload sucess!\n")
}

type IndexController struct {
	web.Controller
}

func (this *IndexController) Index() {
	router := strings.Join(this.UrlKey, "/")
	log.InfoTag(this, router)

	if val, ok := routers[router]; ok {
		sinArr := []string{}
		if val.Session != "" {
			sinArr = strings.Split(val.Session, ":")
		}

		if !this.verfiySession(sinArr) {
			return
		}
		model := &DefaultModel{}
		web.D(model)
		if val.Tpl != "" {
			for key, value := range val.Data {
				this.Data[key] = model.Aws(value, this.Form)
			}
			this.doSession(sinArr)
			log.InfoTag(this, this.Data)
			this.Display(val.Tpl)
		} else {
			this.WriteJson(this.Data)
		}

	} else {
		this.WriteString(router + " do not find!")
	}

}

func (this *IndexController) doSession(s []string) {

	if len(s) < 2 {
		return
	}

	switch s[1] {
	case "save":

		if val, ok := sessions[s[0]]; ok {
			data := (this.Data[val.Name]).(map[string]interface{})
			if (data["code"]).(float64) == 1 {
				log.InfoTag(this, "doSession save", data["result"])
				this.SetSession(s[0], data["result"])
			}
		}
		break
	case "delete":
		this.FlushSession()
		break
	default:

	}
}

func (this *IndexController) verfiySession(s []string) bool {
	if len(s) == 1 {
		if val, ok := sessions[s[0]]; ok {
			if this.GetSession(s[0]) != nil {
				log.InfoTag(this, this.GetSession(s[0]))
				cursession := (this.GetSession(s[0])).(map[string]interface{})
				for k, v := range cursession {
					this.Form[k] = v
				}
			} else {
				this.Redirect(val.Fail)
				return false
			}
		}
	}
	return true
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
