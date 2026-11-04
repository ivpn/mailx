package utils

import (
	"strings"
	"testing"
)

// FuzzVerifyEmailAuth fuzzes parsing of the Authentication-Results /
// ARC-Authentication-Results headers, which a remote mail server fully
// controls the content of (they are attacker-influenced: any relay ahead of
// the trust boundary, or a directly-connecting malicious MTA if the field is
// not stripped/re-stamped, can set them).
func FuzzVerifyEmailAuth(f *testing.F) {
	seeds := []string{
		"From: a@example.com\r\nAuthentication-Results: mx.example.com; dkim=pass header.d=example.com; spf=pass smtp.mailfrom=a@example.com; dmarc=pass header.from=example.com\r\n\r\nBody.",
		"From: a@example.com\r\nAuthentication-Results: mx.example.com; dkim=fail header.d=evil.com\r\n\r\nBody.",
		"From: a@example.com\r\nARC-Authentication-Results: i=1; mx.example.com; spf=pass smtp.mailfrom=a@example.com\r\n\r\nBody.",
		// No Authentication-Results at all.
		"From: a@example.com\r\nMessage-ID: <id@example.com>\r\n\r\nBody.",
		// Malformed From.
		"From: not-an-address\r\nAuthentication-Results: mx; dkim=pass header.d=example.com\r\n\r\nBody.",
		// Domain mismatch.
		"From: a@example.com\r\nAuthentication-Results: mx; dkim=pass header.d=attacker.com; spf=pass smtp.mailfrom=a@attacker.com; dmarc=pass header.from=attacker.com\r\n\r\nBody.",
		// header.d value with no terminator (runs to end of string).
		"From: a@example.com\r\nAuthentication-Results: mx; dkim=pass header.d=example.com",
		// Quoted / weird domain values.
		"From: \"weird\"@example.com\r\nAuthentication-Results: mx; dkim=pass header.d=\"example.com\"\r\n\r\nBody.",
		// Multiple Authentication-Results-like tokens crammed together.
		"From: a@example.com\r\nAuthentication-Results: mx; dkim=pass header.d=a.com; dkim=pass header.d=b.com; spf=pass smtp.mailfrom=c@c.com; dmarc=pass header.from=d.com\r\n\r\nBody.",
		// Non-UTF8 bytes in header.
		"From: a@example.com\r\nAuthentication-Results: mx; dkim=pass header.d=\xff\xfe\r\n\r\nBody.",
		// Extremely nested/garbage separators.
		"From: a@example.com\r\nAuthentication-Results: ;;;;dkim=pass;;;header.d=;;;\r\n\r\nBody.",
		// Empty From domain (bare '@').
		"From: @\r\nAuthentication-Results: mx; spf=pass smtp.mailfrom=@\r\n\r\nBody.",
		// Unbalanced comments and quotes.
		"From: a@example.com\r\nAuthentication-Results: mx.example.com; dmarc=fail (policy=reject header.from=\"example.com\r\n\r\nBody.",
		// ARC chain with a trusted sealer.
		"From: a@example.com\r\nAuthentication-Results: mx.example.com; dmarc=fail header.from=example.com; arc=pass\r\nARC-Seal: i=1; d=google.com; b=x\r\nARC-Authentication-Results: i=1; mx.google.com; dmarc=pass header.from=example.com\r\n\r\nBody.",
	}
	for _, s := range seeds {
		f.Add([]byte(s), "mx.example.com,mx2.example.com", "google.com,microsoft.com")
	}

	// f.Fuzz only supports a fixed set of primitive argument types, so the
	// []string AuthServIDs and ARCTrustedSealers config fields are passed as
	// comma-separated strings and split before use.
	f.Fuzz(func(t *testing.T, data []byte, servIDsCSV string, sealersCSV string) {
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("VerifyEmailAuth panicked on input (len=%d): %v", len(data), r)
			}
		}()
		var sealers []string
		if sealersCSV != "" {
			sealers = strings.Split(sealersCSV, ",")
		}
		cfg := EmailAuthConfig{AuthServIDs: strings.Split(servIDsCSV, ","), ARCTrustedSealers: sealers}
		_, _ = VerifyEmailAuth(data, cfg)
		_, _ = VerifyForwardAuth(data, cfg)
	})
}

// FuzzParseAuthResults fuzzes the lower-level Authentication-Results value
// parser directly (comment, quoted string and property splitting), isolating
// it from the mail.ReadMessage envelope so odd separator/edge-case inputs are
// explored more efficiently than via the full header fuzz above.
func FuzzParseAuthResults(f *testing.F) {
	seeds := []string{
		"dkim=pass header.d=example.com",
		"spf=pass smtp.mailfrom=a@example.com",
		"dmarc=pass header.from=example.com",
		"dkim=pass header.d=",
		"dkim=pass header.d=\r\nspf=pass smtp.mailfrom=",
		"",
		"dkim=passheader.d=nsep",
		"DKIM=PASS HEADER.D=EXAMPLE.COM",
		"i=1; mx; dmarc=fail (policy=reject) reason=\"a; (b\" header.from=x",
		"mx; spf=pass (((unbalanced) smtp.mailfrom=a@b\\",
		// Sender-controlled quoted value inside an Rspamd comment.
		"mx; spf=pass (mx: domain of \"a)b;dmarc=pass header.from=x;(c@d\" designates 1.2.3.4) smtp.mailfrom=\"a)b;(c@d\"; dmarc=fail",
		"mx; spf=pass (mx: domain of \"a\\\\\"b);dmarc=pass\" designates) smtp.mailfrom=x",
	}
	for _, s := range seeds {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, header string) {
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("parseAuthResults panicked on input %q: %v", header, r)
			}
		}()
		_ = parseAuthResults(header)
	})
}
