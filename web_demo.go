package main

import (
	aolog "github.com/artwebs/aogo/log"
	"github.com/artwebs/aogo/web"
)

type TestController struct {
	web.Controller
}

func (this *TestController) Index() {
	if user, ok := this.Form["user"]; ok {
		this.SetSession("user", user)
	}
	aolog.Info(this.Data)
	tmp := make(map[string]interface{})
	tmp["a"] = 1
	tmp["b"] = "jom"
	this.WriteJson(tmp)
	// this.WriteString("artwebs")
	// this.WriteString("artwebs1")
}

func (this *TestController) TestTpl() {
	aolog.Info(this.GetSession("user"))
	this.Data["name"] = "hello"
	this.Display()
}

func (this *TestController) Loginout() {
	this.FlushSession()
}

func (this *TestController) Upload() {
	this.Display()
}

func (this *TestController) Save() {
	err := this.SaveToFile("UpLoadFile", "")
	if err == nil {
		this.WriteString("success")
	} else {
		this.WriteString("fail")
	}
}

func webRun() {
	web.Router("/tet", &TestController{}, "Index")
	web.Router("/tpl", &TestController{}, "TestTpl")
	web.AutoRouter("/", &TestController{})
	ns := web.NewNamespace("/ns",
		web.NSRouter("/demo", &TestController{}, "TestTpl"),
		web.NSRouter("/demo1", &TestController{}, "Index"),
		web.NSNamespace("/demo2",
			web.NSRouter("/demo1", &TestController{}, "Index"),
		),
		web.NSAutoRouter("/", &TestController{}),
	)
	web.AddNamespace(ns)
	web.Run()
}
