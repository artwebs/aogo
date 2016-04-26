package database

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"log"
	"strings"
	"time"
	// "reflect"
	"github.com/artwebs/aogo/cache"
	"github.com/artwebs/aogo/utils"
	aolog "github.com/artwebs/aogo/log"
	"strconv"
)

var drivers = make(map[string]DriverInterface)
var timeOutDuration = 10 * time.Second

func Register(name string, d DriverInterface) {
	drivers[name] = d
}

func Drivers(name string) DriverInterface {
	if drv, ok := drivers[name]; ok {
		return drv
	}
	return nil
}

type Driver struct {
	CacheObj       *cache.Cache
	DBPrifix       string
	TabPrifix      string
	TabName        string
	db             *sql.DB
	DriverName     string
	DataSourceName string

	fields                      []string
	where                       string
	whereArgs                   []interface{}
	limit, order, group, having string
	isCache bool
}

type DriverInterface interface {
	Init(DriverName, DataSourceName, TabPrifix string)
	SetDBPrifix(p string) 
	SetCache(c *cache.Cache)
	IsCache(flag bool)
	Conn()
	Close()
	SetTabName(name string)
	Query(s string, args ...interface{}) ([]map[string]string, error)
	QueryNoConn(s string, args ...interface{}) ([]map[string]string, error)
	QueryRow(s string, args ...interface{}) (map[string]string, error)
	QueryRowNoConn(s string, args ...interface{}) (map[string]string, error)
	Exec(sql string, args ...interface{}) (sql.Result, error)
	ExecNoConn(sql string, args ...interface{}) (sql.Result, error)
	Insert(d DriverInterface, values map[string]interface{}) (int64, error)
	Update(d DriverInterface, values map[string]interface{}) (int64, error)
	Delete(d DriverInterface) (int64, error)
	Find(d DriverInterface) (map[string]string, error)
	Total(d DriverInterface) (int, error)
	Select(d DriverInterface) ([]map[string]string, error)
	Where(w string, args ...interface{}) DriverInterface
	Order(o string) DriverInterface
	Limit(l string) DriverInterface
	Group(g string) DriverInterface
	Having(h string) DriverInterface
	Field(fields ...string) DriverInterface
	ClearCache(args ...string)
}

func (this *Driver) SetTabName(name string) {
	this.TabName = name
}

func (this *Driver) SetCache(c *cache.Cache) {
	this.CacheObj = c
}

func (this *Driver) IsCache(flag bool) {
	this.isCache = flag
}

func (this *Driver) SetDBPrifix(p string) {
	this.DBPrifix = p
}

func (this *Driver) Init(DriverName, DataSourceName, TabPrifix string) {
	this.DriverName = DriverName
	this.DataSourceName = DataSourceName
	this.TabPrifix = TabPrifix
	this.isCache = true
}

func (this *Driver) Conn() {
	var err error
	this.db, err = sql.Open(this.DriverName, this.DataSourceName)
	if err != nil {
		log.Fatalln("Database open fail!")
	}
}

func (this *Driver) Close() {
	if this.db == nil {
		return
	}
	this.db.Close()
}

func (this *Driver) getTabName() string {
	return this.TabPrifix + utils.StrUpperUnderline(this.TabName)
}

func (this *Driver) Reset() {
	this.fields = nil
	this.where = ""
	this.limit = ""
	this.having = ""
	this.order = ""
	this.group = ""
	this.isCache = true
}

func (this *Driver) addWhere(sql string, args []interface{}) (string, []interface{}) {
	if this.where != "" && !strings.Contains(strings.ToLower(sql), " where ") {
		sql += " where " + this.where
		args = append(args, this.whereArgs...)
	}
	return sql, args
}

func (this *Driver) addOrder(sql string) string {
	if this.order != "" && !strings.Contains(strings.ToLower(sql), " order by ") {
		sql += " order by " + this.order
	}
	return sql
}

func (this *Driver) addGroup(sql string) string {
	if this.group != "" && !strings.Contains(strings.ToLower(sql), " group by ") {
		sql += " group by " + this.group
	}
	return sql
}

func (this *Driver) addHaving(sql string) string {
	if this.group != "" && !strings.Contains(strings.ToLower(sql), " having ") {
		sql += " having " + this.having
	}
	return sql
}

func (this *Driver) addLimit(sql string) string {
	if this.limit != "" && !strings.Contains(strings.ToLower(sql), " limit ") {
		sql += " limit " + this.limit
	}
	return sql
}

func (this *Driver) initSelect() string {
	sql := "select "
	if this.fields != nil {
		sql += strings.Join(this.fields, ",")
	} else {
		sql += "*"
	}
	sql += " from " + this.getTabName()
	return sql
}

func (this *Driver) Query(s string, args ...interface{}) ([]map[string]string, error) {
	this.Conn()
	defer this.Close()
	return this.QueryNoConn(s, args...)
}

