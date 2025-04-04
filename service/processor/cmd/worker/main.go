package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/charmingruby/devicio/lib/database"
	"github.com/charmingruby/devicio/lib/messaging/rabbitmq"
	"github.com/charmingruby/devicio/service/processor/config"
	"github.com/charmingruby/devicio/service/processor/internal/device"
	"github.com/charmingruby/devicio/service/processor/internal/device/postgres"
	"github.com/charmingruby/devicio/service/processor/pkg/logger"
	"github.com/charmingruby/devicio/service/processor/pkg/observability"
	"github.com/jmoiron/sqlx"
)

func main() {
	logger.New()

	cfg, err := config.New()
	if err != nil {
		logger.Log.Error(err.Error())
		os.Exit(1)
	}

	if err := observability.NewTracing(cfg.ServiceName); err != nil {
		logger.Log.Error(err.Error())
		os.Exit(1)
	}

	queue, err := rabbitmq.New(logger.Log, observability.Tracing, &rabbitmq.Config{
		URL:       cfg.RabbitMQURL,
		QueueName: cfg.RabbitMQQueueName,
	})
	if err != nil {
		logger.Log.Error(err.Error())
		os.Exit(1)
	}

	db, err := database.NewPostgres(database.PostgresConnectionInput{
		User:         cfg.DatabaseUser,
		Password:     cfg.DatabasePassword,
		Host:         cfg.DatabaseHost,
		DatabaseName: cfg.DatabaseName,
		SSL:          cfg.DatabaseSSL,
	})
	if err != nil {
		logger.Log.Error(err.Error())
		os.Exit(1)
	}

	repo, err := postgres.NewRoutineRepository(db, observability.Tracing)
	if err != nil {
		logger.Log.Error(err.Error())
		os.Exit(1)
	}

	svc := device.NewService(queue, repo)

	go func() {
		if err := queue.Subscribe(context.Background(), svc.ProcessRoutine); err != nil {
			logger.Log.Error(err.Error())
			os.Exit(1)
		}
	}()

	gracefulShutdown(queue, db)
}

func gracefulShutdown(queue *rabbitmq.Client, db *sqlx.DB) {
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	<-stopChan
	logger.Log.Info("shutting down gracefully...")

	queue.Close()

	if err := db.Close(); err != nil {
		logger.Log.Error(err.Error())
		os.Exit(1)
	}

	if err := observability.Tracing.Close(); err != nil {
		logger.Log.Error(err.Error())
		os.Exit(1)
	}

	time.Sleep(2 * time.Second)

	logger.Log.Info("shutdown complete")
}
