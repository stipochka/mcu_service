package producer

import (
	"encoding/json"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/mcu_service/internal/broker"
	"github.com/mcu_service/internal/models"
)

type Producer struct {
	intervalToAsk time.Duration
	*broker.Broker
	devices []models.Record
}

func New(intervalToAsk time.Duration, broker *broker.Broker, devices []models.Record) *Producer {
	return &Producer{
		intervalToAsk: intervalToAsk,
		Broker:        broker,
		devices:       devices,
	}
}

func (p *Producer) Produce(log *slog.Logger) {
	const op = "producer.Produce"
	defer p.Broker.Producer.Close()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for e := range p.Broker.Producer.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					log.Error(
						"failed to send kafka message",
						slog.String("op", op),
						slog.String("error", ev.TopicPartition.Error.Error()),
					)
				}

			}

		}

	}()

	ticker := time.NewTicker(p.intervalToAsk)

	hbTicker := time.NewTicker(p.intervalToAsk / 2)
	defer hbTicker.Stop()
	defer ticker.Stop()

	_, err := p.getData()
	if err != nil {
		log.Error("failed to get data", slog.String("op", op), slog.String("err", err.Error()))
	}

	for {
		select {
		case <-ticker.C:
			devicesInfo, err := p.getData()
			if err != nil {
				log.Error("failed to get data", slog.String("op", op), slog.String("error", err.Error()))
			}

			for _, val := range devicesInfo {
				p.Broker.WriteData(val, log)
			}

		//case <-hbTicker.C:
		//TODO: CHECK MCU
		//err := p.CheckDevices()

		case <-sigChan:
			log.Info("Gracefully stopping service")
			time.Sleep(10 * time.Second)
			return
		}

	}

}

func (p *Producer) getData() ([]string, error) {
	data := []string{}

	for _, val := range p.devices {
		recordBytes, err := json.Marshal(val)
		if err != nil {
			return []string{}, err
		}

		data = append(data, string(recordBytes))
	}

	return data, nil
}
