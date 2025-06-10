package utils

import (
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
			expected:   true,
			name:       "email like from domain",
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

func TestExtractValue(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
		name     string
	}{
		{
			input:    "value;param",
			expected: "value",
			name:     "with semicolon",
		},
		{
			input:    "value param",
			expected: "value",
			name:     "with space",
		},
		{
			input:    "value\nparam",
			expected: "value",
			name:     "with newline",
		},
		{
			input:    "value\rparam",
			expected: "value",
			name:     "with carriage return",
		},
		{
			input:    "value",
			expected: "value",
			name:     "no separator",
		},
		{
			input:    "",
			expected: "",
			name:     "empty string",
		},
		{
			input:    "value; param",
			expected: "value",
			name:     "space after semicolon",
		},
		{
			input:    "value   ",
			expected: "value",
			name:     "trailing spaces",
		},
		{
			input:    "value;param;another",
			expected: "value",
			name:     "multiple semicolons",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := extractValue(tc.input)
			if actual != tc.expected {
				t.Errorf("extractValue(%q) = %q, expected %q", tc.input, actual, tc.expected)
			}
		})
	}
}

func TestParseAuthResults(t *testing.T) {
	testCases := []struct {
		headers  []string
		expected AuthResults
		name     string
	}{
		{
			headers: []string{"DKIM=pass header.d=example.com"},
			expected: AuthResults{
				DKIM:        "pass",
				DKIMDomain:  "example.com",
				SPF:         "",
				SPFDomain:   "",
				DMARC:       "",
				DMARCDomain: "",
			},
			name: "dkim pass",
		},
		{
			headers: []string{"SPF=pass smtp.mailfrom=example.com"},
			expected: AuthResults{
				DKIM:        "",
				DKIMDomain:  "",
				SPF:         "pass",
				SPFDomain:   "example.com",
				DMARC:       "",
				DMARCDomain: "",
			},
			name: "spf pass",
		},
		{
			headers: []string{"DMARC=pass header.from=example.com"},
			expected: AuthResults{
				DKIM:        "",
				DKIMDomain:  "",
				SPF:         "",
				SPFDomain:   "",
				DMARC:       "pass",
				DMARCDomain: "example.com",
			},
			name: "dmarc pass",
		},
		{
			headers: []string{"DKIM=pass header.d=example.com", "SPF=pass smtp.mailfrom=example.net", "DMARC=pass header.from=example.org"},
			expected: AuthResults{
				DKIM:        "pass",
				DKIMDomain:  "example.com",
				SPF:         "pass",
				SPFDomain:   "example.net",
				DMARC:       "pass",
				DMARCDomain: "example.org",
			},
			name: "all pass",
		},
		{
			headers: []string{"DKIM=fail header.d=example.com", "SPF=fail smtp.mailfrom=example.net", "DMARC=fail header.from=example.org"},
			expected: AuthResults{
				DKIM:        "",
				DKIMDomain:  "",
				SPF:         "",
				SPFDomain:   "",
				DMARC:       "",
				DMARCDomain: "",
			},
			name: "all fail",
		},
		{
			headers: []string{"DKIM=pass header.d=example.com;param", "SPF=pass smtp.mailfrom=example.net;param", "DMARC=pass header.from=example.org;param"},
			expected: AuthResults{
				DKIM:        "pass",
				DKIMDomain:  "example.com",
				SPF:         "pass",
				SPFDomain:   "example.net",
				DMARC:       "pass",
				DMARCDomain: "example.org",
			},
			name: "with params",
		},
		{
			headers: []string{"DKIM=pass", "SPF=pass", "DMARC=pass"},
			expected: AuthResults{
				DKIM:        "pass",
				DKIMDomain:  "",
				SPF:         "pass",
				SPFDomain:   "",
				DMARC:       "pass",
				DMARCDomain: "",
			},
			name: "pass without domain",
		},
		{
			headers: []string{},
			expected: AuthResults{
				DKIM:        "",
				DKIMDomain:  "",
				SPF:         "",
				SPFDomain:   "",
				DMARC:       "",
				DMARCDomain: "",
			},
			name: "empty headers",
		},
		{
			headers: []string{"dkim=pass header.d=EXAMPLE.com", "spf=pass smtp.mailfrom=EXAMPLE.net", "dmarc=pass header.from=EXAMPLE.org"},
			expected: AuthResults{
				DKIM:        "pass",
				DKIMDomain:  "example.com",
				SPF:         "pass",
				SPFDomain:   "example.net",
				DMARC:       "pass",
				DMARCDomain: "example.org",
			},
			name: "uppercase",
		},
		{
			headers: []string{"ARC-Authentication-Results: dkim=pass header.d=example.com"},
			expected: AuthResults{
				DKIM:        "pass",
				DKIMDomain:  "example.com",
				SPF:         "",
				SPFDomain:   "",
				DMARC:       "",
				DMARCDomain: "",
			},
			name: "ARC header",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := parseAuthResults(tc.headers)
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Errorf("parseAuthResults(%q) = %v, expected %v", tc.headers, actual, tc.expected)
			}
		})
	}
}

