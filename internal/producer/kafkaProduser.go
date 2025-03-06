package producer

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/creack/pty"
	"github.com/mcu_service/internal/broker"
	"github.com/mcu_service/internal/models"
	"golang.org/x/term"
)

type kafkaProducer struct {
	intervalToAsk time.Duration
	Broker        *broker.KafkaBroker
	devices       []models.Record
}

func NewKafkaProducer(intervalToAsk time.Duration, broker *broker.KafkaBroker, devices []models.Record) *kafkaProducer {
	return &kafkaProducer{
		intervalToAsk: intervalToAsk,
		Broker:        broker,
		devices:       devices,
	}
}

func (p *kafkaProducer) Produce(log *slog.Logger) {
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

	master, slave, err := pty.Open()
	if err != nil {
		log.Error("failed to open pty", slog.String("op", op), slog.String("error", err.Error()))
	}
	oldState, err := term.MakeRaw(int(slave.Fd()))
	if err != nil {
		log.Error("failed to set raw mode", slog.String("error", err.Error()))
	}
	defer term.Restore(int(slave.Fd()), oldState)

	defer master.Close()
	defer slave.Close()

	log.Info("started mater and slave", slog.String("slavePath", slave.Name()))

	_, err = p.askMcu(master)
	if err != nil {
		log.Error("failed to get data", slog.String("op", op), slog.String("err", err.Error()))
	}

	for {
		select {
		case <-ticker.C:
			devicesInfo, err := p.askMcu(master)
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

// iterate over devices and requesting data from each
func (p *kafkaProducer) askMcu(file *os.File) ([]string, error) {
	const op = "producer.askMcu"

	res := make([]string, 0)

	for _, val := range p.devices {
		data, err := sendRequest(file, val.ID)
		if err != nil {
			return res, fmt.Errorf("%s %v", op, err)
		}

		jsonData, err := parseData(data)
		if err != nil {
			return res, err
		}

		res = append(res, jsonData)
	}

	return res, nil
}

// parse data from device response
func parseData(data string) (string, error) {
	const op = "producer.parseData"

	var record models.Record
	valuesFromDevice := strings.Split(data, ";")

	if len(valuesFromDevice) < 2 {
		return "", fmt.Errorf("%s %s", op, "failed to parse data")
	}
	recID, err := strconv.Atoi(valuesFromDevice[0])
	if err != nil {
		return "", fmt.Errorf("%s %s", op, err.Error())
	}

	record.ID = recID
	record.Data = valuesFromDevice[1]

	marshallResp, err := json.Marshal(record)
	if err != nil {
		return "", fmt.Errorf("%s %v", op, err)
	}

	return string(marshallResp), nil
}

// send prototype of request
func sendRequest(file *os.File, id int) (string, error) {
	const op = "producer.sendRequest"

	req := fmt.Sprintf("0%d03\n", id)
	_, err := file.WriteString(req)
	if err != nil {
		return "", fmt.Errorf("%s %v", op, err)
	}

	reader := bufio.NewReader(file)
	line, _, err := reader.ReadLine()
	if err != nil {
		return "", fmt.Errorf("%s %v", op, err)
	}

	return fmt.Sprintf("%d; %s", id, strings.TrimSpace(string(line))), nil
}
