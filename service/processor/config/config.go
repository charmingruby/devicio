package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/charmingruby/devicio/service/processor/pkg/logger"
	"github.com/joho/godotenv"
)

type Config struct {
	RabbitMQURL       string `env:"RABBITMQ_URL"`
	RabbitMQQueueName string `env:"RABBITMQ_QUEUE_NAME"`
	DatabaseUser      string `env:"DATABASE_USER,required"`
	DatabasePassword  string `env:"DATABASE_PASSWORD,required"`
	DatabaseHost      string `env:"DATABASE_HOST,required"`
	DatabaseName      string `env:"DATABASE_NAME,required"`
	DatabaseSSL       string `env:"DATABASE_SSL,required"`
}

func New() (Config, error) {
	if err := godotenv.Load(); err != nil {
		logger.Log.Warn(".env file not found")
	}

	cfg := Config{}

	if err := env.Parse(&cfg); err != nil {
		return Config{}, err
	}

	return cfg, nil
}
