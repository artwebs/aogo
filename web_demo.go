package main

import (
	aolog "github.com/artwebs/aogo/log"
	"github.com/artwebs/aogo/web"
)

type TestController struct {
	web.Controller
}

func (this *TestController) Index() {
	this.Writer.Write([]byte("artwebs"))
}

func (this *TestController) TestTpl() {

	this.Data["name"] = "hello"
	aolog.Debug("hello")
	this.Display("TestTpl.html")
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
	)
	web.AddNamespace(ns)
	web.Run()
}
