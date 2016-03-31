package web

import (
	"fmt"
	"github.com/astaxie/beego/config"
	"log"
)

var (
	confObj     *config.ConfigContainer
	HttpAddress = ""
	HttpPort    = 8080

	ViewsPath   = "views"
	TemplateExt = "html"
)

type AppConfig struct {
}

func InitAppConfig() {
	confObj, err := config.NewConfig("ini", "app.conf")
	if err != nil {
		log.Println("no app.conf")
		return
	}
	HttpAddress = AppConfigString("HttpAddress", "")
	HttpPort = AppConfigInt("HttpPort", 8080)
}

func AppConfigString(key, def string) string {
	if confObj != nil {
		return def
	}
	return confObj.String(key)
}

func AppConfigInt(key string, def int) int {
	if confObj != nil {
		return def
	}
	val, err := confObj.Int(key)
	if err != nil {
		return def
	}
	return val
}
