package utils

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt"
	"ivpn.net/email/api/config"
)

func TestCreateAuthToken(t *testing.T) {
	cfg := config.APIConfig{
		TokenSecret:     "testsecret",
		TokenExpiration: time.Hour,
	}

	userID := "12345"
	email := "test@example.com"

	tokenString, err := CreateAuthToken(cfg, userID, email)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrInvalidKey
		}
		return []byte(cfg.TokenSecret), nil
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !token.Valid {
		t.Fatalf("expected token to be valid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		t.Fatalf("expected claims to be of type jwt.MapClaims")
	}

	if claims["user_id"] != userID {
		t.Errorf("expected user_id to be %v, got %v", userID, claims["user_id"])
	}

	if claims["email"] != email {
		t.Errorf("expected email to be %v, got %v", email, claims["email"])
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		t.Fatalf("expected exp to be a float64")
	}

	if time.Unix(int64(exp), 0).Sub(time.Now()) > cfg.TokenExpiration {
		t.Errorf("expected token expiration to be within %v, got %v", cfg.TokenExpiration, time.Unix(int64(exp), 0).Sub(time.Now()))
	}
}
