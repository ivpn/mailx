package service

import (
	"context"
	"errors"
	"log"

	"gorm.io/gorm"
	"ivpn.net/email-service/internal/client/mailer"
	"ivpn.net/email-service/internal/model"
	"ivpn.net/email-service/internal/utils"
)

var (
	ErrGetUser        = errors.New("could not get user by ID")
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
)

type UserStore interface {
	GetUser(context.Context, string) (model.User, error)
	GetUserByEmail(context.Context, string) (model.User, error)
	PostUser(context.Context, model.User) (model.User, error)
	ActivateUser(context.Context, string) error
	DeleteUser(context.Context, string) error
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

func (s *Service) PostUser(ctx context.Context, user model.User) error {
	err := user.Validate()
	if err != nil {
		log.Printf("error creating user: %s", err.Error())
		return err
	}

	err = user.SetPassword(*user.PasswordPlain)
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
		mailer := mailer.New(s.Cfg.SMTPClient)
		err = mailer.Send(user.Email, "Activate your account", otp.Secret)
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
		mailer := mailer.New(s.Cfg.SMTPClient)
		err = mailer.Send(user.Email, "Activate your account", otp.Secret)
		if err != nil {
			log.Printf("error sending OTP: %s", err.Error())
		}
	})

	return nil
}

func (s *Service) ActivateUser(ctx context.Context, ID string, otp string) error {
	err := utils.ValidateOTP(otp)
	if err != nil {
		log.Printf("error activating user: %s", err.Error())
		return err
	}

	hash, err := s.Cache.Get(ctx, "activation_"+ID)
	if err != nil {
		log.Printf("error activating user: %s", err.Error())
		return ErrExpiredOTP
	}

	if hash != utils.HashOTP(otp) {
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

func (s *Service) DeleteUser(ctx context.Context, ID string) error {
	err := s.Store.DeleteUser(ctx, ID)
	if err != nil {
		log.Printf("error deleting user: %s", err.Error())
		return ErrDeleteUser
	}

	return nil
}
