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

func MatchOTP(secret string, hash string) bool {
	return HashMatches(secret, hash)
}
