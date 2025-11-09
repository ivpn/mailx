package model

import "time"

type Discard struct {
	ID          string    `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	AttemptedAt time.Time `json:"attempted_at"`
	UserID      string    `json:"-"`
	AliasID     string    `json:"-"`
	From        string    `json:"from"`
	Destination string    `json:"destination"`
	Message     string    `json:"message"`
}
