package api

type UserReq struct {
	Email    string `json:"email" validate:"required,emailx"`
	Password string `json:"password" validate:"password"`
	OTP      string `json:"otp" validate:"min=0,max=8"`
}

type EmailReq struct {
	Email string `json:"email" validate:"required,emailx"`
}

type SignupUserReq struct {
	Email    string `json:"email" validate:"required,emailx"`
	Password string `json:"password" validate:"password"`
	SubID    string `json:"subid" validate:"required,uuid"`
}

type SignupEmailReq struct {
	Email string `json:"email" validate:"required,emailx"`
	SubID string `json:"subid" validate:"required,uuid"`
}

type SubscriptionReq struct {
	ID          string `json:"id" validate:"required,uuid"`
	ActiveUntil string `json:"active_until" validate:"required"`
}

type AliasReq struct {
	Description    string `json:"description"`
	Enabled        bool   `json:"enabled"`
	Recipients     string `json:"recipients" validate:"required"`
	FromName       string `json:"from_name"`
	Format         string `json:"format"`
	Domain         string `json:"domain" validate:"required"`
	CatchAllSuffix string `json:"catch_all_suffix" validate:"omitempty,alphanum,min=6,max=12"`
}

type RecipientReq struct {
	ID         string `json:"id" validate:"required,uuid"`
	PGPKey     string `json:"pgp_key" validate:"omitempty,pgp"`
	PGPEnabled bool   `json:"pgp_enabled"`
	PGPInline  bool   `json:"pgp_inline"`
}

type SettingsReq struct {
	ID           string `json:"id" validate:"required,uuid"`
	Domain       string `json:"domain"`
	Recipient    string `json:"recipient"`
	FromName     string `json:"from_name"`
	AliasFormat  string `json:"alias_format"`
	LogBounce    bool   `json:"log_bounce"`
	LogDiscard   bool   `json:"log_discard"`
	RemoveHeader bool   `json:"remove_header"`
}

type DeleteUserReq struct {
	OTP string `json:"otp" validate:"required,len=8"`
}

type ChangePasswordReq struct {
	Password string `json:"password" validate:"password"`
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
