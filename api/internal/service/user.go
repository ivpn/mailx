package service

import (
	"context"
	"crypto/sha256"
	"encoding/base32"
	"encoding/base64"
	"errors"
	"log"
	"strings"

	"slices"

	"github.com/go-sql-driver/mysql"
	"ivpn.net/email/api/internal/client/mailer"
	"ivpn.net/email/api/internal/model"
	"ivpn.net/email/api/internal/utils"
)

var (
	ErrGetUser             = errors.New("Unable to retrieve user by ID.")
	ErrGetUserStats        = errors.New("Unable to retrieve user statistics.")
	ErrSaveUser            = errors.New("Unable to save user.")
	ErrPostUser            = errors.New("Unable to create user.")
	ErrSignupWebhook       = errors.New("Unable to call signup webhook.")
	ErrActivateUser        = errors.New("Unable to activate user.")
	ErrDeleteUser          = errors.New("Unable to delete user.")
	ErrCreateOTP           = errors.New("Unable to generate OTP.")
	ErrSaveOTP             = errors.New("Unable to save OTP.")
	ErrSendOTP             = errors.New("Unable to send OTP.")
	ErrExpiredOTP          = errors.New("The OTP has expired. Please request a new one.")
	ErrIncorrectOTP        = errors.New("The OTP you entered is incorrect.")
	ErrIncorrectEmail      = errors.New("The email address you entered is incorrect.")
	ErrIncorrectPass       = errors.New("The password you entered is incorrect.")
	ErrLogoutUser          = errors.New("Unable to log out.")
	ErrChangePassword      = errors.New("Unable to change password.")
	ErrChangeEmail         = errors.New("Unable to change email.")
	ErrTotpDisabled        = errors.New("Two-factor authentication is disabled.")
	ErrGetTotp             = errors.New("Unable to retrieve 2FA code.")
	ErrTotpBackupUsed      = errors.New("This 2FA backup code has already been used.")
	ErrTotpBackupNotFound  = errors.New("2FA backup code not found.")
	ErrTotpSetBackup       = errors.New("Unable to set 2FA backup.")
	ErrTotpDisable         = errors.New("Unable to disable 2FA.")
	ErrInvalidTOTPCode     = errors.New("The 2FA code you entered is invalid.")
	ErrInvalidSubscription = errors.New("Invalid subscription or signup URL.")
	ErrTokenHashMismatch   = errors.New("Subscription token hash does not match.")
)

type UserStore interface {
	GetUser(context.Context, string) (model.User, error)
	GetUserByEmail(context.Context, string) (model.User, error)
	GetUserByEmailUnfinishedSignup(context.Context, string) (model.User, error)
	PostUser(context.Context, model.User) (model.User, error)
	ActivateUser(context.Context, string) error
	SaveUser(context.Context, model.User) error
	DeleteUser(context.Context, string) error
	GetUserStats(context.Context, string) (model.UserStats, error)
	TotpEnable(context.Context, string, string, string) error
	TotpDisable(context.Context, string) error
	TotpGetBackup(context.Context, string) (string, string, error)
	TotpSetUsedBackup(context.Context, string, string) error
}

func (s *Service) GetUser(ctx context.Context, userID string) (model.User, error) {
	user, err := s.Store.GetUser(ctx, userID)
	if err != nil {
		return model.User{}, ErrGetUser
	}

	return user, nil
}

func (s *Service) GetUserByCredentials(ctx context.Context, email string, password string) (model.User, error) {
	user, err := s.Store.GetUserByEmail(ctx, email)
	if err != nil {
		return model.User{}, ErrIncorrectEmail
	}

	matches := user.Matches(password)
	if !matches {
		return model.User{}, ErrIncorrectPass
	}

	return user, nil
}

func (s *Service) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	user, err := s.Store.GetUserByEmail(ctx, email)
	if err != nil {
		return model.User{}, ErrIncorrectEmail
	}

	return user, nil
}

func (s *Service) GetUnfinishedSignupOrPostUser(ctx context.Context, user model.User, subID string, sessionId string) (model.User, error) {
	email := user.Email
	pass := user.PasswordPlain
	user, err := s.Store.GetUserByEmailUnfinishedSignup(ctx, email)
	if err != nil {
		user := model.User{
			Email:         email,
			PasswordPlain: pass,
			IsActive:      false,
		}
		err = s.PostUser(ctx, user, subID, sessionId)
		if err != nil {
			log.Printf("error creating user: %s", err.Error())
			return model.User{}, ErrPostUser
		}
	}

	return user, nil
}

