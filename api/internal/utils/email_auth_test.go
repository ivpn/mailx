package utils

import (
	"errors"
	"reflect"
	"testing"
)

func TestRelaxedMatch(t *testing.T) {
	testCases := []struct {
		fromDomain string
		authDomain string
		expected   bool
		name       string
	}{
		{
			fromDomain: "example.com",
			authDomain: "example.com",
			expected:   true,
			name:       "exact match",
		},
		{
			fromDomain: "sub.example.com",
			authDomain: "example.com",
			expected:   true,
			name:       "subdomain match",
		},
		{
			fromDomain: "example.net",
			authDomain: "example.com",
			expected:   false,
			name:       "domain mismatch",
		},
		{
			fromDomain: "example.com",
			authDomain: "example.com.net",
			expected:   false,
			name:       "auth domain longer",
		},
		{
			fromDomain: "test@example.com",
			authDomain: "example.com",
			expected:   false,
			name:       "email address is not a bare domain",
		},
		{
			fromDomain: "notexample.com",
			authDomain: "example.com",
			expected:   false,
			name:       "lookalike domain without label boundary must not match",
		},
		{
			fromDomain: "example.com.evil.com",
			authDomain: "example.com",
			expected:   false,
			name:       "victim domain used as prefix must not match",
		},
		{
			fromDomain: "Sub.EXAMPLE.com",
			authDomain: "example.COM",
			expected:   true,
			name:       "case-insensitive subdomain match",
		},
		{
			fromDomain: "",
			authDomain: "example.com",
			expected:   false,
			name:       "empty from domain",
		},
		{
			fromDomain: "example.com",
			authDomain: "",
			expected:   false,
			name:       "empty auth domain",
		},
		{
			fromDomain: "",
			authDomain: "",
			expected:   false,
			name:       "both empty",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := relaxedMatch(tc.fromDomain, tc.authDomain)
			if actual != tc.expected {
				t.Errorf("relaxedMatch(%q, %q) = %v, expected %v", tc.fromDomain, tc.authDomain, actual, tc.expected)
			}
		})
	}
}

func TestExtractDomain(t *testing.T) {
	testCases := []struct {
		email    string
		expected string
		name     string
	}{
		{
			email:    "test@example.com",
			expected: "example.com",
			name:     "valid email",
		},
		{
			email:    "test.user@sub.example.com",
			expected: "sub.example.com",
			name:     "subdomain email",
		},
		{
			email:    "test@example.co.uk",
			expected: "example.co.uk",
			name:     "co.uk email",
		},
		{
			email:    "test",
			expected: "test",
			name:     "no @ symbol",
		},
		{
			email:    "",
			expected: "",
			name:     "empty string",
		},
		{
			email:    "@example.com",
			expected: "example.com",
			name:     "leading @",
		},
		{
			email:    "test@",
			expected: "",
			name:     "trailing @",
		},
		{
			email:    "multiple@at@symbols",
			expected: "symbols",
			name:     "multiple @ symbols",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := extractDomain(tc.email)
			if actual != tc.expected {
				t.Errorf("extractDomain(%q) = %q, expected %q", tc.email, actual, tc.expected)
			}
		})
	}
}

func TestIsTrustedSealer(t *testing.T) {
	trusted := []string{"google.com", " Microsoft.com "}

	testCases := []struct {
		domain   string
		expected bool
		name     string
	}{
		{domain: "google.com", expected: true, name: "exact match"},
		{domain: "arc.google.com", expected: true, name: "subdomain match"},
		{domain: "MICROSOFT.com", expected: true, name: "case-insensitive with trimmed config"},
		{domain: "evilgoogle.com", expected: false, name: "lookalike domain without label boundary"},
		{domain: "com", expected: false, name: "parent of trusted domain"},
		{domain: "google.com.evil.com", expected: false, name: "trusted domain used as prefix"},
		{domain: "", expected: false, name: "empty domain"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := isTrustedSealer(tc.domain, trusted)
			if actual != tc.expected {
				t.Errorf("isTrustedSealer(%q) = %v, expected %v", tc.domain, actual, tc.expected)
			}
		})
	}
}

func TestLatestARCSeal(t *testing.T) {
	testCases := []struct {
		seals            []string
		expectedInstance int
		expectedDomain   string
		name             string
	}{
		{
			seals:            nil,
			expectedInstance: 0,
			expectedDomain:   "",
			name:             "no seals",
		},
		{
			seals:            []string{"i=1; a=rsa-sha256; t=1; cv=none; d=Google.com; s=arc-20240605; b=abc"},
			expectedInstance: 1,
			expectedDomain:   "google.com",
			name:             "single seal",
		},
		{
			seals: []string{
				"i=2; a=rsa-sha256; cv=pass; d=microsoft.com; s=arc; b=abc",
				"i=1; a=rsa-sha256; cv=none; d=google.com; s=arc; b=abc",
			},
			expectedInstance: 2,
			expectedDomain:   "microsoft.com",
			name:             "highest instance wins",
		},
		{
			seals: []string{
				"i=1; d=google.com; b=abc",
				"i=1; d=attacker.example; b=abc",
			},
			expectedInstance: 0,
			expectedDomain:   "",
			name:             "duplicate highest instance is ambiguous",
		},
		{
			seals:            []string{"i=x; d=google.com", "d=google.com", "i=0; d=google.com"},
			expectedInstance: 0,
			expectedDomain:   "",
			name:             "invalid instances are ignored",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			instance, domain := latestARCSeal(tc.seals)
			if instance != tc.expectedInstance || domain != tc.expectedDomain {
				t.Errorf("latestARCSeal(%q) = (%d, %q), expected (%d, %q)", tc.seals, instance, domain, tc.expectedInstance, tc.expectedDomain)
			}
		})
	}
}

