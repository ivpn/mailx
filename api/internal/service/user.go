package service

import (
	"context"
	"errors"
	"log"
	"time"

	"gorm.io/gorm"
	"ivpn.net/email/api/internal/client/mailer"
	"ivpn.net/email/api/internal/model"
	"ivpn.net/email/api/internal/utils"
)

var (
	ErrGetUser        = errors.New("could not get user by ID")
	ErrGetUserStats   = errors.New("could not get user stats")
	ErrPostUser       = errors.New("could not create user")
	ErrActivateUser   = errors.New("could not activate user")
	ErrDeleteUser     = errors.New("could not delete user")
	ErrCreateOTP      = errors.New("could not create OTP")
	ErrSaveOTP        = errors.New("could not save OTP")
	ErrSendOTP        = errors.New("could not send OTP")
	ErrExpiredOTP     = errors.New("expired OTP")
	ErrIncorrectOTP   = errors.New("incorrect OTP")
	ErrIncorrectEmail = errors.New("incorrect email")
	ErrIncorrectPass  = errors.New("incorrect password")
	ErrLogoutUser     = errors.New("could not logout user")
)

type UserStore interface {
	GetUser(context.Context, string) (model.User, error)
	GetUserByEmail(context.Context, string) (model.User, error)
	PostUser(context.Context, model.User) (model.User, error)
	ActivateUser(context.Context, string) error
	DeleteUser(context.Context, string) error
	GetUserStats(context.Context, string) (model.UserStats, error)
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

func (s *Service) PostUser(ctx context.Context, user model.User) error {
	err := user.SetPassword(*user.PasswordPlain)
	if err != nil {
		log.Printf("error creating user: %s", err.Error())
		return err
	}

	user, err = s.Store.PostUser(ctx, user)
	if err != nil {
		log.Printf("error creating user: %s", err.Error())
		switch {
		case errors.Is(err, gorm.ErrDuplicatedKey):
			return model.ErrDuplicateEmail
		default:
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
		data := map[string]interface{}{"otp": otp.Secret}
		mailer := mailer.New(s.Cfg.SMTPClient)
		mailer.Sender = s.Cfg.SMTPClient.Sender
		mailer.SenderName = s.Cfg.SMTPClient.SenderName
		err = mailer.SendTemplate(user.Email, "["+mailer.SenderName+"] OTP to activate account", "otp_account.tmpl", data)
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
		data := map[string]interface{}{"otp": otp.Secret}
		mailer := mailer.New(s.Cfg.SMTPClient)
		mailer.Sender = s.Cfg.SMTPClient.Sender
		mailer.SenderName = s.Cfg.SMTPClient.SenderName
		err = mailer.SendTemplate(user.Email, "["+mailer.SenderName+"] OTP to activate account", "otp_account.tmpl", data)
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

func (s *Service) LogoutUser(ctx context.Context, jwtSignature string, jwtExp time.Duration) error {
	err := s.Cache.Set(ctx, "logout_"+jwtSignature, "true", jwtExp)
	if err != nil {
		log.Printf("error saving jwt: %s", err.Error())
		return ErrLogoutUser
	}

	return nil
}
