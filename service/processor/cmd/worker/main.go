package main

import (
	"context"
	"fmt"
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
	cfg, exists, err := config.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if !exists {
		if err := instrumentation.NewLogger(""); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		instrumentation.Logger.Warn("no config found, using default values")
	}

	if err := instrumentation.NewLogger(cfg.Base.LogLevel); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := instrumentation.NewTracer(cfg.ServiceName); err != nil {
		instrumentation.Logger.Error(err.Error())
		os.Exit(1)
	}

	queue, err := rabbitmq.New(instrumentation.Logger, instrumentation.Tracer, &rabbitmq.Config{
		URL:       cfg.Custom.RabbitMQURL,
		QueueName: cfg.Custom.RabbitMQQueueName,
	})
	if err != nil {
		instrumentation.Logger.Error(err.Error())
		os.Exit(1)
	}

	db, err := database.NewPostgres(database.PostgresConnectionInput{
		User:         cfg.Custom.DatabaseUser,
		Password:     cfg.Custom.DatabasePassword,
		Host:         cfg.Custom.DatabaseHost,
		DatabaseName: cfg.Custom.DatabaseName,
		SSL:          cfg.Custom.DatabaseSSL,
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