func (s *Service) GetUserByPassword(ctx context.Context, userID string, password string) (model.User, error) {
	user, err := s.Store.GetUser(ctx, userID)
	if err != nil {
		return model.User{}, ErrGetUser
	}

	matches := user.Matches(password)
	if !matches {
		return model.User{}, ErrIncorrectPass
	}

	return user, nil
}

func (s *Service) SaveUser(ctx context.Context, user model.User) error {
	err := s.Store.SaveUser(ctx, user)
	if err != nil {
		log.Printf("error saving user: %s", err.Error())
		return ErrSaveUser
	}

	return nil
}

func (s *Service) PostUser(ctx context.Context, user model.User, subID string, sessionId string) error {
	paSession, err := s.GetPASession(ctx, sessionId)
	if err != nil {
		log.Printf("error creating user: %s", err.Error())
		return ErrPASessionNotFound
	}

	preauthId := paSession.PreauthId
	token := paSession.Token
	tokenHash := sha256.Sum256([]byte(token))
	tokenHashStr := base64.StdEncoding.EncodeToString(tokenHash[:])

	preauth, err := s.Http.GetPreauth(preauthId)
	if err != nil {
		log.Printf("error creating user: %s", err.Error())
		return ErrInvalidSubscription
	}

	if preauth.TokenHash != tokenHashStr {
		log.Printf("error creating user: Token hash does not match")
		return ErrTokenHashMismatch
	}

	exists, err := s.Store.CheckDuplicateRecipient(ctx, user.Email)
	if exists || err != nil {
		log.Printf("error creating user: ErrDuplicateEmail")
		return model.ErrDuplicateEmail
	}

	if user.PasswordPlain != nil {
		err = user.SetPassword(*user.PasswordPlain)
		if err != nil {
			log.Printf("error creating user: %s", err.Error())
			return err
		}
	}

	user, err = s.Store.PostUser(ctx, user)
	if err != nil {
		log.Printf("error creating user: %s", err.Error())
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return model.ErrDuplicateEmail
		} else {
			return ErrPostUser
		}
	}

	err = s.PostSubscription(ctx, user.ID, preauth)
	if err != nil {
		log.Printf("error creating user: %s", err.Error())
		return ErrPostUser
	}

	err = s.PostSettings(ctx, user.ID)
	if err != nil {
		log.Printf("error creating user: %s", err.Error())
		return ErrPostUser
	}

	err = s.Http.SignupWebhook(subID)
	if err != nil {
		log.Printf("error creating user: %s", err.Error())
		return ErrSignupWebhook
	}

	return nil
}

func (s *Service) SendUserOTP(ctx context.Context, userID string) error {
	user, err := s.Store.GetUser(ctx, userID)
	if err != nil {
		log.Printf("error sending OTP: %s", err.Error())
		return ErrGetUser
	}

	otp, err := utils.CreateOTP()
	if err != nil {
		log.Printf("error sending OTP: %s", err.Error())
		return ErrCreateOTP
	}

	err = s.Cache.Set(ctx, "activation_"+userID, otp.Hash, s.Cfg.Service.OTPExpiration)
	if err != nil {
		log.Printf("error sending OTP: %s", err.Error())
		return ErrSaveOTP
	}

	utils.Background(func() {
		data := map[string]any{
			"otp":  otp.Secret,
			"from": s.Cfg.SMTPClient.SenderName,
		}
		mailer := mailer.New(s.Cfg.SMTPClient)
		mailer.Sender = s.Cfg.SMTPClient.Sender
		mailer.SenderName = s.Cfg.SMTPClient.SenderName
		err = mailer.SendTemplate(user.Email, "Verify Your Email Address", "otp_account.tmpl", data)
		if err != nil {
			log.Printf("error sending OTP: %s", err.Error())
		}
	})

	return nil
}

