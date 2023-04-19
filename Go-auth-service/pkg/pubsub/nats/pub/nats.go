package pub

import (
	"Golang-practice-2023/internal/domain/logger"
	"github.com/nats-io/nats.go"
)

type NatsPublisher struct {
	Conn   *nats.Conn
	logger logger.Logger
}

func New(natsURL string, logger logger.Logger) (*NatsPublisher, error) {
	conn, err := nats.Connect(natsURL)
	if err != nil {
		return nil, err
	}

	return &NatsPublisher{
		Conn:   conn,
		logger: logger,
	}, nil
}

func (p *NatsPublisher) Publish(topic string, data []byte) error {
	err := p.Conn.Publish(topic, data)
	if err != nil {
		return err
	}

	return nil
}
