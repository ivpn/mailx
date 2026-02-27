package service

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"log"
	"strings"

	"ivpn.net/email/api/internal/model"
	"ivpn.net/email/api/internal/utils"
)

var (
	ErrGetDomains      = errors.New("Unable to retrieve domains.")
	ErrGetDomain       = errors.New("Unable to retrieve domain.")
	ErrGetDomainsCount = errors.New("Unable to retrieve domains count.")
	ErrGetDNSConfig    = errors.New("Unable to retrieve DNS config.")
	ErrPostDomain      = errors.New("Unable to create domain. Please try again.")
	ErrUpdateDomain    = errors.New("Unable to update domain. Please try again.")
	ErrDeleteDomain    = errors.New("Unable to delete domain. Please try again.")
	ErrDNSLookupOwner  = errors.New("Unable to verify domain ownership. Please ensure the correct TXT record is set.")
)

type DomainStore interface {
	GetDomains(context.Context, string) ([]model.Domain, error)
	GetDomain(context.Context, string, string) (model.Domain, error)
	GetDomainsCount(context.Context, string) (int64, error)
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

func (s *Service) GetDomain(ctx context.Context, domainID string, userID string) (model.Domain, error) {
	domain, err := s.Store.GetDomain(ctx, domainID, userID)
	if err != nil {
		log.Printf("error getting domain: %s", err.Error())
		return model.Domain{}, ErrGetDomain
	}

	return domain, nil
}

func (s *Service) GetDomainsCount(ctx context.Context, userId string) (int64, error) {
	count, err := s.Store.GetDomainsCount(ctx, userId)
	if err != nil {
		log.Printf("error getting domains count: %s", err.Error())
		return 0, ErrGetDomainsCount
	}

	return count, nil
}

func (s *Service) GetDNSConfig(ctx context.Context, userId string) (model.DNSConfig, error) {
	count, err := s.GetDomainsCount(ctx, userId)
	if err != nil {
		log.Printf("error getting domains count for DNS config: %s", err.Error())
		return model.DNSConfig{}, ErrGetDNSConfig
	}

	domains := strings.Split(s.Cfg.API.Domains, ",")
	if len(domains) == 0 {
		log.Printf("no domains configured for DNS config")
		return model.DNSConfig{}, ErrGetDNSConfig
	}

	verify := sha256.Sum256([]byte(s.Cfg.API.TokenSecret + userId + fmt.Sprint(count)))
	domain := domains[0]
	dkim := strings.Split(s.Cfg.SMTPClient.DkimSelector, ",")
	hosts := strings.Split(s.Cfg.SMTPClient.Host, ",")

	dnsConfig := model.DNSConfig{
		Verify: fmt.Sprintf("%x", verify),
		Domain: domain,
		DKIM:   dkim,
		Hosts:  hosts,
	}

	return dnsConfig, nil
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

func (s *Service) VerifyDomainOwner(ctx context.Context, domain string, userID string) error {
	dnsConfig, err := s.GetDNSConfig(ctx, userID)
	if err != nil {
		log.Printf("error getting DNS config for domain ownership verification: %s", err.Error())
		return ErrGetDNSConfig
	}

	ok, err := utils.LookupTXTExact(domain, "mailx-verify="+dnsConfig.Verify)
	if err != nil {
		log.Printf("error looking up TXT record for domain ownership verification: %s", err.Error())
		return ErrDNSLookupOwner
	}

	if !ok {
		log.Printf("TXT record not found for domain ownership verification")
		return ErrDNSLookupOwner
	}

	return nil
}
