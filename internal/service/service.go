package service

import "context"

type ServiceStore interface {
	GetRecipient(context.Context, string) (Recipient, error)
	GetRecipients(context.Context, string) ([]Recipient, error)
	PostRecipient(context.Context, Recipient) (Recipient, error)
	UpdateRecipient(context.Context, string, Recipient) (Recipient, error)
	DeleteRecipient(context.Context, string) error
}

type Service struct {
	Store ServiceStore
}

func New(store ServiceStore) *Service {
	return &Service{
		Store: store,
	}
}
