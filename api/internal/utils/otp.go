package utils

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
)

var (
	ErrCreateOTP = errors.New("create OTP failed")
)

const (
	charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	length  = 32
)

type OTP struct {
	Secret string
	Hash   string
}

func CreateOTP() (*OTP, error) {
	bigInt, err := rand.Int(rand.Reader, big.NewInt(900000))
	if err != nil {
		return nil, ErrCreateOTP
	}

	sixDigitNum := bigInt.Int64() + 100000
	sixDigitStr := fmt.Sprintf("%06d", sixDigitNum)
	token := OTP{Secret: sixDigitStr}

	hash, err := Hash(token.Secret)
	if err != nil {
		return nil, ErrCreateOTP
	}

	token.Hash = hash

	return &token, nil
}

func CreateLongOTP() (*OTP, error) {
	str, err := GenerateRandomString(length)
	if err != nil {
		return nil, ErrCreateOTP
	}

	token := OTP{Secret: str}

	hash, err := Hash(token.Secret)
	if err != nil {
		return nil, ErrCreateOTP
	}

	token.Hash = hash

	return &token, nil
}

func MatchOTP(secret string, hash string) bool {
	return HashMatches(secret, hash)
}

func GenerateRandomString(length int) (string, error) {
	result := make([]byte, length)
	charsetLength := big.NewInt(int64(len(charset)))

	for i := range result {
		num, err := rand.Int(rand.Reader, charsetLength)
		if err != nil {
			return "", err
		}
		result[i] = charset[num.Int64()]
	}

	return string(result), nil
}
