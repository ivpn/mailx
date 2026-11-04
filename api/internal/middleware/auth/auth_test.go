package auth

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"slices"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
	"ivpn.net/email/api/config"
	"ivpn.net/email/api/internal/model"
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

// fakeSessionService holds sessions by token and users by ID. It records every token it is asked for.
type fakeSessionService struct {
	sessions map[string]model.Session // by token
	users    map[string]model.User
	lookups  []string
}

func (s *fakeSessionService) GetSession(_ context.Context, token string) (model.Session, bool, error) {
	s.lookups = append(s.lookups, token)
	session, ok := s.sessions[token]
	if !ok {
		// Same as the repository: unknown or expired tokens come back as not found, with an error
		return model.Session{}, false, errors.New("could not get session by token")
	}
	session.Token = token
	return session, true, nil
}

func (s *fakeSessionService) GetUser(_ context.Context, id string) (model.User, error) {
	user, ok := s.users[id]
	if !ok {
		return model.User{}, errors.New("user not found")
	}
	return user, nil
}

// Session tokens are standard base64 (GenSessionToken), so they can contain +, / and =
const (
	tokenA = "a+B/c=="
	tokenB = "b+C/d=="
	// tokenAPI is user-a's session made from an access key, tokenLegacy one saved before scopes
	tokenAPI    = "k+E/y=="
	tokenLegacy = "l+E/g=="
)

func newFakeSessionService() *fakeSessionService {
	return &fakeSessionService{
		sessions: map[string]model.Session{
			tokenA:      {UserID: "user-a", Scope: model.SessionScopeWeb},
			tokenB:      {UserID: "user-b", Scope: model.SessionScopeWeb},
			"orphan":    {UserID: "deleted-user", Scope: model.SessionScopeWeb},
			tokenAPI:    {UserID: "user-a", Scope: model.SessionScopeAPI, AccessKeyID: "key-1"},
			tokenLegacy: {UserID: "user-a"},
		},
		users: map[string]model.User{
			"user-a": {BaseModel: model.BaseModel{ID: "user-a"}},
			"user-b": {BaseModel: model.BaseModel{ID: "user-b"}},
		},
	}
}

func authRequest(cookie, authorization string) *http.Request {
	req := httptest.NewRequest(fiber.MethodGet, "/v1/user", nil)
	if cookie != "" {
		req.AddCookie(&http.Cookie{Name: AUTHN_COOKIE, Value: cookie})
	}
	if authorization != "" {
		req.Header.Set("Authorization", authorization)
	}
	return req
}

func TestNew(t *testing.T) {
	tests := []struct {
		name          string
		cookie        string
		authorization string
		wantStatus    int
		wantUser      string
	}{
		{name: "valid cookie", cookie: tokenA, wantStatus: fiber.StatusOK, wantUser: "user-a"},
		{name: "web session as bearer token", authorization: "Bearer " + tokenA, wantStatus: fiber.StatusUnauthorized},
		{name: "cookie wins over a different bearer token", cookie: tokenA, authorization: "Bearer " + tokenB, wantStatus: fiber.StatusOK, wantUser: "user-a"},
		{name: "invalid cookie does not fall back to a valid bearer token", cookie: "unknown", authorization: "Bearer " + tokenB, wantStatus: fiber.StatusUnauthorized},
		{name: "unknown or expired bearer token", authorization: "Bearer unknown", wantStatus: fiber.StatusUnauthorized},
		{name: "unknown cookie", cookie: "unknown", wantStatus: fiber.StatusUnauthorized},
		{name: "bearer token for a deleted user", authorization: "Bearer orphan", wantStatus: fiber.StatusUnauthorized},
		{name: "empty bearer token", authorization: "Bearer ", wantStatus: fiber.StatusUnauthorized},
		{name: "lowercase scheme", authorization: "bearer " + tokenA, wantStatus: fiber.StatusUnauthorized},
		{name: "other scheme", authorization: "Basic " + tokenA, wantStatus: fiber.StatusUnauthorized},
		{name: "token without scheme", authorization: tokenA, wantStatus: fiber.StatusUnauthorized},
		{name: "extra space is not trimmed", authorization: "Bearer  " + tokenA, wantStatus: fiber.StatusUnauthorized},
		{name: "URL-encoded token is not decoded", authorization: "Bearer a%2BB%2Fc%3D%3D", wantStatus: fiber.StatusUnauthorized},
		{name: "no credentials", wantStatus: fiber.StatusUnauthorized},
		// An access-key session reaches only /v1/api, whichever way it's sent
		{name: "access-key session as bearer token", authorization: "Bearer " + tokenAPI, wantStatus: fiber.StatusUnauthorized},
		{name: "access-key session as cookie", cookie: tokenAPI, wantStatus: fiber.StatusUnauthorized},
		{name: "access-key cookie does not fall back to a web bearer token", cookie: tokenAPI, authorization: "Bearer " + tokenA, wantStatus: fiber.StatusUnauthorized},
		{name: "session without a scope", cookie: tokenLegacy, wantStatus: fiber.StatusUnauthorized},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := newFakeSessionService()
			app := fiber.New()
			app.Use(New(config.APIConfig{}, nil, svc))
			app.Get("/v1/user", func(c *fiber.Ctx) error {
				return c.SendString(GetUserID(c))
			})

			resp, err := app.Test(authRequest(tt.cookie, tt.authorization))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			body, _ := io.ReadAll(resp.Body)

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, resp.StatusCode)
			}
			if tt.wantStatus == fiber.StatusOK && string(body) != tt.wantUser {
				t.Errorf("expected user %q, got %q", tt.wantUser, string(body))
			}
			if slices.Contains(svc.lookups, "") {
				t.Errorf("looked up a session with an empty token")
			}
		})
	}
}

func TestNewAPIAuth(t *testing.T) {
	tests := []struct {
		name          string
		cookie        string
		authorization string
		wantStatus    int
	}{
		{name: "access-key session", authorization: "Bearer " + tokenAPI, wantStatus: fiber.StatusOK},
		{name: "web session as bearer token", authorization: "Bearer " + tokenA, wantStatus: fiber.StatusUnauthorized},
		{name: "session without a scope", authorization: "Bearer " + tokenLegacy, wantStatus: fiber.StatusUnauthorized},
		{name: "access-key session as cookie", cookie: tokenAPI, wantStatus: fiber.StatusUnauthorized},
		{name: "unknown token", authorization: "Bearer unknown", wantStatus: fiber.StatusUnauthorized},
		{name: "no credentials", wantStatus: fiber.StatusUnauthorized},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := newFakeSessionService()
			app := fiber.New()
			app.Use(NewAPIAuth(config.APIConfig{}, svc))
			app.Get("/v1/user", func(c *fiber.Ctx) error {
				return c.SendString(GetUserID(c))
			})

			resp, err := app.Test(authRequest(tt.cookie, tt.authorization))
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			body, _ := io.ReadAll(resp.Body)

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("expected status %d, got %d", tt.wantStatus, resp.StatusCode)
			}
			if tt.wantStatus == fiber.StatusOK && string(body) != "user-a" {
				t.Errorf("expected user %q, got %q", "user-a", string(body))
			}
		})
	}
}