func (this *Driver) QueryNoConn(s string, args ...interface{}) ([]map[string]string, error) {
	defer this.Reset()
	cacheKey := this.getCacheName(s, args...)
	if this.CacheObj != nil && this.CacheObj.IsExist(cacheKey) && this.isCache{
		aolog.InfoTag(this, " get "+cacheKey)
		result := []map[string]string{}
		rbyte, err := base64.StdEncoding.DecodeString(this.CacheObj.GetString(cacheKey))
		if err == nil {
			json.Unmarshal(rbyte, &result)
		}
		return result, nil
	}
	result := []map[string]string{}
	aolog.Info(s, args)
	rows, err := this.db.Query(s, args...)
	if err != nil {
		return result, err
	}
	columns, err := rows.Columns()
	if err != nil {
		return result, err
	}
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			return result, err
		}
		row := map[string]string{}
		for i, col := range values {
			if col == nil {
				row[columns[i]] = "NULL"
			} else {
				row[columns[i]] = string(col)
			}
		}
		result = append(result, row)
	}
	if this.CacheObj != nil && this.isCache{
		aolog.InfoTag(this, " save "+cacheKey)
		rbyte, err := json.Marshal(result)
		if err == nil {
			aolog.InfoTag(this, this.CacheObj.Put(cacheKey, base64.StdEncoding.EncodeToString(rbyte), 600*time.Second))
		}

	}
	return result, nil
}

func (this *Driver) QueryRow(s string, args ...interface{}) (map[string]string, error) {
	this.Conn()
	defer this.Close()
	return this.QueryRowNoConn(s, args...)
}

func (this *Driver) QueryRowNoConn(s string, args ...interface{}) (map[string]string, error) {
	defer this.Reset()
	this.limit = "0,1"
	var result map[string]string
	s = this.addLimit(s)
	rows, err := this.QueryNoConn(s, args...)
	if err != nil {
		return result, err
	}
	if len(rows) > 0 {
		result = rows[0]
	} else {
		result = map[string]string{}
	}
	return result, nil

}

func (this *Driver) Exec(sql string, args ...interface{}) (sql.Result, error) {
	this.Conn()
	defer this.Close()
	return this.ExecNoConn(sql, args...)
}

func (this *Driver) ExecNoConn(sql string, args ...interface{}) (sql.Result, error) {
	defer this.Reset()
	aolog.InfoTag(this, sql, args)
	stmt, err := this.db.Prepare(sql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	return stmt.Exec(args...)
}

func (this *Driver) Insert(d DriverInterface, values map[string]interface{}) (int64, error) {
	d.Conn()
	defer d.Close()
	var fm, vm string
	val := []interface{}{}
	for k, v := range values {
		if fm != "" {
			fm += ","
			vm += ","
		}
		fm += k
		vm += "?"
		val = append(val, v)
	}
	sql := "insert into " + this.getTabName() + " (" + fm + ") VALUES (" + vm + ")"
	result, err := d.ExecNoConn(sql, val...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (this *Driver) Update(d DriverInterface, values map[string]interface{}) (int64, error) {
	d.Conn()
	defer d.Close()
	u := ""
	val := []interface{}{}
	for k, v := range values {
		if u != "" {
			u += ","
		}
		u += k + "=?"
		val = append(val, v)
	}
	sql := "update " + this.getTabName() + " set " + u
	sql, val = this.addWhere(sql, val)
	result, err := d.ExecNoConn(sql, val...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (this *Driver) Delete(d DriverInterface) (int64, error) {
	d.Conn()
	defer d.Close()
	val := []interface{}{}
	sql := "delete from " + this.getTabName()
	sql, val = this.addWhere(sql, val)
	result, err := d.ExecNoConn(sql, val...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (this *Driver) Find(d DriverInterface) (map[string]string, error) {
	d.Conn()
	defer d.Close()
	var args []interface{}
	sql := this.initSelect()
	sql, args = this.addWhere(sql, []interface{}{})
	return d.QueryRowNoConn(sql, args...)

}

func (this *Driver) Total(d DriverInterface) (int, error) {
	d.Conn()
	defer d.Close()
	var args []interface{}
	this.Field("count(*) as c")
	sql := this.initSelect()
	sql, args = this.addWhere(sql, []interface{}{})
	row, err := d.QueryRowNoConn(sql, args...)
	if err != nil {
		return 0, nil
	}
	return strconv.Atoi(string(row["c"]))
}

func (this *Driver) Select(d DriverInterface) ([]map[string]string, error) {
	d.Conn()
	defer d.Close()
	var args []interface{}
	sql := this.initSelect()
	sql, args = this.addWhere(sql, []interface{}{})
	sql = this.addOrder(sql)
	sql = this.addLimit(sql)
	sql = this.addGroup(sql)
	sql = this.addHaving(sql)
	return d.QueryNoConn(sql, args...)

}

func (this *Driver) Where(w string, args ...interface{}) DriverInterface {
	this.where = w
	this.whereArgs = args
	return this
}

func (this *Driver) Order(o string) DriverInterface {
	this.order = o
	return this
}

func (this *Driver) Limit(l string) DriverInterface {
	this.limit = l
	return this
}

func (this *Driver) Group(g string) DriverInterface {
	this.group = g
	return this
}

func (this *Driver) Having(h string) DriverInterface {
	this.having = h
	return this
}

func (this *Driver) Field(fields ...string) DriverInterface {
	this.fields = fields
	return this
}

func (this *Driver) getCacheName(s string, args ...interface{}) string {
	jbyte, _ := json.Marshal(args)
	return base64.StdEncoding.EncodeToString([]byte(this.DBPrifix +" DataBase "+s + string(jbyte)))
}

func (this *Driver) ClearCache(args ...string) {
	if this.CacheObj == nil {
		return
	}
	if len(args) == 0 {
		this.CacheObj.ClearAll()
		return
	}
	for i := range args {
		this.CacheObj.Delete(args[i])
	}
}
