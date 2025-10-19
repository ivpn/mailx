package model

type PASession struct {
	ID        string `json:"id"`
	Token     string `json:"token"`
	PreAuthID string `json:"preauth_id"`
}
