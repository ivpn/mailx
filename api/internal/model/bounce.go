package model

import "time"

type Bounce struct {
	ID             string    `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	AttemptedAt    time.Time `json:"attempted_at"`
	UserID         string    `json:"-"`
	AliasID        string    `json:"-"`
	RemoteMta      string    `json:"remote_mta"`
	Destination    string    `json:"destination"`
	Status         string    `json:"status"`
	DiagnosticCode int       `json:"diagnostic_code"`
}
