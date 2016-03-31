package web

import (
	"database/sql"
)

type Model struct {
	db             *sql.DB
	driverName     string
	dataSourceName string
}

func (this *Model) Conn() {
	this.db = sql.Open(this.driverName, this.dataSourceName)
}

func (this *Model) Close() {
	this.db.Close()
}

func (this *Model) Query() {

}

func (this *Model) Insert() {

}

func (this *Model) Update() {

}

func (this *Model) Delete() {

}

func (this *Model) Find() {

}

func (this *Model) Total() int {
	return 0
}

func (this *Model) Select() {

}

func (this *Model) Where() *Model {
	return this
}

func (this *Model) Order() *Model {
	return this
}

func (this *Model) Limit() *Model {
	return this
}

func (this *Model) Group() *Model {
	return this
}

func (this *Model) Having() *Model {
	return this
}

func (this *Model) Field() *Model {
	return this
}
