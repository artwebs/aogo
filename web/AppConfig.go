package web

import (
	"fmt"
	"github.com/astaxie/beego/config"
	"log"
)

var (
	HttpAddress = ""
	HttpPort    = 8080

	StaticPath  = "views"
	TemplateExt = "html"
)

func InitAppConfig() {
	conf, err := config.NewConfig("ini", "app.conf")
	if err != nil {
		log.Println("no app.conf")
		return
	}
	HttpAddress = conf.String("HttpAddress")
	HttpPort, err = conf.Int("HttpPort")
	fmt.Printf("HttpAddress=%s,HttpPort=%d\n", HttpAddress, HttpPort)
}
