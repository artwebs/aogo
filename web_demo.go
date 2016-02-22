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
	web.Router("/test", &TestController{}, "Index")
	web.Router("/tpl", &TestController{}, "TestTpl")
	web.Run()
}
