package service

import (
	"context"
	"errors"
	"log"

	"ivpn.net/email-service/internal/model"
)

var (
	ErrGetMessagesByUser  = errors.New("could not get messages by user ID")
	ErrGetMessagesByAlias = errors.New("could not get messages by alias ID")
	ErrPostMessage        = errors.New("could not post message")
)

type MessageStore interface {
	GetMessagesByUser(context.Context, string) ([]model.Message, error)
	GetMessagesByAlias(context.Context, string) ([]model.Message, error)
	PostMessage(context.Context, model.Message) error
}

func (s *Service) GetMessagesByUser(ctx context.Context, userID string) ([]model.Message, error) {
	messages, err := s.Store.GetMessagesByUser(ctx, userID)
	if err != nil {
		log.Printf("error getting messages by user ID: %s", err.Error())
		return nil, ErrGetMessagesByUser
	}

	return messages, nil
}

func (s *Service) GetMessagesByAlias(ctx context.Context, aliasID string) ([]model.Message, error) {
	messages, err := s.Store.GetMessagesByAlias(ctx, aliasID)
	if err != nil {
		log.Printf("error getting messages by alias ID: %s", err.Error())
		return nil, ErrGetMessagesByAlias
	}

	return messages, nil
}

func (s *Service) PostMessage(ctx context.Context, message model.Message) error {
	err := s.Store.PostMessage(ctx, message)
	if err != nil {
		log.Printf("error posting message: %s", err.Error())
		return ErrPostMessage
	}

	return nil
}
