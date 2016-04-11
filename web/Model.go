package web

import (
	"database/sql"
	"log"
	"strings"
	"reflect"
	"strconv"
	aolog "github.com/artwebs/aogo/log"
)

func D(model ModelInterface) ModelInterface{
	model.SetTabName(strings.TrimSuffix(reflect.Indirect(reflect.ValueOf(model)).Type().Name(), "Model"))
	return model
}

type ModelInterface interface{
	Conn()
	Close()
	SetTabName(name string)
}

type Model struct {
	DBPrifix       string
	TabPrifix      string
	TabName        string
	db             *sql.DB
	driverName     string
	dataSourceName string

	fields []string
	where string
	whereArgs []interface{}

	limit,order,group,having string
}

func (this *Model)SetTabName(name string) {
	this.TabName =name
}

func (this *Model) Conn() {
	var err error
	conf, err = InitAppConfig()
	if err == nil {
		this.driverName = conf.String(this.DBPrifix+"DataBase::driverName", "")
		this.dataSourceName = conf.String(this.DBPrifix+"DataBase::dataSourceName", "")
		this.TabPrifix = conf.String(this.DBPrifix+"DataBase::tabPrifix", "")
	} else {
		log.Fatalln("AppConfig init fail")
	}
	log.Println("driverName:" + this.driverName)
	log.Println("dataSourceName:" + this.dataSourceName)
	log.Println("TabPrifix:" + this.TabPrifix)
	log.Println("TabName:" + this.TabName)
	this.db, err = sql.Open(this.driverName, this.dataSourceName)
	if err != nil {
		log.Fatalln("Database open fail!")
	}
}

func (this *Model) Close() {
	if this.db == nil{
		return
	}
	this.db.Close()
}

func (this *Model)getTabName() string {
	return this.TabPrifix+strings.ToLower(this.TabName)
}

func (this *Model)Reset() {
	this.fields = nil
	this.where = ""
	this.limit =""
	this.having =""
	this.order =""
	this.group=""
}

func (this *Model)addWhere(sql string,args []interface{}) (string,[]interface{}){
	if this.where !="" && !strings.Contains(strings.ToLower(sql), " where "){
		sql += " where "+ this.where
		args=append(args,this.whereArgs...)
	}
	return sql ,args
}

func (this *Model)addOrder(sql string) string {
	if this.order !="" && !strings.Contains(strings.ToLower(sql), " order by "){
		sql += " order by "+ this.order
	}
	return sql
}

func (this *Model)addGroup(sql string)string {
	if this.group !="" && !strings.Contains(strings.ToLower(sql), " group by "){
		sql += " group by "+this.group
	}
	return sql
}

func (this *Model)addHaving(sql string)string {
	if this.group !="" && !strings.Contains(strings.ToLower(sql), " having "){
		sql +=" having "+ this.having
	}
	return sql
}

func (this *Model)addLimit(sql string)string {
	if this.limit != "" && !strings.Contains(strings.ToLower(sql), " limit "){
		sql += " limit "+ this.limit
	}
	return sql
}

func (this *Model)initSelect()string {
	sql := "select "
	if this.fields !=nil{
		sql += strings.Join(this.fields, ",")
	}else {
		sql += "*"
	}
	sql += " from "+this.getTabName()
	return sql
}

func (this *Model) Query(s string,args ...interface{}) ([]map[string]string ,error) {
	this.Conn()
	defer this.Close()
	return this.QueryNoConn(s, args...)
}

func (this *Model) QueryNoConn(s string,args ...interface{}) ([]map[string]string ,error){
	defer this.Reset()
	result :=[]map[string]string{}
	aolog.Info(s,args)
	rows ,err :=  this.db.Query(s, args...)
	if err !=nil{
		return result,err
	}
	columns, err := rows.Columns()
	if err !=nil{
		return result,err
	}
	values := make([]sql.RawBytes, len(columns))
    scanArgs := make([]interface{}, len(values))
    for i := range values {
        scanArgs[i] = &values[i]
    }
    for rows.Next() {
        err = rows.Scan(scanArgs...)
        if err != nil {
            return result,err
        }
 		row := map[string]string{}
        for i, col := range values {
        	if col ==nil{
        		row[columns[i]] = "NULL"
        	}else{
        		row[columns[i]] = string(col)
        	}
        }
        result = append(result,row)
    }
    return result,nil
}

