package model

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"net/mail"
	"regexp"
	"strings"
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
	// Extract DKIM domains using our custom function
	dkimDomains, err := extractDKIMDomains(data)
	if err != nil {
		return Msg{}, fmt.Errorf("failed to extract DKIM domains: %w", err)
	}

	// Check if any DKIM domains were found
	if len(dkimDomains) == 0 {
		return Msg{}, errors.New("no valid DKIM signature found")
	}

	// Use the first DKIM domain found
	dkimDomain := dkimDomains[0]

	// Parse message headers and body
	msg, err := mail.ReadMessage(bytes.NewReader(data))
	if err != nil {
		return Msg{}, err
	}

	// Subject
	subject := msg.Header.Get("Subject")

	// Parse To
	to := make([]string, 0)
	for _, t := range strings.Split(msg.Header.Get("To"), ",") {
		address, err := mail.ParseAddress(t)
		if err != nil {
			return Msg{}, err
		}

		to = append(to, address.Address)
	}

	// Parse From header
	from, err := mail.ParseAddress(msg.Header.Get("From"))
	if err != nil {
		return Msg{}, err
	}

	fromAddr := from.Address
	fromDomain := strings.ToLower(strings.SplitN(fromAddr, "@", 2)[1])

	// Validate DKIM domain vs From domain
	if !strings.EqualFold(fromDomain, dkimDomain) {
		return Msg{}, fmt.Errorf("DKIM domain mismatch: From domain '%s' does not match DKIM domain '%s'", fromDomain, dkimDomain)
	}

	// Read body
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(msg.Body)
	if err != nil {
		return Msg{}, err
	}
	body := buf.String()
	msgType := Send

	// Determine message type
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

func extractDKIMDomains(data []byte) ([]string, error) {
	var domains []string

	// Use scanner to read headers line by line
	scanner := bufio.NewScanner(bytes.NewReader(data))
	var currentHeader string
	var dkimHeaders []string

	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			break // End of headers
		}

		if strings.HasPrefix(line, " ") || strings.HasPrefix(line, "\t") {
			// Header folding
			currentHeader += line
		} else {
			// New header
			if strings.HasPrefix(currentHeader, "DKIM-Signature:") {
				dkimHeaders = append(dkimHeaders, strings.TrimPrefix(currentHeader, "DKIM-Signature:"))
			}
			currentHeader = line
		}
	}

	// Add last DKIM-Signature if needed
	if strings.HasPrefix(currentHeader, "DKIM-Signature:") {
		dkimHeaders = append(dkimHeaders, strings.TrimPrefix(currentHeader, "DKIM-Signature:"))
	}

	// Extract domains
	re := regexp.MustCompile(`\bd=([^;\s]+)`)
	for _, hdr := range dkimHeaders {
		if match := re.FindStringSubmatch(hdr); len(match) == 2 {
			domains = append(domains, match[1])
		}
	}

	return domains, scanner.Err()
}
