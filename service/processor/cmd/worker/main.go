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
	"github.com/charmingruby/devicio/service/processor/internal/device/client"
	"github.com/charmingruby/devicio/service/processor/internal/device/postgres"
	"github.com/charmingruby/devicio/service/processor/pkg/instrumentation"
	"github.com/jmoiron/sqlx"
)

func main() {
	instrumentation.NewLogger()

	cfg, err := config.New()
	if err != nil {
		instrumentation.Logger.Error(err.Error())
		os.Exit(1)
	}

	if err := instrumentation.NewTracer(cfg.ServiceName); err != nil {
		instrumentation.Logger.Error(err.Error())
		os.Exit(1)
	}

	queue, err := rabbitmq.New(instrumentation.Logger, instrumentation.Tracer, &rabbitmq.Config{
		URL:       cfg.RabbitMQURL,
		QueueName: cfg.RabbitMQQueueName,
	})
	if err != nil {
		instrumentation.Logger.Error(err.Error())
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
		instrumentation.Logger.Error(err.Error())
		os.Exit(1)
	}

	repo, err := postgres.NewRoutineRepository(db, instrumentation.Tracer)
	if err != nil {
		instrumentation.Logger.Error(err.Error())
		os.Exit(1)
	}

	externalAPI := client.NewUnstableAPI(instrumentation.Tracer)

	svc := device.NewService(queue, repo, externalAPI)

	go func() {
		if err := queue.Subscribe(context.Background(), svc.ProcessRoutine); err != nil {
			instrumentation.Logger.Error(err.Error())
			os.Exit(1)
		}
	}()

	gracefulShutdown(queue, db)
}

func gracefulShutdown(queue *rabbitmq.Client, db *sqlx.DB) {
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	<-stopChan
	instrumentation.Logger.Info("shutting down gracefully...")

	queue.Close()

	if err := db.Close(); err != nil {
		instrumentation.Logger.Error(err.Error())
		os.Exit(1)
	}

	if err := instrumentation.Tracer.Close(); err != nil {
		instrumentation.Logger.Error(err.Error())
		os.Exit(1)
	}

	time.Sleep(2 * time.Second)

	instrumentation.Logger.Info("shutdown complete")
}