func TestParseAuthResults(t *testing.T) {
	testCases := []struct {
		header   string
		expected authResults
		name     string
	}{
		{
			header: "mx.example.net; dkim=pass header.d=example.com header.s=sel; spf=pass smtp.mailfrom=a@example.com; dmarc=pass header.from=example.com",
			expected: authResults{
				servID: "mx.example.net",
				results: []resInfo{
					{method: "dkim", result: "pass", props: map[string]string{"header.d": "example.com", "header.s": "sel"}},
					{method: "spf", result: "pass", props: map[string]string{"smtp.mailfrom": "a@example.com"}},
					{method: "dmarc", result: "pass", props: map[string]string{"header.from": "example.com"}},
				},
			},
			name: "basic",
		},
		{
			header: "MX.Example.NET;\r\n\tDKIM=Pass header.d=EXAMPLE.com;\r\n\tdmarc=fail reason=\"No valid SPF; No valid DKIM\" header.from=example.com (policy=reject);\r\n\tspf=pass (mx.example.net: domain of a@example.com designates 192.0.2.1 as permitted sender; (nested)) smtp.mailfrom=a@example.com",
			expected: authResults{
				servID: "mx.example.net",
				results: []resInfo{
					{method: "dkim", result: "pass", props: map[string]string{"header.d": "example.com"}},
					{method: "dmarc", result: "fail", props: map[string]string{"reason": "no valid spf; no valid dkim", "header.from": "example.com"}, comments: " policy=reject"},
					{method: "spf", result: "pass", props: map[string]string{"smtp.mailfrom": "a@example.com"}, comments: " mx.example.net: domain of a@example.com designates 192.0.2.1 as permitted sender; (nested)"},
				},
			},
			name: "rspamd style with comments, quoted strings and folding",
		},
		{
			header: "i=2; mx.google.com; arc=pass (i=1); dkim=pass header.i=@example.com",
			expected: authResults{
				servID:   "mx.google.com",
				instance: 2,
				results: []resInfo{
					{method: "arc", result: "pass", props: map[string]string{}, comments: " i=1"},
					{method: "dkim", result: "pass", props: map[string]string{"header.i": "@example.com"}},
				},
			},
			name: "ARC instance tag",
		},
		{
			header: `mx.example.net; spf=pass (mx.example.net: domain of "a)b;dmarc=pass header.from=victim.example;(c@evil.example" designates 192.0.2.1 as permitted sender) smtp.mailfrom="a)b;dmarc=pass header.from=victim.example;(c@evil.example"; dmarc=fail header.from=victim.example`,
			expected: authResults{
				servID: "mx.example.net",
				results: []resInfo{
					{method: "spf", result: "pass", props: map[string]string{"smtp.mailfrom": "a)b;dmarc=pass header.from=victim.example;(c@evil.example"}, comments: ` mx.example.net: domain of "a)b;dmarc=pass header.from=victim.example;(c@evil.example" designates 192.0.2.1 as permitted sender`},
					{method: "dmarc", result: "fail", props: map[string]string{"header.from": "victim.example"}},
				},
			},
			name: "quoted string inside a comment",
		},
		{
			header: `mx.example.net; spf=pass (domain of "a\"b" designates) smtp.mailfrom="a\"b"`,
			expected: authResults{
				servID: "mx.example.net",
				results: []resInfo{
					{method: "spf", result: "pass", props: map[string]string{"smtp.mailfrom": `a"b`}, comments: ` domain of "a"b" designates`},
				},
			},
			name: "escaped quote inside a quoted string",
		},
		{
			header: `mx.example.net; spf=pass "unterminated`,
			expected: authResults{
				servID:    "mx.example.net",
				results:   []resInfo{{method: "spf", result: "pass", props: map[string]string{}}},
				malformed: "unterminated quoted string or comment",
			},
			name: "unterminated quoted string",
		},
		{
			header: "mx.example.net; spf=pass (unterminated smtp.mailfrom=x",
			expected: authResults{
				servID:    "mx.example.net",
				results:   []resInfo{{method: "spf", result: "pass", props: map[string]string{}, comments: " unterminated smtp.mailfrom=x"}},
				malformed: "unterminated quoted string or comment",
			},
			name: "unterminated comment",
		},
		{
			header: "mx.example.net; dkim=none; garbage; dmarc=none",
			expected: authResults{
				servID: "mx.example.net",
				results: []resInfo{
					{method: "dkim", result: "none", props: map[string]string{}},
					{method: "dmarc", result: "none", props: map[string]string{}},
				},
				malformed: `result without method=value: "garbage"`,
			},
			name: "result without method=value",
		},
		{
			header:   "mx.example.net; none",
			expected: authResults{servID: "mx.example.net"},
			name:     "no results",
		},
		{
			header:   "",
			expected: authResults{},
			name:     "empty header",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := parseAuthResults(tc.header)
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("parseAuthResults(%q) = %+v, expected %+v", tc.header, actual, tc.expected)
			}
		})
	}
}

