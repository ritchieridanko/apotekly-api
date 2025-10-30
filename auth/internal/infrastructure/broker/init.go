package broker

import (
	"strings"

	"github.com/ritchieridanko/apotekly-api/auth/configs"
	"github.com/segmentio/kafka-go"
)

type Broker struct {
	producer *kafka.Writer
}

func NewClient(cfg *configs.Broker) *Broker {
	brokers := strings.Split(cfg.Brokers, ",")

	p := kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Balancer:     &kafka.Hash{},
		RequiredAcks: kafka.RequireAll,
		Async:        false,
		BatchTimeout: cfg.Timeout.Batch,
	}

	return &Broker{producer: &p}
}

func (b *Broker) Producer() *kafka.Writer {
	return b.producer
}

func (b *Broker) Close() error {
	return b.producer.Close()
}
