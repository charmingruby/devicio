package device

import (
	"context"
	"math/rand"
	"time"

	"github.com/charmingruby/devicio/lib/messaging/rabbitmq"
	pb "github.com/charmingruby/devicio/lib/proto/gen/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Service struct {
	queue *rabbitmq.Client
}

var diagnosticOptions = []string{
	"Temperature within normal range",
	"Pressure levels optimal",
	"Flow rate stable",
	"Power consumption normal",
	"System response time acceptable",
}

var statusOptions = []pb.DeviceStatus{
	pb.DeviceStatus_HEALTHY,
	pb.DeviceStatus_WARNING,
	pb.DeviceStatus_ERROR,
	pb.DeviceStatus_CRITICAL,
}

var areas = []string{"A", "B", "C"}

func NewService(queue *rabbitmq.Client) *Service {
	return &Service{queue: queue}
}

func (s *Service) DispatchRoutineMessage(ctx context.Context, device Device) error {
	now := time.Now()
	timestamp := timestamppb.New(now)

	routine := &pb.DeviceRoutine{
		Id:           device.ID,
		Status:       getRandomStatus(),
		Context:      "routine",
		Diagnostics:  getRandomDiagnostic(),
		Area:         getRandomArea(),
		DispatchedAt: timestamp,
	}

	return s.queue.Publish(ctx, routine)
}

func getRandomDiagnostic() string {
	options := make([]string, len(diagnosticOptions))
	copy(options, diagnosticOptions)

	rand.Shuffle(len(options), func(i, j int) {
		options[i], options[j] = options[j], options[i]
	})

	return options[0]
}

func getRandomStatus() pb.DeviceStatus {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return statusOptions[r.Intn(len(statusOptions))]
}

func getRandomArea() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	return areas[r.Intn(len(areas))]
}