func (s *Service) ActivateUser(ctx context.Context, ID string, otp string) error {
	hash, err := s.Cache.Get(ctx, "activation_"+ID)
	if err != nil {
		log.Printf("error activating user: %s", err.Error())
		return ErrExpiredOTP
	}

	if !utils.MatchOTP(otp, hash) {
		return ErrIncorrectOTP
	}

	err = s.Store.ActivateUser(ctx, ID)
	if err != nil {
		log.Printf("error activating user: %s", err.Error())
		return ErrActivateUser
	}

	err = s.Cache.Del(ctx, "activation_"+ID)
	if err != nil {
		log.Printf("error activating user: %s", err.Error())
	}

	user, err := s.Store.GetUser(ctx, ID)
	if err != nil {
		log.Printf("error activating user: %s", err.Error())
		return ErrActivateUser
	}

	sub, err := s.GetSubscription(context.Background(), ID)
	if err != nil {
		log.Printf("error fetching subscription: %s", err.Error())
		return nil
	}

	if !sub.ActiveStatus() {
		log.Println("error creating recipient: subscription is not active")
		return nil
	}

	recipient := model.Recipient{
		UserID:   ID,
		Email:    user.Email,
		IsActive: true,
	}

	_, err = s.Store.PostRecipient(ctx, recipient)
	if err != nil {
		log.Printf("error saving account email as recipient: %s", err.Error())
	}

	return nil
}

func (s *Service) DeleteUserRequest(ctx context.Context, userID string) (string, error) {
	otp, err := utils.GenerateRandomString(8)
	if err != nil {
		log.Printf("error deleting user request: %s", err.Error())
		return "", ErrDeleteUser
	}

	err = s.Cache.Set(ctx, "delete_account_"+userID, otp, s.Cfg.Service.OTPExpiration)
	if err != nil {
		log.Printf("error deleting user request: %s", err.Error())
		return "", ErrSaveOTP
	}

	return otp, nil
}

func (s *Service) DeleteUser(ctx context.Context, userID string, OTP string) error {
	otp, err := s.Cache.Get(ctx, "delete_account_"+userID)
	if err != nil {
		log.Printf("error deleting user: %s", err.Error())
		return ErrIncorrectOTP
	}

	if otp != OTP {
		log.Printf("error deleting user: OTP does not match")
		return ErrIncorrectOTP
	}

	err = s.Store.DeleteAliasByUserID(ctx, userID)
	if err != nil {
		log.Printf("error deleting user: %s", err.Error())
		return ErrDeleteUser
	}

	err = s.Store.DeleteRecipientByUserID(ctx, userID)
	if err != nil {
		log.Printf("error deleting user: %s", err.Error())
		return ErrDeleteUser
	}

	err = s.Store.DeleteMessageByUserID(ctx, userID)
	if err != nil {
		log.Printf("error deleting user: %s", err.Error())
		return ErrDeleteUser
	}

	err = s.Store.DeleteSubscription(ctx, userID)
	if err != nil {
		log.Printf("error deleting user: %s", err.Error())
		return ErrDeleteUser
	}

	err = s.Store.DeleteSettings(ctx, userID)
	if err != nil {
		log.Printf("error deleting user: %s", err.Error())
		return ErrDeleteUser
	}

	err = s.Store.DeleteCredentialByUserID(ctx, userID)
	if err != nil {
		log.Printf("error deleting user: %s", err.Error())
		return ErrDeleteUser
	}

	err = s.Store.DeleteSessionByUserID(ctx, userID)
	if err != nil {
		log.Printf("error deleting user: %s", err.Error())
		return ErrDeleteUser
	}

	err = s.Store.DeleteUser(ctx, userID)
	if err != nil {
		log.Printf("error deleting user: %s", err.Error())
		return ErrDeleteUser
	}

	return nil
}

func (s *Service) GetUserStats(ctx context.Context, userID string) (model.UserStats, error) {
	stats, err := s.Store.GetUserStats(ctx, userID)
	if err != nil {
		log.Printf("error getting user stats: %s", err.Error())
		return model.UserStats{}, ErrGetUserStats
	}

	return stats, nil
}

func (s *Service) LogoutUser(ctx context.Context, authnToken string) error {
	err := s.Store.DeleteSession(ctx, authnToken)
	if err != nil && authnToken != "" {
		log.Printf("error deleting session: %s", err.Error())
		return ErrLogoutUser
	}

	return nil
}

