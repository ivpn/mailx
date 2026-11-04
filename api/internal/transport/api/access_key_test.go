package api

import (
	"testing"
	"time"
)

func TestParseExpiresAt(t *testing.T) {
	cases := []struct {
		in      string
		want    time.Time
		wantErr bool
	}{
		// What the app sends: JSON-serialised JS Date.
		{in: "2027-10-06T10:15:30.123Z", want: time.Date(2027, 10, 6, 10, 15, 30, 123000000, time.UTC)},
		{in: "2027-10-06T10:15:30Z", want: time.Date(2027, 10, 6, 10, 15, 30, 0, time.UTC)},
		{in: "2027-10-06T12:15:30+02:00", want: time.Date(2027, 10, 6, 10, 15, 30, 0, time.UTC)},
		{in: "2027-10-06", want: time.Date(2027, 10, 6, 0, 0, 0, 0, time.UTC)},
		{in: "not a date", wantErr: true},
		{in: "10/06/2027", wantErr: true},
	}

	for _, tc := range cases {
		got, err := parseExpiresAt(tc.in)
		if tc.wantErr {
			if err == nil {
				t.Errorf("parseExpiresAt(%q): expected error, got %v", tc.in, got)
			}
			continue
		}
		if err != nil {
			t.Errorf("parseExpiresAt(%q): unexpected error: %v", tc.in, err)
			continue
		}
		if !got.Equal(tc.want) {
			t.Errorf("parseExpiresAt(%q) = %v, want %v", tc.in, got, tc.want)
		}
	}
}
