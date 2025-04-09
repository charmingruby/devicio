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

	traceID := instrumentation.Tracer.GetTraceIDFromContext(ctx)

	instrumentation.Logger.Debug("Starting to process routine", "traceId", traceID)

	r, ctx, err := s.parseProcessRoutineData(ctx, msg)
	if err != nil {
		instrumentation.Logger.Error("Failed to parse process routine data", "error", err)

		errorsCounterListMetric, err := instrumentation.ErrorsCounterListMetric()
		if err != nil {
			instrumentation.Logger.Error("Failed to get errors counter list metric", "error", err)
		}

		errorsCounterListMetric.WithLabelValues("parse_error").Inc()
		return err
	}

	instrumentation.Logger.Debug("Processing routine", "traceId", traceID, "routineId", r.ID)

	ctx, err = s.externalAPI.VolatileCall(ctx)
	if err != nil {
		instrumentation.Logger.Error("Failed to call external API", "error", err)

		errorsCounterListMetric, err := instrumentation.ErrorsCounterListMetric()
		if err != nil {
			instrumentation.Logger.Error("Failed to get errors counter list metric", "error", err)
		}

		errorsCounterListMetric.WithLabelValues("api_error").Inc()
		return err
	}

	instrumentation.Logger.Debug("External API call completed", "traceId", traceID)

	if _, err := s.repo.Store(ctx, r); err != nil {
		instrumentation.Logger.Error("Failed to store routine", "error", err)

		errorsCounterListMetric, err := instrumentation.ErrorsCounterListMetric()
		if err != nil {
			instrumentation.Logger.Error("Failed to get errors counter list metric", "error", err)
		}

		errorsCounterListMetric.WithLabelValues("store_error").Inc()
		return err
	}

	instrumentation.Logger.Debug("Stored routine", "traceId", traceID, "routineId", r.ID)

	// messageProcessedCounterMetric, err := instrumentation.MessagesProcessedCounterMetric()
	// if err != nil {
	// 	instrumentation.Logger.Error("Failed to get messages processed counter metric", "error", err)
	// }

	// messageProcessedCounterMetric.WithLabelValues(r.ID).Inc()

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
