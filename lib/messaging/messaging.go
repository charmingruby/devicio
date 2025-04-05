package messaging

import (
	"context"

	"google.golang.org/protobuf/proto"
)

type Queue interface {
	Publish(ctx context.Context, msg proto.Message) (context.Context, error)
	Subscribe(ctx context.Context, handler func(context.Context, []byte) error) error
	Close()
}
