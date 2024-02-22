package model

type Alias struct {
	BaseModel
	RecipientID string `json:"recipient_id"`
	Name        string `json:"name"`
	Descripion  string `json:"description"`
}
