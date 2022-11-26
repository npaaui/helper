package rabbitmq

import (
	"fmt"

	"github.com/streadway/amqp"

	"helper/common/logger"
)

type MqProducer struct {
	Uri          string `json:"uri"`           // amqp的地址
	ExchangeName string `json:"exchange_name"` // exchange的名称
	QueueName    string `json:"queue_name"`    // queue的名称
}

var mqProducer *MqProducer

func (m *MqProducer) InitMqProducerConf() {
	mqProducer = m
}

func NewMqProducer() *MqProducer {
	return mqProducer
}

// 发布
func (m *MqProducer) Publish(body string) error {
	//建立连接
	logger.Instance.Debugf("dialing %q", m.Uri)
	connection, err := amqp.Dial(m.Uri)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}
	defer func() {
		_ = connection.Close()
	}()

	//创建一个Channel
	logger.Instance.Debug("got Connection, getting Channel")
	channel, err := connection.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %v", err)
	}
	defer func() {
		_ = channel.Close()
	}()

	//声明exchange
	if err = channel.ExchangeDeclare(
		m.ExchangeName, //name
		"direct",       //exchangeType
		true,           //durable
		false,          //auto-deleted
		false,          //internal
		false,          //noWait
		nil,            //arguments
	); err != nil {
		return fmt.Errorf("exchange Declare: %s", err.Error())
	}

	// Reliable publisher confirms require confirm.select support from the
	// connection.
	logger.Instance.Debug("enabling publishing confirms.")
	if err = channel.Confirm(false); err != nil {
		return fmt.Errorf("channel could not be put into confirm mode: %s", err.Error())
	}

	confirms := channel.NotifyPublish(make(chan amqp.Confirmation, 1))

	defer confirmOne(confirms)

	// Producer只能发送到exchange，它是不能直接发送到queue的。
	// 现在我们使用默认的exchange（名字是空字符）。这个默认的exchange允许我们发送给指定的queue。
	// routing_key就是指定的queue名字。
	if err = channel.Publish(
		m.ExchangeName, // exchange
		m.QueueName,    // routing key
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			Headers:         amqp.Table{},
			ContentType:     "text/plain",
			ContentEncoding: "",
			Body:            []byte(body),
			DeliveryMode:    amqp.Transient, // 1=non-persistent, 2=persistent
			Priority:        0,              // 0-9
		},
	); err != nil {
		return fmt.Errorf("failed to publish a message %s", err.Error())
	}
	logger.Instance.Debugf("published %dB OK", len(body))

	return nil
}

// One would typically keep a channel of publishing, a sequence number, and a
// set of unacknowledged sequence numbers and loop until the publishing channel
// is closed.
func confirmOne(confirms <-chan amqp.Confirmation) {
	logger.Instance.Debug("waiting for confirmation of one publishing")

	if confirmed := <-confirms; confirmed.Ack {
		logger.Instance.Debugf("confirmed delivery with delivery tag: %d", confirmed.DeliveryTag)
	} else {
		logger.Instance.Debugf("failed delivery of delivery tag: %d", confirmed.DeliveryTag)
	}
}
