package service

import (
	"bytes"
	"context"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/araddon/dateparse"
	"github.com/google/uuid"
	"ivpn.net/email/api/internal/model"
)

var (
	ErrGetBouncesByUser     = errors.New("Unable to retrieve bounces for this user.")
	ErrPostBounce           = errors.New("Unable to create bounce.")
	ErrDeleteBounceByUserID = errors.New("Unable to delete bounces for this user.")
)

type BounceStore interface {
	GetBouncesByUser(context.Context, string) ([]model.Bounce, error)
	GetBounce(context.Context, string, string) (model.Bounce, error)
	GetBounceFile(context.Context, string) ([]byte, error)
	PostBounce(context.Context, model.Bounce) error
	DeleteBounceByUserID(context.Context, string) error
	SaveBounceToFile(context.Context, string, []byte) error
}

func (s *Service) GetBounces(ctx context.Context, userID string) ([]model.Bounce, error) {
	bounces, err := s.Store.GetBouncesByUser(ctx, userID)
	if err != nil {
		log.Printf("error getting bounces by user ID: %s", err.Error())
		return nil, ErrGetBouncesByUser
	}

	return bounces, nil
}

func (s *Service) GetBounce(ctx context.Context, bounceID string, userId string) (model.Bounce, error) {
	bounce, err := s.Store.GetBounce(ctx, bounceID, userId)
	if err != nil {
		log.Printf("error getting bounce: %s", err.Error())
		return model.Bounce{}, err
	}

	return bounce, nil
}

func (s *Service) GetBounceFile(ctx context.Context, bounceID string, userId string) ([]byte, error) {
	bounce, err := s.GetBounce(ctx, bounceID, userId)
	if err != nil {
		log.Printf("error getting bounce for file retrieval: %s", err.Error())
		return nil, err
	}

	filename := bounce.ID

	data, err := s.Store.GetBounceFile(ctx, filename)
	if err != nil {
		log.Printf("error getting bounce file: %s", err.Error())
		return nil, err
	}

	return data, nil
}

func (s *Service) PostBounce(ctx context.Context, bounce model.Bounce) error {
	err := s.Store.PostBounce(ctx, bounce)
	if err != nil {
		log.Printf("error posting bounce: %s", err.Error())
		return ErrPostBounce
	}

	return nil
}

func (s *Service) SaveBounceToFile(ctx context.Context, filename string, data []byte) error {
	err := s.Store.SaveBounceToFile(ctx, filename, data)
	if err != nil {
		log.Printf("error saving bounce to file: %s", err.Error())
		return err
	}

	return nil
}

func (s *Service) DeleteBounceByUserID(ctx context.Context, userID string) error {
	err := s.Store.DeleteBounceByUserID(ctx, userID)
	if err != nil {
		log.Printf("error deleting bounces by user ID: %s", err.Error())
		return ErrDeleteBounceByUserID
	}

	return nil
}

func (s *Service) ProcessBounce(userId string, aliasId string, data []byte, msg model.Msg) error {
	settings, err := s.GetSettings(context.Background(), userId)
	if err != nil {
		return err
	}

	if !settings.LogBounce {
		return nil
	}

	// log the entire data email for debugging
	log.Printf("Bounce email raw data: %s", string(data))

	var messageId string
	var to string
	var remoteMta string
	var status string
	var diagnosticCode string
	var date time.Time
	msgType := model.Send
	for _, line := range bytes.Split(data, []byte{'\n'}) {
		if after, ok := bytes.CutPrefix(line, []byte("Message-Id: ")); ok {
			messageId = string(after)
			if start := bytes.IndexByte(after, '<'); start != -1 {
				if end := bytes.IndexByte(after[start:], '@'); end != -1 {
					messageId = string(after[start+1 : start+end])
				}
			}
		}
		if after, ok := bytes.CutPrefix(line, []byte("Original-Recipient: rfc822;")); ok {
			to = string(after)
		}
		if after, ok := bytes.CutPrefix(line, []byte("Remote-MTA:")); ok {
			remoteMta = string(after)
		}
		if after, ok := bytes.CutPrefix(line, []byte("Status:")); ok {
			status = string(after)
		}
		// Capture Diagnostic-Code header and its folded continuation lines.
		if after, ok := bytes.CutPrefix(line, []byte("Diagnostic-Code:")); ok {
			diagnosticCode = strings.TrimSpace(string(after))
			continue
		}
		// RFC 5322 header folding: continuation lines start with space or tab.
		if len(diagnosticCode) > 0 && len(line) > 0 && (line[0] == ' ' || line[0] == '\t') {
			diagnosticCode += " " + strings.TrimSpace(string(line))
		}
		if _, ok := bytes.CutPrefix(line, []byte("In-Reply-To")); ok {
			msgType = model.Reply
		}
		if _, ok := bytes.CutPrefix(line, []byte("References")); ok {
			msgType = model.Reply
		}
		if after, ok := bytes.CutPrefix(line, []byte("Date:")); ok {
			date, err = dateparse.ParseAny(string(after))
			if err != nil {
				log.Println("error parsing bounce date:", err.Error())
			}
		}
	}

	bounce := model.Bounce{
		ID:             uuid.New().String(),
		CreatedAt:      time.Now(),
		AttemptedAt:    date,
		UserID:         userId,
		AliasID:        aliasId,
		From:           msg.From,
		Destination:    to,
		RemoteMta:      remoteMta,
		Status:         status,
		DiagnosticCode: diagnosticCode,
	}

	err = s.SaveBounceToFile(context.Background(), bounce.ID, data)
	if err != nil {
		return err
	}

	err = s.PostBounce(context.Background(), bounce)
	if err != nil {
		return err
	}

	err = s.RemoveLastMessage(context.Background(), userId, userId, msgType)
	if err != nil {
		return err
	}

	log.Printf("Bounce email processed successfully, %v", messageId)

	return nil
}
