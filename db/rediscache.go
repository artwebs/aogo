package db

import (
	"github.com/artwebs/aogo/log"
	"github.com/artwebs/aogo/utils"
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

func (this *RedisCache) Publish(ch, val string) error {
	var err error
	err = this.Conn()
	if err != nil {
		return err
	}
	this.client.Publish(ch, []byte(val))
	return nil
}

func (this *RedisCache) Receive(sub, unsub, psub, punsub string, f func(m redis.Message)) {
	var err error
	err = this.Conn()
	if err != nil {
		println("连接错误")
	}
	subscribe := make(chan string, 1)
	unsubscribe := make(chan string, 0)
	psubscribe := make(chan string, 0)
	punsubscribe := make(chan string, 0)
	messages := make(chan redis.Message, 0)
	go this.client.Subscribe(subscribe, unsubscribe, psubscribe, punsubscribe, messages)

	if sub != "" {
		subscribe <- sub
	}
	if unsub != "" {
		unsubscribe <- unsub
	}
	if psub != "" {
		psubscribe <- psub
	}
	if punsub != "" {
		punsubscribe <- punsub
	}

	forever := make(chan bool)

	go func() {
		for d := range messages {
			f(d)
		}
	}()
	utils.FailOnError(err, "Waiting for messages. To exit press CTRL+C")
	<-forever
	close(subscribe)
	close(unsubscribe)
	close(psubscribe)
	close(punsubscribe)
	close(messages)
}

func (this *RedisCache) Close() error {
	var err error
	this.client = nil
	return err
}
