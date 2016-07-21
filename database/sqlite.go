package database

import (
	_ "github.com/mattn/go-sqlite3"
)

func init() {
	Register("sqlite", &Sqlite{})
}

type Sqlite struct {
	Driver
}
