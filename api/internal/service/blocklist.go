package service

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
	"ivpn.net/email/api/internal/model"
)

var (
	ErrGetBlocklist            = errors.New("Unable to retrieve blocklist.")
	ErrPostBlocklistItem       = errors.New("Unable to add to blocklist. Please try again.")
	ErrPostBlocklistItemLimit  = errors.New("You’ve reached the maximum number of blocklist entries.")
	ErrDuplicateBlocklistItem  = errors.New("This entry is already in your blocklist.")
	ErrDeleteBlocklistItem     = errors.New("Unable to delete blocklist entry. Please try again.")
	ErrDeleteBlocklistByUserID = errors.New("Unable to delete blocklist for this user.")
)

type BlocklistStore interface {
	GetBlocklist(context.Context, string, int, int) ([]model.BlocklistItem, error)
	GetBlocklistCount(context.Context, string) (int, error)
	FindBlocklistMatches(context.Context, string, string, []string) ([]model.BlocklistItem, error)
	PostBlocklistItem(context.Context, model.BlocklistItem) (model.BlocklistItem, error)
	IncrementBlocklistBlocked(context.Context, string) error
	DeleteBlocklistItem(context.Context, string, string) error
	DeleteBlocklistByUserID(context.Context, string) error
}

func (s *Service) GetBlocklist(ctx context.Context, userID string, limit int, page int) (model.BlocklistList, error) {
	offset := (page - 1) * limit
	if page < 1 {
		offset = 0
	}

	items, err := s.Store.GetBlocklist(ctx, userID, limit, offset)
	if err != nil {
		log.Printf("error fetching blocklist: %s", err.Error())
		return model.BlocklistList{}, ErrGetBlocklist
	}

	total, err := s.Store.GetBlocklistCount(ctx, userID)
	if err != nil {
		log.Printf("error fetching blocklist count: %s", err.Error())
		return model.BlocklistList{}, ErrGetBlocklist
	}

	return model.BlocklistList{
		Items: items,
		Total: total,
	}, nil
}

func (s *Service) PostBlocklistItem(ctx context.Context, item model.BlocklistItem) (model.BlocklistItem, error) {
	count, err := s.Store.GetBlocklistCount(ctx, item.UserID)
	if err != nil {
		log.Printf("error fetching blocklist count: %s", err.Error())
		return model.BlocklistItem{}, ErrPostBlocklistItem
	}

	if count >= model.MaxBlocklistItems {
		return model.BlocklistItem{}, ErrPostBlocklistItemLimit
	}

	created, err := s.Store.PostBlocklistItem(ctx, item)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return model.BlocklistItem{}, ErrDuplicateBlocklistItem
		}
		log.Printf("error creating blocklist item: %s", err.Error())
		return model.BlocklistItem{}, ErrPostBlocklistItem
	}

	return created, nil
}

func (s *Service) DeleteBlocklistItem(ctx context.Context, ID string, userID string) error {
	err := s.Store.DeleteBlocklistItem(ctx, ID, userID)
	if err != nil {
		log.Printf("error deleting blocklist item: %s", err.Error())
		return ErrDeleteBlocklistItem
	}

	return nil
}

func (s *Service) DeleteBlocklistByUserID(ctx context.Context, userID string) error {
	err := s.Store.DeleteBlocklistByUserID(ctx, userID)
	if err != nil {
		log.Printf("error deleting blocklist by user ID: %s", err.Error())
		return ErrDeleteBlocklistByUserID
	}

	return nil
}

// FindBlocklistMatch reports whether the user has blocked the sender. A failed lookup is treated
// as no match, so a database hiccup delivers a message that should have been blocked instead of
// silently dropping one that shouldn't.
func (s *Service) FindBlocklistMatch(ctx context.Context, userID string, from string) (model.BlocklistItem, bool) {
	email, domains := model.BlocklistCandidates(from)
	if email == "" {
		return model.BlocklistItem{}, false
	}

	items, err := s.Store.FindBlocklistMatches(ctx, userID, email, domains)
	if err != nil {
		log.Printf("error fetching blocklist matches: %s", err.Error())
		return model.BlocklistItem{}, false
	}

	return pickBlocklistMatch(items, email, domains)
}

// pickBlocklistMatch picks the item credited with the block when several match: an exact address
// match first, then the most specific domain (domains is ordered most specific first).
func pickBlocklistMatch(items []model.BlocklistItem, email string, domains []string) (model.BlocklistItem, bool) {
	for _, item := range items {
		if item.Type == model.BlocklistEmail && item.Value == email {
			return item, true
		}
	}

	for _, domain := range domains {
		for _, item := range items {
			if item.Type == model.BlocklistDomain && item.Value == domain {
				return item, true
			}
		}
	}

	return model.BlocklistItem{}, false
}

// BlockMessage discards a message the way a disabled alias does - the alias records a block in
// its stats and, with diagnostics enabled, the block is logged - and credits the matched item.
func (s *Service) BlockMessage(alias model.Alias, from string, to string, item model.BlocklistItem) {
	ctx := context.Background()

	if err := s.Store.IncrementBlocklistBlocked(ctx, item.ID); err != nil {
		log.Printf("error incrementing blocklist blocked count: %s", err.Error())
	}

	if err := s.SaveMessage(ctx, alias, model.Block); err != nil {
		log.Println("error saving message", err)
	}

	settings, err := s.GetSettings(ctx, alias.UserID)
	if err != nil {
		log.Println("error getting settings", err)
		return
	}

	if settings.LogIssues {
		reason := fmt.Sprintf("Sender is on your blocklist (%s: %s)", item.Type, item.Value)
		if err := s.ProcessDiagnosticLog(alias, from, to, reason, model.BlockedSender); err != nil {
			log.Println("error processing diagnostic log", err)
		}
	}
}
