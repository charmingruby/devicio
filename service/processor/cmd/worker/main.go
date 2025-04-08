package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/charmingruby/devicio/lib/database"
	"github.com/charmingruby/devicio/lib/messaging/rabbitmq"
	"github.com/charmingruby/devicio/lib/observability"
	"github.com/charmingruby/devicio/service/processor/config"
	"github.com/charmingruby/devicio/service/processor/internal/device"
	"github.com/charmingruby/devicio/service/processor/internal/device/client"
	"github.com/charmingruby/devicio/service/processor/internal/device/postgres"
	"github.com/charmingruby/devicio/service/processor/pkg/instrumentation"
	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	instrumentation.NewLogger("")

	cfg, exists, err := config.New()
	if err != nil {
		instrumentation.Logger.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	if !exists {
		instrumentation.Logger.Warn("No configuration file found, using default values")
	} else {
		instrumentation.Logger.Info("Configuration file found, using custom values")
	}

	if cfg.Base.LogLevel != observability.LOG_LEVEL_INFO {
		instrumentation.Logger.Info("Updating log level configuration", "new_level", cfg.Base.LogLevel)

		instrumentation.NewLogger(cfg.Base.LogLevel)

		instrumentation.Logger.Info("New log level successfully configured", "new_level", cfg.Base.LogLevel)
	}

	instrumentation.Logger.Info("Initializing tracing system")

	if err := instrumentation.NewTracer(cfg.ServiceName); err != nil {
		instrumentation.Logger.Error("Failed to initialize tracer", "error", err)
		os.Exit(1)
	}

	instrumentation.Logger.Info("Tracing system initialized successfully")

	instrumentation.Logger.Info("Establishing RabbitMQ connection")

	queue, err := rabbitmq.New(instrumentation.Logger, instrumentation.Tracer, &rabbitmq.Config{
		URL:       cfg.Custom.RabbitMQURL,
		QueueName: cfg.Custom.RabbitMQQueueName,
	})
	if err != nil {
		instrumentation.Logger.Error("Failed to establish RabbitMQ connection", "error", err)
		os.Exit(1)
	}

	instrumentation.Logger.Info("RabbitMQ connection established successfully")

	instrumentation.Logger.Info("Establishing Postgres connection")

	db, err := database.NewPostgres(database.PostgresConnectionInput{
		User:         cfg.Custom.DatabaseUser,
		Password:     cfg.Custom.DatabasePassword,
		Host:         cfg.Custom.DatabaseHost,
		DatabaseName: cfg.Custom.DatabaseName,
		SSL:          cfg.Custom.DatabaseSSL,
	})
	if err != nil {
		instrumentation.Logger.Error("Failed to establish Postgres connection", "error", err)
		os.Exit(1)
	}

	instrumentation.Logger.Info("Postgres connection established successfully")

	repo, err := postgres.NewRoutineRepository(db, instrumentation.Tracer)
	if err != nil {
		instrumentation.Logger.Error("Failed to create routine repository", "error", err)
		os.Exit(1)
	}

	externalAPI := client.NewUnstableAPI(instrumentation.Tracer)

	svc := device.NewService(queue, repo, externalAPI)

	instrumentation.Logger.Info("Subscribing to RabbitMQ queue", "queue", cfg.Custom.RabbitMQQueueName)

	prometheus.MustRegister(instrumentation.MessagesProcessedCounter)
	prometheus.MustRegister(instrumentation.ProcessingTimeHistogram)
	prometheus.MustRegister(instrumentation.ErrorCounter)
	prometheus.MustRegister(instrumentation.QueueLatencyGauge)

	go func() {
		if err := queue.Subscribe(context.Background(), svc.ProcessRoutine); err != nil {
			instrumentation.Logger.Error("Failed to subscribe to RabbitMQ queue", "error", err)
			os.Exit(1)
		}
	}()

	instrumentation.Logger.Info("Subscribed to RabbitMQ queue successfully")

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		if err := http.ListenAndServe(":2112", nil); err != nil {
			instrumentation.Logger.Error("Failed to start Prometheus metrics server", "error", err)
			os.Exit(1)
		}
	}()

	instrumentation.Logger.Info("Prometheus metrics server started on :2112")

	gracefulShutdown(queue, db)
}

func gracefulShutdown(queue *rabbitmq.Client, db *sqlx.DB) {
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	<-stopChan
	instrumentation.Logger.Info("Shutting down gracefully")

	instrumentation.Logger.Info("Closing RabbitMQ connection")

	queue.Close()

	instrumentation.Logger.Info("RabbitMQ connection closed successfully")

	instrumentation.Logger.Info("Closing Postgres connection")

	if err := db.Close(); err != nil {
		instrumentation.Logger.Error("Failed to close Postgres connection", "error", err)
		os.Exit(1)
	}

	instrumentation.Logger.Info("Postgres connection closed successfully")

	instrumentation.Logger.Info("Closing tracing system")

	if err := instrumentation.Tracer.Close(); err != nil {
		instrumentation.Logger.Error(err.Error())
		os.Exit(1)
	}

	instrumentation.Logger.Info("Tracing system closed successfully")

	instrumentation.Logger.Info("Gracefully shutdown")
}
