package main

import (
	"log/slog"
	"os"

	"github.com/mcu_service/internal/broker"
	"github.com/mcu_service/internal/config"
	"github.com/mcu_service/internal/producer"
)

func main() {
	//TODO create config
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)
	log.Info(
		"started mcu-service",
		slog.String("env", cfg.Env),
	)

	broker, err := broker.New(cfg.BrokerAddress, "device-data")
	if err != nil {
		log.Error("failed to initialize broker", slog.String("error", err.Error()))
		return
	}

	producer := producer.New(cfg.IntervalToAsk, *broker, cfg.Devices)

	producer.Produce(log)
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case "local":
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case "dev":
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	case "prod":
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
