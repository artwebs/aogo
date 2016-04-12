package database

import (
	"database/sql"
	aolog "github.com/artwebs/aogo/log"
	_ "github.com/lib/pq"
	"strconv"
	"strings"
)

type Postgresql struct {
	Driver
}

func (this *Postgresql) QueryNoConn(s string, args ...interface{}) ([]map[string]string, error) {
	return this.Driver.QueryNoConn(this.toSql(s, args...), args...)
}

func (this *Postgresql) ExecNoConn(s string, args ...interface{}) (sql.Result, error) {
	return this.Driver.ExecNoConn(this.toSql(s, args...), args...)
}

func (this *Postgresql) toSql(s string, args ...interface{}) string {
	aolog.InfoTag(this, args)
	for i := 0; i < len(args); i++ {
		s = strings.Replace(s, "?", "$"+strconv.Itoa(i+1), 1)
	}
	return s
}
