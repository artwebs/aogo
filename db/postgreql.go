package db

import (
	"database/sql"
	"strconv"
	"strings"

	"github.com/artwebs/aogo/logger"

	_ "github.com/lib/pq"
)

func init() {
	Register("postgres", &Postgresql{})
}

type Postgresql struct {
	Driver
}

func (this *Postgresql) QueryRowNoConn(conn func(), s string, args ...interface{}) (map[string]string, error) {
	defer this.Reset()
	this.limit = "1"
	var result map[string]string
	s = this.addLimit(s)
	rows, err := this.QueryNoConn(conn, s, args...)
	logger.InfoTag(this, rows)
	if err != nil {
		logger.InfoTag(this, err)
		return result, err
	}
	if len(rows) > 0 {
		result = rows[0]
	} else {
		result = map[string]string{}
	}
	return result, nil

}

func (this *Postgresql) QueryNoConn(conn func(), s string, args ...interface{}) ([]map[string]string, error) {
	return this.Driver.QueryNoConn(conn, this.toSql(s, args...), args...)
}

func (this *Postgresql) ExecNoConn(conn func(), s string, args ...interface{}) (sql.Result, error) {
	return this.Driver.ExecNoConn(conn, this.toSql(s, args...), args...)
}

func (this *Postgresql) toSql(s string, args ...interface{}) string {
	for i := 0; i < len(args); i++ {
		s = strings.Replace(s, "?", "$"+strconv.Itoa(i+1), 1)
	}
	return s
}

func (this *Postgresql) getCacheName(s string, args ...interface{}) string {
	for i := 0; i < len(args); i++ {
		s = strings.Replace(s, "$"+strconv.Itoa(i+1), args[i].(string), 1)
	}
	return s
}
