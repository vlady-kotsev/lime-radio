package events

import (
	"time"

	"github.com/ThreeDotsLabs/watermill-nats/v2/pkg/nats"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/google/uuid"
	nc "github.com/nats-io/nats.go"
	sharedconfig "github.com/vlady-kotsev/lime-radio/shared/config"
	sharedevents "github.com/vlady-kotsev/lime-radio/shared/events"
	"go.uber.org/zap"
)

const RequestTopicName string = "request_topic"

type EventPublisher struct {
	p            *nats.Publisher
	requestTopic string
}

var _ EventPublisherer = (*EventPublisher)(nil)

func NewPublisher(logger *zap.Logger, config sharedconfig.EventConfiger) (*EventPublisher, error) {
	// TODO move other values to config
	watermillLogger := sharedevents.NewZapLoggerAdapter(logger)

	marshaler := &nats.GobMarshaler{}

	options := []nc.Option{
		nc.RetryOnFailedConnect(true),
		nc.Timeout(30 * time.Second),
		nc.ReconnectWait(1 * time.Second),
		nc.UserInfo(config.GetEventUsername(), config.GetEventPassword()),
	}
	subscribeOptions := []nc.SubOpt{
		nc.DeliverAll(),
		nc.AckExplicit(),
	}

	jsConfig := nats.JetStreamConfig{
		Disabled:         false,
		AutoProvision:    true,
		ConnectOptions:   nil,
		SubscribeOptions: subscribeOptions,
		PublishOptions:   nil,
		TrackMsgId:       false,
		AckAsync:         false,
		DurablePrefix:    "",
	}

	p, err := nats.NewPublisher(
		nats.PublisherConfig{
			URL:         config.GetBrokerUrl(),
			NatsOptions: options,
			Marshaler:   marshaler,
			JetStream:   jsConfig,
		},
		watermillLogger,
	)
	if err != nil {
		return nil, err
	}

	return &EventPublisher{
		p:            p,
		requestTopic: RequestTopicName,
	}, nil
}

func (ep *EventPublisher) PublishMessage(payload string) error {
	msg := message.NewMessage(uuid.New().String(), []byte(payload))
	return ep.p.Publish(ep.requestTopic, msg)
}
