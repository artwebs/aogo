package web

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"reflect"
	"strings"
)

type Controller struct {
	ctl, fun string
	Writer   http.ResponseWriter
	Request  *http.Request
	Form     map[string]interface{}
	Data     map[string]interface{}
	session  *Session
}

type ControllerInterface interface {
	Init(w http.ResponseWriter, r *http.Request, ctl ControllerInterface, fun string, data []string)
	SetUrl(arr []string)
	Release()
	SetSession(key, value interface{})
	GetSession(key interface{}) interface{}
	FlushSession()
	SaveToFile(fromfile, tofile string) error
}

func (this *Controller) Init(w http.ResponseWriter, r *http.Request, ctl ControllerInterface, fun string, data []string) {
	this.ctl = strings.TrimSuffix(reflect.Indirect(reflect.ValueOf(ctl)).Type().Name(), "Controller")
	this.fun = fun
	this.Writer = w
	this.Request = r
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

	log.Println(this.Form)
}

func (this *Controller) SetUrl(arr []string) {
	this.Data["url"] = strings.Join(arr[:len(arr)-1], "/")
	this.Data["nspace"] = strings.Join(arr[:len(arr)-2], "/")
}

func (this *Controller) Display(args ...string) {
	tpl := ""
	if len(args) == 0 {
		tpl = StaticPath + "/" + this.ctl + "/" + this.fun + "." + TemplateExt
	} else if len(args) == 1 {
		tpl = StaticPath + "/" + this.ctl + "/" + args[0] + "." + TemplateExt
	} else {
		tpl = StaticPath + "/" + args[1] + "/" + args[0] + "." + TemplateExt
	}
	log.Println(tpl)
	t, err := template.ParseFiles(tpl)
	if err != nil {
		log.Println(err)
	}
	t.Execute(this.Writer, this.Data)
}

func (this *Controller) Release() {
	if this.session != nil {
		defer this.session.Release(this.Writer)
	}

}

func (this *Controller) SetSession(key, value interface{}) {
	if this.session == nil {
		this.session = InitSession()
		this.session.Start(this.Writer, this.Request)
	}
	this.session.Set(key, value)
}

func (this *Controller) GetSession(key interface{}) interface{} {
	if this.session == nil {
		this.session = InitSession()
		this.session.Start(this.Writer, this.Request)
	}
	return this.session.Get(key)
}

func (this *Controller) FlushSession() {
	if this.session == nil {
		this.session = InitSession()
		this.session.Start(this.Writer, this.Request)
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
	file, handle, err := this.Request.FormFile(fromfile)
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
