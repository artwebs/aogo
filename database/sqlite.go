package database

import (
	_ "github.com/mattn/go-sqlite3"
)

func init() {
	Register("sqlite3", &Sqlite{})
}

type Sqlite struct {
	Driver
}
