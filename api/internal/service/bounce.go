package service

import (
	"context"
	"errors"
	"log"

	"ivpn.net/email/api/internal/model"
)

var (
	ErrGetBouncesByUser     = errors.New("Unable to retrieve bounces for this user.")
	ErrPostBounce           = errors.New("Unable to create bounce.")
	ErrDeleteBounceByUserID = errors.New("Unable to delete bounces for this user.")
)

type BounceStore interface {
	GetBouncesByUser(context.Context, string) ([]model.Bounce, error)
	PostBounce(context.Context, model.Bounce) error
	DeleteBounceByUserID(context.Context, string) error
	SaveBounceToFile(context.Context, string, []byte) error
	GetBounceFile(context.Context, string) ([]byte, error)
}

func (s *Service) GetBouncesByUser(ctx context.Context, userID string) ([]model.Bounce, error) {
	bounces, err := s.Store.GetBouncesByUser(ctx, userID)
	if err != nil {
		log.Printf("error getting bounces by user ID: %s", err.Error())
		return nil, ErrGetBouncesByUser
	}

	return bounces, nil
}

func (s *Service) PostBounce(ctx context.Context, bounce model.Bounce) error {
	err := s.Store.PostBounce(ctx, bounce)
	if err != nil {
		log.Printf("error posting bounce: %s", err.Error())
		return ErrPostBounce
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

func (s *Service) SaveBounceToFile(ctx context.Context, filename string, data []byte) error {
	err := s.Store.SaveBounceToFile(ctx, filename, data)
	if err != nil {
		log.Printf("error saving bounce to file: %s", err.Error())
		return err
	}

	return nil
}

func (s *Service) GetBounceFile(ctx context.Context, filename string) ([]byte, error) {
	data, err := s.Store.GetBounceFile(ctx, filename)
	if err != nil {
		log.Printf("error getting bounce file: %s", err.Error())
		return nil, err
	}

	return data, nil
}
