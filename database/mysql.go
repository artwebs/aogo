package database

import (
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	Register("mysql", &Mysql{})
}

type Mysql struct {
	Driver
}
