package database

import "github.com/bradfitz/gomemcache/memcache"

type Memcache struct {
	Cstr  string
	mcObj *memcache.Client
}

func (this *Memcache) Conn() error {
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
	vtemp, err1 := this.mcObj.Get(key)
	if vtemp == nil {
		this.mcObj.Set(&memcache.Item{Key: table, Value: []byte(key)})
		this.mcObj.Set(&memcache.Item{Key: key, Value: []byte(value)})
	}
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
	vtemp, _ := this.mcObj.Get(key)
	if vtemp != nil {
		flag = true
	}
	return flag
}

func (this *Memcache) DelCache(table string) error {
	err := this.Conn()
	if err != nil {
		return err
	}
	// vtemp, _ := this.mcObj.Get(table)
	// log.InfoTag(this, string(vtemp.Value))
	return this.mcObj.Delete(table)
}
func (this *Memcache) Close() error {
	this.mcObj = nil
	return nil
}
