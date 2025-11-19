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
	ErrGetLogs    = errors.New("Unable to retrieve logs for this user.")
	ErrPostLog    = errors.New("Unable to create log.")
	ErrDeleteLogs = errors.New("Unable to delete logs for this user.")
)

type LogStore interface {
	GetLogs(context.Context, string) ([]model.Log, error)
	GetLog(context.Context, string, string) (model.Log, error)
	GetLogFile(context.Context, string) ([]byte, error)
	PostLog(context.Context, model.Log) error
	DeleteLogs(context.Context, string) error
	SaveLogToFile(context.Context, string, []byte) error
}

func (s *Service) GetLogs(ctx context.Context, userId string) ([]model.Log, error) {
	logs, err := s.Store.GetLogs(ctx, userId)
	if err != nil {
		log.Printf("error getting logs by user ID: %s", err.Error())
		return nil, ErrGetLogs
	}

	return logs, nil
}

func (s *Service) GetLog(ctx context.Context, logId string, userId string) (model.Log, error) {
	lg, err := s.Store.GetLog(ctx, logId, userId)
	if err != nil {
		log.Printf("error getting log: %s", err.Error())
		return model.Log{}, err
	}

	return lg, nil
}

func (s *Service) GetLogFile(ctx context.Context, logId string, userId string) ([]byte, error) {
	lg, err := s.GetLog(ctx, logId, userId)
	if err != nil {
		log.Printf("error getting log for file retrieval: %s", err.Error())
		return nil, err
	}

	filename := lg.ID
	data, err := s.Store.GetLogFile(ctx, filename)
	if err != nil {
		log.Printf("error getting log file: %s", err.Error())
		return nil, err
	}

	return data, nil
}

func (s *Service) PostLog(ctx context.Context, lg model.Log) error {
	err := s.Store.PostLog(ctx, lg)
	if err != nil {
		log.Printf("error posting log: %s", err.Error())
		return ErrPostLog
	}

	return nil
}

func (s *Service) SaveLogToFile(ctx context.Context, filename string, data []byte) error {
	err := s.Store.SaveLogToFile(ctx, filename, data)
	if err != nil {
		log.Printf("error saving log to file: %s", err.Error())
		return err
	}

	return nil
}

func (s *Service) DeleteLogs(ctx context.Context, userId string) error {
	err := s.Store.DeleteLogs(ctx, userId)
	if err != nil {
		log.Printf("error deleting logs by user ID: %s", err.Error())
		return ErrDeleteLogs
	}

	return nil
}

func (s *Service) ProcessDiscardLog(alias model.Alias, from string, destination string, message string, logType model.LogType) error {
	lg := model.Log{
		ID:          uuid.New().String(),
		CreatedAt:   time.Now(),
		Type:        logType,
		UserID:      alias.UserID,
		AliasID:     alias.ID,
		From:        from,
		Destination: destination,
		Message:     message,
	}

	err := s.PostLog(context.Background(), lg)
	if err != nil {
		log.Printf("error processing discard: %s", err.Error())
		return err
	}

	return nil
}

func (s *Service) ProcessBounceLog(userId string, aliasId string, data []byte, msg model.Msg) error {
	settings, err := s.GetSettings(context.Background(), userId)
	if err != nil {
		return err
	}

	if !settings.LogBounce {
		return nil
	}

	var messageId string
	var to string
	var remoteMta string
	var status string
	var diagnosticCode string
	var date time.Time
	var inDiag = false
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
		// Collect all lines from Diagnostic-Code: until a line starting with "--"
		if inDiag {
			trim := bytes.TrimSpace(line)
			if bytes.HasPrefix(trim, []byte("--")) || len(trim) == 0 {
				inDiag = false
			} else {
				if diagnosticCode != "" {
					diagnosticCode += " "
				}
				diagnosticCode += string(trim)
			}
		}
		if after, ok := bytes.CutPrefix(line, []byte("Diagnostic-Code:")); ok {
			diagnosticCode = string(bytes.TrimSpace(after))
			inDiag = true
		}
		if _, ok := bytes.CutPrefix(line, []byte("In-Reply-To:")); ok {
			msgType = model.Reply
		}
		if _, ok := bytes.CutPrefix(line, []byte("References:")); ok {
			msgType = model.Reply
		}
		if after, ok := bytes.CutPrefix(line, []byte("Subject:")); ok {
			subject := string(after)
			if strings.Contains(subject, "Re:") || strings.Contains(subject, "RE:") {
				msgType = model.Reply
			}
		}
		if after, ok := bytes.CutPrefix(line, []byte("Date:")); ok {
			date, err = dateparse.ParseAny(string(bytes.TrimSpace(after)))
			if err != nil {
				log.Println("error parsing bounce date:", err.Error())
			}
		}
	}

	lg := model.Log{
		ID:          uuid.New().String(),
		CreatedAt:   time.Now(),
		AttemptedAt: date,
		Type:        model.BounceMessage,
		UserID:      userId,
		AliasID:     aliasId,
		From:        msg.From,
		Destination: to,
		RemoteMta:   remoteMta,
		Status:      status,
		Message:     diagnosticCode,
	}

	err = s.SaveLogToFile(context.Background(), lg.ID, data)
	if err != nil {
		return err
	}

	err = s.PostLog(context.Background(), lg)
	if err != nil {
		return err
	}

	err = s.RemoveLastMessage(context.Background(), aliasId, userId, msgType)
	if err != nil {
		return err
	}

	log.Printf("Bounce email processed successfully, %v", messageId)

	return nil
}
