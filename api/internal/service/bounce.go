package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/mnako/letters"
	"ivpn.net/email/api/internal/model"
)

var (
	ErrGetBouncesByUser     = errors.New("Unable to retrieve bounces for this user.")
	ErrPostBounce           = errors.New("Unable to create bounce.")
	ErrDeleteBounceByUserID = errors.New("Unable to delete bounces for this user.")
)

type BounceStore interface {
	GetBouncesByUser(context.Context, string) ([]model.Bounce, error)
	GetBounce(context.Context, string, string) (model.Bounce, error)
	GetBounceFile(context.Context, string) ([]byte, error)
	PostBounce(context.Context, model.Bounce) error
	DeleteBounceByUserID(context.Context, string) error
	SaveBounceToFile(context.Context, string, []byte) error
}

func (s *Service) GetBouncesByUser(ctx context.Context, userID string) ([]model.Bounce, error) {
	bounces, err := s.Store.GetBouncesByUser(ctx, userID)
	if err != nil {
		log.Printf("error getting bounces by user ID: %s", err.Error())
		return nil, ErrGetBouncesByUser
	}

	return bounces, nil
}

func (s *Service) GetBounce(ctx context.Context, bounceID string, userId string) (model.Bounce, error) {
	bounce, err := s.Store.GetBounce(ctx, bounceID, userId)
	if err != nil {
		log.Printf("error getting bounce: %s", err.Error())
		return model.Bounce{}, err
	}

	return bounce, nil
}

func (s *Service) GetBounceFile(ctx context.Context, bounceID string, userId string) ([]byte, error) {
	bounce, err := s.GetBounce(ctx, bounceID, userId)
	if err != nil {
		log.Printf("error getting bounce for file retrieval: %s", err.Error())
		return nil, err
	}

	filename := bounce.ID

	data, err := s.Store.GetBounceFile(ctx, filename)
	if err != nil {
		log.Printf("error getting bounce file: %s", err.Error())
		return nil, err
	}

	return data, nil
}

func (s *Service) PostBounce(ctx context.Context, bounce model.Bounce) error {
	err := s.Store.PostBounce(ctx, bounce)
	if err != nil {
		log.Printf("error posting bounce: %s", err.Error())
		return ErrPostBounce
	}

	return nil
}

func (s *Service) SaveBounceToFile(ctx context.Context, filename string, data []byte) error {
	err := s.Store.SaveBounceToFile(ctx, filename, data)
	if err != nil {
		log.Printf("error saving bounce to file: %s", err.Error())
		return err
	}

	return nil
}

func (s *Service) DeleteBounceByUserID(ctx context.Context, userID string) error {
	err := s.Store.DeleteBounceByUserID(ctx, userID)
	if err != nil {
		log.Printf("error deleting bounces by user ID: %s", err.Error())
		return ErrDeleteBounceByUserID
	}

	return nil
}

func (s *Service) ProcessBounce(userId string, aliasId string, data []byte) error {
	reader := bytes.NewReader(data)
	email, err := letters.ParseEmail(reader)
	if err != nil {
		return err
	}

	remoteMta := ""
	if vals := email.Headers.ExtraHeaders["Received-From-MTA"]; len(vals) > 0 {
		remoteMta = vals[0]
	}

	to := ""
	if len(email.Headers.To) > 0 {
		to = email.Headers.To[0].Address
	}

	status := ""
	if vals := email.Headers.ExtraHeaders["Status"]; len(vals) > 0 {
		status = vals[0]
	}

	diagnosticCode := 0
	if vals := email.Headers.ExtraHeaders["Diagnostic-Code"]; len(vals) > 0 {
		fmt.Sscanf(vals[0], "%d", &diagnosticCode)
	}

	bounce := model.Bounce{
		ID:             uuid.New().String(),
		AttemptedAt:    email.Headers.Date,
		UserID:         userId,
		AliasID:        aliasId,
		RemoteMta:      remoteMta,
		Destination:    to,
		Status:         status,
		DiagnosticCode: diagnosticCode,
	}

	err = s.SaveBounceToFile(context.Background(), bounce.ID, data)
	if err != nil {
		return err
	}

	err = s.PostBounce(context.Background(), bounce)
	if err != nil {
		return err
	}

	return nil
}
