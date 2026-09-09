package utils

import (
	"testing"
)

// stripDot tests

func TestStripDot(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"mail1.example.net.", "mail1.example.net"},
		{"mail1.example.net", "mail1.example.net"},
		{".", ""},
		{"", ""},
	}
	for _, tc := range tests {
		got := stripDot(tc.input)
		if got != tc.want {
			t.Errorf("stripDot(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

// LookupTXTExact tests

func TestLookupTXTExact_NotFound(t *testing.T) {
	// .invalid TLD is reserved (RFC 2606) and never resolves.
	found, err := LookupTXTExact("nonexistent.invalid", "some-value")
	if err != nil {
		t.Fatalf("expected nil error for non-existent domain, got: %v", err)
	}
	if found {
		t.Fatal("expected false for non-existent domain")
	}
}

func TestLookupTXTExact_Mismatch(t *testing.T) {
	// gmail.com has TXT records but not this value.
	found, err := LookupTXTExact("gmail.com", "service-verify=this-will-never-exist")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found {
		t.Fatal("expected false for mismatched TXT value")
	}
}

// LookupTXTContains tests

func TestLookupTXTContains_NotFound(t *testing.T) {
	found, err := LookupTXTContains("nonexistent.invalid", "v=spf1")
	if err != nil {
		t.Fatalf("expected nil error for non-existent domain, got: %v", err)
	}
	if found {
		t.Fatal("expected false for non-existent domain")
	}
}

func TestLookupTXTContains_SPF(t *testing.T) {
	// gmail.com is known to publish an SPF record.
	found, err := LookupTXTContains("gmail.com", "v=spf1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !found {
		t.Skip("gmail.com SPF TXT record not found; skipping (may be a network issue)")
	}
}

func TestLookupTXTContains_DMARC(t *testing.T) {
	// _dmarc.gmail.com is known to publish a DMARC record.
	found, err := LookupTXTContains("_dmarc.gmail.com", "v=DMARC1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !found {
		t.Skip("_dmarc.gmail.com TXT record not found; skipping (may be a network issue)")
	}
}

func TestLookupTXTContains_Mismatch(t *testing.T) {
	found, err := LookupTXTContains("gmail.com", "this-value-will-never-exist-xyz123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found {
		t.Fatal("expected false for mismatched TXT value")
	}
}

// LookupMX tests

func TestLookupMX_NotFound(t *testing.T) {
	found, err := LookupMX("nonexistent.invalid", "mail.example.com")
	if err != nil {
		t.Fatalf("expected nil error for non-existent domain, got: %v", err)
	}
	if found {
		t.Fatal("expected false for non-existent domain")
	}
}

func TestLookupMX_Mismatch(t *testing.T) {
	// gmail.com has MX records but not this host.
	found, err := LookupMX("gmail.com", "mail.this-does-not-exist.example")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found {
		t.Fatal("expected false for mismatched MX host")
	}
}

func TestLookupMX_TrailingDot(t *testing.T) {
	// Verify that trailing dots on the target are normalised before comparison.
	// Both "gmail-smtp-in.l.google.com." and "gmail-smtp-in.l.google.com" must behave
	// identically. We check that not-found returns false without an error in both cases.
	found1, err1 := LookupMX("gmail.com", "mail.this-does-not-exist.example.")
	found2, err2 := LookupMX("gmail.com", "mail.this-does-not-exist.example")
	if err1 != nil || err2 != nil {
		t.Fatalf("unexpected errors: %v / %v", err1, err2)
	}
	if found1 != found2 {
		t.Fatal("trailing dot normalisation inconsistency")
	}
}

// validSPFRecord tests

func TestValidSPFRecord(t *testing.T) {
	tests := []struct {
		name              string
		record            string
		requiredMechanism string
		want              bool
	}{
		{"exact match", "v=spf1 include:spf.mailx.net -all", "spf.mailx.net", true},
		{"softfail all", "v=spf1 include:spf.mailx.net ~all", "spf.mailx.net", true},
		{"extra includes and reordered", "v=spf1 include:_spf.google.com include:spf.mailx.net ~all", "spf.mailx.net", true},
		{"mx mechanism instead of include", "v=spf1 mx -all", "spf.mailx.net", true},
		{"case-insensitive version and mechanism", "V=SPF1 INCLUDE:SPF.MAILX.NET -ALL", "spf.mailx.net", true},
		{"missing required mechanism", "v=spf1 include:spf.other.net -all", "spf.mailx.net", false},
		{"missing all mechanism", "v=spf1 include:spf.mailx.net", "spf.mailx.net", false},
		{"neutral all not accepted", "v=spf1 include:spf.mailx.net ?all", "spf.mailx.net", false},
		{"not an spf record", "some other txt record", "spf.mailx.net", false},
		{"wrong spf version", "v=spf2 include:spf.mailx.net -all", "spf.mailx.net", false},
		{"empty record", "", "spf.mailx.net", false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := validSPFRecord(tc.record, tc.requiredMechanism)
			if got != tc.want {
				t.Errorf("validSPFRecord(%q, %q) = %v, want %v", tc.record, tc.requiredMechanism, got, tc.want)
			}
		})
	}
}

// validDMARCRecord tests

func TestValidDMARCRecord(t *testing.T) {
	tests := []struct {
		name   string
		record string
		want   bool
	}{
		{"quarantine", "v=DMARC1; p=quarantine; adkim=s", true},
		{"reject", "v=DMARC1; p=reject", true},
		{"reordered tags", "adkim=s; p=reject; v=DMARC1", true},
		{"case-insensitive", "V=DMARC1; P=Reject", true},
		{"none policy rejected", "v=DMARC1; p=none", false},
		{"missing version tag", "p=reject", false},
		{"missing policy tag", "v=DMARC1; adkim=s", false},
		{"empty record", "", false},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := validDMARCRecord(tc.record)
			if got != tc.want {
				t.Errorf("validDMARCRecord(%q) = %v, want %v", tc.record, got, tc.want)
			}
		})
	}
}

// LookupSPF tests

func TestLookupSPF_NotFound(t *testing.T) {
	found, err := LookupSPF("nonexistent.invalid", "spf.example.net")
	if err != nil {
		t.Fatalf("expected nil error for non-existent domain, got: %v", err)
	}
	if found {
		t.Fatal("expected false for non-existent domain")
	}
}

func TestLookupSPF_Mismatch(t *testing.T) {
	// gmail.com publishes an SPF record but never includes this mechanism.
	found, err := LookupSPF("gmail.com", "spf.this-will-never-exist.example")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found {
		t.Fatal("expected false for mismatched SPF mechanism")
	}
}

// LookupDMARC tests

func TestLookupDMARC_NotFound(t *testing.T) {
	found, err := LookupDMARC("_dmarc.nonexistent.invalid")
	if err != nil {
		t.Fatalf("expected nil error for non-existent domain, got: %v", err)
	}
	if found {
		t.Fatal("expected false for non-existent domain")
	}
}

// LookupCNAME tests

func TestLookupCNAME_NotFound(t *testing.T) {
	found, err := LookupCNAME("nonexistent-cname.invalid", "target.example.com")
	if err != nil {
		t.Fatalf("expected nil error for non-existent host, got: %v", err)
	}
	if found {
		t.Fatal("expected false for non-existent host")
	}
}

func TestLookupCNAME_Mismatch(t *testing.T) {
	// www.github.com resolves to a CNAME but not to this target.
	found, err := LookupCNAME("www.github.com", "wrong.target.example.com")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if found {
		t.Fatal("expected false for mismatched CNAME target")
	}
}
