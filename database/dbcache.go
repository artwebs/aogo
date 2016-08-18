package database

import (
	"github.com/artwebs/aogo/log"
	"github.com/garyburd/redigo/redis"
)

type DBCache struct {
	Cstr string
}

func (this *DBCache) Conn() (redis.Conn, error) {
	return redis.Dial("tcp", this.Cstr)
}
func (this *DBCache) Close(conn redis.Conn) error {
	var err error
	if conn != nil {
		err = conn.Close()
		conn = nil
	}
	return err
}

func (this *DBCache) Set(key, value string) error {

	conn, err := this.Conn()
	defer this.Close(conn)
	if err != nil {
		return err
	}
	if _, err = conn.Do("SET", key, value, "EX", "600", "NX"); err != nil {
		return err
	}
	return err
}

func (this *DBCache) SetList(name, key, value string) error {
	conn, err := this.Conn()
	defer this.Close(conn)
	if err != nil {
		return err
	}
	val, err := redis.Int64(conn.Do("EXISTS", key))
	if val == 0 {
		_, err = conn.Do("RPUSH", name, key)
		if err != nil {
			return err
		}
		if _, err = conn.Do("SET", key, value, "EX", "600", "NX"); err != nil {
			return err
		}
	}
	return err
}

func (this *DBCache) IsExist(key string) bool {
	conn, err := this.Conn()
	defer this.Close(conn)
	flag := false
	if err != nil {
		return flag
	}
	n, _ := redis.Int(conn.Do("EXISTS", key))
	log.InfoTag(this, "IsExist", n, key)
	if n == 1 {
		return true
	}
	return flag
}

func (this *DBCache) Get(key string) (string, error) {
	conn, err := this.Conn()
	defer this.Close(conn)
	if err != nil {
		return "", err
	}
	value, err := conn.Do("GET", key)
	if err != nil {
		return "", err
	}
	return redis.String(value, err)
}

func (this *DBCache) Delete(key string) {
	conn, err := this.Conn()
	if err != nil {
		return
	}
	defer this.Close(conn)
	conn.Send("DEL", key)
	conn.Do("EXEC")
}

func (this *DBCache) DeleteList(name string) {
	conn, err := this.Conn()
	if err != nil {
		return
	}
	defer this.Close(conn)
	values, _ := redis.Values(conn.Do("lrange", name, "0", "200"))
	for _, v := range values {
		log.InfoTag(this, "Del", string(v.([]byte)))
		conn.Send("DEL", string(v.([]byte)))
	}
	conn.Send("DEL", name)
	conn.Do("EXEC")
}
