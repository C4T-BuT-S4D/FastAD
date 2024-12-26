package centclient

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/centrifugal/centrifuge-go"
	"go.uber.org/zap"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type Producer struct {
	client       *centrifuge.Client
	channel      string
	disconnected chan struct{}

	logger *zap.Logger
}

func NewProducer(client *centrifuge.Client, channel string) *Producer {
	return &Producer{
		client:       client,
		channel:      channel,
		disconnected: make(chan struct{}),
		logger: zap.L().With(
			zap.String("component", "centrifuge-client"),
			zap.String("channel", channel),
		),
	}
}

func (c *Producer) Init() error {
	c.client.OnConnected(func(event centrifuge.ConnectedEvent) {
		c.logger.Debug("connected", zap.String("client_id", event.ClientID))
	})

	c.client.OnError(func(event centrifuge.ErrorEvent) {
		c.logger.Error("error", zap.Error(event.Error))
	})

	c.client.OnDisconnected(func(event centrifuge.DisconnectedEvent) {
		c.logger.Debug(
			"disconnected",
			zap.Uint32("code", event.Code),
			zap.String("reason", event.Reason),
		)
		close(c.disconnected)
	})

	if err := c.client.Connect(); err != nil {
		return fmt.Errorf("connecting to centrifuge: %w", err)
	}

	return nil
}

func (c *Producer) Run(ctx context.Context) error {
	select {
	case <-ctx.Done():
		c.client.Close()
		c.logger.Debug("client closed, waiting for disconnect")
		<-c.disconnected
		return nil
	case <-c.disconnected:
		c.logger.Error("disconnected from centrifuge")
		return errors.New("premature disconnect")
	}
}

func (c *Producer) Publish(ctx context.Context, msg any) error {
	raw, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshaling message: %w", err)
	}
	if _, err := c.client.Publish(ctx, c.channel, raw); err != nil {
		return fmt.Errorf("publishing message: %w", err)
	}
	return nil
}

func (c *Producer) PublishProto(ctx context.Context, msg proto.Message) error {
	raw, err := protojson.MarshalOptions{EmitUnpopulated: true}.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshaling message: %w", err)
	}
	if _, err := c.client.Publish(ctx, c.channel, raw); err != nil {
		return fmt.Errorf("publishing message: %w", err)
	}
	return nil
}
