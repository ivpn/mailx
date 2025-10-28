package model

import "time"

type Bounce struct {
	ID             string    `json:"id"`
	CreatedAt      time.Time `json:"created_at"`
	AttemptedAt    time.Time `json:"attempted_at"`
	UserID         string    `json:"-"`
	AliasID        string    `json:"-"`
	From           string    `json:"from"`
	Destination    string    `json:"destination"`
	Status         string    `json:"status"`
	DiagnosticCode string    `json:"diagnostic_code"`
	RemoteMta      string    `json:"remote_mta"`
}
