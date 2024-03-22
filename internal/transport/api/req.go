package api

type AliasReq struct {
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
	Recipients  string `json:"recipients"`
	FromName    string `json:"from_name"`
	Format      string `json:"format"`
}

type RecipientReq struct {
	Email string `json:"email"`
}

type SubscriptionReq struct {
	ID          string `json:"id"`
	ActiveUntil string `json:"active_until"`
}

type UserReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ActivateReq struct {
	OTP string `json:"otp"`
}
