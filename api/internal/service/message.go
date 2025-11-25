package service

import (
	"context"
	"errors"
	"fmt"
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
	DeleteMessage(context.Context, uint, string) error
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

func (s *Service) SaveMessage(ctx context.Context, alias model.Alias, msgType model.MessageType) error {
	message := model.Message{
		AliasID: alias.ID,
		UserID:  alias.UserID,
		Type:    msgType,
	}

	err := s.Store.PostMessage(ctx, message)
	if err != nil {
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

func (s *Service) RemoveLastMessage(ctx context.Context, aliasId string, userId string, typ model.MessageType) error {
	messages, err := s.Store.GetMessagesByAlias(ctx, aliasId)
	if err != nil {
		log.Printf("error getting messages by alias ID: %s", err.Error())
		return ErrGetMessagesByAlias
	}

	var lastMessageID uint
	for i := len(messages) - 1; i >= 0; i-- {
		if messages[i].Type == typ {
			lastMessageID = messages[i].ID
			break
		}
	}
	if lastMessageID == 0 {
		return fmt.Errorf("no message found to delete for type: %v", typ)
	}

	err = s.Store.DeleteMessage(ctx, lastMessageID, userId)
	if err != nil {
		log.Printf("error deleting message by ID: %s", err.Error())
		return ErrDeleteMessageByUserID
	}

	return nil
}
