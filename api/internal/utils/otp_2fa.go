package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base32"
	"encoding/binary"
	"errors"
	"net/url"
	"strings"
	"time"
)

// OTPConfig holds the configuration for the OTP verification
type OTPConfig struct {
	WindowSize int // Number of time intervals to check
}

// DefaultOTPConfig returns a default configuration
func DefaultOTPConfig() *OTPConfig {
	return &OTPConfig{
		WindowSize: 2,
	}
}

// ValidateSecret checks if the provided secret is valid
func ValidateSecret(secret string) error {
	// Ensure the secret is properly base32 encoded
	secret = strings.TrimSpace(secret)
	secret = strings.ToUpper(secret)

	if len(secret) < 16 {
		return errors.New("secret is too short, must be at least 16 characters")
	}

	_, err := base32.StdEncoding.DecodeString(secret)
	if err != nil {
		return errors.New("invalid base32 encoding in secret")
	}

	return nil
}

// ComputeCode computes the response code for a 64-bit challenge 'value' using the secret 'secret'.
// To avoid breaking compatibility with the previous API, it returns an invalid code (-1) when an error occurs,
// but does not silently ignore them (it forces a mismatch so the code will be rejected).
func computeCode(secret string, value int64) (int, error) {
	key, err := base32.StdEncoding.DecodeString(strings.ToUpper(strings.TrimSpace(secret)))
	if err != nil {
		return -1, err
	}

	h := hmac.New(sha256.New, key)

	if err := binary.Write(h, binary.BigEndian, value); err != nil {
		return -1, err
	}
	hm := h.Sum(nil)

	offset := hm[31] & 0x0f

	truncated := binary.BigEndian.Uint32(hm[offset : offset+4])

	truncated &= 0x7fffffff
	code := truncated % 1000000

	return int(code), nil
}

// Check checks whether specific code is valid given specific secret
// it doesn't check for duplicate use, which have to be checked on a different level
func Check(secret string, code int) (bool, error) {
	return CheckWithConfig(secret, code, DefaultOTPConfig())
}

// CheckWithConfig checks whether specific code is valid with custom configuration
func CheckWithConfig(secret string, code int, config *OTPConfig) (bool, error) {
	if err := ValidateSecret(secret); err != nil {
		return false, err
	}

	if code < 0 || code >= 1000000 {
		return false, errors.New("invalid OTP code")
	}

	t0 := int(time.Now().UTC().Unix() / 30)

	minT := t0 - (config.WindowSize / 2)
	maxT := t0 + (config.WindowSize / 2)

	for t := minT; t <= maxT; t++ {
		c, err := computeCode(secret, int64(t))
		if err != nil {
			return false, err
		}
		// Use constant-time comparison to prevent timing attacks
		if subtle.ConstantTimeEq(int32(c), int32(code)) == 1 {
			return true, nil
		}
	}

	return false, nil
}

// GenerateURI generates a URI that can be turned into a QR code
// to configure a Google Authenticator mobile app. It respects the recommendations
// on how to avoid conflicting accounts.
//
// See https://github.com/google/google-authenticator/wiki/Conflicting-Accounts
func GenerateURI(secret, user string, issuer string) string {
	auth := "totp/"

	q := make(url.Values)
	q.Add("secret", secret)
	if issuer != "" {
		q.Add("issuer", issuer)
		auth += issuer + ":"
	}

	return "otpauth://" + auth + url.QueryEscape(user) + "?" + q.Encode()
}
