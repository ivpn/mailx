package model

type PASession struct {
	ID        string `json:"id"`
	Token     string `json:"token"`
	PreauthId string `json:"preauth_id"`
}
