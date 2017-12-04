package db

import (
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func init() {
	DBRegister("mysql", &DBMysql{})
}

type DBMysql struct {
	DB
}

func (this *DBMysql) getCacheName(s string, args ...interface{}) string {
	for i := 0; i < len(args); i++ {
		s = strings.Replace(s, "?", args[i].(string), 1)
	}
	return s
}
