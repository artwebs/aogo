package main

import (
	"github.com/artwebs/aogo/web"
)

func SubData() string {
	return "subData"
}

type TestController struct {
	web.Controller
}

func (this *TestController) Index() {
	this.Writer.Write([]byte("artwebs"))
}

func (this *TestController) TestTpl() {
	this.Data["name"] = "hello"
	this.Display("TestTpl.html")
}

func main() {
	web.Router("/test", &TestController{}, "Index")
	web.Router("/tpl", &TestController{}, "TestTpl")
	web.Run()
}
