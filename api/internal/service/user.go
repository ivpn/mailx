package service

import (
	"context"
	"encoding/base32"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"ivpn.net/email/api/internal/client/mailer"
	"ivpn.net/email/api/internal/model"
	"ivpn.net/email/api/internal/utils"
)

var (
	ErrGetUser            = errors.New("could not get user by ID")
	ErrGetUserStats       = errors.New("could not get user stats")
	ErrSaveUser           = errors.New("could not save user")
	ErrPostUser           = errors.New("could not create user")
	ErrActivateUser       = errors.New("could not activate user")
	ErrDeleteUser         = errors.New("could not delete user")
	ErrCreateOTP          = errors.New("could not create OTP")
	ErrSaveOTP            = errors.New("could not save OTP")
	ErrSendOTP            = errors.New("could not send OTP")
	ErrExpiredOTP         = errors.New("expired OTP")
	ErrIncorrectOTP       = errors.New("incorrect OTP")
	ErrIncorrectEmail     = errors.New("incorrect email")
	ErrIncorrectPass      = errors.New("incorrect password")
	ErrLogoutUser         = errors.New("could not logout user")
	ErrChangePassword     = errors.New("could not change password")
	ErrTotpDisabled       = errors.New("2FA is disabled")
	ErrGetTotp            = errors.New("could not get 2FA code")
	ErrTotpBackupUsed     = errors.New("2FA backup is already used")
	ErrTotpBackupNotFound = errors.New("2FA backup not found")
	ErrTotpSetBackup      = errors.New("could not set 2FA backup")
	ErrTotpDisable        = errors.New("could not disable 2FA")
	ErrInvalidTOTPCode    = errors.New("invalid 2FA code")
)

type UserStore interface {
	GetUser(context.Context, string) (model.User, error)
	GetUserByEmail(context.Context, string) (model.User, error)
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

func (s *Service) GetOrPostUser(ctx context.Context, user model.User) (model.User, error) {
	user, err := s.Store.GetUserByEmail(ctx, user.Email)
	if err != nil {
		user, err = s.Store.PostUser(ctx, user)
		if err != nil {
			log.Printf("error creating user: %s", err.Error())
			return model.User{}, ErrPostUser
		}
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

func (s *Service) PostUser(ctx context.Context, user model.User) error {
	rcpCount, err := s.Store.GetRecipientsCountByEmail(ctx, user.Email)
	if err != nil || rcpCount > 0 {
		return model.ErrDuplicateEmail
	}

	err = user.SetPassword(*user.PasswordPlain)
	if err != nil {
		log.Printf("error creating user: %s", err.Error())
		return err
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

	err = s.PostSubscription(ctx, user.ID)
	if err != nil {
		log.Printf("error creating user: %s", err.Error())
		return ErrPostUser
	}

	err = s.PostSettings(ctx, user.ID)
	if err != nil {
		log.Printf("error creating user: %s", err.Error())
		return ErrPostUser
	}

	otp, err := utils.CreateOTP()
	if err != nil {
		log.Printf("error creating user: %s", err.Error())
		return ErrCreateOTP
	}

	err = s.Cache.Set(ctx, "activation_"+user.ID, otp.Hash, s.Cfg.Service.OTPExpiration)
	if err != nil {
		log.Printf("error creating user: %s", err.Error())
		return ErrSaveOTP
	}

	utils.Background(func() {
		data := map[string]interface{}{
			"otp":  otp.Secret,
			"from": s.Cfg.SMTPClient.SenderName,
		}
		mailer := mailer.New(s.Cfg.SMTPClient)
		mailer.Sender = s.Cfg.SMTPClient.Sender
		mailer.SenderName = s.Cfg.SMTPClient.SenderName
		err = mailer.SendTemplate(user.Email, "["+mailer.SenderName+"] Verify Account Notification", "otp_account.tmpl", data)
		if err != nil {
			log.Printf("error creating user: %s", err.Error())
		}
	})

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
		data := map[string]interface{}{
			"otp":  otp.Secret,
			"from": s.Cfg.SMTPClient.SenderName,
		}
		mailer := mailer.New(s.Cfg.SMTPClient)
		mailer.Sender = s.Cfg.SMTPClient.Sender
		mailer.SenderName = s.Cfg.SMTPClient.SenderName
		err = mailer.SendTemplate(user.Email, "["+mailer.SenderName+"] Verify Account Notification", "otp_account.tmpl", data)
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

func (s *Service) DeleteUser(ctx context.Context, userID string) error {
	err := s.Store.DeleteAliasByUserID(ctx, userID)
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

func (s *Service) LogoutUser(ctx context.Context, jwtSignature string, jwtExp time.Duration, authnToken string) error {
	err := s.Cache.Set(ctx, "logout_"+jwtSignature, "true", jwtExp)
	if err != nil && jwtSignature != "" {
		log.Printf("error saving jwt: %s", err.Error())
		return ErrLogoutUser
	}

	err = s.Store.DeleteSession(ctx, authnToken)
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
		data := map[string]interface{}{
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
	totpSecret := base32.StdEncoding.EncodeToString(
		[]byte(utils.RandomString(10, utils.AlphaNumericUserFriendlyUppercase)),
	)

	err := s.Cache.Set(ctx, "totp_"+userID, totpSecret, s.Cfg.Service.OTPExpiration)
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
	for i := 0; i < 8; i++ {
		backupCodes = append(backupCodes, utils.RandomString(8, utils.AlphaNumericUserFriendly))
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

	return isValid, nil
}

func (s *Service) TotpUseBackup(ctx context.Context, userID string, backup string) (bool, error) {
	backups, used, err := s.Store.TotpGetBackup(ctx, userID)
	if err != nil {
		return false, ErrGetTotp
	}

	usedSlice := strings.Fields(used)

	for _, code := range usedSlice {
		if backup == code {
			return false, ErrTotpBackupUsed
		}
	}

	found := false

	for _, code := range strings.Fields(backups) {
		if backup == code {
			found = true
			break
		}
	}

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
