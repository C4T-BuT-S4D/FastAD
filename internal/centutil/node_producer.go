package centutil

import (
	"context"
	"fmt"

	"github.com/centrifugal/centrifuge"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var _ Producer = (*NodeProducer)(nil)

type NodeProducer struct {
	channel string
	node    *centrifuge.Node
}

func NewNodeProducer(node *centrifuge.Node, channel string) *NodeProducer {
	return &NodeProducer{
		channel: channel,
		node:    node,
	}
}

func (p *NodeProducer) PublishProto(_ context.Context, msg proto.Message) error {
	raw, err := protojson.MarshalOptions{EmitUnpopulated: true}.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshaling message: %w", err)
	}
	if _, err := p.node.Publish(p.channel, raw); err != nil {
		return fmt.Errorf("publishing message: %w", err)
	}
	return nil
}
