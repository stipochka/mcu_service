package broker

import (
	"log/slog"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaBroker struct {
	Producer kafka.Producer
	Topic    string
}

func NewKafkaBroker(address string, topic string) (*KafkaBroker, error) {
	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": address,
	})
	if err != nil {
		return nil, err
	}
	return &KafkaBroker{
		Producer: *producer,
		Topic:    topic,
	}, err
}

//docker exec -it 07be40d28c27 kafka-console-consumer.sh --bootstrap-server localhost:9092 --topic device-data --from-beginning
//docker exec -it cb622cc33eb7 kafka-topics.sh --bootstrap-server localhost:9092

func (b *KafkaBroker) WriteData(data string, log *slog.Logger) error {
	const op = "broker.WriteData"

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

	log.Info("message sent to kafka", slog.String("op", op), slog.String("message", data))
	return nil
}
