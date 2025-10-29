package model

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"net/mail"
	"strings"

	"ivpn.net/email/api/internal/utils"
)

type Msg struct {
	From     string
	FromName string
	To       []string
	Subject  string
	Body     string
	Type     MessageType
}

func ParseMsg(data []byte) (Msg, error) {
	msg, err := mail.ReadMessage(bytes.NewReader(data))
	if err != nil {
		return Msg{}, err
	}

	subject := msg.Header.Get("Subject")

	to := make([]string, 0)
	for _, t := range strings.Split(msg.Header.Get("To"), ",") {
		address, err := mail.ParseAddress(t)
		if err != nil {
			return Msg{}, err
		}

		to = append(to, address.Address)
	}

	from, err := mail.ParseAddress(msg.Header.Get("From"))
	if err != nil {
		return Msg{}, err
	}
	fromAddress := from.Address
	log.Println("Parsed email from:", fromAddress)

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(msg.Body)
	if err != nil {
		return Msg{}, err
	}
	body := buf.String()
	msgType := Send

	if isReply(msg) {
		msgType = Reply
	}

	if isBounce(msg) {
		msgType = FailBounce
		fromAddress, err = ExtractOriginalFrom(data)
		if err != nil {
			return Msg{}, fmt.Errorf("extract original from bounce: %w", err)
		}
		log.Println("Parsed email from ORIGINAL:", fromAddress)

	} else {
		pass, err := utils.VerifyEmailAuth(data)
		if err != nil {
			log.Println("email authentication failed with error:", err)
		}
		if !pass {
			return Msg{}, errors.New("email authentication failed")
		}
	}

	return Msg{
		From:     fromAddress,
		FromName: from.Name,
		To:       to,
		Subject:  subject,
		Body:     body,
		Type:     msgType,
	}, nil
}

// isReply checks whether the given email is a reply.
func isReply(m *mail.Message) bool {
	if m.Header.Get("In-Reply-To") != "" || m.Header.Get("References") != "" {
		return true
	}

	return false
}

// isBounce checks whether the given email is a bounce (DSN).
func isBounce(m *mail.Message) bool {
	// 1. Check Return-Path header
	if rp := m.Header.Get("Return-Path"); strings.TrimSpace(rp) == "<>" {
		return true
	}

	// 2. Check Content-Type header
	ct := m.Header.Get("Content-Type")
	mediatype, params, err := mime.ParseMediaType(ct)
	if err == nil && strings.EqualFold(mediatype, "multipart/report") {
		if strings.EqualFold(params["report-type"], "delivery-status") {
			return true
		}
	}

	// 3. Optional: check Auto-Submitted
	if strings.EqualFold(m.Header.Get("Auto-Submitted"), "auto-replied") {
		return true
	}

	return false
}

// ExtractOriginalFrom parses a bounce/DSN email and returns the "From:"
// address of the *original* message embedded as "message/rfc822".
// Returns an empty string if not found.
func ExtractOriginalFrom(data []byte) (string, error) {
	msg, err := mail.ReadMessage(bytes.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("read message: %w", err)
	}

	contentType := msg.Header.Get("Content-Type")
	mediaType, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		return "", fmt.Errorf("parse content-type: %w", err)
	}

	// Read entire body first to safely re-use it.
	body, err := io.ReadAll(msg.Body)
	if err != nil {
		return "", fmt.Errorf("read body: %w", err)
	}

	// Must be multipart with a boundary.
	if !strings.HasPrefix(mediaType, "multipart/") || params["boundary"] == "" {
		return "", fmt.Errorf("not a multipart/report message")
	}

	mr := multipart.NewReader(bytes.NewReader(body), params["boundary"])
	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("next part: %w", err)
		}

		pType := part.Header.Get("Content-Type")
		if strings.HasPrefix(strings.ToLower(pType), "message/rfc822") {
			// Found the original message part
			innerData, err := io.ReadAll(part)
			if err != nil {
				return "", fmt.Errorf("read inner part: %w", err)
			}

			innerMsg, err := mail.ReadMessage(bytes.NewReader(innerData))
			if err != nil {
				return "", fmt.Errorf("read inner message: %w", err)
			}

			from := innerMsg.Header.Get("From")
			if from == "" {
				return "", fmt.Errorf("original From not found")
			}
			return from, nil
		}
	}

	return "", fmt.Errorf("no message/rfc822 part found")
}
