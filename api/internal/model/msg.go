package model

import (
	"bytes"
	"errors"
	"log"
	"mime"
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
	pass, err := utils.VerifyEmailAuth(data)
	if err != nil {
		log.Println("email authentication failed with error:", err)
	}
	if !pass {
		return Msg{}, errors.New("email authentication failed")
	}

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

	return Msg{
		From:     from.Address,
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
