package db

import (
	"log"
	"strconv"

	aolog "github.com/artwebs/aogo/log"
	"github.com/astaxie/beego/config"
)

type Table struct {
	prifix, name, limit, order, group, having string
	pageSize, pageCount                       int
	field, where, join                        *DBParamer
	db                                        DBInterface
	defaultParam                              []interface{}
}

func TableNewWithParamer(subDB *DBParamer, args ...string) *Table {
	return TableNew(subDB.GetFormat(), args...).SetDefaultParam(subDB.GetArgs()).SetPrefix("")
}

func TableNew(tb string, args ...string) *Table {
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
		return nil
	}

	driverName = cobj.String(dbPrifix + "DataBase::driverName")
	dataSourceName = cobj.String(dbPrifix + "DataBase::dataSourceName")
	tabPrifix = cobj.String(dbPrifix + "DataBase::tabPrifix")

	if err != nil {
		aolog.Info("db::Table", "dataSourceName", dataSourceName)
	}
	db := Selector(driverName, dataSourceName)
	table := &Table{db: db, name: tb, defaultParam: []interface{}{}}
	table.SetPrefix(tabPrifix)
	return table
}

func (this *Table) SetDefaultParam(p []interface{}) *Table {
	this.defaultParam = append(this.defaultParam, p...)
	return this
}

func (this *Table) Select(args ...interface{}) ([]map[string]string, error) {
	dp := this.ToSelectasSql(args...)
	return this.db.Query(dp.GetFormat(), dp.GetArgs()...)
}

func (this *Table) Total(args ...interface{}) int64 {
	dp := DBParamerNew("select ", this.defaultParam...)
	dp.Append(" count(*) as c from ")
	dp.Append(this.prifix + this.name)
	if this.join != nil {
		dp.Append(this.join.GetFormat(), this.join.GetArgs()...)
	}

	//计算where 条件
	this.composeWhere(dp, args...)

	if this.limit != "" {
		dp.Append(" limit " + this.limit)
	}
	if this.group != "" {
		dp.Append(" group by " + this.group)
	}

	if this.having != "" {
		dp.Append(" having " + this.having)
	}
	aolog.Debug(dp.ToString())
	data, _ := this.db.Query(dp.GetFormat(), dp.GetArgs()...)
	i, err := strconv.ParseInt(data[0]["c"], 10, 64)
	if err != nil {
		return 0
	}
	return i
}

func (this *Table) ToSelectasAlias(alias string, args ...interface{}) *DBParamer {
	dp := this.ToSelectasSql(args...)
	return DBParamerNew("("+dp.GetFormat()+") "+alias, dp.GetArgs()...)
}

func (this *Table) ToSelectasSql(args ...interface{}) *DBParamer {
	dp := DBParamerNew("select ", this.defaultParam...)

	if this.field != nil {
		dp.Append(this.field.GetFormat()+" from ", this.field.GetArgs()...)
	} else {
		dp.Append(" * from ")
	}
	dp.Append(this.prifix + this.name)

	if this.join != nil {
		dp.Append(this.join.GetFormat(), this.join.GetArgs()...)
	}

	//计算where 条件
	this.composeWhere(dp, args...)

	if this.order != "" {
		dp.Append(" order by " + this.order)
	}

	if this.limit != "" {
		dp.Append(" limit " + this.limit)
	}
	if this.group != "" {
		dp.Append(" group by " + this.group)
	}

	if this.having != "" {
		dp.Append(" having " + this.having)
	}
	return dp
}

func (this *Table) Insert(values map[string]interface{}) (int64, error) {
	var dp *DBParamer
	var dpv *DBParamer
	for k, v := range values {
		if dp == nil {
			dp = DBParamerNew("insert into "+this.prifix+this.name+" ( "+k, v)
		} else {
			dp.AppendwithSplit(",", k, v)
		}
		if dpv != nil {
			dpv = DBParamerNew("(?")
		} else {
			dpv.AppendwithSplit(",", "?")
		}
	}
	dp.Append(") values ")
	dpv.Append(")")
	result, err := this.db.Exec(dp.GetFormat()+dpv.GetFormat(), dp.GetArgs()...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (this *Table) Update(values map[string]interface{}, args ...interface{}) (int64, error) {
	dp := DBParamerNew("update " + this.prifix + this.name + " set ")
	for k, v := range values {
		if dp.index > 2 {
			dp.AppendwithSplit(",", k+"=? ", v)
		} else {
			dp.AppendwithSplit(" ", k+"=? ", v)
		}
	}
	this.composeWhere(dp, args...)
	result, err := this.db.Exec(dp.GetFormat(), dp.GetArgs()...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (this *Table) Delete(args ...interface{}) (int64, error) {
	dp := DBParamerNew("delete from " + this.prifix + this.name)
	this.composeWhere(dp, args...)
	result, err := this.db.Exec(dp.GetFormat(), dp.GetArgs()...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
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
	return this.WhereWithSplit("and", format, args...)
}

func (this *Table) WhereOr(format string, args ...interface{}) *Table {
	return this.WhereWithSplit("or", format, args...)
}

func (this *Table) WhereWithSplit(split, format string, args ...interface{}) *Table {
	if this.where == nil {
		this.where = DBParamerNew(format, args...)
	} else {
		this.where.AppendwithParenthesis(split, format, args...)
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

func (this *Table) composeWhere(dp *DBParamer, args ...interface{}) {
	var w *DBParamer
	if len(args) > 0 && this.where != nil {
		w = DBParamerNew(this.where.GetFormat(), this.where.GetArgs()...)
		w.AppendwithParenthesis(" and ", args[0].(string), args[1:]...)
	} else if len(args) > 0 {
		w = DBParamerNew(args[0].(string), args[1:]...)
	} else if this.where != nil {
		w = DBParamerNew(this.where.GetFormat(), this.where.GetArgs()...)
	}
	if w != nil {
		dp.Append(" where "+w.GetFormat(), w.GetArgs()...)
	}
}
