package cache

import (
	bgcache "github.com/astaxie/beego/cache"
	_ "github.com/astaxie/beego/cache/memcache"
	_ "github.com/astaxie/beego/cache/redis"
	_ "github.com/astaxie/beego/cache/ssdb"
	redigo "github.com/garyburd/redigo/redis"
	"time"
	// "log"
)

type Cache struct {
	name    string
	adapter bgcache.Cache
}

func NewCache(name, config string) (*Cache, error) {
	adp, err := bgcache.NewCache(name, config)
	if err != nil {
		return nil, err
	}
	return &Cache{adapter: adp, name: name}, nil
}

func (this *Cache) Get(key string) interface{} {
	return this.adapter.Get(key)
}

func (this *Cache) GetMulti(keys []string) []interface{} {
	return this.adapter.GetMulti(keys)
}

func (this *Cache) Put(key string, val interface{}, timeout time.Duration) error {
	return this.adapter.Put(key, val, timeout)
}

func (this *Cache) Delete(key string) error {
	return this.adapter.Delete(key)
}

func (this *Cache) Incr(key string) error {
	return this.adapter.Incr(key)
}

func (this *Cache) Decr(key string) error {
	return this.adapter.Decr(key)
}

func (this *Cache) IsExist(key string) bool {
	return this.adapter.IsExist(key)
}

func (this *Cache) ClearAll() error {
	return this.adapter.ClearAll()
}

func (this *Cache) StartAndGC(config string) error {
	return this.adapter.StartAndGC(config)
}

func (this *Cache) GetString(key string) string {
	val := this.Get(key)
	var result string
	switch this.name {
	case "redis":
		result, _ = redigo.String(val, nil)
		break
	default:
		result = val.(string)
		break
	}
	return result
}
