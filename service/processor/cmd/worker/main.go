package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/charmingruby/devicio/lib/pkg/messaging/rabbitmq"
	"github.com/charmingruby/devicio/service/processor/config"
	"github.com/charmingruby/devicio/service/processor/internal/device"
	"github.com/charmingruby/devicio/service/processor/internal/device/postgres"
	"github.com/charmingruby/devicio/service/processor/pkg/logger"
	"github.com/charmingruby/devicio/service/processor/pkg/observability"
	"github.com/charmingruby/devicio/service/processor/pkg/pg"
	"github.com/jmoiron/sqlx"
)

func main() {
	logger.New()

	cfg, err := config.New()
	if err != nil {
		logger.Log.Error(err.Error())
		os.Exit(1)
	}

	cleanupTrace, err := observability.NewTracer(cfg.ServiceName)
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

	db, err := pg.New(pg.ConnectionInput{
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

	repo, err := postgres.NewRoutineRepository(db)
	if err != nil {
		logger.Log.Error(err.Error())
		os.Exit(1)
	}

	svc := device.NewService(queue, repo)

	ctx, span := observability.Tracer.Start(context.Background(), "main")
	defer span.End()

	go func() {
		if err := queue.Subscribe(ctx, svc.ProcessRoutine); err != nil {
			logger.Log.Error(err.Error())
			os.Exit(1)
		}
	}()

	gracefulShutdown(queue, db, cleanupTrace)
}

func gracefulShutdown(queue *rabbitmq.Client, db *sqlx.DB, cleanupTrace func() error) {
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	<-stopChan
	logger.Log.Info("shutting down gracefully...")

	queue.Close()

	if err := db.Close(); err != nil {
		logger.Log.Error(err.Error())
		os.Exit(1)
	}

	if err := cleanupTrace(); err != nil {
		logger.Log.Error(err.Error())
		os.Exit(1)
	}

	time.Sleep(2 * time.Second)

	logger.Log.Info("shutdown complete")
}
