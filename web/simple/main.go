package main

import (
	"encoding/json"
	"io/ioutil"
	"strconv"

	"github.com/artwebs/aogo/log"
	"github.com/artwebs/aogo/utils"
	"github.com/artwebs/aogo/web"
	// "net/http"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
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
	reg, _ := regexp.Compile("^[^rw_]\\w.+\\.json$")
	// 遍历目录
	filepath.Walk("./router",
		func(path string, f os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if f.IsDir() {
				return nil
			}
			// 匹配目录
			matched := reg.MatchString(f.Name())
			if matched {
				log.Info("reload", path)
				sstr = readSimple(path)
				temp := make(map[string]*Router)
				if sstr != "" {
					json.Unmarshal([]byte(sstr), &temp)
					for k, v := range temp {
						if strings.HasPrefix(k, "/") {
							routers[k] = v
						} else {
							routers["/"+strings.TrimSuffix(f.Name(), ".json")+"/"+k] = v
						}
					}
				}
			}
			return nil
		})
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
		if !this.verfiySession(val.Session) {
			return
		}
		model := &DefaultModel{}
		web.D(model)
		if val.Data != nil {
			for key, value := range val.Data {
				switch value {
				case "SaveFile":
					file, err := this.SaveToFile("File", "")
					if err == nil {
						this.Data[key] = map[string]interface{}{"code": 1, "message": "上传成功", "result": file}
					} else {
						this.Data[key] = map[string]interface{}{"code": 0, "message": err, "result": ""}
					}

					break
				case "VerfiyCode":
					d := make([]byte, 4)
					s := utils.NewLen(4)
					ss := ""
					d = []byte(s)
					for v := range d {
						d[v] %= 10
						ss += strconv.FormatInt(int64(d[v]), 32)
					}
					this.SetSession(key, ss)
					this.WriteImage(utils.NewImage(d, 100, 40))
					return
				default:
					this.Data[key] = model.Aws(value, this.Form)
				}

			}
		}
		this.doSession(val.Session)
		tpl := router
		if val.Tpl != "" {
			tpl = val.Tpl
		}
		this.Data["Requst"] = this.Form
		if tpl == "json" {
			this.WriteJson(this.Data)
		} else {
			this.Display(tpl)
		}

	} else {
		this.WriteString(router + " do not find!")
	}

}

func (this *IndexController) doSession(sin string) {
	for _, ss := range strings.Split(sin, ",") {
		s := strings.Split(ss, ":")
		if len(s) < 2 {
			continue
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
}

func (this *IndexController) verfiySession(sin string) bool {
	for _, ss := range strings.Split(sin, ",") {
		s := strings.Split(ss, ":")
		if len(s) == 1 {
			if val, ok := sessions[s[0]]; ok {
				if this.GetSession(s[0]) != nil {
					log.InfoTag(this, this.GetSession(s[0]))
					if val.Verfiy != "" {
						if this.Form[val.Verfiy] != this.GetSession(s[0]) {
							this.Redirect(val.Fail)
							return false
						}
					}
					cursession := (this.GetSession(s[0])).(map[string]interface{})
					for k, v := range cursession {
						if _, tok := this.Form[k]; tok {
							this.Form["_"+k] = v
						} else {
							this.Form[k] = v
						}
					}
					continue
				}
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
	Name, Fail, Verfiy string
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
