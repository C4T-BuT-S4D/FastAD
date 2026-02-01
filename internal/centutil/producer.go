package centutil

import (
	"context"

	"google.golang.org/protobuf/proto"
)

type Producer interface {
	PublishProto(ctx context.Context, msg proto.Message) error
}
