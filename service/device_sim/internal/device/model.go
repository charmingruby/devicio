package device

import "time"

type Device struct {
	ID   string
	Name string
	//... any other fields
	Area      string
	CreatedAt time.Time
}
