package config

import (
	"github.com/charmingruby/devicio/lib/config"
)

type CustomConfig struct {
	RabbitMQURL       string `env:"RABBITMQ_URL"`
	RabbitMQQueueName string `env:"RABBITMQ_QUEUE_NAME"`
	DatabaseUser      string `env:"DATABASE_USER,required"`
	DatabasePassword  string `env:"DATABASE_PASSWORD,required"`
	DatabaseHost      string `env:"DATABASE_HOST,required"`
	DatabaseName      string `env:"DATABASE_NAME,required"`
	DatabaseSSL       string `env:"DATABASE_SSL,required"`
}

func New() (config.Config[CustomConfig], bool, error) {
	return config.New[CustomConfig]()
}
