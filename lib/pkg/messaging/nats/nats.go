package nats

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
)

type Client struct {
	conn *nats.Conn
	js   nats.JetStreamContext
}

func New(cfg *Config) (*Client, error) {
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
		conn: conn,
		js:   js,
	}, nil
}

func (c *Client) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}

func (c *Client) Publish(ctx context.Context, subject string, msg proto.Message) error {
	data, err := proto.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal protobuf message: %w", err)
	}

	_, err = c.js.Publish(subject, data)
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	return nil
}

func (c *Client) Subscribe(subject string, msgType proto.Message, handler func(proto.Message) error) (*nats.Subscription, error) {
	sub, err := c.js.Subscribe(subject, func(msg *nats.Msg) {
		protoMsg := proto.Clone(msgType)

		if err := proto.Unmarshal(msg.Data, protoMsg); err != nil {
			msg.Nak()
			return
		}

		if err := handler(protoMsg); err != nil {
			msg.Nak()
			return
		}

		msg.Ack()
	})

	if err != nil {
		return nil, fmt.Errorf("failed to subscribe: %w", err)
	}

	return sub, nil
}
