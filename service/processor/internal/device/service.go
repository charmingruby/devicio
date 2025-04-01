package device

import (
	"github.com/charmingruby/devicio/lib/pkg/messaging/rabbitmq"
	"github.com/charmingruby/devicio/lib/proto/gen/pb"
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

func (s *Service) ProcessRoutine(msg []byte) error {
	var protoRoutine pb.DeviceRoutine

	if err := proto.Unmarshal(msg, &protoRoutine); err != nil {
		return err
	}

	r := &Routine{}
	r.MapFromProto(&protoRoutine)

	if err := s.repo.Store(r); err != nil {
		return err
	}

	return nil
}
