package main

import (
	"encoding/json"
	"errors"
	"flag"
	"reflect"
	"strings"

	"github.com/artwebs/aogo/log"
	"github.com/artwebs/aogo/security"
	"github.com/artwebs/aogo/utils"
	"github.com/artwebs/aogo/web"
)

var (
	filter               []string
	runmode, securitykey string
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
		} else {
			filter = []string{"login"}
		}

		runmode = conf.String("RunMode", "dev")
		securitykey = conf.String("SecurityKey", "Y8gyxetKJ68N3d35Lass72GP")
	}

}

type IndexController struct {
	web.Controller
	isEncrypt bool
}

func data(code, count int, message string, obj interface{}) map[string]interface{} {
	return map[string]interface{}{"code": code, "count": count, "message": message, "result": obj}
}

func data_error(message string) map[string]interface{} {
	return data(0, 0, message, "")
}

func (this *IndexController) Index() {
	this.isEncrypt = false
	log.InfoTag(this, this.UrlVal)
	this.parseQuery()
	if len(this.UrlVal) >= 1 {
		key := this.UrlVal[0]
		if val := reflect.ValueOf(this).MethodByName(key); val.IsValid() || key == strings.Title(key) {
			val.Call([]reflect.Value{})
		} else {
			this.Normal()
		}

	} else {
		this.WriteJson(data_error("非法请求，已经进行了记录"))
	}
}

func (this *IndexController) Login() {
	key := this.UrlVal[0]
	if _, ok := this.Form["appId"]; !ok {
		goto fail
	}

	if _, ok := this.Form["clientId"]; !ok {
		goto fail
	}

	if _, ok := this.Form["clientVersion"]; !ok {
		goto fail
	}

	this.SetSession(key, this.Form)
	this.write(data(1, 0, "登录成功", ""))
	return
fail:
	{
		this.write(data_error("非法登录！"))
		return
	}
}

func (this *IndexController) Normal() {
	if err := this.verfiySession(); err != nil {
		this.write(data_error(err.Error()))
		return
	}
	model := &DefaultModel{}
	web.D(model)
	key := this.UrlVal[0]
	this.write(model.Aws(key, this.Form))
}

func (this *IndexController) Upload() {
	if err := this.verfiySession(); err != nil {
		this.WriteJson(data_error(err.Error()))
		return
	}
	file, err := this.SaveToFile("File", "")
	if err == nil {
		this.write(data(1, 0, "文件上传成功", map[string]interface{}{"file": file}))
	} else {
		this.write(data_error("文件删除失败"))
	}
}

func (this *IndexController) Download() {
	if err := this.verfiySession(); err != nil {
		this.write(data_error(err.Error()))
		return
	}
	file := strings.Join(this.UrlVal[2:], "/")
	this.ServeFile(strings.TrimPrefix(file, "/"))
}

func (this *IndexController) verfiySession() error {
	client := this.GetSession("Login")
	if client == nil {
		return errors.New("非法请求，已经进行了记录")
	}
	this.addForm(client.(map[string]interface{}))
	return nil
}

func (this *IndexController) parseQuery() error {
	if val, ok := this.Form["cmd"]; ok {
		aesObj := security.NewSecurityAES()
		str, err := aesObj.DecryptString(securitykey, val.(string))
		if err == nil {
			json.Unmarshal([]byte(str), &this.Form)
			delete(this.Form, "cmd")
			this.isEncrypt = true
		} else {
			return errors.New("非法数据请求，已经进行了记录")
		}
	} else if runmode == "product" {
		return errors.New("数据非法请求，已经进行了记录")
	}
	return nil

}

func (this *IndexController) Encode() {
	str, err := this.encode(this.Form)
	if err == nil {
		this.WriteString(utils.UrlEncode(str))
	} else {
		this.WriteString(err.Error())
	}

}

func (this *IndexController) Decode() {
	this.WriteJson(this.Form)
}

func (this *IndexController) write(d map[string]interface{}) {
	if runmode == "product" || this.isEncrypt {
		data, _ := this.encode(d)
		this.WriteString(data)
	} else {
		this.WriteJson(d)
	}

}

func (this *IndexController) encode(d map[string]interface{}) (string, error) {
	data, _ := json.Marshal(d)
	aesObj := security.NewSecurityAES()
	return aesObj.EncryptString(securitykey, string(data))
}

func (this *IndexController) addForm(obj map[string]interface{}) {
	for k, v := range obj {
		if _, tok := this.Form[k]; tok {
			this.Form["_"+k] = v
		} else {
			this.Form[k] = v
		}
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