var testServIDs = []string{"mail1.example.net", " mx.example.net "}

var testSealers = []string{"google.com", "microsoft.com"}

func TestVerifyEmailAuth(t *testing.T) {
	tests := []struct {
		name       string
		emailRaw   string
		servIDs    []string
		want       bool
		wantErrIs  error
		wantNilErr bool
	}{
		{
			name: "DMARC pass",
			emailRaw: `From: sender@example.com
To: recipient@example.net
Authentication-Results: mx.example.net; dmarc=pass header.from=example.com

This is a test email.`,
			want:       true,
			wantNilErr: true,
		},
		{
			name: "DKIM pass with aligned domain",
			emailRaw: `From: sender@example.com
To: recipient@example.net
Authentication-Results: mx.example.net; dkim=pass header.d=example.com

This is a test email.`,
			want:       true,
			wantNilErr: true,
		},
		{
			name: "SPF pass with aligned domain",
			emailRaw: `From: sender@example.com
To: recipient@example.net
Authentication-Results: mx.example.net; spf=pass smtp.mailfrom=bounce@example.com

This is a test email.`,
			want:       true,
			wantNilErr: true,
		},
		{
			name: "Rspamd style header with comments and folding",
			emailRaw: "From: Sender <sender@example.com>\r\n" +
				"To: recipient@example.net\r\n" +
				"Authentication-Results: mx.example.net;\r\n" +
				"\tdkim=pass header.d=example.com header.s=sel header.b=AbCd;\r\n" +
				"\tdmarc=pass (policy=reject) header.from=example.com;\r\n" +
				"\tspf=pass (mx.example.net: domain of sender@example.com designates 192.0.2.1 as permitted sender) smtp.mailfrom=sender@example.com\r\n" +
				"\r\nThis is a test email.",
			want:       true,
			wantNilErr: true,
		},
		{
			name: "authserv-id matched case-insensitively",
			emailRaw: `From: sender@example.com
To: recipient@example.net
Authentication-Results: MX.Example.NET; dmarc=pass header.from=example.com

This is a test email.`,
			want:       true,
			wantNilErr: true,
		},
		{
			name: "subdomain in From aligned with organizational DKIM domain",
			emailRaw: `From: sender@sub.example.com
To: recipient@example.net
Authentication-Results: mx.example.net; dkim=pass header.d=example.com

This is a test email.`,
			want:       true,
			wantNilErr: true,
		},
		{
			name: "Gmail style DKIM result with header.i only",
			emailRaw: `From: sender@example.com
To: recipient@example.net
Authentication-Results: mx.example.net; dkim=pass header.i=@example.com header.s=20230601

This is a test email.`,
			want:       true,
			wantNilErr: true,
		},
		{
			name: "email service provider: aligned DKIM, misaligned SPF",
			emailRaw: `From: news@example.com
To: recipient@example.net
Authentication-Results: mx.example.net; dkim=pass header.d=example.com; spf=pass smtp.mailfrom=bounces+123@bounces.sendgrid.net

This is a test email.`,
			want:       true,
			wantNilErr: true,
		},
		{
			name: "Google Workspace fallback DKIM signature with DMARC pass",
			emailRaw: `From: sender@example.com
To: recipient@example.net
Authentication-Results: mx.example.net; dkim=pass header.d=example-com.20230601.gappssmtp.com; spf=pass smtp.mailfrom=sender@example.com; dmarc=pass header.from=example.com

This is a test email.`,
			want:       true,
			wantNilErr: true,
		},
		{
			name: "aligned DKIM signature after a misaligned one",
			emailRaw: `From: sender@example.com
To: recipient@example.net
Authentication-Results: mx.example.net; dkim=pass header.d=relay.example; dkim=pass header.d=example.com

This is a test email.`,
			want:       true,
			wantNilErr: true,
		},
		{
			name: "DMARC pass for a different domain",
			emailRaw: `From: sender@example.com
To: recipient@example.net
Authentication-Results: mx.example.net; dmarc=pass header.from=different.org

This is a test email.`,
			want:      false,
			wantErrIs: ErrNoAlignedAuthentication,
		},
		{
			name: "DKIM pass with misaligned domain",
			emailRaw: `From: sender@example.com
To: recipient@example.net
Authentication-Results: mx.example.net; dkim=pass header.d=attacker.example; dmarc=none header.from=example.com

This is a test email.`,
			want:      false,
			wantErrIs: ErrNoAlignedAuthentication,
		},
		{
			name: "SPF pass with misaligned domain",
			emailRaw: `From: sender@example.com
To: recipient@example.net
Authentication-Results: mx.example.net; spf=pass smtp.mailfrom=a@attacker.example

This is a test email.`,
			want:      false,
			wantErrIs: ErrNoAlignedAuthentication,
		},
		{
			name: "spoof signed by sender's own domain with DMARC fail",
			emailRaw: `From: owner@example.com
To: recipient@example.net
Authentication-Results: mx.example.net; dkim=pass header.d=attacker.example; spf=pass smtp.mailfrom=a@attacker.example; dmarc=fail header.from=example.com (policy=none)

This is a test email.`,
			want:      false,
			wantErrIs: ErrDMARCFail,
		},
		{
			name: "DMARC fail decides over an aligned DKIM pass",
			emailRaw: `From: sender@example.com
To: recipient@example.net
Authentication-Results: mx.example.net; dkim=pass header.d=example.com; dmarc=fail header.from=example.com

This is a test email.`,
			want:      false,
			wantErrIs: ErrDMARCFail,
		},
		{
			name: "all auth fails",
			emailRaw: `From: sender@example.com
To: recipient@example.net
Authentication-Results: mx.example.net; dkim=fail; spf=fail

This is a test email.`,
			want:      false,
			wantErrIs: ErrNoAlignedAuthentication,
		},
		{
			name: "report reproduction: forged ARC-Authentication-Results",
			emailRaw: "Authentication-Results: mx.example.net;\r\n\tdkim=none;\r\n\tspf=softfail smtp.mailfrom=owner@victim.example;\r\n\tdmarc=fail header.from=victim.example (policy=none)\r\n" +
				"ARC-Authentication-Results: i=1; mx.google.com; dkim=pass header.d=victim.example; spf=pass smtp.mailfrom=a@victim.example; dmarc=pass header.from=victim.example\r\n" +
				"From: Victim <owner@victim.example>\r\nTo: alias@example.net\r\nSubject: forged\r\n\r\nhi\r\n",
			want:      false,
			wantErrIs: ErrDMARCFail,
		},
		{
			name: "pass from another authserv-id is ignored",
			emailRaw: `From: sender@example.com
To: recipient@example.net
Authentication-Results: mx.attacker.example; dmarc=pass header.from=example.com
Authentication-Results: mx.example.net; dmarc=fail header.from=example.com

This is a test email.`,
			want:      false,
			wantErrIs: ErrDMARCFail,
		},
		{
			name: "only another authserv-id",
			emailRaw: `From: sender@example.com
To: recipient@example.net
Authentication-Results: mx.attacker.example; dmarc=pass header.from=example.com
ARC-Authentication-Results: i=1; mx.example.net; dmarc=pass header.from=example.com

This is a test email.`,
			want:      false,
			wantErrIs: ErrNoOwnAuthResults,
		},
		{
			name: "forged second header with own authserv-id",
			emailRaw: `From: sender@example.com
To: recipient@example.net
Authentication-Results: mx.example.net; dmarc=pass header.from=example.com
Authentication-Results: mx.example.net; dmarc=fail header.from=example.com

This is a test email.`,
			want:      false,
			wantErrIs: ErrDuplicateOwnAuthResults,
		},
		{
			name: "header from another of our mail servers",
			emailRaw: `From: sender@example.com
To: recipient@example.net
Authentication-Results: mail1.example.net; dmarc=pass header.from=example.com

This is a test email.`,
			want:       true,
			wantNilErr: true,
		},
		{
			name: "headers from two of our mail servers",
			emailRaw: `From: sender@example.com
To: recipient@example.net
Authentication-Results: mail1.example.net; dmarc=pass header.from=example.com
Authentication-Results: mx.example.net; dmarc=fail header.from=example.com

This is a test email.`,
			want:      false,
			wantErrIs: ErrDuplicateOwnAuthResults,
		},
		{
			name: "empty authserv-id entry does not match a header without one",
			emailRaw: `From: sender@example.com
To: recipient@example.net
Authentication-Results: ; dmarc=pass header.from=example.com
Authentication-Results: mx.example.net; dmarc=fail header.from=example.com

This is a test email.`,
			servIDs:   []string{"mx.example.net", ""},
			want:      false,
			wantErrIs: ErrDMARCFail,
		},
		{
			name: "authserv-id not configured",
			emailRaw: `From: sender@example.com
To: recipient@example.net
Authentication-Results: mx.example.net; dmarc=pass header.from=example.com

This is a test email.`,
			servIDs:   []string{" ", ""},
			want:      false,
			wantErrIs: ErrAuthServIDNotConfigured,
		},
		{
			name: "no Authentication-Results",
			emailRaw: `From: sender@example.com
To: recipient@example.net

This is a test email.`,
			want:      false,
			wantErrIs: ErrNoOwnAuthResults,
		},
		{
			name: "no From header",
			emailRaw: `To: recipient@example.net
Authentication-Results: mx.example.net; dkim=pass header.d=example.com

This is a test email.`,
			want: false,
		},
		{
			name: "malformed From address",
			emailRaw: `From: invalid-email
To: recipient@example.net
Authentication-Results: mx.example.net; dkim=pass header.d=example.com

This is a test email.`,
			want: false,
		},
		{
			name: "invalid email format",
			emailRaw: `This is not a valid email format
Just some random text`,
			want: false,
		},
		{
			name: "Rspamd header with a quoted sender address",
			emailRaw: "From: sender@example.com\r\n" +
				"To: recipient@example.net\r\n" +
				"Authentication-Results: mx.example.net;\r\n" +
				"\tdkim=pass header.d=example.com header.s=sel header.b=AbCd;\r\n" +
				"\tspf=pass (mx.example.net: domain of \"john (smith)@example.com\" designates 192.0.2.1 as permitted sender) smtp.mailfrom=\"john (smith)@example.com\";\r\n" +
				"\tdmarc=pass (policy=reject) header.from=example.com\r\n" +
				"\r\nThis is a test email.",
			want:       true,
			wantNilErr: true,
		},
		{
			name: "malformed: backslash in own header",
			emailRaw: "From: sender@example.com\r\n" +
				"To: recipient@example.net\r\n" +
				"Authentication-Results: mx.example.net; dkim=pass header.d=example.com; spf=pass smtp.mailfrom=\"\\\"john smith\\\"@example.com\"; dmarc=pass header.from=example.com\r\n" +
				"\r\nThis is a test email.",
			want:      false,
			wantErrIs: ErrMalformedOwnAuthResults,
		},
		{
			// Rspamd escapes '"' but not '\', so DKIM b= ending in a backslash
			// and a '"' in the sender shift the quotes and expose the injection.
			name: "injection through DKIM b= ending in a backslash",
			emailRaw: "From: Victim <owner@victim.example>\r\n" +
				"To: alias@example.net\r\n" +
				"Authentication-Results: mx.example.net;\r\n" +
				"\tdkim=permerror header.d=evil.example header.s=sel header.b=\"abcdefg\\\";\r\n" +
				"\tdmarc=none;\r\n" +
				"\tspf=pass (mx.example.net: domain of \";dkim=pass header.d=victim.example;dkim=none \\\"\" designates 192.0.2.1 as permitted sender) smtp.mailfrom=\";dkim=pass header.d=victim.example;dkim=none \\\"\"\r\n" +
				"\r\nhi\r\n",
			want:      false,
			wantErrIs: ErrMalformedOwnAuthResults,
		},
		{
			name: "several DKIM results before the other results",
			emailRaw: `From: sender@example.com
To: recipient@example.net
Authentication-Results: mx.example.net; dkim=pass header.d=relay.example; dkim=pass header.d=example.com; dmarc=pass header.from=example.com; spf=pass smtp.mailfrom=a@relay.example; auth=pass smtp.auth=user

This is a test email.`,
			want:       true,
			wantNilErr: true,
		},
		{
			name: "injection through the SPF comment: dmarc=pass before the real dmarc=fail",
			emailRaw: "From: Victim <owner@victim.example>\r\n" +
				"To: alias@example.net\r\n" +
				"Authentication-Results: mx.example.net;\r\n" +
				"\tdkim=none;\r\n" +
				"\tspf=pass (mx.example.net: domain of \"a)b;dmarc=pass header.from=victim.example;(c@evil.example\" designates 192.0.2.1 as permitted sender) smtp.mailfrom=\"a)b;dmarc=pass header.from=victim.example;(c@evil.example\";\r\n" +
				"\tdmarc=fail reason=\"No valid SPF, No valid DKIM\" header.from=victim.example (policy=none)\r\n" +
				"\r\nhi\r\n",
			want:      false,
			wantErrIs: ErrDMARCFail,
		},
		{
			name: "injection through the SPF comment: aligned dkim=pass when DMARC is none",
			emailRaw: "From: Victim <owner@victim.example>\r\n" +
				"To: alias@example.net\r\n" +
				"Authentication-Results: mx.example.net;\r\n" +
				"\tdkim=none;\r\n" +
				"\tdmarc=none;\r\n" +
				"\tspf=pass (mx.example.net: domain of \"a)b;dkim=pass header.d=victim.example;dkim=none (c@evil.example\" designates 192.0.2.1 as permitted sender) smtp.mailfrom=\"a)b;dkim=pass header.d=victim.example;dkim=none (c@evil.example\"\r\n" +
				"\r\nhi\r\n",
			want:      false,
			wantErrIs: ErrNoAlignedAuthentication,
		},
		{
			name: "malformed: second dmarc result",
			emailRaw: `From: owner@victim.example
To: alias@example.net
Authentication-Results: mx.example.net; dkim=none; spf=pass smtp.mailfrom=a@evil.example; dmarc=pass header.from=victim.example; dmarc=fail header.from=victim.example (policy=none)

hi`,
			want:      false,
			wantErrIs: ErrMalformedOwnAuthResults,
		},
		{
			name: "malformed: second spf result",
			emailRaw: `From: owner@victim.example
To: alias@example.net
Authentication-Results: mx.example.net; dkim=none; spf=fail smtp.mailfrom=a@evil.example; spf=pass smtp.mailfrom=a@victim.example; dmarc=none

hi`,
			want:      false,
			wantErrIs: ErrMalformedOwnAuthResults,
		},
		{
			name: "malformed: dkim result after spf",
			emailRaw: `From: owner@victim.example
To: alias@example.net
Authentication-Results: mx.example.net; dkim=none; dmarc=none; spf=pass smtp.mailfrom=a@evil.example; dkim=pass header.d=victim.example

hi`,
			want:      false,
			wantErrIs: ErrMalformedOwnAuthResults,
		},
		{
			name: "malformed: unknown method",
			emailRaw: `From: owner@victim.example
To: alias@example.net
Authentication-Results: mx.example.net; dkim=none; spf=pass smtp.mailfrom=a@evil.example; dmarc=pass header.from=victim.example; smtp.mailfrom=x

hi`,
			want:      false,
			wantErrIs: ErrMalformedOwnAuthResults,
		},
		{
			name: "malformed: result without method=value",
			emailRaw: `From: owner@victim.example
To: alias@example.net
Authentication-Results: mx.example.net; dkim=none; spf=pass b; dmarc=pass header.from=victim.example; c

hi`,
			want:      false,
			wantErrIs: ErrMalformedOwnAuthResults,
		},
		{
			name: "malformed: unterminated quoted string",
			emailRaw: `From: owner@victim.example
To: alias@example.net
Authentication-Results: mx.example.net; dkim=pass header.d=victim.example; spf=pass smtp.mailfrom="a@evil.example; dmarc=fail header.from=victim.example

hi`,
			want:      false,
			wantErrIs: ErrMalformedOwnAuthResults,
		},
		{
			name: "ARC: trusted sealer vouches for a forwarded message",
			emailRaw: `From: sender@example.com
To: recipient@example.net
Authentication-Results: mx.example.net; dkim=fail header.d=example.com; spf=fail smtp.mailfrom=a@forwarder.example; dmarc=fail header.from=example.com (policy=reject); arc=pass ("google.com:s=arc-20240605:i=1")
ARC-Seal: i=1; a=rsa-sha256; t=1; cv=none; d=google.com; s=arc-20240605; b=abc
ARC-Authentication-Results: i=1; mx.google.com; dkim=pass header.i=@example.com; spf=pass smtp.mailfrom=sender@example.com; dmarc=pass (p=REJECT) header.from=example.com

This is a test email.`,
			want:       true,
			wantNilErr: true,
		},
		{
			name: "ARC: sealer not trusted",
			emailRaw: `From: sender@example.com
To: recipient@example.net
Authentication-Results: mx.example.net; dmarc=fail header.from=example.com; arc=pass
ARC-Seal: i=1; a=rsa-sha256; cv=none; d=attacker.example; s=arc; b=abc
ARC-Authentication-Results: i=1; mx.attacker.example; dmarc=pass header.from=example.com

This is a test email.`,
			want:      false,
			wantErrIs: ErrDMARCFail,
		},
		{
			name: "ARC: chain not verified by own mail server",
			emailRaw: `From: sender@example.com
To: recipient@example.net
Authentication-Results: mx.example.net; dmarc=fail header.from=example.com; arc=fail
ARC-Seal: i=1; a=rsa-sha256; cv=none; d=google.com; s=arc; b=abc
ARC-Authentication-Results: i=1; mx.google.com; dmarc=pass header.from=example.com

This is a test email.`,
			want:      false,
			wantErrIs: ErrDMARCFail,
		},
		{
			name: "ARC: no ARC result from own mail server",
			emailRaw: `From: sender@example.com
To: recipient@example.net
Authentication-Results: mx.example.net; dmarc=fail header.from=example.com
ARC-Seal: i=1; a=rsa-sha256; cv=none; d=google.com; s=arc; b=abc
ARC-Authentication-Results: i=1; mx.google.com; dmarc=pass header.from=example.com

This is a test email.`,
			want:      false,
			wantErrIs: ErrDMARCFail,
		},
		{
			name: "ARC: latest seal is from an untrusted sealer",
			emailRaw: `From: sender@example.com
To: recipient@example.net
Authentication-Results: mx.example.net; dmarc=fail header.from=example.com; arc=pass
ARC-Seal: i=2; a=rsa-sha256; cv=pass; d=attacker.example; s=arc; b=abc
ARC-Authentication-Results: i=2; mx.attacker.example; dmarc=pass header.from=example.com
ARC-Seal: i=1; a=rsa-sha256; cv=none; d=google.com; s=arc; b=abc
ARC-Authentication-Results: i=1; mx.google.com; dmarc=fail header.from=example.com

This is a test email.`,
			want:      false,
			wantErrIs: ErrDMARCFail,
		},
		{
			name: "ARC: trusted sealer recorded a failure",
			emailRaw: `From: sender@example.com
To: recipient@example.net
Authentication-Results: mx.example.net; dmarc=fail header.from=example.com; arc=pass
ARC-Seal: i=1; a=rsa-sha256; cv=none; d=google.com; s=arc; b=abc
ARC-Authentication-Results: i=1; mx.google.com; dkim=pass header.d=attacker.example; dmarc=fail header.from=example.com

This is a test email.`,
			want:      false,
			wantErrIs: ErrDMARCFail,
		},
		{
			name: "ARC: no results for the latest instance",
			emailRaw: `From: sender@example.com
To: recipient@example.net
Authentication-Results: mx.example.net; dmarc=fail header.from=example.com; arc=pass
ARC-Seal: i=2; a=rsa-sha256; cv=pass; d=google.com; s=arc; b=abc
ARC-Authentication-Results: i=1; mx.google.com; dmarc=pass header.from=example.com

This is a test email.`,
			want:      false,
			wantErrIs: ErrDMARCFail,
		},
		{
			name: "ARC: duplicate results for the latest instance",
			emailRaw: `From: sender@example.com
To: recipient@example.net
Authentication-Results: mx.example.net; dmarc=fail header.from=example.com; arc=pass
ARC-Seal: i=1; a=rsa-sha256; cv=none; d=google.com; s=arc; b=abc
ARC-Authentication-Results: i=1; mx.google.com; dmarc=pass header.from=example.com
ARC-Authentication-Results: i=1; mx.google.com; dmarc=fail header.from=example.com

This is a test email.`,
			want:      false,
			wantErrIs: ErrDMARCFail,
		},
		{
			name: "ARC: one arc result per verified seal",
			emailRaw: `From: sender@example.com
To: recipient@example.net
Authentication-Results: mx.example.net; dmarc=fail header.from=example.com; arc=pass ("google.com:s=arc-20240605:i=1"); arc=pass ("microsoft.com:s=arcselector10001:i=2")
ARC-Seal: i=2; a=rsa-sha256; cv=pass; d=microsoft.com; s=arcselector10001; b=abc
ARC-Authentication-Results: i=2; mx.microsoft.com 1; dmarc=pass header.from=example.com
ARC-Seal: i=1; a=rsa-sha256; cv=none; d=google.com; s=arc-20240605; b=abc
ARC-Authentication-Results: i=1; mx.google.com; dmarc=pass header.from=example.com

This is a test email.`,
			want:       true,
			wantNilErr: true,
		},
		{
			name: "ARC: arc=pass next to arc=fail",
			emailRaw: `From: sender@example.com
To: recipient@example.net
Authentication-Results: mx.example.net; dmarc=fail header.from=example.com; arc=fail; arc=pass
ARC-Seal: i=1; a=rsa-sha256; cv=none; d=google.com; s=arc; b=abc
ARC-Authentication-Results: i=1; mx.google.com; dmarc=pass header.from=example.com

This is a test email.`,
			want:      false,
			wantErrIs: ErrDMARCFail,
		},
		{
			name: "ARC: trusted sealer's dmarc=fail decides over an earlier dmarc=pass",
			emailRaw: `From: sender@example.com
To: recipient@example.net
Authentication-Results: mx.example.net; dmarc=fail header.from=example.com; arc=pass
ARC-Seal: i=1; a=rsa-sha256; cv=none; d=google.com; s=arc; b=abc
ARC-Authentication-Results: i=1; mx.google.com; dmarc=pass header.from=example.com; dmarc=fail header.from=example.com

This is a test email.`,
			want:      false,
			wantErrIs: ErrDMARCFail,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := EmailAuthConfig{AuthServIDs: testServIDs, ARCTrustedSealers: testSealers}
			if tt.servIDs != nil {
				cfg.AuthServIDs = tt.servIDs
			}
			got, err := VerifyEmailAuth([]byte(tt.emailRaw), cfg)
			if got != tt.want {
				t.Errorf("VerifyEmailAuth() = %v, want %v (err: %v)", got, tt.want, err)
			}
			if tt.wantNilErr && err != nil {
				t.Errorf("VerifyEmailAuth() error = %v, want nil", err)
			}
			if !tt.want && err == nil {
				t.Errorf("VerifyEmailAuth() rejected without an error")
			}
			if tt.wantErrIs != nil && !errors.Is(err, tt.wantErrIs) {
				t.Errorf("VerifyEmailAuth() error = %v, want %v", err, tt.wantErrIs)
			}
		})
	}
}

