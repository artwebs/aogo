package cache

import (
	"sync"

	"github.com/artwebs/aogo/log"
	"github.com/artwebs/aogo/utils"
	_ "github.com/astaxie/beego/cache/memcache"
	_ "github.com/astaxie/beego/cache/redis"
	_ "github.com/astaxie/beego/cache/ssdb"
	"github.com/hoisie/redis"
	// "log"
)

var redisLockobj *sync.RWMutex

type Cache struct {
	name   string
	client *redis.Client
}

func NewCache(name, config string) (*Cache, error) {
	data, err := utils.MapFromString(config)
	if err != nil {
		return nil, err
	}
	if redisLockobj == nil {
		redisLockobj = new(sync.RWMutex)
	}
	return &Cache{name: name, client: &redis.Client{Addr: data["conn"].(string)}}, nil

}

func (this *Cache) Get(key string) []byte {
	redisLockobj.RLock()
	rs, err := this.client.Get(key)
	redisLockobj.RUnlock()
	if err != nil {
		log.ErrorTag(this, err)
		return nil
	}
	return rs
}

func (this *Cache) Put(key string, val []byte, timeout int64) (bool, error) {
	var err error
	flag := false
	redisLockobj.Lock()
	err = this.client.Set(key, val)
	redisLockobj.Unlock()
	if err != nil {
		return flag, err
	}
	redisLockobj.Lock()
	flag, err = this.client.Expire(key, timeout)
	redisLockobj.Unlock()
	if err != nil {
		return flag, err
	}
	return flag, err
}

func (this *Cache) Set(key string, val []byte) (bool, error) {
	var err error
	flag := false
	redisLockobj.Lock()
	err = this.client.Set(key, val)
	redisLockobj.Unlock()
	if err != nil {
		return flag, err
	}
	return flag, err
}

func (this *Cache) Delete(key string) (bool, error) {
	redisLockobj.Lock()
	rs, err := this.client.Del(key)
	redisLockobj.Unlock()
	return rs, err
}

func (this *Cache) Incr(key string) (int64, error) {
	redisLockobj.Lock()
	rs, err := this.client.Incr(key)
	redisLockobj.Unlock()
	return rs, err
}

func (this *Cache) Decr(key string) (int64, error) {
	redisLockobj.Lock()
	rs, err := this.client.Decr(key)
	redisLockobj.Unlock()
	return rs, err
}

func (this *Cache) IsExist(key string) (bool, error) {
	redisLockobj.RLock()
	flag, err := this.client.Exists(key)
	redisLockobj.Unlock()
	return flag, err
}

func (this *Cache) ClearAll() error {
	redisLockobj.RLock()
	err := this.client.Flush(false)
	redisLockobj.Unlock()
	return err
}

func (this *Cache) GetString(key string) string {
	val := this.Get(key)
	return string(val)
}
