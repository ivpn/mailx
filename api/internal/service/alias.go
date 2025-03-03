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
	ErrGetAlias            = errors.New("could not get alias by ID")
	ErrGetAliases          = errors.New("could not get aliass by recipient ID")
	ErrGetAliasByName      = errors.New("could not get alias by name")
	ErrPostAlias           = errors.New("could not create alias, try again")
	ErrPostAliasLimit      = errors.New("maximum number of aliases reached")
	ErrUpdateAlias         = errors.New("could not update alias")
	ErrDeleteAlias         = errors.New("could not delete alias")
	ErrDeleteAliasByUserID = errors.New("could not delete alias by user ID")
)

type AliasStore interface {
	GetAlias(context.Context, string, string) (model.Alias, error)
	GetAliases(context.Context, string, int, int, string, string, string) ([]model.Alias, error)
	GetAliasCount(context.Context, string, string) (int, error)
	GetAliasDailyCount(context.Context, string) (int, error)
	GetAliasByName(string) (model.Alias, error)
	PostAlias(context.Context, model.Alias) error
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

func (s *Service) GetAliases(ctx context.Context, userID string, limit int, page int, sortBy string, sortOrder string, catchAll string) (model.AliasList, error) {
	offset := (page - 1) * limit
	if page < 1 {
		offset = 0
	}

	aliases, err := s.Store.GetAliases(ctx, userID, limit, offset, sortBy, sortOrder, catchAll)
	if err != nil {
		log.Printf("error fetching aliass: %s", err.Error())
		return model.AliasList{}, ErrGetAliases
	}

	total, err := s.Store.GetAliasCount(ctx, userID, catchAll)
	if err != nil {
		log.Printf("error fetching alias count: %s", err.Error())
		return model.AliasList{}, ErrGetAliases
	}

	return model.AliasList{
		Aliases: aliases,
		Total:   total,
	}, nil
}

func (s *Service) GetAliasByName(name string) (model.Alias, error) {
	alias, err := s.Store.GetAliasByName(name)
	if err != nil {
		log.Printf("error fetching alias %s: %s", name, err.Error())
		return model.Alias{}, ErrGetAliasByName
	}

	return alias, nil
}

func (s *Service) PostAlias(ctx context.Context, alias model.Alias, format string, domain string, sufix string) error {
	count, err := s.Store.GetAliasDailyCount(ctx, alias.UserID)
	if err != nil {
		log.Printf("error creating alias: %s", err.Error())
		return ErrPostAlias
	}

	if count >= s.Cfg.Service.MaxDailyAliases {
		return ErrPostAliasLimit
	}

	// Catch-all alias
	if format == model.AliasFormatCatchAll {
		userAliases, err := s.Store.GetAliases(ctx, alias.UserID, 0, 0, "", "", "true")
		if err != nil {
			log.Printf("error fetching user aliases: %s", err.Error())
			return ErrPostAlias
		}

		for _, userAlias := range userAliases {
			if strings.Contains(userAlias.Name, domain) {
				return model.ErrDuplicateAliasDomain
			}
		}

		alias.Name = model.GenerateAlias(format, sufix) + "@" + domain
		alias.CatchAll = true
		err = s.Store.PostAlias(ctx, alias)
		if err != nil {
			log.Printf("error creating catch-all alias: %s", err.Error())
			return ErrPostAlias
		}

		return nil
	}

	// Standard alias
	for range 5 {
		alias.Name = model.GenerateAlias(format, "") + "@" + domain
		err := s.Store.PostAlias(ctx, alias)
		if err != nil {
			log.Printf("error creating standard alias: %s", err.Error())
			var mysqlErr *mysql.MySQLError
			if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
				continue
			} else {
				return ErrPostAlias
			}
		}
		break
	}

	return nil
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
