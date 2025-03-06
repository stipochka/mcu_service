package producer

import (
	"log/slog"
	"time"

	"github.com/mcu_service/internal/broker"
	"github.com/mcu_service/internal/models"
)

type KafkaProducer interface {
	Produce(log *slog.Logger)
}

type Producer struct {
	KafkaProducer
}

func NewProducer(intervalToAsk time.Duration, broker *broker.KafkaBroker, devices []models.Record) Producer {
	return Producer{
		KafkaProducer: NewKafkaProducer(intervalToAsk, broker, devices),
	}
}
