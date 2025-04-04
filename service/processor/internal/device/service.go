package device

import (
	"context"
	"fmt"
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

	instrumentation.Logger.Info(fmt.Sprintf("started processing routine with traceId=%s", traceID))

	r, ctx, err := s.parseProcessRoutineData(ctx, msg)
	if err != nil {
		return err
	}

	instrumentation.Logger.Info(fmt.Sprintf("parsed routine with id=%s,traceId=%s", r.ID, traceID))

	ctx, err = s.externalAPI.VolatileCall(ctx)
	if err != nil {
		return err
	}

	if _, err := s.repo.Store(ctx, r); err != nil {
		return err
	}

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
