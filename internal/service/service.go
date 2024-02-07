package service

type Service struct {
	Store ServiceStore
}

type ServiceStore interface {
	RecipienteStore
	AliasStore
}

func New(store ServiceStore) *Service {
	return &Service{
		Store: store,
	}
}
