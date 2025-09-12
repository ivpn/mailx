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
	subscription, err := s.Store.GetSubscription(ctx, userID)
	if err != nil {
		return model.Subscription{}, ErrGetSubscription
	}

	return subscription, nil
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

func (s *Service) UpdateSubscription(ctx context.Context, subscription model.Subscription) error {
	err := s.Store.UpdateSubscription(ctx, subscription)
	if err != nil {
		log.Printf("error updating subscription: %s", err.Error())
		return ErrUpdateSubscription
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
