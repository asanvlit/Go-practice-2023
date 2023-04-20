package sub

import (
	"Golang-practice-2023/internal/domain/logger"
	"github.com/nats-io/nats.go"
)

type NatsSubscriber struct {
	Conn   *nats.Conn
	logger logger.Logger
}

func New(natsURL string, logger logger.Logger) (*NatsSubscriber, error) {
	conn, err := nats.Connect(natsURL)
	if err != nil {
		return nil, err
	}

	return &NatsSubscriber{
		Conn:   conn,
		logger: logger,
	}, nil
}

func (s *NatsSubscriber) Subscribe(topic string, callback func(msg *nats.Msg)) (*nats.Subscription, error) {
	return s.Conn.Subscribe(topic, callback)
}

func (s *NatsSubscriber) Unsubscribe(subscription *nats.Subscription) error {
	err := subscription.Unsubscribe()
	if err != nil {
		return err
	}
	return nil
}
