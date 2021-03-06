package main

import (
	aolog "github.com/artwebs/aogo/log"
	"github.com/artwebs/aogo/web"
)

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

type TestController struct {
	web.Controller
}

func (this *TestController) Index() {
	// model := web.D(new(TestModel))
	model := &UserModel{}
	// model.DBPrifix = "PG"
	// model := NewPgUser()
	web.D(model, "PG")
	// aolog.Info(model.Insert(map[string]interface{}{"name": "test"}))
	// aolog.Info(model.Where("id=?", 1).Update(map[string]interface{}{"name": "test1"}))
	// aolog.Info(model.Where("id=?",16).Delete())
	// aolog.Info(model.Query("select * from user"))
	// aolog.Info(model.Where("id=?",1).Find())
	// aolog.Info(model.Total())
	// aolog.Info(model.Where("id=?",1).Total())
	aolog.Info(model.Order("id desc").Select())
	if user, ok := this.Form["user"]; ok {
		this.SetSession("user", user)
	}
	aolog.Info(this.Data)
	tmp := make(map[string]interface{})
	tmp["a"] = 1
	tmp["b"] = "jom"
	this.WriteJson(tmp)
	// this.WriteString("artwebs")
	// this.WriteString("artwebs1")
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
	_, err := this.SaveToFile("UpLoadFile", "")
	if err == nil {
		this.WriteString("success")
	} else {
		this.WriteString("fail")
	}
}

type UserModel struct {
	// web.Model
	web.Model
}

func (this *UserModel) DoTest() {
	aolog.Info("DoTest")

}
