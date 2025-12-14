package utils

import (
	"bytes"
	"errors"
	"net/mail"
	"strings"
)

type AuthResults struct {
	DKIM        string
	SPF         string
	DMARC       string
	DKIMDomain  string
	SPFDomain   string
	DMARCDomain string
}

func VerifyEmailAuth(data []byte) (bool, error) {
	msg, err := mail.ReadMessage(bytes.NewReader(data))
	if err != nil {
		return false, err
	}

	headers := msg.Header
	authResults := []string{}

	if ar := headers.Get("Authentication-Results"); ar != "" {
		authResults = append(authResults, ar)
	}

	if arcAr := headers.Get("ARC-Authentication-Results"); arcAr != "" {
		authResults = append(authResults, arcAr)
	}

	if len(authResults) == 0 {
		msgID := headers.Get("Message-ID")
		return false, errors.New("no Authentication-Results headers found, Message-ID: " + msgID)
	}

	parsed := parseAuthResults(authResults)

	// fromAddr, err := mail.ParseAddress(headers.Get("From"))
	// if err != nil {
	// 	return false, err
	// }
	// fromDomain := extractDomain(fromAddr.Address)

	// Check for domain mismatches when authentication method is present
	// TODO: Re-enable DKIM domain check if needed, this check prevents real case scenarios when 3rd party services sending email on behalf of some domain
	// if parsed.DKIM != "" && !relaxedMatch(fromDomain, parsed.DKIMDomain) {
	// 	return false, errors.New("DKIM domain mismatch, fromDomain: " + fromDomain + ", DKIM domain: " + parsed.DKIMDomain)
	// }

	// if parsed.SPF != "" && !relaxedMatch(fromDomain, parsed.SPFDomain) {
	// 	return false, errors.New("SPF domain mismatch, fromDomain: " + fromDomain + ", SPF domain: " + parsed.SPFDomain)
	// }

	// if parsed.DMARC != "" && !relaxedMatch(fromDomain, parsed.DMARCDomain) {
	// 	return false, errors.New("DMARC domain mismatch, fromDomain: " + fromDomain + ", DMARC domain: " + parsed.DMARCDomain)
	// }

	// Continue with original verification checks
	switch {
	case parsed.DMARC == "pass":
		return true, nil
	case parsed.DKIM == "pass":
		return true, nil
	case parsed.SPF == "pass":
		return true, nil
	default:
		return false, nil
	}
}

func parseAuthResults(headers []string) AuthResults {
	result := AuthResults{}
	for _, header := range headers {
		h := strings.ToLower(header)

		if strings.Contains(h, "dkim=pass") {
			result.DKIM = "pass"
			if _, after, ok := strings.Cut(h, "header.d="); ok {
				domain := extractValue(after)
				result.DKIMDomain = domain
			}
		}
		if strings.Contains(h, "spf=pass") {
			result.SPF = "pass"
			if _, after, ok := strings.Cut(h, "smtp.mailfrom="); ok {
				domain := extractDomain(extractValue(after))
				result.SPFDomain = domain
			}
		}
		if strings.Contains(h, "dmarc=pass") {
			result.DMARC = "pass"
			if _, after, ok := strings.Cut(h, "header.from="); ok {
				domain := extractValue(after)
				result.DMARCDomain = domain
			}
		}
	}
	return result
}

func extractValue(s string) string {
	end := strings.IndexAny(s, " ;\n\r")
	if end != -1 {
		return s[:end]
	}
	return s
}

func extractDomain(email string) string {
	// Trim whitespace
	email = strings.TrimSpace(email)

	// Trim surrounding quotes if present
	if len(email) > 1 && email[0] == '"' && email[len(email)-1] == '"' {
		email = email[1 : len(email)-1]
	}

	// Find the last '@' character and return the domain part
	if at := strings.LastIndex(email, "@"); at != -1 {
		return email[at+1:]
	}

	return email
}

func relaxedMatch(fromDomain, authDomain string) bool {
	if fromDomain == "" || authDomain == "" {
		return false
	}

	return strings.HasSuffix(fromDomain, authDomain) || strings.HasSuffix(authDomain, fromDomain)
}
