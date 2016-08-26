package database

import (
	"strconv"

	"github.com/garyburd/redigo/redis"
)

type RedisCache struct {
	Cstr  string
	rdObj redis.Conn
}

func (this *RedisCache) Conn() error {
	var err error
	if this.rdObj == nil {
		this.rdObj, err = redis.Dial("tcp", this.Cstr)
	}
	return err
}

func (this *RedisCache) AddCache(table, key, value string) error {
	err := this.Conn()
	if err != nil {
		return err
	}
	val, err1 := redis.Int64(this.rdObj.Do("EXISTS", key))
	if err1 != nil {
		return err1
	}
	if val == 0 {
		this.rdObj.Do("SET", table, key, "EX", strconv.Itoa(timeOutDuration))
		this.rdObj.Do("SET", key, value, "EX", strconv.Itoa(timeOutDuration))
	}
	return nil
}

func (this *RedisCache) GetCache(key string) (string, error) {
	var value string
	err := this.Conn()
	if err != nil {
		return "", err
	}
	value, err = redis.String(this.rdObj.Do("GET", key))
	if err != nil {
		return value, err
	}
	return value, err
}

func (this *RedisCache) IsExist(key string) bool {
	flag := false
	err := this.Conn()
	if err != nil {
		return flag
	}

	val, _ := redis.Int64(this.rdObj.Do("EXISTS", key))
	if val > 0 {
		flag = true
	}
	return flag
}

func (this *RedisCache) DelCache(table string) error {
	var err error
	err = this.Conn()
	if err != nil {
		return err
	}
	var values []interface{}
	values, err = redis.Values(this.rdObj.Do("GET", table))
	for _, v := range values {
		this.rdObj.Send("DEL", string(v.([]byte)))
	}
	this.rdObj.Send("DEL", table)
	this.rdObj.Do("EXEC")
	return err
}

func (this *RedisCache) Close() error {
	var err error
	if this.rdObj != nil {
		err = this.rdObj.Close()
		this.rdObj = nil
	}
	return err
}
