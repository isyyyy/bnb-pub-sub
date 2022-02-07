package mq

import (
	"github.com/isyyyy/bnb-pub-sub/internal/config"
	"github.com/streadway/amqp"
)

type MQConnection struct {
	Channel      *amqp.Channel
	ExchangeName string
	QueueName    string
}

func NewRabbitMQConnection(config config.RabbitMQConfig) (*MQConnection, error) {

	conn, err := amqp.Dial(config.URL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	err = ch.ExchangeDeclare(config.ExchangeName, config.ExchangeKind, true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	queue, err := ch.QueueDeclare(config.QueueName, false, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	sub := &MQConnection{
		Channel:      ch,
		ExchangeName: config.ExchangeName,
		QueueName:    queue.Name,
	}
	return sub, nil
}

