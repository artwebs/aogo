package db

import (
	"log"

	aolog "github.com/artwebs/aogo/log"
	"github.com/astaxie/beego/config"
)

type Table struct {
	prifix, name, limit, order, group, having string
	pageSize, pageCount                       int
	field, where, join                        *DBParamer
	db                                        DBInterface
}

func TableNew(tb string, args ...string) (*Table, error) {
	var err error
	var dbPrifix, driverName, dataSourceName, tabPrifix string
	for i, v := range args {
		switch i {
		case 0:
			tabPrifix = v
			break
		case 1:
			dbPrifix = v
			break
		}
	}
	cobj, err := config.NewConfig("ini", "app.conf")
	if err != nil {
		log.Fatalln("no app.conf")
		return nil, err
	}

	driverName = cobj.String(dbPrifix + "DataBase::driverName")
	dataSourceName = cobj.String(dbPrifix + "DataBase::dataSourceName")
	tabPrifix = cobj.String(dbPrifix + "DataBase::tabPrifix")

	if err != nil {
		aolog.Info("db::Table", "dataSourceName", dataSourceName)
	}
	db := Selector(driverName, dataSourceName)
	table := &Table{db: db, name: tb}
	table.SetPrefix(tabPrifix)
	return table, nil
}

func (this *Table) Select() ([]map[string]string, error) {
	sql, param := this.ToSelectasSql()
	println(sql)
	println(param)
	return this.db.Query(sql, param...)
}

func (this *Table) ToSelectasAlias(alias string) (string, []interface{}) {
	sql, param := this.ToSelectasSql()
	return "(" + sql + ") " + alias, param
}

func (this *Table) ToSelectasSql() (string, []interface{}) {
	param := []interface{}{}
	sql := "select "
	if this.field != nil {
		sql = sql + this.field.GetFormat() + " from "
		param = append(param, this.field.GetArgs()...)
	} else {
		sql = sql + " * from "
	}
	sql = sql + this.prifix + this.name
	if this.where != nil {
		sql = sql + " where " + this.where.GetFormat()
		param = append(param, this.where.GetArgs()...)
	}
	if this.order != "" {
		sql = sql + " order by " + this.order
	}

	if this.limit != "" {
		sql = sql + " limit " + this.limit
	}
	if this.group != "" {
		sql = sql + " group by " + this.group
	}

	if this.having != "" {
		sql = sql + " having " + this.having
	}
	return sql, param
}

func (this *Table) Update(values map[string]interface{}) (int64, error) {
	return 0, nil
}

func (this *Table) SetPrefix(prefix string) *Table {
	this.prifix = prefix
	return this
}

func (this *Table) SetName(name string) *Table {
	this.name = name
	return this
}

func (this *Table) Field(fild string) *Table {
	if this.field == nil {
		this.field = DBParamerNew(fild)
	} else {
		this.field.AppendwithSplit(",", fild)
	}

	return this
}

func (this *Table) Where(format string, args ...interface{}) *Table {
	return this.WhereWithSplit("and", format, args)
}

func (this *Table) WhereOr(format string, args ...interface{}) *Table {
	return this.WhereWithSplit("or", format, args)
}

func (this *Table) WhereWithSplit(split, format string, args ...interface{}) *Table {
	if this.where == nil {
		this.where = DBParamerNew(format, args)
	} else {
		this.where.AppendwithSplit(split, format, args)
	}
	return this
}

func (this *Table) Join(format string, args ...interface{}) *Table {
	if this.join == nil {
		this.join = DBParamerNew(format, args)
	} else {
		this.join.Append(format, args)
	}
	return this
}

func (this *Table) Limit(limit string) *Table {
	this.limit = limit
	return this
}

func (this *Table) Order(order string) *Table {
	this.order = order
	return this
}

func (this *Table) Group(group string) *Table {
	this.group = group
	return this
}

func (this *Table) Having(having string) *Table {
	this.having = having
	return this
}
