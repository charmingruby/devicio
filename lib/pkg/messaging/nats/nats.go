package nats

import (
	"context"
	"fmt"

	"github.com/charmingruby/devicio/lib/pkg/logger"
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

type Client struct {
	conn   *nats.Conn
	js     nats.JetStreamContext
	cfg    *Config
	logger *logger.Logger
}

func New(logger *logger.Logger, cfg *Config) (*Client, error) {
	conn, err := nats.Connect(cfg.URL, cfg.Options...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to NATS: %w", err)
	}

	js, err := conn.JetStream()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to create JetStream context: %w", err)
	}

	return &Client{
		conn:   conn,
		js:     js,
		logger: logger,
		cfg:    cfg,
	}, nil
}

func (c *Client) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}

func (c *Client) Publish(ctx context.Context, subject string, msg proto.Message) error {
	c.logger.Info(fmt.Sprintf("publishing message to subject %s", subject))

	data, err := proto.Marshal(msg)
	if err != nil {
		c.logger.Error(fmt.Sprintf("failed to marshal protobuf message: %v", err))
		return fmt.Errorf("failed to marshal protobuf message: %w", err)
	}

	_, err = c.js.Publish(subject, data)
	if err != nil {
		c.logger.Error(fmt.Sprintf("failed to publish message: %v", err))
		return fmt.Errorf("failed to publish message: %w", err)
	}

	c.logger.Info(fmt.Sprintf("message published to subject %s", subject))

	return nil
}

func (c *Client) Subscribe(subject string, msgType proto.Message, handler func(proto.Message) error) (*nats.Subscription, error) {
	sub, err := c.js.QueueSubscribe(subject, c.cfg.QueueName, func(msg *nats.Msg) {
		protoMsg := proto.Clone(msgType)

		if err := proto.Unmarshal(msg.Data, protoMsg); err != nil {
			msg.Nak()
			c.logger.Error(fmt.Sprintf("failed to unmarshal protobuf message: %v", err))
			return
		}

		if err := handler(protoMsg); err != nil {
			msg.Nak()
			c.logger.Error(fmt.Sprintf("failed to handle message: %v", err))
			return
		}

		msg.Ack()
	})

	if err != nil {
		return nil, fmt.Errorf("failed to subscribe: %w", err)
	}

	return sub, nil
}
