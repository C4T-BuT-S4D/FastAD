package centutil

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

var _ Producer = (*ClientProducer)(nil)

type ClientProducer struct {
	client       *centrifuge.Client
	channel      string
	disconnected chan struct{}

	logger *zap.Logger
}

func NewClientProducer(address, channel, name, token string) (*ClientProducer, error) {
	rawData, err := json.Marshal(clientData{Name: name})
	if err != nil {
		return nil, fmt.Errorf("marshaling client data: %w", err)
	}

	client := centrifuge.NewJsonClient(address, centrifuge.Config{
		Token:             token,
		Data:              rawData,
		Name:              name,
		EnableCompression: true,
	})

	p := &ClientProducer{
		client:       client,
		channel:      channel,
		disconnected: make(chan struct{}),
		logger: zap.L().Named("centrifuge_client").With(
			zap.String("channel", channel),
		),
	}
	if err := p.init(); err != nil {
		return nil, fmt.Errorf("initializing ClientProducer: %w", err)
	}

	return p, nil
}

func (p *ClientProducer) Run(ctx context.Context) error {
	select {
	case <-ctx.Done():
		p.client.Close()
		p.logger.Debug("client closed, waiting for disconnect")
		<-p.disconnected
		return nil
	case <-p.disconnected:
		p.logger.Error("disconnected from centrifuge")
		return errors.New("premature disconnect")
	}
}

func (p *ClientProducer) PublishProto(ctx context.Context, msg proto.Message) error {
	raw, err := protojson.MarshalOptions{EmitUnpopulated: true}.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshaling message: %w", err)
	}
	if _, err := p.client.Publish(ctx, p.channel, raw); err != nil {
		return fmt.Errorf("publishing message: %w", err)
	}
	return nil
}

func (p *ClientProducer) init() error {
	p.client.OnError(func(event centrifuge.ErrorEvent) {
		p.logger.Error("error", zap.Error(event.Error))
	})

	p.client.OnDisconnected(func(event centrifuge.DisconnectedEvent) {
		p.logger.Debug(
			"disconnected",
			zap.Uint32("code", event.Code),
			zap.String("reason", event.Reason),
		)
		close(p.disconnected)
	})

	if err := p.client.Connect(); err != nil {
		return fmt.Errorf("connecting to centrifuge: %w", err)
	}

	return nil
}

type clientData struct {
	Name string `json:"name"`
}

func ClientNameFromData(data []byte) (string, error) {
	var d clientData
	if err := json.Unmarshal(data, &d); err != nil {
		return "", fmt.Errorf("unmarshaling client data: %w", err)
	}
	return d.Name, nil
}
