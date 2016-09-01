package database

import (
	"github.com/artwebs/aogo/log"
	"github.com/hoisie/redis"
)

type RedisCache struct {
	Cstr   string
	client *redis.Client
}

func (this *RedisCache) Conn() error {
	var err error

	if this.client == nil {
		this.client = &redis.Client{Addr: this.Cstr}
	}
	return err
}

func (this *RedisCache) AddCache(table, key, value string) error {
	err := this.Conn()
	if err != nil {
		return err
	}

	if ok, err1 := this.client.Exists(key); !ok {
		this.client.Rpush(table, []byte(key))
		this.client.Set(key, []byte(value))
	} else {
		log.ErrorTag(this, err1)
	}
	return nil
}

func (this *RedisCache) GetCache(key string) (string, error) {
	var value string
	err := this.Conn()
	if err != nil {
		return "", err
	}
	val, err1 := this.client.Get(key)

	if err1 != nil {
		return value, err1
	} else {
		value = string(val)
	}
	return value, err
}

func (this *RedisCache) IsExist(key string) bool {
	err := this.Conn()
	if err != nil {
		return false
	}
	flag, err := this.client.Exists(key)
	if err != nil {
		log.ErrorTag(this, err)
	}
	return flag
}

func (this *RedisCache) DelCache(table string) error {
	var err error
	err = this.Conn()
	if err != nil {
		return err
	}
	vals, _ := this.client.Lrange(table, 0, 500)
	for _, v := range vals {
		this.client.Del(string(v))
	}
	this.client.Del(table)

	return err
}

func (this *RedisCache) Close() error {
	var err error
	this.client = nil
	return err
}
