package api

type AliasReq struct {
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
	Recipients  string `json:"recipients"`
	FromName    string `json:"from_name"`
	Format      string `json:"format"`
	Domain      string `json:"domain"`
}

type RecipientReq struct {
	Email string `json:"email" validate:"required,email"`
}

type SubscriptionReq struct {
	ID          string `json:"id"`
	ActiveUntil string `json:"active_until"`
}

type SettingsReq struct {
	ID        string `json:"id"`
	Domain    string `json:"domain"`
	Recipient string `json:"recipient"`
	FromName  string `json:"from_name"`
}

type UserReq struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=64"`
}

type DeleteUserReq struct {
	Password string `json:"password" validate:"required,min=8,max=64"`
}

type ActivateReq struct {
	OTP string `json:"otp" validate:"required,len=6"`
}
