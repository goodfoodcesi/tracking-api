package main

import (
	"encoding/json"
	"fmt"
	"github.com/goodfoodcesi/api-utils-go/pkg/message"
	"github.com/goodfoodcesi/api-utils-go/pkg/order"
	amqp "github.com/rabbitmq/amqp091-go"
	"golang.org/x/net/context"
)

func main() {
	pub, err := NewPublisher("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}

	newOrder := order.NewOrder("order-1", "customer-1", "shop-1", "delivery-1")
	orderMessage := message.NewMessage(order.OrderCreated, "tools", newOrder)

	err = pub.PublishMessage("orders", orderMessage)
	if err != nil {
		panic(err)
	}
}

type Publisher struct {
	channel *amqp.Channel
}

func NewPublisher(amqpURL string) (*Publisher, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	return &Publisher{channel: ch}, nil
}

func (p *Publisher) PublishMessage(queueName string, message any) error {
	_, err := p.channel.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	data, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	return p.channel.PublishWithContext(
		context.Background(),
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        data,
		},
	)
}
