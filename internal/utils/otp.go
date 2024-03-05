package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"
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
	hash := sha256.Sum256([]byte(token.Secret))
	token.Hash = fmt.Sprintf("%x\n", hash)

	return &token, nil
}

func HashOTP(secret string) string {
	hash := sha256.Sum256([]byte(secret))
	return fmt.Sprintf("%x\n", hash)
}
