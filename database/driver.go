package database

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"log"
	"strings"
	// "reflect"
	"strconv"

	"github.com/artwebs/aogo/cache"
	aolog "github.com/artwebs/aogo/log"
	"github.com/artwebs/aogo/utils"
)

var drivers = make(map[string]DriverInterface)
var timeOutDuration = 10 * 60

func Register(name string, d DriverInterface) {
	drivers[name] = d
}

func Drivers(name string) DriverInterface {
	if drv, ok := drivers[name]; ok {
		return drv
	}
	// switch name {
	// case "mysql":
	// 	return &Mysql{}
	// case "postgres":
	// 	return &Postgresql{}
	// default:
	//
	// }
	return nil
}

type Driver struct {
	dbCache        DBCache
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
	cacheKey                    string
}

type DriverInterface interface {
	Init(DriverName, DataSourceName, TabPrifix string)
	SetDBPrifix(p string)
	SetDBCache(c DBCache)
	Conn()
	Close()
	SetTabName(name string)
	Db() *sql.DB
	Query(s string, args ...interface{}) ([]map[string]string, error)
	QueryNoConn(conn func(), s string, args ...interface{}) ([]map[string]string, error)
	QueryRow(s string, args ...interface{}) (map[string]string, error)
	QueryRowNoConn(conn func(), s string, args ...interface{}) (map[string]string, error)
	Exec(sql string, args ...interface{}) (sql.Result, error)
	ExecNoConn(conn func(), sql string, args ...interface{}) (sql.Result, error)
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
}

func (this *Driver) SetTabName(name string) {
	this.TabName = name
}

func (this *Driver) SetCache(c *cache.Cache) {
	this.CacheObj = c
}

func (this *Driver) SetDBCache(c DBCache) {
	this.dbCache = c
}

func (this *Driver) SetDBPrifix(p string) {
	this.DBPrifix = p
}

func (this *Driver) Db() *sql.DB {
	return this.db
}

func (this *Driver) Init(DriverName, DataSourceName, TabPrifix string) {
	this.DriverName = DriverName
	this.DataSourceName = DataSourceName
	this.TabPrifix = TabPrifix
}

func (this *Driver) Conn() {
	var err error
	if this.db == nil {
		this.db, err = sql.Open(this.DriverName, this.DataSourceName)
		if err != nil {
			log.Fatalln("Database open fail!")
		}
	}

}

func (this *Driver) Close() {
	// if this.db != nil {
	// 	this.db.Close()
	// }

	// if this.dbCache != nil {
	// 	this.dbCache.Close()
	// }
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
	return this.QueryNoConn(this.Conn, s, args...)
}

func (this *Driver) QueryNoConn(conn func(), s string, args ...interface{}) ([]map[string]string, error) {
	defer this.Reset()
	this.cacheKey = this.getCacheName(s, args...)
	result := []map[string]string{}
	if this.dbCache != nil && this.dbCache.IsExist(this.cacheKey) {
		val, err := this.dbCache.GetCache(this.cacheKey)
		if err == nil {
			// aolog.InfoTag(this, " get =>"+this.cacheKey)
			rbyte, err := base64.StdEncoding.DecodeString(val)
			if err == nil {
				err := json.Unmarshal(rbyte, &result)
				// aolog.InfoTag(this, " get =>", result)
				if err == nil {
					return result, nil
				} else {
					aolog.InfoTag(this, err, val)
				}

			}
		}
	}
	aolog.Info(s, args)
	conn()
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
	if this.dbCache != nil {
		// aolog.InfoTag(this, " save "+this.cacheKey)
		rbyte, err := json.Marshal(result)
		if err == nil {
			// aolog.InfoTag(this, " save =>", result)
			this.dbCache.AddCache(strings.ToLower(this.TabPrifix+this.TabName), this.cacheKey, base64.StdEncoding.EncodeToString(rbyte))
			// aolog.InfoTag(this, this.dbCache.AddCache(this.TabName, this.cacheKey, base64.StdEncoding.EncodeToString(rbyte)))
		}
	}
	return result, nil
}

func (this *Driver) QueryRow(s string, args ...interface{}) (map[string]string, error) {
	this.Conn()
	return this.QueryRowNoConn(this.Conn, s, args...)
}

func (this *Driver) QueryRowNoConn(conn func(), s string, args ...interface{}) (map[string]string, error) {
	defer this.Reset()
	this.limit = "0,1"
	var result map[string]string
	s = this.addLimit(s)
	rows, err := this.QueryNoConn(conn, s, args...)
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
	return this.ExecNoConn(this.Conn, sql, args...)
}

func (this *Driver) ExecNoConn(conn func(), sql string, args ...interface{}) (sql.Result, error) {
	defer this.Reset()
	aolog.InfoTag(this, sql, args, this.dbCache)
	if this.dbCache != nil {
		this.dbCache.DelCache(strings.ToLower(this.TabPrifix + this.TabName))
	}
	conn()
	stmt, err := this.db.Prepare(sql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	return stmt.Exec(args...)
}

func (this *Driver) Insert(d DriverInterface, values map[string]interface{}) (int64, error) {
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
	result, err := d.ExecNoConn(d.Conn, sql, val...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (this *Driver) Update(d DriverInterface, values map[string]interface{}) (int64, error) {
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
	result, err := d.ExecNoConn(d.Conn, sql, val...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (this *Driver) Delete(d DriverInterface) (int64, error) {
	val := []interface{}{}
	sql := "delete from " + this.getTabName()
	sql, val = this.addWhere(sql, val)
	result, err := d.ExecNoConn(d.Conn, sql, val...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (this *Driver) Find(d DriverInterface) (map[string]string, error) {
	var args []interface{}
	sql := this.initSelect()
	sql, args = this.addWhere(sql, []interface{}{})
	return d.QueryRowNoConn(d.Conn, sql, args...)

}

func (this *Driver) Total(d DriverInterface) (int, error) {
	var args []interface{}
	this.Field("count(*) as c")
	sql := this.initSelect()
	sql, args = this.addWhere(sql, []interface{}{})
	row, err := d.QueryRowNoConn(d.Conn, sql, args...)
	if err != nil {
		return 0, nil
	}
	return strconv.Atoi(string(row["c"]))
}

func (this *Driver) Select(d DriverInterface) ([]map[string]string, error) {

	var args []interface{}
	sql := this.initSelect()
	sql, args = this.addWhere(sql, []interface{}{})
	sql = this.addOrder(sql)
	sql = this.addLimit(sql)
	sql = this.addGroup(sql)
	sql = this.addHaving(sql)
	return d.QueryNoConn(d.Conn, sql, args...)

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
	return base64.StdEncoding.EncodeToString([]byte(this.DBPrifix + " DataBase " + s + string(jbyte)))
	// return this.DBPrifix + " DataBase " + s + string(jbyte)
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

func (this *Driver) DeleteCache() {
	this.ClearCache(this.cacheKey)
}
