package web

import (
	// "fmt"
	"log"

	"github.com/astaxie/beego/config"
)

var (
	conf        *AppConfig
	HttpAddress = ""
	HttpPort    = 8080

	ViewsPath   = "views"
	TemplateExt = "html"
	UploadPath  = "files"
	RouterParam = 0
	Debug       = 1
)

type AppConfig struct {
	obj config.Configer
}

func InitAppConfig() (*AppConfig, error) {
	if conf == nil {
		conf = &AppConfig{}
	}

	if conf.obj == nil {
		var err error
		conf.obj, err = config.NewConfig("ini", "app.conf")
		if err != nil {
			log.Println("no app.conf")
			return nil, err
		}

		HttpAddress = conf.String("HttpAddress", "")
		HttpPort = conf.Int("HttpPort", 8080)
		ViewsPath = conf.String("ViewsPath", "views")
		TemplateExt = conf.String("TemplateExt", "html")
		UploadPath = conf.String("TemplateExt", "files")
		RouterParam = conf.Int("RouterParam", 0)
		Debug = conf.Int("Debug", 1)
	}
	return conf, nil

}

func (this *AppConfig) String(key, def string) string {
	if this.obj == nil {
		return def
	}
	val := this.obj.String(key)
	if val == "" {
		return def
	}
	return val
}

func (this *AppConfig) Int(key string, def int) int {
	if this.obj == nil {
		return def
	}
	val, err := this.obj.Int(key)
	if err != nil {
		return def
	}
	return val
}
