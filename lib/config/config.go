package config

import (
	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
)

type Base struct {
	LogLevel    string `env:"LOG_LEVEL"`
	ServiceName string `env:"SERVICE_NAME,required"`
}

type Config[T any] struct {
	Base
	Custom T
}

func New[T any]() (Config[T], bool, error) {
	if err := godotenv.Load(); err != nil {
		return Config[T]{}, false, err
	}

	cfg := Config[T]{}
	if err := env.Parse(&cfg); err != nil {
		return Config[T]{}, false, err
	}

	return cfg, true, nil
}
