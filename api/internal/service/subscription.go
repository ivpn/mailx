package service

import (
	"context"
	"errors"
	"log"

	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"ivpn.net/email/api/internal/model"
)

var (
	ErrGetSubscription    = errors.New("Unable to retrieve subscription by user ID.")
	ErrAddSubscription    = errors.New("Unable to add subscription.")
	ErrPostSubscription   = errors.New("Unable to create subscription.")
	ErrUpdateSubscription = errors.New("Unable to update subscription.")
	ErrDeleteSubscription = errors.New("Unable to delete subscription.")
)

type SubscriptionStore interface {
	GetSubscription(context.Context, string) (model.Subscription, error)
	PostSubscription(context.Context, model.Subscription) error
	UpdateSubscription(context.Context, model.Subscription) error
	DeleteSubscription(context.Context, string) error
}

func (s *Service) GetSubscription(ctx context.Context, userID string) (model.Subscription, error) {
	sub, err := s.Store.GetSubscription(ctx, userID)
	if err != nil {
		return model.Subscription{}, ErrGetSubscription
	}

	sub.Status = sub.GetStatus()

	return sub, nil
}

func (s *Service) PostSubscription(ctx context.Context, userID string, preauth model.Preauth) error {
	sub := model.Subscription{
		UserID:      userID,
		ActiveUntil: preauth.ActiveUntil,
		IsActive:    preauth.IsActive,
		Tier:        preauth.Tier,
		TokenHash:   preauth.TokenHash,
	}
	sub.ID = uuid.New().String()

	err := s.Store.PostSubscription(ctx, sub)
	if err != nil {
		log.Printf("error posting subscription: %s", err.Error())
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return model.ErrDuplicateSubscription
		} else {
			return ErrPostSubscription
		}
	}

	return nil
}

func (s *Service) AddSubscription(ctx context.Context, subscription model.Subscription, activeUntil string) error {
	err := s.Cache.Set(ctx, "sub_"+subscription.ID, activeUntil, s.Cfg.Service.OTPExpiration)
	if err != nil {
		log.Printf("error adding subscription: %s", err.Error())
		return ErrAddSubscription
	}

	return nil
}

func (s *Service) UpdateSubscription(ctx context.Context, sub model.Subscription, subID string, preauthID string, preauthTokenHash string) error {
	preauth, err := s.Http.GetPreauth(preauthID)
	if err != nil {
		log.Printf("error creating user: %s", err.Error())
		return ErrInvalidSubscription
	}

	if preauth.TokenHash != preauthTokenHash {
		log.Printf("error creating user: Token hash does not match")
		return ErrTokenHashMismatch
	}

	sub.ActiveUntil = preauth.ActiveUntil
	sub.IsActive = preauth.IsActive
	sub.Tier = preauth.Tier
	sub.TokenHash = preauth.TokenHash

	if sub.ID == "" {
		log.Printf("error updating subscription: Subscription ID is required")
		return ErrInvalidSubscription
	}

	err = s.Store.UpdateSubscription(ctx, sub)
	if err != nil {
		log.Printf("error updating subscription: %s", err.Error())
		return ErrUpdateSubscription
	}

	err = s.Http.SignupWebhook(subID)
	if err != nil {
		log.Printf("error creating user: %s", err.Error())
		return ErrSignupWebhook
	}

	return nil
}

func (s *Service) DeleteSubscription(ctx context.Context, userID string) error {
	err := s.Store.DeleteSubscription(ctx, userID)
	if err != nil {
		log.Printf("error deleting subscription: %s", err.Error())
		return ErrDeleteSubscription
	}

	return nil
}
