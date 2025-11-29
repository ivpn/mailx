package auth

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
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
