package service

type Service struct {
	Store ServiceStore
}

func New(store ServiceStore) *Service {
	return &Service{
		Store: store,
	}
}
