package model

type Settings struct {
	BaseModel
	UserID       string `json:"-"`
	Domain       string `json:"domain"`
	Recipient    string `json:"recipient"`
	FromName     string `json:"from_name"`
	AliasFormat  string `json:"alias_format"`
	LogBounce    bool   `json:"log_bounce"`
	LogDiscard   bool   `json:"log_discard"`
	RemoveHeader bool   `json:"remove_header"`
}
