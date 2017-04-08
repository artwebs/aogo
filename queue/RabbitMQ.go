package queue

import (
	"strings"

	"github.com/artwebs/amqp"
	"github.com/artwebs/aogo/log"
	"github.com/artwebs/aogo/utils"
)

type RabbitMQ struct {
	connstr   string
	exchange  string
	routerkey string
	qtype     string //direct fanout topic
	conn      *amqp.Connection
	ch        *amqp.Channel
}

func InitQueue(connstr, exchange, routerkey string) *RabbitMQ {
	var q = &RabbitMQ{connstr: connstr, exchange: exchange, routerkey: routerkey}
	return q
}

func (this *RabbitMQ) SetType(t string) *RabbitMQ {
	this.qtype = t
	return this
}

func (this *RabbitMQ) Open() {
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
		defer this.conn.Close()
	}
	if this.ch != nil {
		defer this.ch.Close()
	}
}

func (this *RabbitMQ) Send(key, tp, msg string) {
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
	if err := this.ch.Publish(this.exchange, this.routerkey+"."+key, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Type:        tp,
		Body:        []byte(msg),
	}); err != nil {
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
	queue := utils.Identity()
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
	err = this.ch.QueueBind(q.Name, this.routerkey, this.exchange, false, nil)
	if err != nil {
		log.Error(this, err)
		return
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			if d.RoutingKey == this.routerkey {
				f(d)
			}
		}
	}()
	utils.FailOnError(err, "Waiting for messages. To exit press CTRL+C")
	<-forever
}

func (this *RabbitMQ) Routerkeys() []string {
	arr := strings.Split(this.routerkey, ".")
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
