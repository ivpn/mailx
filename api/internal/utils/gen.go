package utils

import (
	"math/rand"
	"strings"
	"time"
)

// AlphaNumericUserFriendly contains letters without those which looks similar to one another
const AlphaNumericUserFriendly = "ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnpqrstuvwxyz23456789"
const AlphaNumericUserFriendlyUppercase = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"

const ALPHANUMERICLOWERUPPER = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
const ALPHALOWER = "abcdefghijkmnopqrstuvwxyz"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

// RandomString Generate a random string fast
// Why so complicated? See https://stackoverflow.com/a/31832326
func RandomString(n int, letterBytes string) string {
	sb := strings.Builder{}
	sb.Grow(n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			sb.WriteByte(letterBytes[idx])
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return sb.String()
}
