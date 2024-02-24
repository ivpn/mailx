package service

import (
	"context"
	"errors"
	"log"
	"time"

	"gorm.io/gorm"
	"ivpn.net/email-service/internal/client/mailer"
	"ivpn.net/email-service/internal/model"
	"ivpn.net/email-service/internal/utils"
)

var (
	ErrPostUser    = errors.New("could not save user")
	ErrGenerateOTP = errors.New("could not generate OTP")
	ErrSaveOTP     = errors.New("could not save OTP")
	ErrSendOTP     = errors.New("could not send OTP")
)

type UserStore interface {
	PostUser(context.Context, model.User) error
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

	otp, err := utils.GenerateOTP()
	if err != nil {
		log.Printf("error creating user: %s", err.Error())
		return ErrGenerateOTP
	}

	err = s.Cache.Set(ctx, "activation_"+user.ID, otp.Hash, 10*time.Minute)
	if err != nil {
		log.Printf("error creating user: %s", err.Error())
		return ErrSaveOTP
	}

	utils.Background(func() {
		err = mailer.Send("form@example.net", user.Email, "Activate your account", utils.FormatOTP(otp.Secret))
		if err != nil {
			log.Printf("error creating user: %s", err.Error())
		}
	})

	return nil
}