func (s *Service) ChangePassword(ctx context.Context, userID string, password string) error {
	user, err := s.Store.GetUser(ctx, userID)
	if err != nil {
		log.Printf("error changing password: %s", err.Error())
		return ErrGetUser
	}

	err = user.SetPassword(password)
	if err != nil {
		log.Printf("error changing password: %s", err.Error())
		return ErrChangePassword
	}

	err = s.Store.SaveUser(ctx, user)
	if err != nil {
		log.Printf("error changing password: %s", err.Error())
		return ErrChangePassword
	}

	return nil
}

func (s *Service) ChangeEmail(ctx context.Context, userID string, email string) error {
	user, err := s.Store.GetUser(ctx, userID)
	if err != nil {
		log.Printf("error changing email: %s", err.Error())
		return ErrGetUser
	}

	user.Email = email
	user.IsActive = false

	err = s.Store.SaveUser(ctx, user)
	if err != nil {
		log.Printf("error changing email: %s", err.Error())
		return ErrChangeEmail
	}

	return nil
}

func (s *Service) InitiatePasswordReset(ctx context.Context, email string) error {
	user, err := s.Store.GetUserByEmail(ctx, email)
	if err != nil {
		log.Printf("error initiating password reset: %s", err.Error())
		return ErrIncorrectEmail
	}

	otp, err := utils.CreateLongOTP()
	if err != nil {
		log.Printf("error initiating password reset: %s", err.Error())
		return ErrCreateOTP
	}

	err = s.Cache.Set(ctx, "reset_"+otp.Secret, email, s.Cfg.Service.OTPExpiration)
	if err != nil {
		log.Printf("error initiating password reset: %s", err.Error())
		return ErrSaveOTP
	}

	utils.Background(func() {
		data := map[string]any{
			"otp":        otp.Secret,
			"from":       s.Cfg.SMTPClient.SenderName,
			"origin":     s.Cfg.API.ApiAllowOrigin,
			"expiration": s.Cfg.Service.OTPExpiration.Minutes(),
		}
		mailer := mailer.New(s.Cfg.SMTPClient)
		mailer.Sender = s.Cfg.SMTPClient.Sender
		mailer.SenderName = s.Cfg.SMTPClient.SenderName
		err = mailer.SendTemplate(user.Email, "["+mailer.SenderName+"] Reset Password Notification", "password_reset.tmpl", data)
		if err != nil {
			log.Printf("error initiating password reset: %s", err.Error())
		}
	})

	return nil
}

func (s *Service) ResetPassword(ctx context.Context, otp string, password string) error {
	email, err := s.Cache.Get(ctx, "reset_"+otp)
	if err != nil {
		log.Printf("error resetting password: %s", err.Error())
		return ErrExpiredOTP
	}

	err = s.Cache.Del(ctx, "reset_"+otp)
	if err != nil {
		log.Printf("error resetting password: %s", err.Error())
	}

	user, err := s.Store.GetUserByEmail(ctx, email)
	if err != nil {
		log.Printf("error resetting password: %s", err.Error())
		return ErrIncorrectEmail
	}

	err = user.SetPassword(password)
	if err != nil {
		log.Printf("error resetting password: %s", err.Error())
		return ErrChangePassword
	}

	err = s.Store.SaveUser(ctx, user)
	if err != nil {
		log.Printf("error resetting password: %s", err.Error())
		return ErrChangePassword
	}

	return nil
}

func (s *Service) TotpEnable(ctx context.Context, userID string) (model.TOTPNew, error) {
	random, err := utils.RandomString(10, utils.AlphaNumericUserFriendlyUppercase)
	if err != nil {
		log.Printf("error enabling TOTP: %s", err.Error())
		return model.TOTPNew{}, ErrCreateOTP
	}

	totpSecret := base32.StdEncoding.EncodeToString(
		[]byte(random),
	)

	err = s.Cache.Set(ctx, "totp_"+userID, totpSecret, s.Cfg.Service.OTPExpiration)
	if err != nil {
		log.Printf("error enabling TOTP: %s", err.Error())
		return model.TOTPNew{}, ErrSaveOTP
	}

	user, err := s.Store.GetUser(ctx, userID)
	if err != nil {
		log.Printf("error enabling TOTP: %s", err.Error())
		return model.TOTPNew{}, ErrGetUser
	}

	return model.TOTPNew{
		Secret:  totpSecret,
		Account: user.Email,
		URI:     utils.GenerateURI(totpSecret, user.Email, s.Cfg.SMTPClient.SenderName),
	}, nil
}

