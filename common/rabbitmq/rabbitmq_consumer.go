package rabbitmq

import (
	"fmt"

	"github.com/streadway/amqp"

	"helper/common/logger"
)

type MqConsumer struct {
	Uri          string `json:"uri"`
	ExchangeName string `json:"exchange_name"`
	ExchangeType string `json:"exchange_type"`
	QueueName    string `json:"queue_name"`
	Tag          string `json:"tag"`
	Key          string `json:"key"`

	Conn       *amqp.Connection
	Channel    *amqp.Channel
	Deliveries <-chan amqp.Delivery
	Done       chan error
}

func (c *MqConsumer) NewConsumer() (*MqConsumer, error) {
	var err error

	logger.Instance.Debugf("dialing %q", c.Uri)
	c.Conn, err = amqp.Dial(c.Uri)
	if err != nil {
		return nil, fmt.Errorf("dial: %s", err)
	}

	go func() {
		logger.Instance.Debugf("closing: %s", <-c.Conn.NotifyClose(make(chan *amqp.Error)))
	}()

	logger.Instance.Debug("got Connection, getting Channel")
	c.Channel, err = c.Conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("channel: %s", err)
	}

	logger.Instance.Debugf("got Channel, declaring Exchange (%q)", c.ExchangeName)
	if err = c.Channel.ExchangeDeclare(
		c.ExchangeName, // name of the exchange
		c.ExchangeType, // type
		true,           // durable
		false,          // delete when complete
		false,          // internal
		false,          // noWait
		nil,            // arguments
	); err != nil {
		return nil, fmt.Errorf("exchange Declare: %s", err)
	}

	logger.Instance.Debugf("declared Exchange, declaring Queue %q", c.QueueName)
	queue, err := c.Channel.QueueDeclare(
		c.QueueName, // name of the queue
		true,        // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // noWait
		nil,         // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("queue Declare: %s", err)
	}

	logger.Instance.Debugf("declared Queue (%q %d messages, %d consumers), binding to Exchange (key %q)",
		queue.Name, queue.Messages, queue.Consumers, c.Key)

	if err = c.Channel.QueueBind(
		queue.Name,     // name of the queue
		c.Key,          // bindingKey
		c.ExchangeName, // sourceExchange
		false,          // noWait
		nil,            // arguments
	); err != nil {
		return nil, fmt.Errorf("queue Bind: %s", err)
	}

	logger.Instance.Debugf("queue bound to Exchange, starting Consume (consumer tag %q)", c.Tag)
	c.Deliveries, err = c.Channel.Consume(
		queue.Name, // name
		c.Tag,      // consumerTag,
		false,      // noAck
		false,      // exclusive
		false,      // noLocal
		false,      // noWait
		nil,        // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("queue Consume: %s", err)
	}

	return c, nil
}

func (c *MqConsumer) Shutdown() error {
	// will close() the deliveries channel
	if err := c.Channel.Cancel(c.Tag, true); err != nil {
		return fmt.Errorf("consumer cancel failed: %s", err)
	}

	if err := c.Conn.Close(); err != nil {
		return fmt.Errorf("AMQP connection close error: %s", err)
	}

	defer logger.Instance.Debug("AMQP shutdown OK")

	// wait for handle() to exit
	return <-c.Done
}
