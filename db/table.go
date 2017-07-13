package db

import (
	"log"

	aolog "github.com/artwebs/aogo/log"
	"github.com/astaxie/beego/config"
)

func Table(tb string, args ...string) DriverInterface {
	var err error
	var dbPrifix, driverName, dataSourceName, tabPrifix string
	for i, v := range args {
		switch i {
		case 0:
			dbPrifix = v
		}
	}
	cobj, err := config.NewConfig("ini", "app.conf")
	if err != nil {
		log.Fatalln("no app.conf")
		return nil
	}

	driverName = cobj.String(dbPrifix + "DataBase::driverName")
	dataSourceName = cobj.String(dbPrifix + "DataBase::dataSourceName")
	tabPrifix = cobj.String(dbPrifix + "DataBase::tabPrifix")

	CobjName := cobj.String("DBCache::name")
	CobjConfig := cobj.String("DBCache::config")
	var Cobj DBCache
	if CobjName != "" && CobjConfig != "" {
		Cobj = OpenDBCache(CobjName, CobjConfig)
	}
	if err != nil {
		aolog.Info("db::Table", "dataSourceName", dataSourceName)
	}
	drv := Drivers(driverName)
	drv.Init(driverName, dataSourceName, tabPrifix)
	drv.SetDBCache(Cobj)
	drv.SetDBPrifix(dbPrifix)
	return drv
}
