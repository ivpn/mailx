package model

import (
	"bytes"
	"encoding/base64"
	"errors"
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

type Attachment struct {
	Filename    string
	ContentType string
	Content     []byte
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

func isReply(m *mail.Message) bool {
	if m.Header.Get("In-Reply-To") != "" || m.Header.Get("References") != "" {
		return true
	}

	return false
}

// ExtractPGPSignatures scans raw email data and returns all application/pgp-signature attachments.
func ExtractPGPSignatures(data []byte) ([]Attachment, error) {
	var results []Attachment

	msg, err := mail.ReadMessage(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	ct := msg.Header.Get("Content-Type")
	mediaType, params, err := mime.ParseMediaType(ct)
	if err != nil {
		return nil, err
	}

	if !strings.HasPrefix(mediaType, "multipart/") {
		// Not multipart, no attachments to look for
		return results, nil
	}

	mr := multipart.NewReader(msg.Body, params["boundary"])
	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		partCT := part.Header.Get("Content-Type")
		if strings.HasPrefix(partCT, "application/pgp-signature") {
			content, err := io.ReadAll(part)
			if err != nil {
				return nil, err
			}

			filename := part.FileName()
			if filename == "" {
				if _, p, err := mime.ParseMediaType(part.Header.Get("Content-Disposition")); err == nil {
					filename = p["filename"]
				}
				if filename == "" {
					filename = "signature.asc" // fallback
				}
			}

			results = append(results, Attachment{
				Filename:    filename,
				ContentType: partCT,
				Content:     content,
			})
		}
	}

	return results, nil
}

// ExtractPGPKeys scans raw email data and returns all application/pgp-keys parts as decoded attachments.
func ExtractPGPKeys(data []byte) ([]Attachment, error) {
	var results []Attachment

	msg, err := mail.ReadMessage(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	ct := msg.Header.Get("Content-Type")
	mediaType, params, err := mime.ParseMediaType(ct)
	if err != nil {
		return nil, err
	}

	if !strings.HasPrefix(mediaType, "multipart/") {
		return results, nil // not a multipart message, nothing to scan
	}

	mr := multipart.NewReader(msg.Body, params["boundary"])
	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		partCT := part.Header.Get("Content-Type")
		if strings.HasPrefix(partCT, "application/pgp-keys") {
			var content []byte
			if strings.EqualFold(part.Header.Get("Content-Transfer-Encoding"), "base64") {
				encoded, _ := io.ReadAll(part)
				decoded := make([]byte, base64.StdEncoding.DecodedLen(len(encoded)))
				n, err := base64.StdEncoding.Decode(decoded, bytes.TrimSpace(encoded))
				if err != nil {
					return nil, err
				}
				content = decoded[:n]
			} else {
				content, err = io.ReadAll(part)
				if err != nil {
					return nil, err
				}
			}

			filename := part.FileName()
			if filename == "" {
				if _, p, err := mime.ParseMediaType(part.Header.Get("Content-Disposition")); err == nil {
					filename = p["filename"]
				}
				if filename == "" {
					filename = "publickey.asc" // fallback
				}
			}

			results = append(results, Attachment{
				Filename:    filename,
				ContentType: partCT,
				Content:     content,
			})
		}
	}

	return results, nil
}
