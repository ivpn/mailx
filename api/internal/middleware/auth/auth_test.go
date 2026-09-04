package auth

import (
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
	"ivpn.net/email/api/config"
)

func TestGetUserID(t *testing.T) {
	tests := []struct {
		name           string
		userID         any
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
func TestGetAuthToken(t *testing.T) {
	tests := []struct {
		name           string
		authorization  string
		expectedResult string
	}{
		{
			name:           "Valid Bearer token",
			authorization:  "Bearer abc123token",
			expectedResult: "abc123token",
		},
		{
			name:           "Valid Bearer token with spaces",
			authorization:  "Bearer token with spaces",
			expectedResult: "token with spaces",
		},
		{
			name:           "Empty Bearer token",
			authorization:  "Bearer ",
			expectedResult: "",
		},
		{
			name:           "No Bearer prefix",
			authorization:  "abc123token",
			expectedResult: "",
		},
		{
			name:           "Different auth scheme",
			authorization:  "Basic abc123",
			expectedResult: "",
		},
		{
			name:           "Empty authorization header",
			authorization:  "",
			expectedResult: "",
		},
		{
			name:           "Bearer with lowercase",
			authorization:  "bearer abc123token",
			expectedResult: "",
		},
		{
			name:           "Just Bearer",
			authorization:  "Bearer",
			expectedResult: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			req := &fasthttp.RequestCtx{}
			c := app.AcquireCtx(req)

			if tt.authorization != "" {
				c.Request().Header.Set("Authorization", tt.authorization)
			}

			result := GetAuthToken(c)
			if result != tt.expectedResult {
				t.Errorf("expected %q, got %q", tt.expectedResult, result)
			}
		})
	}
}
func TestGetAuthnCookie(t *testing.T) {
	tests := []struct {
		name           string
		cookieValue    string
		expectedResult string
	}{
		{
			name:           "Valid authn cookie",
			cookieValue:    "session123token",
			expectedResult: "session123token",
		},
		{
			name:           "Empty authn cookie",
			cookieValue:    "",
			expectedResult: "",
		},
		{
			name:           "Cookie with special characters",
			cookieValue:    "token!@#$%^&*()",
			expectedResult: "token!@#$%^&*()",
		},
		{
			name:           "Long cookie value",
			cookieValue:    "very_long_session_token_with_lots_of_characters_12345678901234567890",
			expectedResult: "very_long_session_token_with_lots_of_characters_12345678901234567890",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			req := &fasthttp.RequestCtx{}
			c := app.AcquireCtx(req)

			if tt.cookieValue != "" {
				c.Request().Header.SetCookie(AUTHN_COOKIE, tt.cookieValue)
			}

			result := GetAuthnCookie(c)
			if result != tt.expectedResult {
				t.Errorf("expected %q, got %q", tt.expectedResult, result)
			}
		})
	}
}

func TestNewPSK(t *testing.T) {
	tests := []struct {
		name           string
		configuredPSK  string
		authorization  string
		expectedStatus int
	}{
		{
			name:           "valid PSK",
			configuredPSK:  "supersecretpsk",
			authorization:  "Bearer supersecretpsk",
			expectedStatus: fiber.StatusOK,
		},
		{
			name:           "wrong PSK",
			configuredPSK:  "supersecretpsk",
			authorization:  "Bearer wrongpsk",
			expectedStatus: fiber.StatusUnauthorized,
		},
		{
			name:           "missing Authorization header",
			configuredPSK:  "supersecretpsk",
			authorization:  "",
			expectedStatus: fiber.StatusUnauthorized,
		},
		{
			// regression test: constant-time compare of "" == "" would otherwise report a match
			name:           "empty configured PSK, no Authorization header",
			configuredPSK:  "",
			authorization:  "",
			expectedStatus: fiber.StatusUnauthorized,
		},
		{
			name:           "empty configured PSK, empty Bearer token",
			configuredPSK:  "",
			authorization:  "Bearer ",
			expectedStatus: fiber.StatusUnauthorized,
		},
		{
			name:           "case-sensitive mismatch",
			configuredPSK:  "SuperSecretPSK",
			authorization:  "Bearer supersecretpsk",
			expectedStatus: fiber.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			app.Use(NewPSK(tt.configuredPSK))
			app.Get("/protected", func(c *fiber.Ctx) error {
				return c.SendStatus(fiber.StatusOK)
			})

			req := httptest.NewRequest(fiber.MethodGet, "/protected", nil)
			if tt.authorization != "" {
				req.Header.Set("Authorization", tt.authorization)
			}

			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, resp.StatusCode)
			}
		})
	}
}

func TestNewCookieTempAuthn(t *testing.T) {
	tests := []struct {
		name            string
		tokenExpiration time.Duration
	}{
		{name: "short TokenExpiration", tokenExpiration: time.Second},
		{name: "long TokenExpiration", tokenExpiration: 168 * time.Hour},
	}

	wantMaxAge := int(WebAuthnCeremonyExpiration.Seconds())

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config.APIConfig{TokenExpiration: tt.tokenExpiration}

			before := time.Now()
			cookie := NewCookieTempAuthn("sometoken", "/some/path", cfg)
			after := time.Now()

			if cookie.Name != AUTHN_TEMP_COOKIE {
				t.Errorf("expected cookie name %q, got %q", AUTHN_TEMP_COOKIE, cookie.Name)
			}

			// regression: MaxAge must stay pinned to the ceremony window regardless of TokenExpiration
			if cookie.MaxAge != wantMaxAge {
				t.Errorf("expected MaxAge %d, got %d", wantMaxAge, cookie.MaxAge)
			}

			minExpires := before.Add(WebAuthnCeremonyExpiration)
			maxExpires := after.Add(WebAuthnCeremonyExpiration)
			if cookie.Expires.Before(minExpires) || cookie.Expires.After(maxExpires) {
				t.Errorf("expected Expires within [%v, %v], got %v", minExpires, maxExpires, cookie.Expires)
			}
		})
	}
}
