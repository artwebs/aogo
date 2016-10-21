package cache

import (
	"github.com/artwebs/aogo/log"
	"github.com/artwebs/aogo/utils"
	_ "github.com/astaxie/beego/cache/memcache"
	_ "github.com/astaxie/beego/cache/redis"
	_ "github.com/astaxie/beego/cache/ssdb"
	"github.com/hoisie/redis"
	// "log"
)

type Cache struct {
	name   string
	client *redis.Client
}

func NewCache(name, config string) (*Cache, error) {
	data, err := utils.MapFromString(config)
	if err != nil {
		return nil, err
	}
	return &Cache{name: name, client: &redis.Client{Addr: data["conn"].(string)}}, nil

}

func (this *Cache) Get(key string) []byte {
	rs, err := this.client.Get(key)
	if err != nil {
		log.ErrorTag(this, err)
		return nil
	}
	return rs
}

func (this *Cache) Put(key string, val []byte, timeout int64) (bool, error) {
	var err error
	flag := false
	err = this.client.Set(key, val)
	if err != nil {
		return flag, err
	}
	flag, err = this.client.Expire(key, timeout)
	if err != nil {
		return flag, err
	}
	return flag, err
}

func (this *Cache) Set(key string, val []byte) (bool, error) {
	var err error
	flag := false
	err = this.client.Set(key, val)
	if err != nil {
		return flag, err
	}
	return flag, err
}

func (this *Cache) Delete(key string) (bool, error) {
	this.client.Del(key)
	return this.client.Del(key)
}

func (this *Cache) Incr(key string) (int64, error) {
	return this.client.Incr(key)
}

func (this *Cache) Decr(key string) (int64, error) {
	return this.client.Decr(key)
}

func (this *Cache) IsExist(key string) (bool, error) {
	this.client.Exists(key)
	return this.client.Exists(key)
}

func (this *Cache) ClearAll() error {

	return nil
}

func (this *Cache) GetString(key string) string {
	val := this.Get(key)
	return string(val)
}
