package database

import (
	"database/sql"
	"strconv"
	"strings"

	aolog "github.com/artwebs/aogo/log"

	_ "github.com/lib/pq"
)

func init() {
	Register("postgres", &Postgresql{})
}

type Postgresql struct {
	Driver
}

func (this *Postgresql) QueryRowNoConn(s string, args ...interface{}) (map[string]string, error) {
	defer this.Reset()
	this.limit = "1"
	var result map[string]string
	s = this.addLimit(s)
	rows, err := this.QueryNoConn(s, args...)
	aolog.InfoTag(this, rows)
	if err != nil {
		aolog.InfoTag(this, err)
		return result, err
	}
	if len(rows) > 0 {
		result = rows[0]
	} else {
		result = map[string]string{}
	}
	return result, nil

}

func (this *Postgresql) QueryNoConn(s string, args ...interface{}) ([]map[string]string, error) {
	return this.Driver.QueryNoConn(this.toSql(s, args...), args...)
}

func (this *Postgresql) ExecNoConn(s string, args ...interface{}) (sql.Result, error) {
	return this.Driver.ExecNoConn(this.toSql(s, args...), args...)
}

func (this *Postgresql) toSql(s string, args ...interface{}) string {
	for i := 0; i < len(args); i++ {
		s = strings.Replace(s, "?", "$"+strconv.Itoa(i+1), 1)
	}
	return s
}
