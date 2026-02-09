package model

import "time"

type LogType string

const (
	BounceMessage        LogType = "bounce"
	DisabledAlias        LogType = "disabled_alias"
	UnauthorisedSend     LogType = "unauthorised_send"
	InactiveSubscription LogType = "inactive_subscription"
)

type Log struct {
	ID          string    `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	AttemptedAt time.Time `json:"attempted_at"`
	Type        LogType   `json:"log_type"`
	UserID      string    `json:"-"`
	AliasID     string    `json:"-"`
	From        string    `json:"from"`
	Destination string    `json:"destination"`
	Message     string    `json:"message"`
	Status      string    `json:"status"`
	RemoteMta   string    `json:"remote_mta"`
}
