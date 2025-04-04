package rabbitmq

import (
	"context"
	"fmt"

	"github.com/charmingruby/devicio/lib/logger"
	"github.com/charmingruby/devicio/lib/observability"
	"github.com/streadway/amqp"
	"google.golang.org/protobuf/proto"
)

type Client struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	cfg     *Config
	logger  *logger.Logger
	tracer  observability.Tracer
}

type Config struct {
	URL       string
	QueueName string
}

func New(logger *logger.Logger, tracer observability.Tracer, cfg *Config) (*Client, error) {
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
		tracer:  tracer,
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

func (c *Client) Subscribe(ctx context.Context, handler func(context.Context, []byte) error) error {
	msgs, err := c.channel.Consume(
		c.cfg.QueueName,
		"",    // consumer
		false, // autoAck
		false, // exclusive
		false, // noLocal
		false, // noWait
		nil,   // args
	)
	if err != nil {
		return fmt.Errorf("failed to subscribe: %w", err)
	}

	go func() {
		for msg := range msgs {
			ctx, complete := c.tracer.Span(ctx, "rabbitmq.Client.Subscribe.Handler")

			if err := handler(ctx, msg.Body); err != nil {
				c.logger.Error(fmt.Sprintf("failed to handle message: %v", err))

				if err := msg.Nack(false, true); err != nil {
					c.logger.Error(fmt.Sprintf("failed to nack message: %v", err))
				}

				complete()
				continue
			}

			if err := msg.Ack(false); err != nil {
				c.logger.Error(fmt.Sprintf("failed to ack message: %v", err))
			}

			complete()
		}
	}()

	return nil
}
