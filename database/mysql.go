package database

import (
	"encoding/base64"
	"encoding/json"
	"strings"
	"time"

	aolog "github.com/artwebs/aogo/log"

	_ "github.com/go-sql-driver/mysql"
)

func init() {
	Register("mysql", &Mysql{})
}

type Mysql struct {
	Driver
}

func (this *Mysql) QueryNoConn(s string, args ...interface{}) ([]map[string]string, error) {
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
	result, _ := this.Driver.QueryNoConn(s, args)
	if this.CacheObj != nil && this.isCache {
		aolog.InfoTag(this, " save "+this.cacheKey)
		rbyte, err := json.Marshal(result)
		if err == nil {
			aolog.InfoTag(this, this.CacheObj.Put(this.cacheKey, base64.StdEncoding.EncodeToString(rbyte), 600*time.Second))
		}
	}
	return result, nil
}

func (this *Mysql) getCacheName(s string, args ...interface{}) string {
	for i := 0; i < len(args); i++ {
		s = strings.Replace(s, "?", args[i].(string), 1)
	}
	return s
}
