package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/goodfoodcesi/api-utils-go/pkg/message"
	"github.com/goodfoodcesi/api-utils-go/pkg/order"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

type QueueConfig struct {
	Name     string
	Handlers map[string]func([]byte) error
}

type MessageHandler struct {
	channel *amqp.Channel
	redis   *redis.Client
	ctx     context.Context
	queues  map[string]QueueConfig
}

func NewMessageHandler(ctx context.Context, amqpURL string, redisClient *redis.Client) (*MessageHandler, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	return &MessageHandler{
		channel: ch,
		redis:   redisClient,
		ctx:     ctx,
		queues:  make(map[string]QueueConfig),
	}, nil
}

func (mh *MessageHandler) RegisterQueue(queueName string) {
	if _, exists := mh.queues[queueName]; !exists {
		mh.queues[queueName] = QueueConfig{
			Name:     queueName,
			Handlers: make(map[string]func([]byte) error),
		}
	}
	fmt.Printf("Queue %s registered\n", queueName)
}

func (mh *MessageHandler) RegisterHandler(queueName, messageType string, handler func([]byte) error) error {
	queue, exists := mh.queues[queueName]
	if !exists {
		return fmt.Errorf("queue %s not registered", queueName)
	}
	queue.Handlers[messageType] = handler
	mh.queues[queueName] = queue
	return nil
}

func (mh *MessageHandler) StartListening() error {
	for _, qConfig := range mh.queues {
		q, err := mh.channel.QueueDeclare(
			qConfig.Name,
			false,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			return fmt.Errorf("failed to declare queue %s: %w", qConfig.Name, err)
		}

		msgs, err := mh.channel.Consume(
			q.Name,
			"",
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			return fmt.Errorf("failed to register consumer for queue %s: %w", qConfig.Name, err)
		}

		go mh.handleMessages(msgs, qConfig.Handlers)
	}

	return nil
}

func (mh *MessageHandler) handleMessages(msgs <-chan amqp.Delivery, handlers map[string]func([]byte) error) {
	for d := range msgs {
		var baseMsg struct {
			Type string `json:"type"`
		}
		if err := json.Unmarshal(d.Body, &baseMsg); err != nil {
			fmt.Printf("Error parsing message type: %v\n", err)
			continue
		}

		handler, exists := handlers[baseMsg.Type]
		if !exists {
			fmt.Printf("No handler for message type: %s\n", baseMsg.Type)
			continue
		}

		if err := handler(d.Body); err != nil {
			fmt.Printf("Error handling message: %v\n", err)
		}
	}
}

func main() {
	ctx := context.Background()
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	handler, err := NewMessageHandler(ctx, "amqp://guest:guest@localhost:5672/", rdb)
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to RabbitMQ")

	handler.RegisterQueue("orders")
	handler.RegisterQueue("payments")

	err = handler.RegisterHandler("orders", order.OrderCreated, func(data []byte) error {
		newMessage, err := message.UnmarshalMessage(data)
		if err != nil {
			return fmt.Errorf("failed to unmarshal message: %w", err)
		}

		var orderPayload order.Order
		if err := newMessage.ExtractPayloadEntity(&orderPayload); err != nil {
			return fmt.Errorf("failed to extract payload: %w", err)
		}

		fmt.Printf("Order received: %+v\n", orderPayload)
		return nil
	})
	if err != nil {
		return
	}

	if err := handler.StartListening(); err != nil {
		panic(err)
	}

	<-make(chan struct{})
}