func (this *Model)QueryRow(s string,args ...interface{})(map[string]string ,error) {
	this.Conn()
	defer this.Close()
	return this.QueryRowNoConn(s, args...)
}

func (this *Model)QueryRowNoConn(s string,args ...interface{})(map[string]string ,error) {
	defer this.Reset()
	this.limit = "0,1"
	var result map[string]string
	s =  this.addLimit(s)
	rows , err := this.QueryNoConn(s, args...)
	if err !=nil{
		return result,err
	}
	if len(rows)>0{
		result = rows[0]
	}else{
		result = map[string]string{}
	}
	return result , nil

}

func (this *Model) Exec(sql string,args ...interface{})(sql.Result,error) {
	this.Conn()
	defer this.Close()
	return this.ExecNoConn(sql, args...)
}

func (this *Model) ExecNoConn(sql string,args ...interface{})(sql.Result,error) {
	defer this.Reset()
	aolog.InfoTag(this,sql,args)
	stmt,err := this.db.Prepare(sql)
	if err !=nil{
		return nil,err
	}
	defer stmt.Close()
	return stmt.Exec(args...)
}

func (this *Model) Insert(values map[string]interface{}) (int64,error) {
	this.Conn()
	defer this.Close()
	var fm ,vm string
	val := []interface{}{}
	for k , v := range values {
		if fm != ""{
			fm += ","
			vm += ","
		}
		fm += k
		vm += "?"
		val = append(val,v)
	}
	sql := "insert into "+this.getTabName()+" ("+fm+") VALUES ("+vm+")"
	result , err := this.ExecNoConn(sql, val...)
	if err !=nil{
		return 0,err
	}
	return result.RowsAffected()
}

func (this *Model) Update(values map[string]interface{})(int64,error) {
	this.Conn()
	defer this.Close()
	u := ""
	val :=[]interface{}{}
	for k,v :=range values{
		if u !="" {
			u +=","
		}
		u += k +"=?"
		val = append(val,v)
	}
	sql := "update "+this.getTabName() +" set " + u
	sql ,val = this.addWhere(sql, val)
	result , err := this.ExecNoConn(sql, val...)
	if err !=nil{
		return 0,err
	}
	return result.RowsAffected()
}

func (this *Model) Delete()(int64,error) {
	this.Conn()
	defer this.Close()
	val :=[]interface{}{}
	sql := "delete from "+this.getTabName()
	sql ,val =this.addWhere(sql, val)
	result , err := this.ExecNoConn(sql, val...)
	if err !=nil{
		return 0,err 
	}
	return result.RowsAffected()
}


func (this *Model) Find()(map[string]string,error) {
	this.Conn()
	defer this.Close()
	var args []interface{}
	sql := this.initSelect()
	sql,args =  this.addWhere(sql,[]interface{}{})
	return this.QueryRowNoConn(sql,args...)

}

func (this *Model) Total() (int,error) {
	this.Conn()
	defer this.Close()
	var args []interface{}
	this.Field("count(*) as c")
	sql := this.initSelect()
	sql,args = this.addWhere(sql, []interface{}{})
	row , err := this.QueryRowNoConn(sql,args...)
	if err != nil{
		return 0,nil
	}
	return strconv.Atoi(string(row["c"]))
}

func (this *Model) Select()([]map[string]string,error) {
	this.Conn()
	defer this.Close()
	var args []interface{}
	sql := this.initSelect()
	sql,args= this.addWhere(sql, []interface{}{})
	sql = this.addOrder(sql)
	sql = this.addLimit(sql)
	sql = this.addGroup(sql)
	sql = this.addHaving(sql)
	return this.QueryNoConn(sql,args...)

}

func (this *Model) Where(w string,args ...interface{}) *Model {
	this.where = w
	this.whereArgs = args
	return this
}

func (this *Model) Order(o string) *Model {
	this.order = o
	return this
}

func (this *Model) Limit(l string) *Model {
	this.limit = l
	return this
}

func (this *Model) Group(g string) *Model {
	this.group = g
	return this
}

func (this *Model) Having(h string) *Model {
	this.having = h
	return this
}

func (this *Model) Field(fields ...string) *Model {
	this.fields = fields
	return this
}
