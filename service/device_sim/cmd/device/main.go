package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"sync"

	"github.com/charmingruby/devicio/lib/messaging/rabbitmq"
	"github.com/charmingruby/devicio/service/device_sim/config"
	"github.com/charmingruby/devicio/service/device_sim/internal/device"
	"github.com/charmingruby/devicio/service/device_sim/pkg/logger"
	"github.com/charmingruby/devicio/service/device_sim/pkg/observability"
)

func main() {
	recordsAmount := flag.Int("records", 10, "Amount of records to dispatch")
	concurrency := flag.Int("concurrency", 5, "Amount of workers")
	flag.Parse()

	logger.New()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.New()
	if err != nil {
		logger.Log.Error(err.Error())
		os.Exit(1)
	}

	if err := observability.NewTracer(cfg.ServiceName); err != nil {
		logger.Log.Error(err.Error())
		os.Exit(1)
	}

	queue, err := rabbitmq.New(logger.Log, observability.Tracer, &rabbitmq.Config{
		URL:       cfg.RabbitMQURL,
		QueueName: cfg.RabbitMQQueueName,
	})
	if err != nil {
		logger.Log.Error(err.Error())
		os.Exit(1)
	}

	svc := device.NewService(queue)

	if err := runWorkerPool(ctx, svc, *recordsAmount, *concurrency); err != nil {
		logger.Log.Error(err.Error())
		os.Exit(1)
	}
}

func runWorkerPool(ctx context.Context, svc *device.Service, recordsAmount, concurrency int) error {
	var wg sync.WaitGroup
	jobs := make(chan int, recordsAmount)
	results := make(chan workerResult, recordsAmount)

	for i := range concurrency {
		wg.Add(1)
		go worker(ctx, &wg, i, svc, jobs, results)
	}

	go func() {
		for i := 1; i <= recordsAmount; i++ {
			select {
			case jobs <- i:
			case <-ctx.Done():
				close(jobs)
				return
			}
		}
		close(jobs)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	var errorCount int
	for result := range results {
		if result.err != nil {
			errorCount++
			logger.Log.Error(result.err.Error())
		}
	}

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
	defer wg.Done()

	for recordID := range jobs {
		select {
		case <-ctx.Done():
			return
		default:
			err := svc.DispatchRoutineMessage(ctx, device.Device{ID: fmt.Sprintf("device-%d", recordID)})
			results <- workerResult{
				workerID: workerID,
				recordID: recordID,
				err:      err,
			}
		}
	}
}
