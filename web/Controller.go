package web

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"reflect"
	"strings"

	"github.com/artwebs/aogo/utils"

	aolog "github.com/artwebs/aogo/log"
)

type Controller struct {
}

type Context struct {
	Ctl, Fun       string
	w              http.ResponseWriter
	r              *http.Request
	Form           map[string]interface{}
	Data           map[string]interface{}
	UrlKey, UrlVal []string
	session        *Session
	Url            string
}

type ControllerInterface interface {
	Init(ctx *Context)
	WillDid(ctx *Context) bool
	Release(ctx *Context)
}

func (this *Controller) Init(ctx *Context) {

}

func (this *Controller) Release(ctx *Context) {
	if ctx.session != nil {
		defer ctx.session.Release(ctx.w)
	}
}

func (this *Context) Init(w http.ResponseWriter, r *http.Request, ctl ControllerInterface, fun string, data []string) {
	this.Ctl = strings.TrimSuffix(reflect.Indirect(reflect.ValueOf(ctl)).Type().Name(), "Controller")
	this.Fun = fun
	this.w = w
	this.r = r
	this.Url = r.URL.String()
	this.Data = make(map[string]interface{})
	this.Form = make(map[string]interface{})
	this.UrlVal = data[:]
	if RouterParam == 1 && len(data)%2 == 0 {
		index := 0
		for {
			if index >= len(data) {
				break
			}
			this.Form[data[index]] = data[index+1]
			index += 2
		}
	}
	// aolog.DebugTag(this, "r.Form ", r.Form)
	r.ParseForm()
	for k, v := range r.Form {
		if len(v) > 0 {
			this.Form[k] = v[0]
		} else {
			this.Form[k] = v
		}
	}
	// aolog.DebugTag(this, "Form ", this.Form)

}

func (this *Controller) WillDid(ctx *Context) bool {
	return true
}

func (this *Context) GetRequst() *http.Request {
	return this.r
}

func (this *Context) GetRespone() http.ResponseWriter {
	return this.w
}

func (this *Context) SetUrl(arr []string) {
	this.UrlKey = arr[:]

}

func (this *Context) Redirect(url string) {
	http.Redirect(this.w, this.r, url, http.StatusFound)
}

func (this *Context) Header(code int) {
	this.w.WriteHeader(code)
}

func (this *Context) SetHeader(key, value string) {
	this.w.Header().Set(key, value)
}

func (this *Context) WriteString(str string) {
	this.w.Write([]byte(str))
}

func (this *Context) WriteImage(img *utils.Image) {
	this.w.Header().Set("Content-Type", "image/png")
	img.WriteTo(this.w)
}

func (this *Context) WriteJson(data interface{}) {
	content, err := json.Marshal(data)
	if err != nil {
		log.Fatal("WriteJson Fail")
		return
	}
	this.w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
	this.w.Write(content)
}

func (this *Context) Template(args ...string) string {
	tpl := ""
	root := ViewsPath
	if v, ok := this.Data["nspace"]; ok {
		root += v.(string)
	}

	if len(args) == 0 {
		tpl = root + "/" + this.Ctl + "/" + this.Fun
	} else if len(args) == 1 {
		if strings.HasPrefix(args[0], "/") {
			tpl = root + args[0]
		} else {
			tpl = root + "/" + this.Ctl + "/" + args[0]
		}
	} else {
		tpl = root + "/" + args[1] + "/" + args[0]
	}

	if !strings.HasSuffix(tpl, "."+TemplateExt) {
		tpl += "." + TemplateExt
	}
	return tpl
}

func (this *Context) Display(args ...string) {
	tpl := this.Template(args...)
	aolog.InfoTag(this, tpl)
	if _, err := os.Stat(tpl); err != nil {
		aolog.ErrorTag(this, "file "+tpl+" do not exist")
		return
	}

	t, err := template.ParseFiles(tpl)
	if err != nil {
		log.Println(err)
	}

	this.Data["url"] = ""
	this.Data["nspace"] = ""
	if len(this.UrlKey) > 1 {
		this.Data["url"] = "/" + strings.Join(this.UrlKey[:len(this.UrlKey)-1], "/")
	}
	if len(this.UrlKey) > 2 {
		this.Data["nspace"] = "/" + strings.Join(this.UrlKey[:len(this.UrlKey)-2], "/")
	}
	this.Data["res"] = this.Data["nspace"]
	t.Execute(this.w, this.Data)
}

func (this *Context) SetSession(key, value interface{}) {
	if this.session == nil {
		this.session = InitSession()
		this.session.Start(this.w, this.r)
	}
	this.session.Set(key, value)
}

func (this *Context) GetSession(key interface{}) interface{} {
	if this.session == nil {
		this.session = InitSession()
		this.session.Start(this.w, this.r)
	}
	return this.session.Get(key)
}

func (this *Context) FlushSession() {
	if this.session == nil {
		this.session = InitSession()
		this.session.Start(this.w, this.r)
	}
	this.session.Flush()
}

// SaveToFile saves uploaded file to new path.
// it only operates the first one of mutil-upload form file field.
// /data/[file].[ext]
func (this *Context) SaveToFile(fromfile, tofile string) (string, error) {
	if tofile == "" {
		tofile = UploadPath + "/[file].[ext]"
	}
	file, handle, err := this.r.FormFile(fromfile)
	if err != nil {
		return tofile, err
	}
	fileNameArr := strings.Split(handle.Filename, ".")
	defer file.Close()
	tofile = strings.Replace(tofile, "[file]", fileNameArr[0], -1)
	tofile = strings.Replace(tofile, "[ext]", fileNameArr[1], -1)
	f, err := os.OpenFile(tofile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return tofile, err
	}
	defer f.Close()
	io.Copy(f, file)
	return tofile, nil
}

func (this *Context) ServeFile(file string) {
	http.ServeFile(this.w, this.r, file)
}

func (this *Context) GetClientIP(r http.Request) (string, string, error) {
	var ip string
	var port string
	var err error
	ip, port, err = net.SplitHostPort(this.r.RemoteAddr)
	if err != nil {
		return ip, port, fmt.Errorf("userip: %q is not IP:port", this.r.RemoteAddr)

	}

	userIP := net.ParseIP(ip)
	if userIP == nil {
		return ip, port, fmt.Errorf("userip: %q is not IP:port", this.r.RemoteAddr)
	}
	return ip, port, nil
}

func (this *Context) IsMobile() bool {
	if strings.Index(this.r.Header.Get("User-Agent"), "Mobile") > -1 {
		return true
	}
	return false
}

func (this *Context) GetUserAgent() string {
	return this.r.Header.Get("User-Agent")
}
