package web

import (
	"encoding/json"
	aolog "github.com/artwebs/aogo/log"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"
)

type Controller struct {
	Ctl, Fun string
	w        http.ResponseWriter
	r        *http.Request
	Form     map[string]interface{}
	Data     map[string]interface{}
	session  *Session
}

type ControllerInterface interface {
	Init(w http.ResponseWriter, r *http.Request, ctl ControllerInterface, fun string, data []string)
	WillDid() bool
	SetUrl(arr []string)
	Redirect(url string)
	WriteString(str string)
	WriteJson(obj interface{})
	Release()
	SetSession(key, value interface{})
	GetSession(key interface{}) interface{}
	FlushSession()
	SaveToFile(fromfile, tofile string) error
}

func (this *Controller) Init(w http.ResponseWriter, r *http.Request, ctl ControllerInterface, fun string, data []string) {
	this.Ctl = strings.TrimSuffix(reflect.Indirect(reflect.ValueOf(ctl)).Type().Name(), "Controller")
	this.Fun = fun
	this.w = w
	this.r = r
	this.Data = make(map[string]interface{})
	this.Form = make(map[string]interface{})
	if len(data)%2 == 0 {
		index := 0
		for {
			if index >= len(data) {
				break
			}
			this.Form[data[index]] = data[index+1]
			index += 2
		}
	}

	for k, v := range r.Form {
		if len(v) > 0 {
			this.Form[k] = v[0]
		} else {
			this.Form[k] = v
		}
	}
	aolog.DebugTag(this, "Form ", this.Form)

}

func (this *Controller) WillDid() bool {
	return true
}

func (this *Controller) SetUrl(arr []string) {
	this.Data["url"] = strings.Join(arr[:len(arr)-1], "/")
	this.Data["nspace"] = strings.Join(arr[:len(arr)-2], "/")
}

func (this *Controller) Redirect(url string) {
	http.Redirect(this.w, this.r, url, http.StatusFound)
}

func (this *Controller) WriteString(str string) {
	this.w.Write([]byte(str))
}

func (this *Controller) WriteJson(data interface{}) {

	content, err := json.Marshal(data)
	if err != nil {
		log.Fatal("WriteJson Fail")
		return
	}
	this.w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	this.w.Write(content)
}

func (this *Controller) Display(args ...string) {
	tpl := ""
	if len(args) == 0 {
		tpl = ViewsPath + "/" + this.Ctl + "/" + this.Fun + "." + TemplateExt
	} else if len(args) == 1 {
		tpl = ViewsPath + "/" + this.Ctl + "/" + args[0] + "." + TemplateExt
	} else {
		tpl = ViewsPath + "/" + args[1] + "/" + args[0] + "." + TemplateExt
	}
	log.Println(tpl)
	t, err := template.ParseFiles(tpl)
	if err != nil {
		log.Println(err)
	}
	t.Execute(this.w, this.Data)
}

func (this *Controller) Release() {
	if this.session != nil {
		defer this.session.Release(this.w)
	}

}

func (this *Controller) SetSession(key, value interface{}) {
	if this.session == nil {
		this.session = InitSession()
		this.session.Start(this.w, this.r)
	}
	this.session.Set(key, value)
}

func (this *Controller) GetSession(key interface{}) interface{} {
	if this.session == nil {
		this.session = InitSession()
		this.session.Start(this.w, this.r)
	}
	return this.session.Get(key)
}

func (this *Controller) FlushSession() {
	if this.session == nil {
		this.session = InitSession()
		this.session.Start(this.w, this.r)
	}
	this.session.Flush()
}

// SaveToFile saves uploaded file to new path.
// it only operates the first one of mutil-upload form file field.
// /data/[file].[ext]
func (this *Controller) SaveToFile(fromfile, tofile string) error {
	if tofile == "" {
		tofile = "[file].[ext]"
	}
	file, handle, err := this.r.FormFile(fromfile)
	if err != nil {
		return err
	}
	fileNameArr := strings.Split(handle.Filename, ".")
	defer file.Close()
	tofile = strings.Replace(tofile, "[file]", fileNameArr[0], -1)
	tofile = strings.Replace(tofile, "[ext]", fileNameArr[1], -1)
	f, err := os.OpenFile(tofile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer f.Close()
	io.Copy(f, file)
	return nil
}
