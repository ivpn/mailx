package utils

import (
	"crypto/rand"
	"errors"
	"io"
	"math/big"
	"strings"
)

const (
	AlphaNumericUserFriendly          = "ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnpqrstuvwxyz23456789"
	AlphaNumericUserFriendlyUppercase = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	AlphaNumericLowerUpper            = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	AlphaLower                        = "abcdefghijkmnopqrstuvwxyz"
)

var (
	ErrInvalidLength       = errors.New("invalid random string length")
	ErrEmptyCharset        = errors.New("empty charset")
	ErrInsufficientEntropy = errors.New("failed to generate secure random string: insufficient entropy")
)

func RandomString(n int, charset string) (string, error) {
	if n <= 0 {
		return "", ErrInvalidLength
	}

	if len(charset) == 0 {
		return "", ErrEmptyCharset
	}

	sb := strings.Builder{}
	sb.Grow(n)

	maxIdx := big.NewInt(int64(len(charset)))

	for range n {
		idx, err := rand.Int(rand.Reader, maxIdx)
		if err != nil {
			if err == io.EOF {
				return "", ErrInsufficientEntropy
			}
			return "", err
		}
		sb.WriteByte(charset[idx.Int64()])
	}

	return sb.String(), nil
}
