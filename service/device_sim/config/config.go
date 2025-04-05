package config

import (
	"github.com/charmingruby/devicio/lib/config"
)

type CustomConfig struct {
	RabbitMQURL       string `env:"RABBITMQ_URL"`
	RabbitMQQueueName string `env:"RABBITMQ_QUEUE_NAME"`
}

func New() (config.Config[CustomConfig], bool, error) {
	return config.New[CustomConfig]()
}
