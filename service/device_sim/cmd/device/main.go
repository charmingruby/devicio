package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"sync"

	"github.com/charmingruby/devicio/lib/messaging/rabbitmq"
	"github.com/charmingruby/devicio/lib/observability"
	"github.com/charmingruby/devicio/service/device_sim/config"
	"github.com/charmingruby/devicio/service/device_sim/internal/device"
	"github.com/charmingruby/devicio/service/device_sim/pkg/instrumentation"
)

func main() {
	instrumentation.NewLogger("")

	instrumentation.Logger.Info("Application started with default log level", "level", observability.LOG_LEVEL_DEFAULT)

	recordsAmount := flag.Int("records", 10, "Amount of records to dispatch")
	concurrency := flag.Int("concurrency", 5, "Amount of workers")
	flag.Parse()

	instrumentation.Logger.Info("Starting device simulator with configuration",
		"records", *recordsAmount,
		"concurrency", *concurrency,
	)

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

	svc := device.NewService(queue)

	instrumentation.Logger.Info("Starting worker pool execution")

	if err := runWorkerPool(svc, *recordsAmount, *concurrency); err != nil {
		instrumentation.Logger.Error("Worker pool execution failed", "error", err)
		os.Exit(1)
	}

	instrumentation.Logger.Info("Worker pool execution completed successfully")

	// Keep the main thread alive to see the traces in the UI
	select {}
}

func runWorkerPool(svc *device.Service, recordsAmount, concurrency int) error {
	ctx, complete := instrumentation.Tracer.Span(context.Background(), "main.runWorkerPool")
	defer complete()

	instrumentation.Logger.Debug("Initializing worker pool",
		"total_workers", concurrency,
		"total_records", recordsAmount,
	)

	var wg sync.WaitGroup
	jobs := make(chan int, recordsAmount)
	results := make(chan workerResult, recordsAmount)

	for i := range concurrency {
		wg.Add(1)
		instrumentation.Logger.Debug("Starting worker", "worker_id", i)
		go worker(ctx, &wg, i, svc, jobs, results)
	}

	go func() {
		instrumentation.Logger.Info("Starting job distribution")
		for i := 1; i <= recordsAmount; i++ {
			select {
			case jobs <- i:
				instrumentation.Logger.Debug("Job dispatched", "job_id", i)
			case <-ctx.Done():
				instrumentation.Logger.Warn("Context cancelled during job distribution")
				close(jobs)
				return
			}
		}
		close(jobs)
		instrumentation.Logger.Info("Job distribution completed")
	}()

	go func() {
		wg.Wait()
		close(results)
		instrumentation.Logger.Info("All workers completed their tasks")
	}()

	var errorCount int
	var successCount int
	for result := range results {
		if result.err != nil {
			errorCount++
			instrumentation.Logger.Warn("Worker encountered error",
				"worker_id", result.workerID,
				"record_id", result.recordID,
				"error", result.err)
		} else {
			successCount++
			instrumentation.Logger.Debug("Worker completed job successfully",
				"worker_id", result.workerID,
				"record_id", result.recordID)
		}
	}

	instrumentation.Logger.Info("Worker pool execution summary",
		"total_jobs", recordsAmount,
		"successful_jobs", successCount,
		"failed_jobs", errorCount,
	)

	if errorCount > 0 {
		return fmt.Errorf("encountered %d errors during processing", errorCount)
	}

	return nil
}

type workerResult struct {
	workerID int
	recordID int
	err      error
}

func worker(ctx context.Context, wg *sync.WaitGroup, workerID int, svc *device.Service, jobs <-chan int, results chan<- workerResult) {
	ctx, complete := instrumentation.Tracer.Span(ctx, "main.worker")
	defer complete()

	defer wg.Done()

	instrumentation.Logger.Debug("Worker started", "worker_id", workerID)

	for recordID := range jobs {
		select {
		case <-ctx.Done():
			instrumentation.Logger.Warn("Worker received cancellation signal", "worker_id", workerID)
			return
		default:
			instrumentation.Logger.Debug("Worker processing job",
				"worker_id", workerID,
				"record_id", recordID,
			)

			err := svc.DispatchRoutineMessage(ctx, device.Device{ID: fmt.Sprintf("device-%d", recordID)})
			results <- workerResult{
				workerID: workerID,
				recordID: recordID,
				err:      err,
			}
		}
	}

	instrumentation.Logger.Debug("Worker completed all jobs", "worker_id", workerID)
}
