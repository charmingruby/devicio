package device

import (
	"context"
	"time"

	"github.com/charmingruby/devicio/lib/pkg/core/id"
	"github.com/charmingruby/devicio/lib/pkg/messaging/rabbitmq"
	"github.com/charmingruby/devicio/lib/proto/gen/pb"
	"github.com/charmingruby/devicio/service/processor/pkg/observability"
	"google.golang.org/protobuf/proto"
)

type Service struct {
	queue *rabbitmq.Client
	repo  RoutineRepository
}

func NewService(queue *rabbitmq.Client, repo RoutineRepository) *Service {
	return &Service{
		queue: queue,
		repo:  repo,
	}
}

func (s *Service) ProcessRoutine(ctx context.Context, msg []byte) error {
	ctx, span := observability.Tracer.Start(ctx, "service.Service.ProcessRoutine")
	defer span.End()

	ctx, r, err := s.parseProcessRoutineData(ctx, msg)
	if err != nil {
		return err
	}

	if _, err := s.repo.Store(ctx, r); err != nil {
		return err
	}

	return nil
}

func (s *Service) parseProcessRoutineData(ctx context.Context, b []byte) (context.Context, *Routine, error) {
	ctx, span := observability.Tracer.Start(ctx, "service.Service.parseProcessRoutineData")
	defer span.End()

	var p pb.DeviceRoutine

	if err := proto.Unmarshal(b, &p); err != nil {
		return ctx, nil, err
	}

	r := &Routine{}

	r.ID = id.New()
	r.DeviceID = p.GetId()
	r.Status = p.GetStatus().String()
	r.Context = p.GetContext()
	r.Area = p.GetArea()
	r.Diagnostics = p.GetDiagnostics()
	r.DispatchedAt = p.GetDispatchedAt().AsTime()
	r.CreatedAt = time.Now()

	return ctx, r, nil
}
