package broker

import (
	"log/slog"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type Broker struct {
	Producer kafka.Producer
	Topic    string
}

func New(address string, topic string) (*Broker, error) {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": address,
	})
	if err != nil {
		return nil, err
	}
	return &Broker{
		Producer: *producer,
		Topic:    topic,
	}, err
}

//docker exec -it 442f74811fc9 kafka-console-consumer.sh --bootstrap-server localhost:9092 --topic device-data --from-beginning
//

func (b *Broker) WriteData(data string, log *slog.Logger) error {
	const op = "broker.WriteData"

	go func() {
		for e := range b.Producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					log.Error(
						"failed to send kafka message",
						slog.String("op", op),
						slog.String("message", data),
						slog.String("error", ev.TopicPartition.Error.Error()),
					)
				} else {
					log.Info(
						"sended kafka message",
						slog.String("op", op),
						slog.String("message", data),
					)
				}
			}
		}

	}()

	message := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &b.Topic,
			Partition: kafka.PartitionAny,
		},
		Value: []byte(data),
	}

	err := b.Producer.Produce(message, nil)
	if err != nil {
		log.Error("failed to produce message", slog.String("op", op), slog.String("error", err.Error()))
		return err
	}

	log.Info("message sent to kafka", slog.String("op", op))
	return nil
}
