// Package signedurl builds and verifies expiring, HMAC-signed URLs, used for endpoints that
// are called without a session (e.g. List-Unsubscribe one-click requests from mail providers).
package signedurl

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"maps"
	"net/url"
	"strconv"
	"time"
)

const (
	paramExpires   = "expires"
	paramSignature = "signature"
)

type Signer struct {
	secret  []byte
	baseURL string
	now     func() time.Time
}

// New returns a Signer, or nil when the secret or base URL isn't configured, so callers can
// treat signed URLs as unavailable instead of producing ones that can't be verified.
func New(secret string, baseURL string) *Signer {
	if secret == "" || baseURL == "" {
		return nil
	}

	return &Signer{secret: []byte(secret), baseURL: baseURL, now: time.Now}
}

// Sign returns baseURL+path?params&expires=...&signature=... The signature covers the path and
// every query param, so none of them can be changed without invalidating the URL.
func (s *Signer) Sign(path string, params url.Values, ttl time.Duration) string {
	q := url.Values{}
	maps.Copy(q, params)
	q.Set(paramExpires, strconv.FormatInt(s.now().Add(ttl).Unix(), 10))
	q.Set(paramSignature, s.mac(path, q))

	return s.baseURL + path + "?" + q.Encode()
}

// Verify checks the signature and expiry of a request's path and raw query string.
func (s *Signer) Verify(path string, rawQuery string) bool {
	q, err := url.ParseQuery(rawQuery)
	if err != nil {
		return false
	}

	sig, err := hex.DecodeString(q.Get(paramSignature))
	if err != nil || len(sig) == 0 {
		return false
	}
	q.Del(paramSignature)

	exp, err := strconv.ParseInt(q.Get(paramExpires), 10, 64)
	if err != nil || s.now().Unix() > exp {
		return false
	}

	expected, err := hex.DecodeString(s.mac(path, q))
	if err != nil {
		return false
	}

	return hmac.Equal(sig, expected)
}

func (s *Signer) mac(path string, q url.Values) string {
	m := hmac.New(sha256.New, s.secret)
	m.Write([]byte(path + "?" + q.Encode())) // Encode sorts by key, so the input is deterministic
	return hex.EncodeToString(m.Sum(nil))
}
