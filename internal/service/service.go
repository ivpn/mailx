package service

type Store interface {
	RecipienteStore
	AliasStore
	UserStore
}

type Service struct {
	Store Store
}

func New(store Store) *Service {
	return &Service{
		Store: store,
	}
}
