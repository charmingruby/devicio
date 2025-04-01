package rabbitmq

import (
	"context"
	"fmt"

	"github.com/charmingruby/devicio/lib/pkg/logger"
	"github.com/streadway/amqp"
	"google.golang.org/protobuf/proto"
)

type Client struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	cfg     *Config
	logger  *logger.Logger
}

type Config struct {
	URL       string
	QueueName string
}

func New(logger *logger.Logger, cfg *Config) (*Client, error) {
	conn, err := amqp.Dial(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to create channel: %w", err)
	}

	_, err = ch.QueueDeclare(cfg.QueueName, true, false, false, false, nil)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	return &Client{
		conn:    conn,
		channel: ch,
		logger:  logger,
		cfg:     cfg,
	}, nil
}

func (c *Client) Close() {
	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
}

func (c *Client) Publish(ctx context.Context, msg proto.Message) error {
	data, err := proto.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal protobuf message: %w", err)
	}

	err = c.channel.Publish("", c.cfg.QueueName, false, false, amqp.Publishing{
		ContentType: "application/protobuf",
		Body:        data,
	})
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}

func (c *Client) Subscribe(handler func(proto.Message) error, msgType proto.Message) error {
	msgs, err := c.channel.Consume(c.cfg.QueueName, "", true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to subscribe: %w", err)
	}

	go func() {
		for d := range msgs {
			protoMsg := proto.Clone(msgType)
			if err := proto.Unmarshal(d.Body, protoMsg); err != nil {
				c.logger.Error(fmt.Sprintf("failed to unmarshal message: %v", err))
				continue
			}

			if err := handler(protoMsg); err != nil {
				c.logger.Error(fmt.Sprintf("failed to handle message: %v", err))
			}
		}
	}()

	return nil
}
