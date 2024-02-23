package service

import (
	"context"
	"errors"
	"log"

	"ivpn.net/email-service/internal/model"
)

var (
	ErrPostUser = errors.New("could not post user")
)

type UserStore interface {
	PostUser(context.Context, model.User) error
}

func (s *Service) PostUser(ctx context.Context, user model.User) error {
	err := s.Store.PostUser(ctx, user)
	if err != nil {
		log.Printf("an error occurred creating user: %s", err.Error())
		return ErrPostUser
	}

	return nil
}
