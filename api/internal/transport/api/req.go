package api

type AliasReq struct {
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
	Recipients  string `json:"recipients" validate:"required"`
	FromName    string `json:"from_name"`
	Format      string `json:"format"`
	Domain      string `json:"domain" validate:"required"`
}

type RecipientReq struct {
	Email string `json:"email" validate:"required,email"`
}

type SubscriptionReq struct {
	ID          string `json:"id" validate:"required,uuid"`
	ActiveUntil string `json:"active_until" validate:"required"`
}

type SettingsReq struct {
	ID        string `json:"id" validate:"required,uuid"`
	Domain    string `json:"domain"`
	Recipient string `json:"recipient"`
	FromName  string `json:"from_name"`
}

type UserReq struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"password"`
	OTP      string `json:"otp" validate:"min=0,max=8"`
}

type DeleteUserReq struct {
	Password string `json:"password" validate:"password"`
}

type ChangePasswordReq struct {
	Password string `json:"password" validate:"password"`
}

type InitiatePasswordResetReq struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordReq struct {
	OTP      string `json:"otp" validate:"required,len=32"`
	Password string `json:"password" validate:"password"`
}

type ActivateReq struct {
	OTP string `json:"otp" validate:"required,len=6"`
}

type TotpReq struct {
	OTP string `json:"otp" validate:"required,min=6,max=8"`
}
