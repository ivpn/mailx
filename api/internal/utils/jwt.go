package utils

import (
	"time"

	"github.com/golang-jwt/jwt"
	"ivpn.net/email/api/config"
)

func CreateAuthToken(cfg config.APIConfig, userID string) (string, error) {
	claims := jwt.MapClaims{}
	claims["user_id"] = userID
	claims["exp"] = time.Now().Add(cfg.TokenExpiration).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.TokenSecret))
}
