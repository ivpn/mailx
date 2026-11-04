package api

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"ivpn.net/email/api/config"
	"ivpn.net/email/api/internal/middleware/auth"
	"ivpn.net/email/api/internal/model"
	"ivpn.net/email/api/internal/utils"
	"ivpn.net/email/api/internal/utils/signedurl"
)

// scopeTestService has one web session and one access-key session for the same user, and records
// the sessions Authenticate saves. The embedded Service is nil, so a request that gets past the auth
// middleware to any other call panics and fails with a 500, never a 401.
type scopeTestService struct {
	Service
	saved []model.Session
}

const (
	scopeTestWebToken = "w+E/b=="
	scopeTestAPIToken = "a+P/i=="
	scopeTestKeyID    = "6f1c2a9e-4b7d-4c3e-9a51-0d2e8f7b6c4a"
)

func (s *scopeTestService) OneClickSigner() *signedurl.Signer { return nil }

func (s *scopeTestService) GetSession(_ context.Context, token string) (model.Session, bool, error) {
	switch token {
	case scopeTestWebToken:
		return model.Session{UserID: "user-a", Token: token, Scope: model.SessionScopeWeb}, true, nil
	case scopeTestAPIToken:
		return model.Session{UserID: "user-a", Token: token, Scope: model.SessionScopeAPI, AccessKeyID: scopeTestKeyID}, true, nil
	}
	return model.Session{}, false, errors.New("could not get session by token")
}

func (s *scopeTestService) GetUser(context.Context, string) (model.User, error) {
	return model.User{BaseModel: model.BaseModel{ID: "user-a"}}, nil
}

func (s *scopeTestService) GetAccessKey(context.Context, string) (model.AccessKey, error) {
	return model.AccessKey{BaseModel: model.BaseModel{ID: scopeTestKeyID}, UserId: "user-a"}, nil
}

func (s *scopeTestService) GetDefaults(context.Context, string) (model.Settings, []string, error) {
	return model.Settings{}, nil, nil
}

func (s *scopeTestService) GetVerifiedDomains(context.Context, string) ([]model.Domain, error) {
	return nil, nil
}

func (s *scopeTestService) SaveSession(_ context.Context, session model.Session) error {
	s.saved = append(s.saved, session)
	return nil
}

func newScopeTestApp(svc *scopeTestService) *fiber.App {
	h := &Handler{
		Service:   svc,
		Server:    fiber.New(),
		Validator: utils.NewValidator(),
		Cfg:       config.APIConfig{ApiTokenExpiration: 24 * time.Hour},
	}
	h.SetupRoutes(config.APIConfig{ApiAllowOrigin: "https://app.mailx.net"})
	return h.Server
}

// Routes an access-key token could use to take over the account or read its data
var webOnlyRoutes = []struct {
	method string
	path   string
}{
	{http.MethodPost, "/v1/register/add"},
	{http.MethodPost, "/v1/register/add/finish"},
	{http.MethodGet, "/v1/accesskeys"},
	{http.MethodPost, "/v1/accesskeys"},
	{http.MethodPost, "/v1/recipient"},
	{http.MethodGet, "/v1/user"},
	{http.MethodGet, "/v1/logs"},
	{http.MethodPut, "/v1/user/changepassword"},
	{http.MethodPost, "/v1/user/logout"},
}

func TestAccessKeySessionRejectedOnWebRoutes(t *testing.T) {
	for _, r := range webOnlyRoutes {
		for _, via := range []string{"bearer", "cookie"} {
			t.Run(via+" "+r.method+" "+r.path, func(t *testing.T) {
				app := newScopeTestApp(&scopeTestService{})

				req := httptest.NewRequest(r.method, r.path, strings.NewReader(`{}`))
				req.Header.Set("Content-Type", "application/json")
				if via == "bearer" {
					req.Header.Set("Authorization", "Bearer "+scopeTestAPIToken)
				} else {
					req.AddCookie(&http.Cookie{Name: auth.AUTHN_COOKIE, Value: scopeTestAPIToken})
				}
				res, err := app.Test(req)
				require.NoError(t, err)

				assert.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
			})
		}
	}
}

func TestWebSessionRejectedOnAccessKeyAPI(t *testing.T) {
	routes := []struct {
		method string
		path   string
	}{
		{http.MethodGet, "/v1/api/aliases"},
		{http.MethodPost, "/v1/api/alias"},
		{http.MethodGet, "/v1/api/defaults"},
		{http.MethodPost, "/v1/api/logout"},
	}

	for _, r := range routes {
		t.Run(r.method+" "+r.path, func(t *testing.T) {
			app := newScopeTestApp(&scopeTestService{})

			req := httptest.NewRequest(r.method, r.path, strings.NewReader(`{}`))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+scopeTestWebToken)
			res, err := app.Test(req)
			require.NoError(t, err)

			assert.Equal(t, fiber.StatusUnauthorized, res.StatusCode)
		})
	}
}

func TestAuthenticateSavesAccessKeySession(t *testing.T) {
	svc := &scopeTestService{}
	app := newScopeTestApp(svc)

	accessKey := scopeTestKeyID + "." + strings.Repeat("k", 48)
	req := httptest.NewRequest(http.MethodPost, "/v1/api/authenticate", strings.NewReader(`{"access_key":"`+accessKey+`"}`))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)
	require.NoError(t, err)

	require.Equal(t, fiber.StatusOK, res.StatusCode)
	require.Len(t, svc.saved, 1)
	assert.Equal(t, model.SessionScopeAPI, svc.saved[0].Scope)
	assert.Equal(t, scopeTestKeyID, svc.saved[0].AccessKeyID, "linked to its key, so revoking the key ends it")
	assert.Equal(t, "user-a", svc.saved[0].UserID)
	assert.False(t, svc.saved[0].Remember)
}
