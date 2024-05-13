package utils

import "github.com/alexedwards/argon2id"

func Hash(secret string) (string, error) {
	return argon2id.CreateHash(secret, argon2id.DefaultParams)
}

func HashMatches(secret string, hash string) bool {
	matches, err := argon2id.ComparePasswordAndHash(secret, hash)
	return matches && err == nil
}
