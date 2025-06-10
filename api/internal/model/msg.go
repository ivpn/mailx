package model

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	"net/mail"
	"net/textproto"
	"strings"

	"github.com/OfimaticSRL/parsemail"
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

func isReply(m *mail.Message) bool {
	if m.Header.Get("In-Reply-To") != "" || m.Header.Get("References") != "" {
		return true
	}

	return false
}

// ExtractPGPSignatures parses the raw email data and extracts PGP signatures
// (Content-Type: application/pgp-signature), returning them with filename
// and decoded content.
func ExtractPGPSignatures(data []byte) ([]parsemail.Attachment, error) {
	msg, err := mail.ReadMessage(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	contentType := msg.Header.Get("Content-Type")
	if contentType == "" {
		return nil, errors.New("missing Content-Type header")
	}

	mediaType, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		return nil, err
	}

	var attachments []parsemail.Attachment

	if strings.HasPrefix(mediaType, "multipart/") {
		mr := multipart.NewReader(msg.Body, params["boundary"])
		for {
			p, err := mr.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				return nil, err
			}

			ct := p.Header.Get("Content-Type")
			filename := p.FileName()
			if filename == "" {
				continue
			}

			// if filename does not end with .asc, skip it
			if !strings.HasSuffix(filename, ".asc") {
				continue
			}

			bodyBytes, err := io.ReadAll(p)
			if err != nil {
				return nil, err
			}

			enc := strings.ToLower(p.Header.Get("Content-Transfer-Encoding"))
			var decodedData []byte

			switch enc {
			case "base64":
				decodedData, err = base64.StdEncoding.DecodeString(string(bodyBytes))
				if err != nil {
					return nil, err
				}
			case "quoted-printable":
				decodedData, err = io.ReadAll(quotedprintable.NewReader(bytes.NewReader(bodyBytes)))
				if err != nil {
					return nil, err
				}
			default:
				decodedData = bodyBytes
			}

			attachments = append(attachments, parsemail.Attachment{
				Filename:    filename,
				ContentType: ct,
				Data:        bytes.NewReader(decodedData),
			})
		}
	} else {
		// Not multipart
		if mediaType == "application/pgp-signature" {
			bodyBytes, err := io.ReadAll(msg.Body)
			if err != nil {
				return nil, err
			}
			attachments = append(attachments, parsemail.Attachment{
				Filename:    "signature.asc",
				ContentType: mediaType,
				Data:        bytes.NewReader(bodyBytes),
			})
		}
	}

	return attachments, nil
}

// ExtractPGPKeys parses the raw email data and extracts attachments
// that are PGP keys (Content-Type: application/pgp-keys), returning
// them with filename and decoded content.
func ExtractPGPKeys(data []byte) ([]parsemail.Attachment, error) {
	msg, err := mail.ReadMessage(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to parse email: %w", err)
	}

	ct := msg.Header.Get("Content-Type")
	mediaType, params, err := mime.ParseMediaType(ct)
	if err != nil {
		// No valid Content-Type header, treat whole body as one PGP key (unlikely)
		return []parsemail.Attachment{}, nil
	}

	if !strings.HasPrefix(mediaType, "multipart/") {
		// Not multipart, treat whole body as one PGP key (maybe)
		return []parsemail.Attachment{}, nil
	}

	// Multipart message: parse parts recursively
	mr := multipart.NewReader(msg.Body, params["boundary"])
	var results []parsemail.Attachment

	err = walkParts(mr, &results)
	if err != nil {
		return nil, err
	}

	return results, nil
}

// walkParts recursively walks MIME multipart parts and extracts PGP keys.
func walkParts(mr *multipart.Reader, results *[]parsemail.Attachment) error {
	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		ct := part.Header.Get("Content-Type")
		mediaType, params, err := mime.ParseMediaType(ct)
		if err != nil {
			mediaType = "application/octet-stream" // fallback
		}

		if strings.HasPrefix(mediaType, "multipart/") {
			// Recursive multipart part
			nestedMR := multipart.NewReader(part, params["boundary"])
			if err := walkParts(nestedMR, results); err != nil {
				return err
			}
			continue
		}

		// Check if part is a PGP key attachment
		if mediaType == "application/pgp-keys" {
			data, err := decodePart(part)
			if err != nil {
				return err
			}

			filename := getFilename(part.Header)
			if filename == "" {
				filename = "publickey.asc"
			}

			*results = append(*results, parsemail.Attachment{
				Filename:    filename,
				ContentType: mediaType,
				Data:        bytes.NewReader(data),
			})
		}
	}

	return nil
}

