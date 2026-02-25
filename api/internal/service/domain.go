package service

import (
	"context"
	"errors"
	"log"

	"ivpn.net/email/api/internal/model"
)

var (
	ErrGetDomains   = errors.New("Unable to retrieve domains.")
	ErrPostDomain   = errors.New("Unable to create domain. Please try again.")
	ErrUpdateDomain = errors.New("Unable to update domain. Please try again.")
	ErrDeleteDomain = errors.New("Unable to delete domain. Please try again.")
)

type DomainStore interface {
	GetDomains(context.Context, string) ([]model.Domain, error)
	PostDomain(context.Context, model.Domain) (model.Domain, error)
	UpdateDomain(context.Context, model.Domain) error
	DeleteDomain(context.Context, string, string) error
	DeleteDomainsByUserID(context.Context, string) error
}

func (s *Service) GetDomains(ctx context.Context, userId string) ([]model.Domain, error) {
	domains, err := s.Store.GetDomains(ctx, userId)
	if err != nil {
		log.Printf("error getting domains: %s", err.Error())
		return nil, ErrGetDomains
	}

	return domains, nil
}

func (s *Service) PostDomain(ctx context.Context, domain model.Domain) (model.Domain, error) {
	createdDomain, err := s.Store.PostDomain(ctx, domain)
	if err != nil {
		log.Printf("error posting domain: %s", err.Error())
		return model.Domain{}, ErrPostDomain
	}

	return createdDomain, nil
}

func (s *Service) DeleteDomain(ctx context.Context, domainID string, userID string) error {
	err := s.Store.DeleteDomain(ctx, domainID, userID)
	if err != nil {
		log.Printf("error deleting domain: %s", err.Error())
		return ErrDeleteDomain
	}

	return nil
}

func (s *Service) UpdateDomain(ctx context.Context, domain model.Domain) error {
	err := s.Store.UpdateDomain(ctx, domain)
	if err != nil {
		log.Printf("error updating domain: %s", err.Error())
		return ErrUpdateDomain
	}

	return nil
}

func (s *Service) DeleteDomainsByUserID(ctx context.Context, userID string) error {
	err := s.Store.DeleteDomainsByUserID(ctx, userID)
	if err != nil {
		log.Printf("error deleting domains by user ID: %s", err.Error())
		return ErrDeleteDomain
	}

	return nil
}
