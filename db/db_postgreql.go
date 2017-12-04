package db

import (
	"database/sql"
	"strconv"
	"strings"

	_ "github.com/lib/pq"
)

func init() {
	DBRegister("postgres", &DBPostgresql{})
}

type DBPostgresql struct {
	DB
}

func (this *DBPostgresql) QueryNoConn(conn func(), s string, args ...interface{}) ([]map[string]string, error) {
	return this.DB.QueryNoConn(conn, this.toSql(s, args...), args...)
}

func (this *DBPostgresql) ExecNoConn(conn func(), s string, args ...interface{}) (sql.Result, error) {
	return this.DB.ExecNoConn(conn, this.toSql(s, args...), args...)
}

func (this *DBPostgresql) toSql(s string, args ...interface{}) string {
	for i := 0; i < len(args); i++ {
		s = strings.Replace(s, "?", "$"+strconv.Itoa(i+1), 1)
	}
	return s
}

func (this *DBPostgresql) getCacheName(s string, args ...interface{}) string {
	for i := 0; i < len(args); i++ {
		s = strings.Replace(s, "$"+strconv.Itoa(i+1), args[i].(string), 1)
	}
	return s
}
