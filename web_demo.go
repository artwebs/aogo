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
	this.Writer.Write([]byte("artwebs"))
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
		this.Writer.Write([]byte("success"))
	} else {
		this.Writer.Write([]byte("fail"))
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
