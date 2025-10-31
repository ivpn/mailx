package service

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log"
	"time"

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
	ErrPANotFound         = errors.New("Pre-auth entry not found.")
	ErrPASessionNotFound  = errors.New("Pre-auth session not found.")
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
	sub.Outage = sub.IsOutage()

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

func (s *Service) UpdateSubscription(ctx context.Context, sub model.Subscription, subID string, sessionId string) error {
	paSession, err := s.GetPASession(ctx, sessionId)
	if err != nil {
		log.Printf("error updating subscription: %s", err.Error())
		return ErrPASessionNotFound
	}

	preauthId := paSession.PreauthId
	token := paSession.Token
	tokenHash := sha256.Sum256([]byte(token))
	tokenHashStr := base64.StdEncoding.EncodeToString(tokenHash[:])

	preauth, err := s.Http.GetPreauth(preauthId)
	if err != nil {
		log.Printf("error updating subscription: %s", err.Error())
		return ErrPANotFound
	}

	if preauth.TokenHash != tokenHashStr {
		log.Printf("error updating subscription: Token hash does not match")
		return ErrTokenHashMismatch
	}

	sub.ActiveUntil = preauth.ActiveUntil
	sub.IsActive = preauth.IsActive
	sub.Tier = preauth.Tier
	sub.TokenHash = preauth.TokenHash

	if sub.ID == "" || sub.UserID == "" {
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
		log.Printf("error updating subscription: %s", err.Error())
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

func (s *Service) AddPASession(ctx context.Context, paSession model.PASession) error {
	data, err := json.Marshal(paSession)
	if err != nil {
		log.Println("failed to marshal pre-auth session to JSON:", err)
		return err
	}

	err = s.Cache.Set(ctx, "pasession_"+paSession.ID, string(data), 15*time.Minute)
	if err != nil {
		log.Println("failed to set pre-auth session in Redis:", err)
		return err
	}

	return nil
}

func (s *Service) GetPASession(ctx context.Context, id string) (model.PASession, error) {
	data, err := s.Cache.Get(ctx, "pasession_"+id)
	if err != nil {
		log.Println("failed to get pre-auth session from Redis:", err)
		return model.PASession{}, err
	}

	var paSession model.PASession
	err = json.Unmarshal([]byte(data), &paSession)
	if err != nil {
		log.Println("failed to unmarshal pre-auth session JSON:", err)
		return model.PASession{}, err
	}

	return paSession, nil
}

func (s *Service) RotatePASessionId(ctx context.Context, id string) (string, error) {
	paSession, err := s.GetPASession(ctx, id)
	if err != nil {
		log.Println("failed to get pre-auth session for rotation:", err)
		return "", err
	}

	newID := uuid.New().String()
	paSession.ID = newID

	data, err := json.Marshal(paSession)
	if err != nil {
		log.Println("failed to marshal rotated pre-auth session to JSON:", err)
		return "", err
	}

	err = s.Cache.Set(ctx, "pasession_"+newID, string(data), 15*time.Minute)
	if err != nil {
		log.Println("failed to set rotated pre-auth session in Redis:", err)
		return "", err
	}

	err = s.Cache.Del(ctx, "pasession_"+id)
	if err != nil {
		log.Println("failed to delete old pre-auth session from Redis:", err)
		return "", err
	}

	return newID, nil
}
