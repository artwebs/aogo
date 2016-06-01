package main

import (
	"github.com/artwebs/aogo/log"
	"github.com/artwebs/aogo/utils"
	"github.com/artwebs/aogo/web"
)

var (
	filter []string
)

func init() {
	reload()
}

func main() {
	web.Router("/", &IndexController{}, "Index")
	web.Run()
}

type IndexController struct {
	web.Controller
}

func (this *IndexController) Index() {
	log.InfoTag(this, this.UrlVal)
	this.WriteString("hello world")
	if len(this.UrlVal) == 1 {
		key := this.UrlVal[0]
		if utils.InSlice(key, filter) {

		}
	}
}

func reload() {
	filter = []string{"login"}
}
