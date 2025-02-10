package producer

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mcu_service/internal/broker"
)

type Producer struct {
	intervalToAsk time.Duration
	broker.Broker
	devices []string
}

func New(intervalToAsk time.Duration, broker broker.Broker, devices []string) *Producer {
	return &Producer{
		intervalToAsk: intervalToAsk,
		Broker:        broker,
		devices:       devices,
	}
}

func (p *Producer) Produce(log *slog.Logger) {
	const op = "producer.Produce"

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

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
	return p.devices, nil
}
