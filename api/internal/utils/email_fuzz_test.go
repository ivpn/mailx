package utils

import "testing"

// FuzzPreprocessEmailData fuzzes the raw-bytes normalisation pass every
// inbound message goes through before net/mail and the letters MIME parser
// see it (called from model.ParseMsg and mailer.Reply/Forward on
// attacker-supplied bytes).
func FuzzPreprocessEmailData(f *testing.F) {
	seeds := []string{
		"From: a@example.com\r\nTo: b@example.com\r\nSubject: =??Q?broken?=\r\n\r\nBody.",
		"From: a@example.com\r\nSubject: =??B?dGVzdA==?= =??Q?more?=\r\n\r\nBody.",
		"From: a@example.com\r\n\r\n",
		"",
		"not-a-valid-header-at-all",
		"From: a@example.com\r\nSubject: \xff\xfe\x00\r\n\r\n\x80\x81",
		// Header value containing many empty-charset encoded words back to back.
		"From: a@example.com\r\nSubject: " + repeatEW(200) + "\r\n\r\nBody.",
	}
	for _, s := range seeds {
		f.Add([]byte(s))
	}

	f.Fuzz(func(t *testing.T, data []byte) {
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("PreprocessEmailData panicked on input (len=%d): %v", len(data), r)
			}
		}()
		_, _ = PreprocessEmailData(data)
	})
}

func repeatEW(n int) string {
	s := ""
	for i := 0; i < n; i++ {
		s += "=??q?x?="
	}
	return s
}

// FuzzDecodeHeaderWithCharset fuzzes RFC 2047 decoding with an attacker
// controlled charset label, which is looked up via golang.org/x/net/html/charset
// (x/net's htmlindex/ianaindex table lookups on arbitrary strings).
func FuzzDecodeHeaderWithCharset(f *testing.F) {
	seeds := []string{
		"=?UTF-8?B?SGVsbG8=?=",
		"=?ISO-8859-1?Q?Caf=E9?=",
		"=?windows-1251?B?zdC50L3Rgg==?=",
		"=?koi8-r?Q?test?=",
		"=??Q?empty_charset?=",
		"=?totally-bogus-charset-xyz?B?AA==?=",
		"plain text, no encoding",
		"=?UTF-8?B?not-valid-base64!!!?=",
		"=?UTF-8?Q?%invalid%escape?=",
		"=?UTF-8?B?" + "QQ==" + "?==?UTF-8?B?QQ==?=", // adjacent encoded words
		"",
	}
	for _, s := range seeds {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, s string) {
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("DecodeHeaderWithCharset panicked on input %q: %v", s, r)
			}
		}()
		_ = DecodeHeaderWithCharset(s)
	})
}

// FuzzSafeDecodeAddressName fuzzes the RFC 2047 display-name fallback
// decoder used on parsed From/To addresses.
func FuzzSafeDecodeAddressName(f *testing.F) {
	seeds := []string{
		"=?UTF-8?B?SGVsbG8=?=",
		"=??Q?empty?=",
		"=?bogus?B?AA==?=",
		"plain name",
		"",
	}
	for _, s := range seeds {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, s string) {
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("SafeDecodeAddressName panicked on input %q: %v", s, r)
			}
		}()
		_ = SafeDecodeAddressName(s)
	})
}

// FuzzNormalizeAddressSeparators fuzzes the Outlook-style ';' separator
// normalisation applied to To/Cc header values before mail.ParseAddressList.
func FuzzNormalizeAddressSeparators(f *testing.F) {
	seeds := []string{
		"a@b.com; c@d.com",
		"a@b.com;",
		";;;",
		"",
		"a@b.com, b@c.com; c@d.com,",
	}
	for _, s := range seeds {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, s string) {
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("NormalizeAddressSeparators panicked on input %q: %v", s, r)
			}
		}()
		_ = NormalizeAddressSeparators(s)
	})
}
