package database

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"strconv"
	"strings"
	"time"

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
	this.cacheKey = this.getCacheName(s, args...)
	if this.CacheObj != nil && this.CacheObj.IsExist(this.cacheKey) && this.isCache {
		aolog.InfoTag(this, " get "+this.cacheKey)
		result := []map[string]string{}
		rbyte, err := base64.StdEncoding.DecodeString(this.CacheObj.GetString(this.cacheKey))
		if err == nil {
			json.Unmarshal(rbyte, &result)
		}
		return result, nil
	}
	result, _ := this.Driver.QueryNoConn(this.toSql(s, args...), args...)
	if this.CacheObj != nil && this.isCache {
		aolog.InfoTag(this, " save "+this.cacheKey)
		rbyte, err := json.Marshal(result)
		if err == nil {
			aolog.InfoTag(this, this.CacheObj.Put(this.cacheKey, base64.StdEncoding.EncodeToString(rbyte), 600*time.Second))
		}
	}
	return result, nil
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

func (this *Postgresql) getCacheName(s string, args ...interface{}) string {
	for i := 0; i < len(args); i++ {
		s = strings.Replace(s, "$"+strconv.Itoa(i+1), args[i].(string), 1)
	}
	return s
}
