package utils

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"

	"github.com/alexedwards/argon2id"
)

var (
	ErrHashFailed = errors.New("hash failed")
)

type OTP struct {
	Secret string
	Hash   string
}

func CreateOTP() (*OTP, error) {
	bigInt, err := rand.Int(rand.Reader, big.NewInt(900000))
	if err != nil {
		return nil, err
	}

	sixDigitNum := bigInt.Int64() + 100000
	sixDigitStr := fmt.Sprintf("%06d", sixDigitNum)
	token := OTP{Secret: sixDigitStr}

	hash, err := argon2id.CreateHash(token.Secret, argon2id.DefaultParams)
	if err != nil {
		return nil, ErrHashFailed
	}

	token.Hash = hash

	return &token, nil
}

func MatchOTP(secret string, hash string) bool {
	match, err := argon2id.ComparePasswordAndHash(secret, hash)
	if err != nil {
		return false
	}

	return match
}
