package device

import "time"

type Routine struct {
	ID           string
	DeviceID     string
	Status       string
	Context      string
	Area         string
	Diagnostics  string
	DispatchedAt time.Time
	CreatedAt    time.Time
}
