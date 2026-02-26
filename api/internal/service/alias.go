package service

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/go-sql-driver/mysql"
	"ivpn.net/email/api/internal/model"
)

var (
	ErrGetAlias            = errors.New("Unable to retrieve alias by ID.")
	ErrGetAliases          = errors.New("Unable to retrieve aliases.")
	ErrGetAliasByName      = errors.New("alias not found:")
	ErrDisabledAlias       = errors.New("alias disabled:")
	ErrPostAlias           = errors.New("Unable to create alias. Please try again.")
	ErrPostAliasLimit      = errors.New("You’ve reached the maximum number of allowed aliases.")
	ErrUpdateAlias         = errors.New("Unable to update alias. Please try again.")
	ErrDeleteAlias         = errors.New("Unable to delete alias. Please try again.")
	ErrDeleteAliasByUserID = errors.New("Unable to delete aliases for this user.")
)

type AliasStore interface {
	GetAlias(context.Context, string, string) (model.Alias, error)
	GetAliases(context.Context, string, int, int, string, string, string, string) ([]model.Alias, error)
	GetAllAliases(context.Context, string) ([]model.Alias, error)
	GetAliasCount(context.Context, string, string, string) (int, error)
	GetAliasDailyCount(context.Context, string) (int, error)
	GetAliasByName(string) (model.Alias, error)
	PostAlias(context.Context, model.Alias) (model.Alias, error)
	UpdateAlias(context.Context, model.Alias) error
	DeleteAlias(context.Context, string, string) error
	DeleteAliasByUserID(context.Context, string) error
}

func (s *Service) GetAlias(ctx context.Context, ID string, userID string) (model.Alias, error) {
	alias, err := s.Store.GetAlias(ctx, ID, userID)
	if err != nil {
		log.Printf("error fetching alias: %s", err.Error())
		return model.Alias{}, ErrGetAlias
	}

	return alias, nil
}

func (s *Service) GetAliases(ctx context.Context, userID string, limit int, page int, sortBy string, sortOrder string, catchAll string, search string) (model.AliasList, error) {
	offset := (page - 1) * limit
	if page < 1 {
		offset = 0
	}

	aliases, err := s.Store.GetAliases(ctx, userID, limit, offset, sortBy, sortOrder, catchAll, search)
	if err != nil {
		log.Printf("error fetching aliases: %s", err.Error())
		return model.AliasList{}, ErrGetAliases
	}

	total, err := s.Store.GetAliasCount(ctx, userID, catchAll, search)
	if err != nil {
		log.Printf("error fetching alias count: %s", err.Error())
		return model.AliasList{}, ErrGetAliases
	}

	return model.AliasList{
		Aliases: aliases,
		Total:   total,
	}, nil
}

func (s *Service) GetAllAliases(ctx context.Context, userID string) ([]model.Alias, error) {
	aliases, err := s.Store.GetAllAliases(ctx, userID)
	if err != nil {
		log.Printf("error fetching all aliases: %s", err.Error())
		return nil, ErrGetAliases
	}

	return aliases, nil
}

func (s *Service) GetAliasByName(name string) (model.Alias, error) {
	alias, err := s.Store.GetAliasByName(name)
	if err != nil {
		return model.Alias{Name: name}, ErrGetAliasByName
	}

	return alias, nil
}

func (s *Service) PostAlias(ctx context.Context, alias model.Alias, format string, domain string, sufix string) (model.Alias, error) {
	sub, err := s.GetSubscription(context.Background(), alias.UserID)
	if err != nil {
		log.Printf("error fetching subscription: %s", err.Error())
		return model.Alias{}, ErrPostAlias
	}

	if !sub.ActiveStatus() {
		return model.Alias{}, ErrPostAlias
	}

	count, err := s.Store.GetAliasDailyCount(ctx, alias.UserID)
	if err != nil {
		log.Printf("error creating alias: %s", err.Error())
		return model.Alias{}, ErrPostAlias
	}

	if count >= s.Cfg.Service.MaxDailyAliases {
		return model.Alias{}, ErrPostAliasLimit
	}

	// Catch-all alias
	if format == model.AliasFormatCatchAll {
		userAliases, err := s.Store.GetAliases(ctx, alias.UserID, 0, 0, "created_at", "DESC", "true", "")
		if err != nil {
			log.Printf("error fetching user aliases: %s", err.Error())
			return model.Alias{}, ErrPostAlias
		}

		// Count how many catch-all aliases the user already has for this domain
		domainAliasCount := 0
		for _, userAlias := range userAliases {
			if strings.Contains(userAlias.Name, domain) {
				domainAliasCount++
				if domainAliasCount >= 2 {
					return model.Alias{}, model.ErrDuplicateAliasDomain
				}
			}
		}

		alias.Name = model.GenerateAlias(format, sufix) + "@" + domain
		alias.CatchAll = true
		alias, err = s.Store.PostAlias(ctx, alias)
		if err != nil {
			log.Printf("error creating catch-all alias: %s", err.Error())
			return model.Alias{}, ErrPostAlias
		}

		return alias, nil
	}

	// Standard alias
	for range 5 {
		alias.Name = model.GenerateAlias(format, "") + "@" + domain
		alias, err = s.Store.PostAlias(ctx, alias)
		if err != nil {
			log.Printf("error creating standard alias: %s", err.Error())
			var mysqlErr *mysql.MySQLError
			if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
				continue
			} else {
				return model.Alias{}, ErrPostAlias
			}
		}
		break
	}

	return alias, nil
}

func (s *Service) UpdateAlias(ctx context.Context, alias model.Alias) error {
	err := s.Store.UpdateAlias(ctx, alias)
	if err != nil {
		log.Printf("error updating alias: %s", err.Error())
		return ErrUpdateAlias
	}

	return nil
}

func (s *Service) DeleteAlias(ctx context.Context, ID string, userID string) error {
	err := s.Store.DeleteAlias(ctx, ID, userID)
	if err != nil {
		log.Printf("error deleting alias: %s", err.Error())
		return ErrDeleteAlias
	}

	return nil
}

func (s *Service) DeleteAliasByUserID(ctx context.Context, userID string) error {
	err := s.Store.DeleteAliasByUserID(ctx, userID)
	if err != nil {
		log.Printf("error deleting alias: %s", err.Error())
		return ErrDeleteAliasByUserID
	}

	return nil
}

func (s *Service) FindAlias(email string) (model.Alias, error) {
	name, _ := model.ParseReplyTo(email)
	alias, err := s.GetAliasByName(name)
	if err != nil {
		return model.Alias{Name: name}, err
	}

	return alias, nil
}
