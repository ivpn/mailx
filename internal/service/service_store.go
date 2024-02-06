package service

import "context"

type RecipienteStore interface {
	GetRecipient(context.Context, string) (Recipient, error)
	GetRecipients(context.Context, string) ([]Recipient, error)
	PostRecipient(context.Context, Recipient) (Recipient, error)
	UpdateRecipient(context.Context, string, Recipient) (Recipient, error)
	DeleteRecipient(context.Context, string) error
}

type AliasStore interface {
	GetAlias(context.Context, string) (Alias, error)
	GetAliases(context.Context, string) ([]Alias, error)
	PostAlias(context.Context, Alias) (Alias, error)
	UpdateAlias(context.Context, string, Alias) (Alias, error)
	DeleteAlias(context.Context, string) error
}

type ServiceStore interface {
	RecipienteStore
	AliasStore
}