// decodePart reads the part data and applies Content-Transfer-Encoding decoding if needed.
func decodePart(part *multipart.Part) ([]byte, error) {
	cte := strings.ToLower(part.Header.Get("Content-Transfer-Encoding"))
	var reader io.Reader = part
	switch cte {
	case "base64":
		reader = base64.NewDecoder(base64.StdEncoding, part)
	case "quoted-printable":
		// Use the standard library quoted-printable decoder
		reader = quotedprintable.NewReader(part)
	default:
		// no decoding
	}

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// getFilename extracts the filename from Content-Disposition or Content-Type headers.
func getFilename(header any) string {
	var cd, ct string
	switch h := header.(type) {
	case mail.Header:
		cd = h.Get("Content-Disposition")
		ct = h.Get("Content-Type")
	case textproto.MIMEHeader:
		cd = h.Get("Content-Disposition")
		ct = h.Get("Content-Type")
	default:
		return ""
	}

	_, cdParams, err := mime.ParseMediaType(cd)
	if err == nil {
		if filename := cdParams["filename"]; filename != "" {
			return filename
		}
	}

	_, ctParams, err := mime.ParseMediaType(ct)
	if err == nil {
		if name := ctParams["name"]; name != "" {
			return name
		}
	}

	return ""

}

// ExtractTextBody extracts the plain text body from an email.
func ExtractTextBody(data []byte) (string, error) {
	return extractBodyRecursive(data, "text/plain")
}

// ExtractHTMLBody extracts the HTML body from an email.
func ExtractHTMLBody(data []byte) (string, error) {
	return extractBodyRecursive(data, "text/html")
}

// extractBodyRecursive recursively walks MIME parts to find the first matching content type.
func extractBodyRecursive(data []byte, contentType string) (string, error) {
	msg, err := mail.ReadMessage(bytes.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("failed to parse email: %w", err)
	}

	ctHeader := msg.Header.Get("Content-Type")
	mediaType, params, err := mime.ParseMediaType(ctHeader)
	if err != nil {
		return "", fmt.Errorf("failed to parse Content-Type: %w", err)
	}

	body, err := io.ReadAll(msg.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read message body: %w", err)
	}

	return walkMIMEParts(mediaType, params, body, contentType)
}

// walkMIMEParts recursively searches for a part of the desired content type.
func walkMIMEParts(mediaType string, params map[string]string, body []byte, targetType string) (string, error) {
	if strings.HasPrefix(mediaType, "multipart/") {
		boundary := params["boundary"]
		mr := multipart.NewReader(bytes.NewReader(body), boundary)
		for {
			p, err := mr.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				return "", fmt.Errorf("error reading multipart: %w", err)
			}

			pMediaType, pParams, _ := mime.ParseMediaType(p.Header.Get("Content-Type"))
			encoding := strings.ToLower(p.Header.Get("Content-Transfer-Encoding"))
			decodedBody, err := decodeBody(p, encoding)
			if err != nil {
				return "", fmt.Errorf("error decoding part body: %w", err)
			}

			result, err := walkMIMEParts(pMediaType, pParams, decodedBody, targetType)
			if err == nil && result != "" {
				return result, nil
			}
		}
	} else if mediaType == targetType {
		return string(body), nil
	}

	return "", fmt.Errorf("%s not found", targetType)
}

func decodeBody(r io.Reader, encoding string) ([]byte, error) {
	switch encoding {
	case "base64":
		return io.ReadAll(base64.NewDecoder(base64.StdEncoding, r))
	case "quoted-printable":
		return io.ReadAll(quotedprintable.NewReader(r))
	case "", "7bit", "8bit", "binary":
		return io.ReadAll(r)
	default:
		return nil, fmt.Errorf("unsupported encoding: %s", encoding)
	}
}
