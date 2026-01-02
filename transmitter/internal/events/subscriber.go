package events

import (
	"context"
	"time"

	"github.com/ThreeDotsLabs/watermill-nats/v2/pkg/nats"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/google/uuid"
	nc "github.com/nats-io/nats.go"
	sharedconfig "github.com/vlady-kotsev/lime-radio/shared/config"
	sharedevents "github.com/vlady-kotsev/lime-radio/shared/events"
	"github.com/vlady-kotsev/lime-radio/transmitter/internal/service/radio"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

const RequestTopicName string = "request_topic"

type EventSubscriber struct {
	logger     *zap.Logger
	pl         radio.PlaylistServicer
	subscriber *nats.Subscriber
	ctx        context.Context
	cancel     context.CancelFunc
}

var _ EventSubscriberer = (*EventSubscriber)(nil)

func NewSubscriber(lc fx.Lifecycle, logger *zap.Logger, config sharedconfig.EventConfiger, pl radio.PlaylistServicer) (*EventSubscriber, error) {
	watermillLogger := sharedevents.NewZapLoggerAdapter(logger)

	marshaler := &nats.GobMarshaler{}

	options := []nc.Option{
		nc.RetryOnFailedConnect(true),
		nc.Timeout(30 * time.Second),
		nc.ReconnectWait(1 * time.Second),
		nc.UserInfo(config.GetEventUsername(), config.GetEventPassword()),
	}
	subscribeOptions := []nc.SubOpt{
		nc.DeliverNew(),
		nc.AckExplicit(),
		nc.Durable("transmitter-consumer"),
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

	subscriber, err := nats.NewSubscriber(
		nats.SubscriberConfig{
			URL:         config.GetBrokerUrl(),
			NatsOptions: options,
			Unmarshaler: marshaler,
			JetStream:   jsConfig,
		},
		watermillLogger,
	)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())

	es := &EventSubscriber{
		logger:     logger,
		pl:         pl,
		subscriber: subscriber,
		ctx:        ctx,
		cancel:     cancel,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return es.Start()
		},
		OnStop: func(ctx context.Context) error {
			return es.Close()
		},
	})

	return es, nil
}

func (es *EventSubscriber) Start() error {
	messages, err := es.subscriber.Subscribe(es.ctx, RequestTopicName)
	if err != nil {
		return err
	}

	go es.handleMessages(messages)
	return nil
}

func (es *EventSubscriber) handleMessages(messages <-chan *message.Message) {
	for {
		select {
		case msg := <-messages:
			if msg == nil {
				return
			}
			es.logger.Info("Received message", zap.String("payload", string(msg.Payload)), zap.String("uuid", msg.UUID))
			msg.Ack()
			songID, err := uuid.ParseBytes(msg.Payload)
			if err != nil {
				es.logger.Error("Invalid message payload", zap.String("payload", string(msg.Payload)))
			}
			err = es.pl.EnqueueSong(songID)
			if err != nil {
				es.logger.Error("Song request enqueue failed", zap.String("songID", songID.String()), zap.Error(err))
			}
		case <-es.ctx.Done():
			return
		}
	}
}

func (es *EventSubscriber) Close() error {
	es.cancel()
	return es.subscriber.Close()
}
