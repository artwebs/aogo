package db

import (
	"database/sql"
	"log"

	"github.com/artwebs/aogo/cache"
)

var dbs = make(map[string]DBInterface)
var dbTimeOutDuration = 10 * 60

func DBRegister(name string, d DBInterface) {
	dbs[name] = d
}

type DB struct {
	dbCache        DBCache
	CacheObj       *cache.Cache
	DBPrifix       string
	sqlDb          *sql.DB
	DriverName     string
	DataSourceName string
}
type DBInterface interface {
	Init(DriverName, DataSourceName string)
	Conn()
	Close()
	Query(s string, args ...interface{}) ([]map[string]string, error)
	QueryNoConn(conn func(), s string, args ...interface{}) ([]map[string]string, error)
	Exec(sql string, args ...interface{}) (sql.Result, error)
	ExecNoConn(conn func(), sql string, args ...interface{}) (sql.Result, error)
}

func Selector(DriverName, DataSourceName string) DBInterface {
	if db, ok := dbs[DriverName]; ok {
		db.Init(DriverName, DataSourceName)
		return db
	}
	return nil
}

func (this *DB) Init(DriverName, DataSourceName string) {
	this.DriverName = DriverName
	this.DataSourceName = DataSourceName
}

func (this *DB) Conn() {
	var err error
	if this.sqlDb == nil {
		this.sqlDb, err = sql.Open(this.DriverName, this.DataSourceName)
		if err != nil {
			log.Fatalln("Database open fail!")
		}
	}
}

func (this *DB) Close() {
	if this.sqlDb != nil {
		this.sqlDb.Close()
		this.sqlDb = nil
	}

	if this.dbCache != nil {
		this.dbCache.Close()
		this.dbCache = nil
	}
}

func (this *DB) Query(s string, args ...interface{}) ([]map[string]string, error) {
	return this.QueryNoConn(this.Conn, s, args...)
}

func (this *DB) QueryNoConn(conn func(), s string, args ...interface{}) ([]map[string]string, error) {
	result := []map[string]string{}
	conn()
	rows, err := this.sqlDb.Query(s, args...)
	if err != nil {
		return result, err
	}
	columns, err := rows.Columns()
	if err != nil {
		return result, err
	}
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			return result, err
		}
		row := map[string]string{}
		for i, col := range values {
			if col == nil {
				row[columns[i]] = "NULL"
			} else {
				row[columns[i]] = string(col)
			}
		}
		result = append(result, row)
	}
	rows.Close()
	return result, nil
}

func (this *DB) Exec(sql string, args ...interface{}) (sql.Result, error) {
	return this.ExecNoConn(this.Conn, sql, args...)
}

func (this *DB) ExecNoConn(conn func(), sql string, args ...interface{}) (sql.Result, error) {
	conn()
	stmt, err := this.sqlDb.Prepare(sql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	return stmt.Exec(args...)
}
