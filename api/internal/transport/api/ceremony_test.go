package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/go-webauthn/webauthn/webauthn"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"ivpn.net/email/api/internal/middleware/auth"
	"ivpn.net/email/api/internal/model"
	"ivpn.net/email/api/internal/service"
)

type memCache struct {
	mu   sync.Mutex
	data map[string]string
}

func (c *memCache) Set(_ context.Context, key string, value any, _ time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = value.(string)
	return nil
}

func (c *memCache) Get(_ context.Context, key string) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	v, ok := c.data[key]
	if !ok {
		return "", errors.New("not found")
	}
	return v, nil
}

func (c *memCache) GetDel(ctx context.Context, key string) (string, error) {
	v, err := c.Get(ctx, key)
	if err == nil {
		_ = c.Del(ctx, key)
	}
	return v, err
}

func (c *memCache) Del(_ context.Context, key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
	return nil
}

func (c *memCache) Incr(context.Context, string, time.Duration) error { return nil }

// ceremonyTestService runs the real ceremony storage and records account access. The embedded
// Service is nil, so any other call panics and fails the request with a 500.
type ceremonyTestService struct {
	Service
	ceremonies *service.Service
	user       model.User
	calls      []string
}

func (s *ceremonyTestService) SaveCeremony(ctx context.Context, c model.Ceremony, ttl time.Duration) (string, error) {
	return s.ceremonies.SaveCeremony(ctx, c, ttl)
}

func (s *ceremonyTestService) ConsumeCeremony(ctx context.Context, token string, t model.CeremonyType) (model.Ceremony, error) {
	return s.ceremonies.ConsumeCeremony(ctx, token, t)
}

func (s *ceremonyTestService) GetUser(context.Context, string) (model.User, error) {
	s.calls = append(s.calls, "GetUser")
	return s.user, nil
}

func (s *ceremonyTestService) SaveCredential(context.Context, webauthn.Credential, string) error {
	s.calls = append(s.calls, "SaveCredential")
	return nil
}

func (s *ceremonyTestService) UpdateCredential(context.Context, webauthn.Credential, string) error {
	s.calls = append(s.calls, "UpdateCredential")
	return nil
}

func (s *ceremonyTestService) SaveSession(context.Context, webauthn.SessionData, string, string, time.Time, bool) error {
	s.calls = append(s.calls, "SaveSession")
	return nil
}

func (s *ceremonyTestService) SetStepUp(context.Context, string) error {
	s.calls = append(s.calls, "SetStepUp")
	return nil
}

func newCeremonyTestApp(svc *ceremonyTestService) *fiber.App {
	h := &Handler{Service: svc}
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = c.SendStatus(fiber.StatusInternalServerError)
			}
		}()
		c.Locals(auth.USER_ID, "victim")
		return c.Next()
	})
	app.Post("/register/finish", h.FinishRegistration)
	app.Post("/register/add/finish", h.FinishAddPasskey)
	app.Post("/login/finish", h.FinishLogin)
	app.Post("/login/passkey/finish", h.FinishPasskeyLogin)
	app.Post("/stepup/passkey/finish", h.StepUpPasskeyFinish)

	return app
}

func postWithCeremony(t *testing.T, app *fiber.App, path, token string) (int, string) {
	req := httptest.NewRequest(http.MethodPost, path, strings.NewReader("{}"))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(&http.Cookie{Name: auth.AUTHN_TEMP_COOKIE, Value: token})
	res, err := app.Test(req)
	assert.NoError(t, err)

	var body struct {
		Error string `json:"error"`
	}
	_ = json.NewDecoder(res.Body).Decode(&body)

	return res.StatusCode, body.Error
}

func TestFinishHandlersRejectOtherCeremonyTypes(t *testing.T) {
	finishes := map[string]model.CeremonyType{
		"/register/finish":       model.CeremonyRegistration,
		"/register/add/finish":   model.CeremonyAddPasskey,
		"/login/finish":          model.CeremonyLogin,
		"/login/passkey/finish":  model.CeremonyPasskeyLogin,
		"/stepup/passkey/finish": model.CeremonyStepUp,
	}

	for path, want := range finishes {
		for _, saved := range finishes {
			if saved == want {
				continue
			}

			t.Run(path+"<-"+string(saved), func(t *testing.T) {
				svc := &ceremonyTestService{ceremonies: &service.Service{Cache: &memCache{data: map[string]string{}}}}
				token, err := svc.SaveCeremony(context.Background(), model.Ceremony{
					Type:        saved,
					UserID:      "victim",
					SessionData: webauthn.SessionData{Challenge: "challenge", UserID: []byte("victim")},
				}, time.Minute)
				assert.NoError(t, err)

				status, msg := postWithCeremony(t, newCeremonyTestApp(svc), path, token)

				assert.Equal(t, fiber.StatusBadRequest, status)
				assert.Equal(t, ErrGetSession, msg)
				assert.Empty(t, svc.calls, "request must stop before touching the account")
			})
		}
	}
}

func TestFinishRegistrationRejectsEstablishedAccount(t *testing.T) {
	established := []model.User{
		{PasswordHash: "hash"},
		{Creds: []webauthn.Credential{{ID: []byte("cred")}}},
		// ChangeEmail clears IsActive on established accounts, so it isn't a signup signal
		{IsActive: false, Creds: []webauthn.Credential{{ID: []byte("cred")}}},
	}

	for _, user := range established {
		svc := &ceremonyTestService{
			ceremonies: &service.Service{Cache: &memCache{data: map[string]string{}}},
			user:       user,
		}
		token, err := svc.SaveCeremony(context.Background(), model.Ceremony{
			Type:        model.CeremonyRegistration,
			UserID:      "victim",
			SessionData: webauthn.SessionData{Challenge: "challenge", UserID: []byte("victim")},
		}, time.Minute)
		assert.NoError(t, err)

		status, msg := postWithCeremony(t, newCeremonyTestApp(svc), "/register/finish", token)

		// Rejected by the signup check, before the attestation is even parsed
		assert.Equal(t, fiber.StatusBadRequest, status)
		assert.Equal(t, ErrFinishRegistration, msg)
		assert.Equal(t, []string{"GetUser"}, svc.calls)
	}
}
