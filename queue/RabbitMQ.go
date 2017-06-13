package queue

import (
	"strconv"
	"strings"
	"time"

	"github.com/artwebs/amqp"
	"github.com/artwebs/aogo/log"
	"github.com/artwebs/aogo/utils"
)

type RabbitMQ struct {
	connstr   string
	exchange  string
	Routerkey string
	qtype     string //direct fanout topic
	conn      *amqp.Connection
	ch        *amqp.Channel
}

func InitQueue(connstr, exchange, routerkey string) *RabbitMQ {
	var q = &RabbitMQ{connstr: connstr, exchange: exchange, Routerkey: routerkey}
	return q
}

func (this *RabbitMQ) SetType(t string) *RabbitMQ {
	this.qtype = t
	return this
}

func (this *RabbitMQ) Open() {
	defer func() {
		if err := recover(); err != nil {
			this.Close()
			log.Error("Open", err)
		}
	}()
	if this.qtype == "" {
		this.qtype = "direct"
	}
	if this.conn == nil || this.ch == nil {
		var err error
		this.conn, err = amqp.Dial(this.connstr)
		utils.FailOnError(err, "Failed to connect to RabbitMQ")
		this.ch, err = this.conn.Channel()
		utils.FailOnError(err, "Failed to open Channel")
	}

}

func (this *RabbitMQ) Close() {
	if this.conn != nil {
		this.conn.Close()
		this.conn = nil
	}
	if this.ch != nil {
		this.ch.Close()
		this.ch = nil
	}
}

func (this *RabbitMQ) Send(key, tp, msg string) {
	this.Publish(this.Routerkey+"."+key, amqp.Publishing{
		ContentType: "text/plain",
		Type:        tp,
		Body:        []byte(msg),
	})

}

func (this *RabbitMQ) Publish(key string, obj amqp.Publishing) {
	defer func() {
		if err := recover(); err != nil {
			this.Close()
			log.Error("Send", err)
		}
	}()
	this.Open()
	if err := this.ch.ExchangeDeclare(
		this.exchange, // name
		this.qtype,    // type
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // noWait
		nil,           // arguments
	); err != nil {
		utils.FailOnError(err, "Exchange Declare")
	}
	if err := this.ch.Publish(this.exchange, key, false, false, obj); err != nil {
		utils.FailOnError(err, "send ok")
	}
}

func (this *RabbitMQ) Revice(f func(d amqp.Delivery)) {
	defer func() {
		if err := recover(); err != nil {
			this.Close()
			log.Error("Revice", err)
		}
	}()
	this.Open()
	if err := this.ch.ExchangeDeclare(
		this.exchange, // name
		this.qtype,    // type
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // noWait
		nil,           // arguments
	); err != nil {
		utils.FailOnError(err, "Exchange Declare")
	}

	queue := utils.Identity() + strconv.FormatInt(time.Now().UnixNano(), 10)
	q, err := this.ch.QueueDeclare(
		queue,
		true,
		true,
		false,
		false,
		nil,
	)
	utils.FailOnError(err, "Failed to declare a queue")

	msgs, err := this.ch.Consume(q.Name, "", true, false, false, false, nil)
	utils.FailOnError(err, "Failed to register a consumer")
	err = this.ch.QueueBind(q.Name, this.Routerkey, this.exchange, false, nil)
	if err != nil {
		log.Error(this, err)
		return
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			f(d)
		}
	}()
	utils.FailOnError(err, "Waiting for messages. To exit press CTRL+C")
	<-forever
}

func (this *RabbitMQ) Routerkeys() []string {
	arr := strings.Split(this.Routerkey, ".")
	rs := make([]string, len(arr))
	temp := ""
	for i, val := range arr {
		if temp != "" {
			temp += "."
		}
		temp += val
		rs[i] = temp
	}
	return rs
}
