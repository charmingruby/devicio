package device

import (
	"time"

	"github.com/charmingruby/devicio/lib/pkg/core/id"
	"github.com/charmingruby/devicio/lib/proto/gen/pb"
)

func (r *Routine) MapFromProto(p *pb.DeviceRoutine) {
	r.ID = id.New()
	r.DeviceID = p.GetId()
	r.Status = p.GetStatus().String()
	r.Context = p.GetContext()
	r.Area = p.GetArea()
	r.Diagnostics = p.GetDiagnostics()
	r.DispatchedAt = p.GetDispatchedAt().AsTime()
	r.CreatedAt = time.Now()
}
