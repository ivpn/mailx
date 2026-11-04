package api

type SuccessRes struct {
	Message string `json:"message"`
}

type AliasPauseRes struct {
	Message     string `json:"message"`
	PausedUntil string `json:"paused_until" example:"2026-10-15T08:00:00Z"`
}

type ErrorRes struct {
	Error string `json:"error"`
}
