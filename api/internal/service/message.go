package service

import (
	"context"
	"errors"
	"log"

	"ivpn.net/email/api/internal/model"
)

var (
	ErrGetMessagesByUser     = errors.New("Unable to retrieve messages for this user.")
	ErrGetMessagesByAlias    = errors.New("Unable to retrieve messages for this alias.")
	ErrPostMessage           = errors.New("Unable to create message.")
	ErrDeleteMessageByUserID = errors.New("Unable to delete messages for this user.")
)

type MessageStore interface {
	GetMessagesByUser(context.Context, string) ([]model.Message, error)
	GetMessagesByAlias(context.Context, string) ([]model.Message, error)
	PostMessage(context.Context, model.Message) error
	DeleteMessageByUserID(context.Context, string) error
	SendReplyDailyCount(context.Context, string) (int, error)
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

func (s *Service) DeleteMessageByUserID(ctx context.Context, userID string) error {
	err := s.Store.DeleteMessageByUserID(ctx, userID)
	if err != nil {
		log.Printf("error deleting messages by user ID: %s", err.Error())
		return ErrDeleteMessageByUserID
	}

	return nil
}

func (s *Service) ValidateSendReplyDailyCount(ctx context.Context, userID string) error {
	count, err := s.Store.SendReplyDailyCount(ctx, userID)
	if err != nil {
		log.Printf("error getting send/reply daily count: %s", err.Error())
		return ErrGetMessagesByUser
	}

	if count >= s.Cfg.Service.MaxDailySendReply {
		return errors.New("daily limit reached")
	}

	return nil
}
