package main

import (
	"os"

	"github.com/charmingruby/devicio/lib/pkg/messaging/rabbitmq"
	"github.com/charmingruby/devicio/service/processor/config"
	"github.com/charmingruby/devicio/service/processor/internal/device"
	"github.com/charmingruby/devicio/service/processor/pkg/logger"
)

func main() {
	logger.New()

	cfg, err := config.New()
	if err != nil {
		logger.Log.Error(err.Error())
		os.Exit(1)
	}

	queue, err := rabbitmq.New(logger.Log, &rabbitmq.Config{
		URL:       cfg.RabbitMQURL,
		QueueName: cfg.RabbitMQQueueName,
	})
	if err != nil {
		logger.Log.Error(err.Error())
		os.Exit(1)
	}

	device.NewService(queue)
}