func TestVerifyEmailAuth(t *testing.T) {
	tests := []struct {
		name     string
		emailRaw string
		want     bool
		wantErr  bool
	}{
		{
			name: "valid email with DMARC pass",
			emailRaw: `From: sender@example.com
To: recipient@example.net
Subject: Test Email
Date: Thu, 22 Aug 2023 12:00:00 -0700
Authentication-Results: mx.example.net; dmarc=pass header.from=example.com

This is a test email.`,
			want:    true,
			wantErr: false,
		},
		{
			name: "valid email with DKIM pass and matching domain",
			emailRaw: `From: sender@example.com
To: recipient@example.net
Subject: Test Email
Date: Thu, 22 Aug 2023 12:00:00 -0700
Authentication-Results: mx.example.net; dkim=pass header.d=example.com

This is a test email.`,
			want:    true,
			wantErr: false,
		},
		{
			name: "valid email with SPF pass and matching domain",
			emailRaw: `From: sender@example.com
To: recipient@example.net
Subject: Test Email
Date: Thu, 22 Aug 2023 12:00:00 -0700
Authentication-Results: mx.example.net; spf=pass smtp.mailfrom=example.com

This is a test email.`,
			want:    true,
			wantErr: false,
		},
		{
			name: "valid email with DMARC pass but domain mismatch",
			emailRaw: `From: sender@example.com
To: recipient@example.net
Subject: Test Email
Date: Thu, 22 Aug 2023 12:00:00 -0700
Authentication-Results: mx.example.net; dmarc=pass header.from=different.org

This is a test email.`,
			want:    false,
			wantErr: true,
		},
		{
			name: "valid email with all auth passes",
			emailRaw: `From: sender@example.com
To: recipient@example.net
Subject: Test Email
Date: Thu, 22 Aug 2023 12:00:00 -0700
Authentication-Results: mx.example.net; dkim=pass header.d=example.com; spf=pass smtp.mailfrom=example.com; dmarc=pass header.from=example.com

This is a test email.`,
			want:    true,
			wantErr: false,
		},
		{
			name: "valid email with DKIM pass but domain mismatch",
			emailRaw: `From: sender@example.com
To: recipient@example.net
Subject: Test Email
Date: Thu, 22 Aug 2023 12:00:00 -0700
Authentication-Results: mx.example.net; dkim=pass header.d=different.org

This is a test email.`,
			want:    false,
			wantErr: true,
		},
		{
			name: "valid email with SPF pass but domain mismatch",
			emailRaw: `From: sender@example.com
To: recipient@example.net
Subject: Test Email
Date: Thu, 22 Aug 2023 12:00:00 -0700
Authentication-Results: mx.example.net; spf=pass smtp.mailfrom=different.org

This is a test email.`,
			want:    false,
			wantErr: true,
		},
		{
			name: "valid email with all auth fails",
			emailRaw: `From: sender@example.com
To: recipient@example.net
Subject: Test Email
Date: Thu, 22 Aug 2023 12:00:00 -0700
Authentication-Results: mx.example.net; dkim=fail; spf=fail; dmarc=fail

This is a test email.`,
			want:    false,
			wantErr: false,
		},
		{
			name: "valid email with ARC authentication results",
			emailRaw: `From: sender@example.com
To: recipient@example.net
Subject: Test Email
Date: Thu, 22 Aug 2023 12:00:00 -0700
ARC-Authentication-Results: mx.example.net; dkim=pass header.d=example.com

This is a test email.`,
			want:    true,
			wantErr: false,
		},
		{
			name: "valid email with no authentication results",
			emailRaw: `From: sender@example.com
To: recipient@example.net
Subject: Test Email
Date: Thu, 22 Aug 2023 12:00:00 -0700

This is a test email.`,
			want:    false,
			wantErr: true,
		},
		{
			name: "invalid email with no From header",
			emailRaw: `To: recipient@example.net
Subject: Test Email
Date: Thu, 22 Aug 2023 12:00:00 -0700
Authentication-Results: mx.example.net; dkim=pass header.d=example.com

This is a test email.`,
			want:    false,
			wantErr: true,
		},
		{
			name: "invalid email format",
			emailRaw: `This is not a valid email format
Just some random text`,
			want:    false,
			wantErr: true,
		},
		{
			name: "email with both Authentication-Results and ARC-Authentication-Results",
			emailRaw: `From: sender@example.com
To: recipient@example.net
Subject: Test Email
Date: Thu, 22 Aug 2023 12:00:00 -0700
Authentication-Results: mx.example.net; dkim=fail; spf=fail
ARC-Authentication-Results: mx.example.net; dkim=pass header.d=example.com

This is a test email.`,
			want:    true,
			wantErr: false,
		},
		{
			name: "multiple conflicting Authentication-Results headers",
			emailRaw: `From: sender@example.com
To: recipient@example.net
Subject: Test Email
Date: Thu, 22 Aug 2023 12:00:00 -0700
Authentication-Results: mx.example.net; dkim=fail; spf=fail
Authentication-Results: mx.example.net; dkim=pass header.d=example.com

This is a test email.`,
			want:    false,
			wantErr: false,
		},
		{
			name: "multiple non-conflicting Authentication-Results headers",
			emailRaw: `From: sender@example.com
To: recipient@example.net
Subject: Test Email
Date: Thu, 22 Aug 2023 12:00:00 -0700
Authentication-Results: mx.example.net; dkim=pass header.d=example.com
Authentication-Results: mx.example.net; spf=pass smtp.mailfrom=example.com

This is a test email.`,
			want:    true,
			wantErr: false,
		},
		{
			name: "multiple ARC-Authentication-Results headers with different results",
			emailRaw: `From: sender@example.com
To: recipient@example.net
Subject: Test Email
Date: Thu, 22 Aug 2023 12:00:00 -0700
ARC-Authentication-Results: mx.example.net; dkim=fail
ARC-Authentication-Results: mx.example.net; dkim=pass header.d=example.com

This is a test email.`,
			want:    false,
			wantErr: false,
		},
		{
			name: "subdomain in From matching broader auth domain",
			emailRaw: `From: sender@sub.example.com
To: recipient@example.net
Subject: Test Email
Date: Thu, 22 Aug 2023 12:00:00 -0700
Authentication-Results: mx.example.net; dkim=pass header.d=example.com

This is a test email.`,
			want:    true,
			wantErr: false,
		},
		{
			name: "malformed From address",
			emailRaw: `From: invalid-email
To: recipient@example.net
Subject: Test Email
Date: Thu, 22 Aug 2023 12:00:00 -0700
Authentication-Results: mx.example.net; dkim=pass header.d=example.com

This is a test email.`,
			want:    false,
			wantErr: true,
		},
		{
			name: "email with Authentication-Results without dkim/spf/dmarc",
			emailRaw: `From: sender@example.com
To: recipient@example.net
Subject: Test Email
Date: Thu, 22 Aug 2023 12:00:00 -0700
Authentication-Results: phl-mx-05.messagingengine.com; arc=none (no signatures found)
ARC-Authentication-Results: mx.example.net; dkim=pass header.d=example.com

This is a test email.`,
			want:    true,
			wantErr: false,
		},
		{
			name: "email with Authentication-Results with some other methods but no fails",
			emailRaw: `From: sender@example.com
To: recipient@example.net
Subject: Test Email
Date: Thu, 22 Aug 2023 12:00:00 -0700
Authentication-Results: mx.example.net; iprev=pass; x-custom=neutral; auth=pass
ARC-Authentication-Results: mx.example.net; dkim=pass header.d=example.com

This is a test email.`,
			want:    true,
			wantErr: false,
		},
		{
			name: "email with Authentication-Results containing other passes and one fail",
			emailRaw: `From: sender@example.com
To: recipient@example.net
Subject: Test Email
Date: Thu, 22 Aug 2023 12:00:00 -0700
Authentication-Results: mx.example.net; iprev=pass; x-custom=pass; dkim=fail
Authentication-Results: mx.example.net; dkim=pass header.d=example.com

This is a test email.`,
			want:    false,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := VerifyEmailAuth([]byte(tt.emailRaw))
			if (err != nil) != tt.wantErr {
				t.Errorf("VerifyEmailAuth() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("VerifyEmailAuth() = %v, want %v", got, tt.want)
			}
		})
	}
}
