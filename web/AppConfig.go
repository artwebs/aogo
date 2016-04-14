package web

import (
	// "fmt"
	"github.com/astaxie/beego/config"
	"log"
)

var (
	conf        *AppConfig
	HttpAddress = ""
	HttpPort    = 8080

	ViewsPath   = "views"
	TemplateExt = "html"
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
	}
	return conf, nil

}

func (this *AppConfig) String(key, def string) string {
	if this.obj == nil {
		return def
	}
	return this.obj.String(key)
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
