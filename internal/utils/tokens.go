package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"
	"strings"
)

type Token struct {
	Secret string
	Hash   string
}

func GenerateOTP() (*Token, error) {
	bigInt, err := rand.Int(rand.Reader, big.NewInt(900000))
	if err != nil {
		return nil, err
	}

	sixDigitNum := bigInt.Int64() + 100000

	// Convert the integer to a string and get the first 6 characters
	sixDigitStr := fmt.Sprintf("%06d", sixDigitNum)

	token := Token{
		Secret: sixDigitStr,
	}

	hash := sha256.Sum256([]byte(token.Secret))

	token.Hash = fmt.Sprintf("%x\n", hash)

	return &token, nil
}

func FormatOTP(s string) string {
	length := len(s)
	half := length / 2
	firstHalf := s[:half]
	secondHalf := s[half:]
	words := []string{firstHalf, secondHalf}
	return strings.Join(words, " ")
}

func ValidateToken(secret string) error {
	if secret == "" {
		return fmt.Errorf("token must be provided")
	}

	if len(secret) != 6 {
		return fmt.Errorf("token must be 6 bytes long")
	}

	return nil
}
