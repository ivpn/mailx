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

// lookupTXTRecords looks up TXT records for host. Not-found/permanent DNS
// errors are treated as a plain empty result rather than an error.
func lookupTXTRecords(host string) ([]string, error) {
	records, err := net.LookupTXT(host)
	if err != nil {
		var dnsErr *net.DNSError
		if errors.As(err, &dnsErr) && (dnsErr.IsNotFound || (!dnsErr.IsTimeout && !dnsErr.IsTemporary)) {
			return nil, nil
		}
		return nil, ErrLookupTXT
	}
	return records, nil
}

// LookupTXTExact looks up TXT records for host and returns true if any record
// is an exact match to value (trailing dots stripped before comparison).
//
// Example use: verify ownership TXT record
//
//	LookupTXTExact("example.com", "service-verify=9487e243822f333d782eabe1115302643b222ef55072c8e77abf75335950a61a")
func LookupTXTExact(host, value string) (bool, error) {
	records, err := lookupTXTRecords(host)
	if err != nil {
		return false, err
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
	records, err := lookupTXTRecords(host)
	if err != nil {
		return false, err
	}

	want := stripDot(value)
	for _, r := range records {
		if strings.Contains(stripDot(r), want) {
			return true, nil
		}
	}
	return false, nil
}

// validSPFRecord reports whether record is a v=spf1 record that authorizes
// requiredMechanism (an "include:" target, e.g. "spf.example.net", or the bare
// "mx" mechanism) and terminates in a "-all"/"~all" mechanism. Mechanism order
// and any additional mechanisms present are ignored.
func validSPFRecord(record, requiredMechanism string) bool {
	fields := strings.Fields(record)
	if len(fields) == 0 || !strings.EqualFold(fields[0], "v=spf1") {
		return false
	}

	wantInclude := "include:" + strings.ToLower(requiredMechanism)
	hasMechanism := false
	hasAll := false
	for _, f := range fields[1:] {
		switch strings.ToLower(f) {
		case wantInclude, "mx":
			hasMechanism = true
		case "-all", "~all":
			hasAll = true
		}
	}

	return hasMechanism && hasAll
}

// validDMARCRecord reports whether record is a v=DMARC1 record with a
// p=quarantine or p=reject policy tag, regardless of tag order.
func validDMARCRecord(record string) bool {
	values := make(map[string]string)
	for tag := range strings.SplitSeq(record, ";") {
		tag = strings.TrimSpace(tag)
		key, value, ok := strings.Cut(tag, "=")
		if !ok {
			continue
		}
		values[strings.ToLower(strings.TrimSpace(key))] = strings.ToLower(strings.TrimSpace(value))
	}

	if values["v"] != "dmarc1" {
		return false
	}
	return values["p"] == "quarantine" || values["p"] == "reject"
}

// LookupSPF looks up the SPF TXT record for host and returns true if it
// authorizes requiredMechanism (an "include:" target or the bare "mx"
// mechanism) and ends in a "-all"/"~all" mechanism, regardless of mechanism
// order or additional mechanisms present.
//
// Example use:
//
//	LookupSPF("example.com", "spf.example.net")
func LookupSPF(host, requiredMechanism string) (bool, error) {
	records, err := lookupTXTRecords(host)
	if err != nil {
		return false, err
	}

	for _, r := range records {
		if validSPFRecord(stripDot(r), requiredMechanism) {
			return true, nil
		}
	}
	return false, nil
}

// LookupDMARC looks up the DMARC TXT record for host (typically
// "_dmarc."+domain) and returns true if it has a p=quarantine or p=reject
// policy, regardless of tag order.
//
// Example use:
//
//	LookupDMARC("_dmarc.example.com")
func LookupDMARC(host string) (bool, error) {
	records, err := lookupTXTRecords(host)
	if err != nil {
		return false, err
	}

	for _, r := range records {
		if validDMARCRecord(stripDot(r)) {
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
