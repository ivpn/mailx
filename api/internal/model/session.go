package model

type Session struct {
	BaseModel
	UserID      string `json:"-"`
	Token       string `gorm:"unique" json:"token"`
	SessionData []byte `gorm:"type:blob" json:"-"`
}
