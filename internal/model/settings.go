package model

type Settings struct {
	BaseModel
	UserID    string `json:"user_id"`
	Domain    string `json:"domain"`
	Recipient string `json:"recipient"`
	FromName  string `json:"from_name"`
}
