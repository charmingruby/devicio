package main

import (
	"os"

	"github.com/charmingruby/devicio/lib/pkg/messaging/nats"
	"github.com/charmingruby/devicio/service/device/config"
	"github.com/charmingruby/devicio/service/device/pkg/logger"
)

func main() {
	logger.New()

	cfg, err := config.New()
	if err != nil {
		logger.Log.Error(err.Error())
		os.Exit(1)
	}

	nats, err := nats.New(
		logger.Log,
		nats.NewConfig(
			nats.WithURL(cfg.NatsURL),
			nats.WithQueueName(cfg.NatsQueueName),
		),
	)
	if err != nil {
		logger.Log.Error(err.Error())
		os.Exit(1)
	}
	defer nats.Close()
}
