package web

import (
	"database/sql"
	"log"
	"reflect"
	"strings"
	// "strconv"
	"github.com/artwebs/aogo/cache"
	"github.com/artwebs/aogo/database"
	// aolog "github.com/artwebs/aogo/log"
)

func D(model ModelInterface, args ...string) ModelInterface {
	model.Init(args...)
	model.SetTabName(strings.TrimSuffix(reflect.Indirect(reflect.ValueOf(model)).Type().Name(), "Model"))
	return model
}

type ModelInterface interface {
	SetTabName(name string)
	Init(args ...string)
}

type Model struct {
	Drv database.DriverInterface
}

func (this *Model) SetTabName(name string) {
	this.Drv.SetTabName(name)
}

func (this *Model) Init(args ...string) {
	var err error
	var dbPrifix, driverName, dataSourceName, tabPrifix string
	for i, v := range args {
		switch i {
		case 0:
			dbPrifix = v
		}
	}
	conf, err = InitAppConfig()
	if err == nil {
		driverName = conf.String(dbPrifix+"DataBase::driverName", "")
		dataSourceName = conf.String(dbPrifix+"DataBase::dataSourceName", "")
		tabPrifix = conf.String(dbPrifix+"DataBase::tabPrifix", "")
	} else {
		log.Fatalln("AppConfig init fail")
	}

	CobjName := conf.String("Cache::name", "")
	CobjConfig := conf.String("Cache::config", "")
	var Cobj *cache.Cache
	if CobjName != "" && CobjConfig != "" {
		Cobj, err = cache.NewCache(CobjName, CobjConfig)
	}
	this.Drv = database.Drivers(driverName)
	this.Drv.Init(driverName, dataSourceName, tabPrifix)
	this.Drv.SetCache(Cobj)
	this.Drv.SetDBPrifix(dbPrifix)
}

func (this *Model) Query(s string, args ...interface{}) ([]map[string]string, error) {
	return this.Drv.Query(s, args...)
}

func (this *Model) QueryNoConn(s string, args ...interface{}) ([]map[string]string, error) {
	return this.Drv.QueryNoConn(s, args...)
}

func (this *Model) QueryRow(s string, args ...interface{}) (map[string]string, error) {
	return this.Drv.QueryRow(s, args...)
}

func (this *Model) QueryRowNoConn(s string, args ...interface{}) (map[string]string, error) {
	return this.Drv.QueryRowNoConn(s, args...)

}

func (this *Model) Exec(sql string, args ...interface{}) (sql.Result, error) {
	return this.Drv.Exec(sql, args...)
}

func (this *Model) ExecNoConn(sql string, args ...interface{}) (sql.Result, error) {
	return this.Drv.ExecNoConn(sql, args...)
}

func (this *Model) Insert(values map[string]interface{}) (int64, error) {
	return this.Drv.Insert(this.Drv, values)
}

func (this *Model) Update(values map[string]interface{}) (int64, error) {
	return this.Drv.Update(this.Drv, values)
}

func (this *Model) Delete() (int64, error) {
	return this.Drv.Delete(this.Drv)
}

func (this *Model) Find() (map[string]string, error) {
	return this.Drv.Find(this.Drv)

}

func (this *Model) Total() (int, error) {
	return this.Drv.Total(this.Drv)
}

func (this *Model) Select() ([]map[string]string, error) {
	return this.Drv.Select(this.Drv)

}

func (this *Model) Where(w string, args ...interface{}) *Model {
	this.Drv.Where(w, args...)
	return this
}

func (this *Model) Order(o string) *Model {
	this.Drv.Order(o)
	return this
}

func (this *Model) Limit(l string) *Model {
	this.Drv.Limit(l)
	return this
}

func (this *Model) Group(g string) *Model {
	this.Drv.Group(g)
	return this
}

func (this *Model) Having(h string) *Model {
	this.Drv.Having(h)
	return this
}

func (this *Model) Field(fields ...string) *Model {
	this.Drv.Field(fields...)
	return this
}

func (this *Model) IsCache(flag bool) *Model {
	this.Drv.IsCache(flag)
	return this
}
