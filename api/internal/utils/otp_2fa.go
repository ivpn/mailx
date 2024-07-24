package utils

// Inspired by:
// https://github.com/dgryski/dgoogauth/blob/master/googauth.go

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"net/url"
	"time"
)

const windowSize = 2

// Much of this code assumes int == int64, which probably is not the case.

// ComputeCode computes the response code for a 64-bit challenge 'value' using the secret 'secret'.
// To avoid breaking compatibility with the previous API, it returns an invalid code (-1) when an error occurs,
// but does not silently ignore them (it forces a mismatch so the code will be rejected).
func computeCode(secret string, value int64) (int, error) {

	key, err := base32.StdEncoding.DecodeString(secret)
	if err != nil {
		return -1, err
	}

	hash := hmac.New(sha1.New, key)

	if err := binary.Write(hash, binary.BigEndian, value); err != nil {
		return -1, err
	}
	h := hash.Sum(nil)

	offset := h[19] & 0x0f

	truncated := binary.BigEndian.Uint32(h[offset : offset+4])

	truncated &= 0x7fffffff
	code := truncated % 1000000

	return int(code), nil
}

// Check checks whether specific code is valid given specific secret
// it doesn't check for duplicate use, which have to be checked on a different level
func Check(secret string, code int) (bool, error) {

	// t0 := int(time.Now().UTC().Unix() / 30)
	t0 := int(time.Now().UTC().Unix() / 30)

	minT := t0 - (windowSize / 2)
	maxT := t0 + (windowSize / 2)
	for t := minT; t <= maxT; t++ {
		c, err := computeCode(secret, int64(t))
		if err != nil {
			return false, err
		}
		if c == code {
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
