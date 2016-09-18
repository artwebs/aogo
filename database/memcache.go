package database

import (
	"github.com/artwebs/aogo/log"
	"github.com/bradfitz/gomemcache/memcache"
)

type Memcache struct {
	Cstr  string
	mcObj *memcache.Client
}

func init() {
	RegisterDBCache("memcache", &Memcache{})
}

func (this *Memcache) Conn() error {
	this.Cstr = dbcachecstr
	if this.mcObj == nil {
		this.mcObj = memcache.New(this.Cstr)
	}
	return nil
}

func (this *Memcache) AddCache(table, key, value string) error {
	err := this.Conn()
	if err != nil {
		return err
	}
	_, err1 := this.mcObj.Get(key)
	if err1 != nil {
		log.InfoTag(this, "AddCache", err1)
		// log.InfoTag(this, "AddCache", table, key, value)
		err := this.mcObj.Set(&memcache.Item{Key: table, Value: []byte(key), Expiration: int32(timeOutDuration)})
		if err != nil {
			log.InfoTag(this, "AddCache", err)
		}
		err = this.mcObj.Set(&memcache.Item{Key: key, Value: []byte(value)})
		if err != nil {
			log.InfoTag(this, "AddCache", err)
		}
	}
	// if vtemp == nil {
	// 	log.InfoTag(this, "AddCache", table, key, value)
	// 	this.mcObj.Set(&memcache.Item{Key: table, Value: []byte(key)})
	// 	this.mcObj.Set(&memcache.Item{Key: key, Value: []byte(value)})
	// }
	return err1
}

func (this *Memcache) GetCache(key string) (string, error) {
	err := this.Conn()
	if err != nil {
		return "", err
	}
	it, err := this.mcObj.Get(key)
	if err == nil {
		return string(it.Value), err
	}
	return "", nil
}

func (this *Memcache) IsExist(key string) bool {
	flag := false
	err := this.Conn()
	if err != nil {
		return false
	}
	_, err = this.mcObj.Get(key)
	if err == nil {
		flag = true
	} else {
		log.InfoTag(this, "IsExist", err)
	}
	return flag
}

func (this *Memcache) DelCache(table string) error {
	err := this.Conn()
	if err != nil {
		return err
	}

	this.mcObj.DeleteAll()
	// vtemp, err := this.mcObj.Get(table)
	// if err != nil {
	// 	log.InfoTag(this, "DelCache", err)
	// 	return err
	// }
	// log.InfoTag(this, "DelCache", string(vtemp.Value))
	// this.mcObj.Delete(table)
	// return this.mcObj.Delete(table)
	return nil
}
func (this *Memcache) Close() error {
	this.mcObj = nil
	return nil
}