func (s *Service) TotpEnableConfirm(ctx context.Context, userID string, otp string) (model.TOTPBackup, error) {
	secret, err := s.Cache.Get(ctx, "totp_"+userID)
	if err != nil {
		log.Printf("error enabling TOTP: %s", err.Error())
		return model.TOTPBackup{}, ErrExpiredOTP
	}

	user, err := s.Store.GetUser(ctx, userID)
	if err != nil {
		log.Printf("error enabling TOTP: %s", err.Error())
		return model.TOTPBackup{}, ErrGetUser
	}

	user.TotpSecret = secret

	isValid, err := user.VerifyTotp(otp)
	if !isValid || err != nil {
		return model.TOTPBackup{}, ErrIncorrectOTP
	}

	backupCodes := []string{}

	for range 8 {
		random, err := utils.RandomString(8, utils.AlphaNumericUserFriendly)
		if err != nil {
			log.Printf("error enabling TOTP: %s", err.Error())
			return model.TOTPBackup{}, ErrCreateOTP
		}

		backupCodes = append(backupCodes, random)
	}

	totpBackup := strings.Join(backupCodes, " ")

	err = s.Store.TotpEnable(ctx, userID, secret, totpBackup)
	if err != nil {
		log.Printf("error enabling TOTP: %s", err.Error())
		return model.TOTPBackup{}, ErrSaveOTP
	}

	return model.TOTPBackup{
		Backup: totpBackup,
	}, nil
}

func (s *Service) TotpDisable(ctx context.Context, userID string, otp string) error {
	isValid, err := s.VerifyTotp(ctx, userID, otp)
	if err != nil {
		log.Printf("error disabling TOTP: %s", err.Error())
		return err
	}

	if isValid {
		err = s.Store.TotpDisable(ctx, userID)
		if err != nil {
			return ErrTotpDisable
		}

		return nil
	}

	return ErrTotpDisable
}

func (s *Service) VerifyTotp(ctx context.Context, userID string, otp string) (bool, error) {
	isValid, err := s.TotpUseBackup(ctx, userID, otp)
	if err != nil {
		log.Printf("error disabling TOTP: %s", err.Error())
		return false, err
	}

	if !isValid {
		user, err := s.Store.GetUser(ctx, userID)
		if err != nil {
			return false, ErrGetUser
		}

		isValid, err = user.VerifyTotp(otp)
		if err != nil {
			return false, ErrInvalidTOTPCode
		}
	}

	idLimiter := utils.IDLimiter{
		ID:    userID,
		Label: "totp_fails",
		Max:   s.Cfg.Service.IdLimiterMax,
		Exp:   s.Cfg.Service.IdLimiterExpiration,
		Cache: s.Cache,
	}

	if !isValid {
		err = idLimiter.Tick()
		if err != nil {
			log.Printf("error ticking ID limiter: %s", err.Error())
			return false, ErrInvalidTOTPCode
		}
	}

	if !idLimiter.IsAllowed() {
		log.Printf("error verifying TOTP: too many failed attempts")
		return false, ErrInvalidTOTPCode
	}

	return isValid, nil
}

func (s *Service) TotpUseBackup(ctx context.Context, userID string, backup string) (bool, error) {
	backups, used, err := s.Store.TotpGetBackup(ctx, userID)
	if err != nil {
		return false, ErrGetTotp
	}

	usedSlice := strings.Fields(used)

	if slices.Contains(usedSlice, backup) {
		return false, ErrTotpBackupUsed
	}

	found := slices.Contains(strings.Fields(backups), backup)

	if !found {
		return false, nil
	}

	usedSlice = append(usedSlice, backup)
	used = strings.Join(usedSlice, " ")

	err = s.Store.TotpSetUsedBackup(ctx, userID, used)
	if err != nil {
		return false, ErrTotpSetBackup
	}

	return true, nil
}
