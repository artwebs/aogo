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

func (this *TestController) Index(ctx *web.Context) {
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
	if user, ok := ctx.Form["user"]; ok {
		ctx.SetSession("user", user)
	}
	aolog.Info(ctx.Data)
	tmp := make(map[string]interface{})
	tmp["a"] = 1
	tmp["b"] = "jom"
	ctx.WriteJson(tmp)
	// this.WriteString("artwebs")
	// this.WriteString("artwebs1")
}

func (this *TestController) TestTpl(ctx *web.Context) {
	aolog.Info(ctx.GetSession("user"))
	ctx.Data["name"] = "hello"
	ctx.Display()
}

func (this *TestController) Loginout(ctx *web.Context) {
	ctx.FlushSession()
}

func (this *TestController) Upload(ctx *web.Context) {
	ctx.Display()
}

func (this *TestController) Save(ctx *web.Context) {
	_, err := ctx.SaveToFile("UpLoadFile", "")
	if err == nil {
		ctx.WriteString("success")
	} else {
		ctx.WriteString("fail")
	}
}

type UserModel struct {
	// web.Model
	web.Model
}

func (this *UserModel) DoTest() {
	aolog.Info("DoTest")

}
