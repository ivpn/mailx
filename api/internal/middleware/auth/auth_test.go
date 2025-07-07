package auth

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

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
