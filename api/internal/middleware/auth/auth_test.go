package auth

import (
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/valyala/fasthttp"
	"ivpn.net/email/api/config"
)

func TestGetTokenSignature(t *testing.T) {
	tests := []struct {
		name     string
		jwt      string
		expected string
	}{
		{
			name:     "Valid JWT",
			jwt:      "header.payload.signature",
			expected: "signature",
		},
		{
			name:     "Invalid JWT with less parts",
			jwt:      "header.payload",
			expected: "",
		},
		{
			name:     "Invalid JWT with more parts",
			jwt:      "header.payload.signature.extra",
			expected: "",
		},
		{
			name:     "Empty JWT",
			jwt:      "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetTokenSignature(tt.jwt)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestGetTokenExp(t *testing.T) {
	tests := []struct {
		name        string
		token       string
		tokenSecret string
		expected    time.Duration
		expectError bool
	}{
		{
			name:        "Valid token",
			token:       generateToken(t, "secret", time.Now().Add(15*time.Hour).Unix()),
			tokenSecret: "secret",
			expected:    15 * time.Hour,
			expectError: false,
		},
		{
			name:        "Invalid token secret",
			token:       generateToken(t, "secret", time.Now().Add(1*time.Hour).Unix()),
			tokenSecret: "wrongsecret",
			expected:    0,
			expectError: true,
		},
		{
			name:        "No exp claim",
			token:       generateTokenWithoutExp(t, "secret"),
			tokenSecret: "secret",
			expected:    0,
			expectError: true,
		},
		{
			name:        "Invalid token format",
			token:       "invalid.token.format",
			tokenSecret: "secret",
			expected:    0,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.APIConfig{TokenSecret: tt.tokenSecret}
			app := fiber.New()
			req := &fasthttp.RequestCtx{}
			c := app.AcquireCtx(req)
			c.Request().Header.Set("Authorization", "Bearer "+tt.token)

			result, err := GetTokenExp(cfg, c)
			if (err != nil) != tt.expectError {
				t.Errorf("expected error: %v, got: %v", tt.expectError, err)
			}
			if !tt.expectError && result.Round(time.Hour) != tt.expected.Round(time.Hour) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestGetToken(t *testing.T) {
	tests := []struct {
		name           string
		authorization  string
		cookie         string
		expectedResult string
	}{
		{
			name:           "Valid Bearer token in Authorization header",
			authorization:  "Bearer validtoken",
			cookie:         "",
			expectedResult: "validtoken",
		},
		{
			name:           "Valid token in cookie",
			authorization:  "",
			cookie:         "validtoken",
			expectedResult: "validtoken",
		},
		{
			name:           "No token in Authorization header or cookie",
			authorization:  "",
			cookie:         "",
			expectedResult: "",
		},
		{
			name:           "Invalid Authorization header format",
			authorization:  "Invalid validtoken",
			cookie:         "",
			expectedResult: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			req := &fasthttp.RequestCtx{}
			c := app.AcquireCtx(req)
			c.Request().Header.Set("Authorization", tt.authorization)
			c.Request().Header.SetCookie(AUTH_COOKIE, tt.cookie)

			result := GetToken(c)
			if result != tt.expectedResult {
				t.Errorf("expected %s, got %s", tt.expectedResult, result)
			}
		})
	}
}

func TestGetUserID(t *testing.T) {
	tests := []struct {
		name           string
		userID         interface{}
		expectedResult string
		expectError    bool
	}{
		{
			name:           "Valid user ID",
			userID:         "12345",
			expectedResult: "12345",
			expectError:    false,
		},
		{
			name:           "Invalid user ID type",
			userID:         12345,
			expectedResult: "",
			expectError:    true,
		},
		{
			name:           "No user ID",
			userID:         nil,
			expectedResult: "",
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			req := &fasthttp.RequestCtx{}
			c := app.AcquireCtx(req)
			c.Locals(USER_ID, tt.userID)

			defer func() {
				if r := recover(); r != nil {
					if !tt.expectError {
						t.Errorf("unexpected panic: %v", r)
					}
				}
			}()

			result := GetUserID(c)
			if result != tt.expectedResult && !tt.expectError {
				t.Errorf("expected %s, got %s", tt.expectedResult, result)
			}
		})
	}
}

func TestGetAuthnToken(t *testing.T) {
	tests := []struct {
		name           string
		cookie         string
		expectedResult string
	}{
		{
			name:           "Valid authn token in cookie",
			cookie:         "validauthtoken",
			expectedResult: "validauthtoken",
		},
		{
			name:           "No authn token in cookie",
			cookie:         "",
			expectedResult: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			req := &fasthttp.RequestCtx{}
			c := app.AcquireCtx(req)
			c.Request().Header.SetCookie(AUTHN_COOKIE, tt.cookie)

			result := GetAuthnToken(c)
			if result != tt.expectedResult {
				t.Errorf("expected %s, got %s", tt.expectedResult, result)
			}
		})
	}
}

func generateToken(t *testing.T, secret string, exp int64) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": exp,
	})
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}
	return tokenString
}

func generateTokenWithoutExp(t *testing.T, secret string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{})
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}
	return tokenString
}
