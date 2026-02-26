package utils

import (
	"errors"
	"net"
	"strings"
)

var (
	ErrLookupTXT   = errors.New("failed to lookup TXT records")
	ErrLookupMX    = errors.New("failed to lookup MX records")
	ErrLookupCNAME = errors.New("failed to lookup CNAME record")
)

// stripDot removes a trailing dot from a DNS hostname.
func stripDot(s string) string {
	return strings.TrimSuffix(s, ".")
}

// LookupTXTExact looks up TXT records for host and returns true if any record
// is an exact match to value (trailing dots stripped before comparison).
//
// Example use: verify ownership TXT record
//
//	LookupTXTExact("example.com", "service-verify=9487e243822f333d782eabe1115302643b222ef55072c8e77abf75335950a61a")
func LookupTXTExact(host, value string) (bool, error) {
	records, err := net.LookupTXT(host)
	if err != nil {
		var dnsErr *net.DNSError
		if errors.As(err, &dnsErr) && (dnsErr.IsNotFound || (!dnsErr.IsTimeout && !dnsErr.IsTemporary)) {
			return false, nil
		}
		return false, ErrLookupTXT
	}

	want := stripDot(value)
	for _, r := range records {
		if stripDot(r) == want {
			return true, nil
		}
	}
	return false, nil
}

// LookupTXTContains looks up TXT records for host and returns true if any record
// contains value as a substring (trailing dots stripped before comparison).
//
// Example uses:
//
//	LookupTXTContains("example.com", "v=spf1 include:spf.example.net -all")
//	LookupTXTContains("_dmarc.example.com", "v=DMARC1; p=quarantine; adkim=s")
func LookupTXTContains(host, value string) (bool, error) {
	records, err := net.LookupTXT(host)
	if err != nil {
		var dnsErr *net.DNSError
		if errors.As(err, &dnsErr) && (dnsErr.IsNotFound || (!dnsErr.IsTimeout && !dnsErr.IsTemporary)) {
			return false, nil
		}
		return false, ErrLookupTXT
	}

	want := stripDot(value)
	for _, r := range records {
		if strings.Contains(stripDot(r), want) {
			return true, nil
		}
	}
	return false, nil
}

// LookupMX looks up MX records for host and returns true if any MX entry's
// hostname matches target (trailing dots stripped, case-insensitive).
// The MX priority/preference value is not checked.
//
// Example use:
//
//	LookupMX("example.com", "mail1.example.net.")
func LookupMX(host, target string) (bool, error) {
	records, err := net.LookupMX(host)
	if err != nil {
		var dnsErr *net.DNSError
		if errors.As(err, &dnsErr) && (dnsErr.IsNotFound || (!dnsErr.IsTimeout && !dnsErr.IsTemporary)) {
			return false, nil
		}
		return false, ErrLookupMX
	}

	want := strings.ToLower(stripDot(target))
	for _, r := range records {
		if strings.ToLower(stripDot(r.Host)) == want {
			return true, nil
		}
	}
	return false, nil
}

// LookupCNAME looks up the canonical name for host and returns true if it matches
// target (trailing dots stripped, case-insensitive).
//
// Example use:
//
//	LookupCNAME("mail._domainkey.example.com", "mail._domainkey.example.net.")
func LookupCNAME(host, target string) (bool, error) {
	cname, err := net.LookupCNAME(host)
	if err != nil {
		var dnsErr *net.DNSError
		if errors.As(err, &dnsErr) && (dnsErr.IsNotFound || (!dnsErr.IsTimeout && !dnsErr.IsTemporary)) {
			return false, nil
		}
		return false, ErrLookupCNAME
	}

	return strings.EqualFold(stripDot(cname), stripDot(target)), nil
}
