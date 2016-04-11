package database

import (
	"strings"
	"strconv"
	"database/sql"
	_ "github.com/lib/pq"
)

type Postgresql struct{
	Driver
}

func (this *Postgresql)QueryNoConn(s string,args ...interface{}) ([]map[string]string ,error){
	return this.Driver.QueryNoConn(this.tranSql(s, args...), args...)
}

func (this *Postgresql) ExecNoConn(s string,args ...interface{})(sql.Result,error){
	return this.Driver.ExecNoConn(this.tranSql(s, args...), args...)
}

func (this *Postgresql)tranSql(s string,args ...interface{})string {
	aolog.InfoTag(this,args)
	for i:=0 ; i<len(args); i++ {
		s = strings.Replace(s, "?", "$"+strconv.Itoa(i+1), 1)
	}
	return s
}