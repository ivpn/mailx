package utils

import "testing"

// FuzzValidSPFRecord and FuzzValidDMARCRecord fuzz parsing of TXT record
// content returned by DNS lookups during domain ownership/SPF/DMARC
// verification (internal/service domain verification flow). The record
// text is controlled by whoever controls DNS for the domain being
// verified, which during a verification flow is not yet a trusted party.
func FuzzValidSPFRecord(f *testing.F) {
	seeds := []struct {
		record string
		mech   string
	}{
		{"v=spf1 include:spf.example.net -all", "spf.example.net"},
		{"v=spf1 mx ~all", "spf.example.net"},
		{"v=spf1", "spf.example.net"},
		{"", "spf.example.net"},
		{"v=spf1 include: -all", ""},
		{"not-an-spf-record", "spf.example.net"},
		{"v=spf1 " + repeatField("include:x.com ", 500) + "-all", "x.com"},
	}
	for _, s := range seeds {
		f.Add(s.record, s.mech)
	}

	f.Fuzz(func(t *testing.T, record, mechanism string) {
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("validSPFRecord panicked on record=%q mechanism=%q: %v", record, mechanism, r)
			}
		}()
		_ = validSPFRecord(record, mechanism)
	})
}

func repeatField(s string, n int) string {
	out := ""
	for i := 0; i < n; i++ {
		out += s
	}
	return out
}

func FuzzValidDMARCRecord(f *testing.F) {
	seeds := []string{
		"v=DMARC1; p=reject",
		"v=DMARC1; p=quarantine; adkim=s",
		"v=DMARC1; p=none",
		"",
		"v=dmarc1;p=reject",
		";;;;",
		"p=reject; v=DMARC1",
		"v=DMARC1" + string([]byte{0}) + "; p=reject",
	}
	for _, s := range seeds {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, record string) {
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("validDMARCRecord panicked on record=%q: %v", record, r)
			}
		}()
		_ = validDMARCRecord(record)
	})
}
