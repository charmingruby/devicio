package device

import (
	"context"
	"time"

	"github.com/charmingruby/devicio/lib/core/id"
	"github.com/charmingruby/devicio/lib/messaging"
	"github.com/charmingruby/devicio/lib/messaging/rabbitmq"
	"github.com/charmingruby/devicio/lib/proto/gen/pb"
	"github.com/charmingruby/devicio/service/processor/internal/device/client"
	"github.com/charmingruby/devicio/service/processor/pkg/instrumentation"
	"google.golang.org/protobuf/proto"
)

type Service struct {
	queue       messaging.Queue
	repo        RoutineRepository
	externalAPI client.UnstableAPI
}

func NewService(queue *rabbitmq.Client, repo RoutineRepository, externalAPI client.UnstableAPI) *Service {
	return &Service{
		queue:       queue,
		repo:        repo,
		externalAPI: externalAPI,
	}
}

func (s *Service) ProcessRoutine(ctx context.Context, msg []byte) error {
	ctx, complete := instrumentation.Tracer.Span(ctx, "service.Service.ProcessRoutine")
	defer complete()

	startTime := time.Now()
	traceID := instrumentation.Tracer.GetTraceIDFromContext(ctx)

	instrumentation.Logger.Debug("Starting to process routine", "traceId", traceID)

	r, ctx, err := s.parseProcessRoutineData(ctx, msg)
	if err != nil {
		instrumentation.ErrorCounter.WithLabelValues("parse_error").Inc()
		return err
	}

	instrumentation.Logger.Debug("Processing routine", "traceId", traceID, "routineId", r.ID)

	ctx, err = s.externalAPI.VolatileCall(ctx)
	if err != nil {
		instrumentation.ErrorCounter.WithLabelValues("api_error").Inc()
		return err
	}

	instrumentation.Logger.Debug("External API call completed", "traceId", traceID)

	if _, err := s.repo.Store(ctx, r); err != nil {
		instrumentation.ErrorCounter.WithLabelValues("store_error").Inc()
		return err
	}

	instrumentation.Logger.Debug("Stored routine", "traceId", traceID, "routineId", r.ID)

	// Record metrics
	instrumentation.MessagesProcessedCounter.Inc()
	processingTime := time.Since(startTime).Seconds()
	instrumentation.ProcessingTimeHistogram.Observe(processingTime)
	instrumentation.QueueLatencyGauge.Set(time.Since(r.DispatchedAt).Seconds())

	return nil
}

func (s *Service) parseProcessRoutineData(ctx context.Context, b []byte) (Routine, context.Context, error) {
	ctx, complete := instrumentation.Tracer.Span(ctx, "service.Service.parseProcessRoutineData")
	defer complete()

	var p pb.DeviceRoutine

	if err := proto.Unmarshal(b, &p); err != nil {
		return Routine{}, ctx, err
	}

	r := Routine{}

	r.ID = id.New()
	r.DeviceID = p.GetId()
	r.Status = p.GetStatus().String()
	r.Context = p.GetContext()
	r.Area = p.GetArea()
	r.Diagnostics = p.GetDiagnostics()
	r.DispatchedAt = p.GetDispatchedAt().AsTime()
	r.CreatedAt = time.Now()

	return r, ctx, nil
}
