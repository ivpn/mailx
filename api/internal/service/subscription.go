package service

import (
	"context"
	"errors"
	"log"
	"time"

	"ivpn.net/email/api/internal/model"
)

var (
	ErrGetSubscription    = errors.New("could not get subscription by user ID")
	ErrPostSubscription   = errors.New("could not create subscription")
	ErrUpdateSubscription = errors.New("could not update subscription")
	ErrDeleteSubscription = errors.New("could not delete subscription")
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

func (s *Service) PostSubscription(ctx context.Context, userID string) error {
	activeUntil := time.Now().AddDate(0, -1, 0)

	if s.Cfg.Service.SubscriptionType == string(model.Managed) {
		activeUntil = time.Now().AddDate(1, 0, 0)
	}

	sub := model.Subscription{
		Type:        model.Free,
		UserID:      userID,
		ActiveUntil: activeUntil,
	}

	err := s.Store.PostSubscription(ctx, sub)
	if err != nil {
		log.Printf("error posting subscription: %s", err.Error())
		return ErrPostSubscription
	}

	return nil
}

func (s *Service) AddSubscription(ctx context.Context, subscription model.Subscription) error {
	return nil
}

func (s *Service) UpdateSubscription(ctx context.Context, subscription model.Subscription) error {
	subscription.Type = model.Managed
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
