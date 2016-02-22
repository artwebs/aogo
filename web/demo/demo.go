package demo

import (
	"github.com/artwebs/aogo/log"
	"github.com/artwebs/aogo/web"
)

func SubData() string {
	return "subData"
}

type TestController struct {
	web.Controller
}

func init() {
	log.NewLogger(1000)
	log.SetLogger("console", "")
	log.SetLevel(log.LevelDebug)
}
func (this *TestController) Index() {
	this.Writer.Write([]byte("artwebs"))
}

func (this *TestController) TestTpl() {

	this.Data["name"] = "hello"
	log.Debug("hello")
	this.Display("TestTpl.html")
}

func main() {
	web.Router("/test", &TestController{}, "Index")
	web.Router("/tpl", &TestController{}, "TestTpl")
	web.Run()
}
