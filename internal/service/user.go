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
	ErrPostUser       = errors.New("could not save user")
	ErrCreateOTP      = errors.New("could not create OTP")
	ErrSaveOTP        = errors.New("could not save OTP")
	ErrSendOTP        = errors.New("could not send OTP")
	ErrIncorrectEmail = errors.New("incorrect email")
	ErrIncorrectPass  = errors.New("incorrect password")
)

type UserStore interface {
	PostUser(context.Context, model.User) error
	UpdateUser(context.Context, model.User) error
	GetUserByEmail(context.Context, string) (model.User, error)
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

	err = s.Store.PostUser(ctx, user)
	if err != nil {
		log.Printf("error creating user: %s", err.Error())
		switch {
		case errors.Is(err, gorm.ErrDuplicatedKey):
			return model.ErrDuplicateEmail
		default:
			return ErrPostUser
		}
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
		err = mailer.Send(user.Email, "Activate your account", utils.FormatOTP(otp.Secret))
		if err != nil {
			log.Printf("error creating user: %s", err.Error())
		}
	})

	return nil
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
