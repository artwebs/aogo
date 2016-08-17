package database

import "github.com/garyburd/redigo/redis"

type DBCache struct {
	conn redis.Conn
	cstr string
}

func (this *DBCache) Conn() error {
	var err error
	if this.conn == nil {
		this.conn, err = redis.Dial("tcp", this.cstr)
		return err
	}
	return nil
}
func (this *DBCache) Close() error {
	var err error
	if this.conn != nil {
		err = this.conn.Close()
		this.conn = nil
	}
	return err
}

func (this *DBCache) Set(key, value string) error {
	defer this.Close()
	err := this.Conn()
	if err != nil {
		return err
	}
	this.conn.Do("SET", value)
	return err
}

func (this *DBCache) SetList(name, key, value string) error {
	defer this.Close()
	err := this.Conn()
	if err != nil {
		return err
	}
	_, err := this.conn.Do("GET", key)
	if err != nil {
		this.conn.Do("lpush", name, key)
	}
	this.conn.Do("SET", key, value)
	return err
}

func (this *DBCache) Get(key string) (string, error) {
	defer this.Close()
	err := this.Conn()
	if err != nil {
		return "", err
	}
	value, err := this.conn.Do("GET", key)
	if err != nil {
		return "", err
	}
	return value.(string), nil
}

func (this *DBCache) Delete(key string) {
	this.conn.Send("DEL", key)
	this.conn.Do("EXEC")
}

func (this *DBCache) DeleteList(name string) {
	values, _ := redis.Values(this.conn.Do("lrange", name, "0", "200"))
	for _, v := range values {
		this.conn.Send("DEL", string(v.([]byte)))
	}
	this.conn.Do("EXEC")
}
