package main

import (
	"encoding/json"
	"flag"
	"strings"

	"github.com/artwebs/aogo/log"
	"github.com/artwebs/aogo/utils"
	"github.com/artwebs/aogo/web"
)

var (
	filter []string
)

func main() {
	version := flag.Bool("version", false, "--version")
	flag.Parse()
	log.Info("version:v1.0.0")
	if *version {
		return
	}
	reload()
	web.Router("/", &IndexController{}, "Index")
	web.Run()

}

func reload() {
	conf, err := web.InitAppConfig()
	if err == nil {
		sin := conf.String("Session::login", "")
		if sin != "" {
			filter = strings.Split(sin, ":")
			return
		}
	}

	filter = []string{"login"}
}

type IndexController struct {
	web.Controller
}

func data(code, count int, message string, obj interface{}) map[string]interface{} {
	return map[string]interface{}{"code": code, "count": count, "message": message, "result": obj}
}

func data_error(message string) map[string]interface{} {
	return data(0, 0, message, "")
}

func (this *IndexController) Index() {
	log.InfoTag(this, this.UrlVal)

	if len(this.UrlVal) >= 1 {
		model := &DefaultModel{}
		web.D(model)
		key := this.UrlVal[0]
		if !utils.InSlice(key, filter) {
			if len(this.UrlVal) < 2 {
				this.WriteJson(data_error("请先进行登录"))
				return
			}

			sin := this.UrlVal[1]
			val := this.GetSession(sin)
			if val != nil {
				if key == "upload" {
					file, err := this.SaveToFile("File", "")
					if err == nil {
						this.WriteJson(data(1, 0, "文件上传成功", map[string]interface{}{"file": file}))
					} else {
						data_error("文件删除失败")
					}
					return
				}
				if key == "download" {
					file := strings.Join(this.UrlVal[2:], "/")
					this.ServeFile(strings.TrimPrefix(file, "/"))
					return
				}

				cursession := (this.GetSession(sin)).(map[string]interface{})
				for k, v := range cursession {
					if _, tok := this.Form[k]; tok {
						this.Form["_"+k] = v
					} else {
						this.Form[k] = v
					}
				}
				this.WriteJson(model.Aws(key, this.Form))
			} else {
				this.WriteJson(data_error("非法登录"))
			}
		} else {
			data := model.Aws(key, this.Form)
			if code, ok := data["code"]; ok {
				if code.(float64) > 0 {
					this.SetSession(key, data["result"])
				}
				this.WriteJson(data)
			} else {
				this.WriteJson(data_error("验证失败"))
			}
		}
	} else {
		this.WriteJson(data_error("非法请求，已经进行了记录"))
	}
}

type DefaultModel struct {
	web.Model
}

func (this *DefaultModel) Aws(name string, args map[string]interface{}) map[string]interface{} {
	data, _ := json.Marshal(args)
	var notused string
	log.InfoTag(this, "drv", this.Drv)
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