func TestVerifyForwardAuth(t *testing.T) {
	tests := []struct {
		name      string
		emailRaw  string
		servIDs   []string
		want      bool
		wantErr   bool
		wantErrIs error
	}{
		{
			name: "aligned pass",
			emailRaw: `From: sender@example.com
To: alias@example.net
Authentication-Results: mx.example.net; dmarc=pass header.from=example.com

This is a test email.`,
			want: true,
		},
		{
			name: "misaligned email service provider signature",
			emailRaw: `From: news@brand.example
To: alias@example.net
Authentication-Results: mx.example.net; dkim=pass header.d=sendgrid.net; spf=pass smtp.mailfrom=a@sendgrid.net; dmarc=none header.from=brand.example

This is a test email.`,
			want: true,
		},
		{
			name: "no authentication at all",
			emailRaw: `From: sender@example.com
To: alias@example.net
Authentication-Results: mx.example.net; dkim=none; spf=none; dmarc=none

This is a test email.`,
			want: true,
		},
		{
			name: "DMARC fail with policy=none",
			emailRaw: `From: sender@example.com
To: alias@example.net
Authentication-Results: mx.example.net; dmarc=fail header.from=example.com (policy=none)

This is a test email.`,
			want: true,
		},
		{
			name: "DMARC fail with policy=quarantine",
			emailRaw: `From: sender@example.com
To: alias@example.net
Authentication-Results: mx.example.net; dmarc=fail header.from=example.com (policy=quarantine)

This is a test email.`,
			want: true,
		},
		{
			name: "DMARC fail with policy=reject",
			emailRaw: `From: sender@bank.example
To: alias@example.net
Authentication-Results: mx.example.net; dkim=pass header.d=attacker.example; dmarc=fail reason="No valid SPF" header.from=bank.example (policy=reject)

This is a test email.`,
			want:      false,
			wantErr:   true,
			wantErrIs: ErrDMARCPolicyReject,
		},
		{
			name: "DMARC fail with policy=reject only claimed by another authserv-id",
			emailRaw: `From: sender@example.com
To: alias@example.net
Authentication-Results: mx.attacker.example; dmarc=fail header.from=example.com (policy=reject)
Authentication-Results: mx.example.net; dmarc=pass header.from=example.com

This is a test email.`,
			want: true,
		},
		{
			name: "DMARC policy=reject overridden by trusted ARC sealer",
			emailRaw: `From: sender@example.com
To: alias@example.net
Authentication-Results: mx.example.net; dmarc=fail header.from=example.com (policy=reject); arc=pass
ARC-Seal: i=1; a=rsa-sha256; cv=none; d=microsoft.com; s=arcselector10001; b=abc
ARC-Authentication-Results: i=1; mx.microsoft.com 1; spf=pass smtp.mailfrom=example.com; dmarc=pass action=none header.from=example.com; dkim=pass header.d=example.com

This is a test email.`,
			want: true,
		},
		{
			name: "own header missing does not block forwarding",
			emailRaw: `From: sender@example.com
To: alias@example.net

This is a test email.`,
			want:      true,
			wantErr:   true,
			wantErrIs: ErrNoOwnAuthResults,
		},
		{
			name: "authserv-id not configured does not block forwarding",
			emailRaw: `From: sender@example.com
To: alias@example.net
Authentication-Results: mx.example.net; dmarc=fail header.from=example.com (policy=reject)

This is a test email.`,
			servIDs:   []string{" ", ""},
			want:      true,
			wantErr:   true,
			wantErrIs: ErrAuthServIDNotConfigured,
		},
		{
			name: "forged second header with own authserv-id",
			emailRaw: `From: sender@example.com
To: alias@example.net
Authentication-Results: mx.example.net; dmarc=pass header.from=example.com
Authentication-Results: mx.example.net; dmarc=fail header.from=example.com (policy=reject)

This is a test email.`,
			want:      false,
			wantErr:   true,
			wantErrIs: ErrDuplicateOwnAuthResults,
		},
		{
			name: "malformed own header hiding policy=reject",
			emailRaw: `From: sender@bank.example
To: alias@example.net
Authentication-Results: mx.example.net; dkim=none; spf=pass smtp.mailfrom="a@evil.example; dmarc=fail header.from=bank.example (policy=reject)

This is a test email.`,
			want:      false,
			wantErr:   true,
			wantErrIs: ErrMalformedOwnAuthResults,
		},
		{
			name: "Rspamd header with a quoted sender address",
			emailRaw: `From: sender@example.com
To: alias@example.net
Authentication-Results: mx.example.net; dkim=none; spf=pass (mx.example.net: domain of "a(b)c@example.com" designates 192.0.2.1 as permitted sender) smtp.mailfrom="a(b)c@example.com"; dmarc=none

This is a test email.`,
			want: true,
		},
		{
			name: "invalid email format",
			emailRaw: `This is not a valid email format
Just some random text`,
			want:    false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := EmailAuthConfig{AuthServIDs: testServIDs, ARCTrustedSealers: testSealers}
			if tt.servIDs != nil {
				cfg.AuthServIDs = tt.servIDs
			}
			got, err := VerifyForwardAuth([]byte(tt.emailRaw), cfg)
			if got != tt.want {
				t.Errorf("VerifyForwardAuth() = %v, want %v (err: %v)", got, tt.want, err)
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("VerifyForwardAuth() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErrIs != nil && !errors.Is(err, tt.wantErrIs) {
				t.Errorf("VerifyForwardAuth() error = %v, want %v", err, tt.wantErrIs)
			}
		})
	}
}
