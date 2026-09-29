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

type EmailAuthConfig struct {
	TrustedRelayDomains []string
	DisableAlignment    bool
	LogMismatch         bool
}

func VerifyEmailAuth(data []byte, cfg EmailAuthConfig) (bool, error) {
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

	fromAddr, err := mail.ParseAddress(headers.Get("From"))
	if err != nil {
		return false, err
	}
	fromDomain := extractDomain(fromAddr.Address)

	reject, mismatchErr := checkAlignment(fromDomain, parsed, cfg)
	if reject {
		return false, mismatchErr
	}

	switch {
	case parsed.DMARC == "pass":
		return true, mismatchErr
	case parsed.DKIM == "pass":
		return true, mismatchErr
	case parsed.SPF == "pass":
		return true, mismatchErr
	default:
		return false, mismatchErr
	}
}

// checkAlignment compares each present DKIM/SPF/DMARC result's authenticated domain against
// fromDomain. It reports whether the message should be rejected under cfg's enforcement
// settings, plus a diagnostic error describing any mismatch found (nil unless a mismatch
// occurred and either it caused rejection or cfg.LogMismatch requested visibility into
// mismatches that were bypassed via TrustedRelayDomains or cfg.DisableAlignment).
func checkAlignment(fromDomain string, parsed AuthResults, cfg EmailAuthConfig) (bool, error) {
	mechanisms := []struct {
		name       string
		result     string
		authDomain string
	}{
		{"DKIM", parsed.DKIM, parsed.DKIMDomain},
		{"SPF", parsed.SPF, parsed.SPFDomain},
		{"DMARC", parsed.DMARC, parsed.DMARCDomain},
	}

	var reject bool
	var errs []error

	for _, m := range mechanisms {
		if m.result == "" || relaxedMatch(fromDomain, m.authDomain) {
			continue
		}

		trusted := isTrustedRelayDomain(m.authDomain, cfg.TrustedRelayDomains)
		mechReject := !trusted && !cfg.DisableAlignment
		if mechReject {
			reject = true
		}

		if mechReject || cfg.LogMismatch {
			errs = append(errs, errors.New(m.name+" domain mismatch, fromDomain: "+fromDomain+", "+m.name+" domain: "+m.authDomain))
		}
	}

	return reject, errors.Join(errs...)
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

	fromDomain = strings.ToLower(fromDomain)
	authDomain = strings.ToLower(authDomain)

	if fromDomain == authDomain {
		return true
	}

	// Require a label boundary so e.g. "notexample.com" does not match "example.com".
	return strings.HasSuffix(fromDomain, "."+authDomain) || strings.HasSuffix(authDomain, "."+fromDomain)
}

// isTrustedRelayDomain reports whether domain is (or is a subdomain of) one of the
// operator-configured relay signing domains explicitly trusted to send authenticated
// mail on behalf of other From: domains (e.g. transactional email providers).
func isTrustedRelayDomain(domain string, trustedRelayDomains []string) bool {
	if domain == "" {
		return false
	}

	for _, trusted := range trustedRelayDomains {
		trusted = strings.TrimSpace(trusted)
		if trusted != "" && relaxedMatch(domain, trusted) {
			return true
		}
	}

	return false
}
