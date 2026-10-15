package api

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"ivpn.net/email/api/config"
	"ivpn.net/email/api/internal/middleware/auth"
	"ivpn.net/email/api/internal/model"
	"ivpn.net/email/api/internal/utils"
	"ivpn.net/email/api/internal/utils/signedurl"
)

// fakeOneClickService embeds Service (as a nil interface) so only the methods under test need
// implementing; anything else panics if invoked unexpectedly.
type fakeOneClickService struct {
	Service

	signer    *signedurl.Signer
	err       error
	calls     []string
	behaviour map[string]model.ListUnsubscribeBehaviour
}

func (f *fakeOneClickService) OneClickSigner() *signedurl.Signer {
	return f.signer
}

func (f *fakeOneClickService) OneClickDeactivateAlias(ctx context.Context, aliasID string, req model.OneClickRequest) error {
	f.calls = append(f.calls, "deactivate:"+aliasID)
	return f.err
}

func (f *fakeOneClickService) OneClickDeleteAlias(ctx context.Context, aliasID string, req model.OneClickRequest) error {
	f.calls = append(f.calls, "delete:"+aliasID)
	return f.err
}

func (f *fakeOneClickService) OneClickBlock(ctx context.Context, aliasID string, itemType model.BlocklistType, value string) error {
	f.calls = append(f.calls, "block-"+string(itemType)+":"+aliasID+":"+value)
	return f.err
}

func (f *fakeOneClickService) UpdateListUnsubscribeBehaviour(ctx context.Context, userID string, behaviour model.ListUnsubscribeBehaviour) error {
	f.behaviour[userID] = behaviour
	return nil
}

func newOneClickTestApp(service *fakeOneClickService) *fiber.App {
	h := &Handler{
		Service:   service,
		Server:    fiber.New(),
		Validator: utils.NewValidator(),
	}
	h.SetupRoutes(config.APIConfig{ApiAllowOrigin: "https://app.mailx.net"})
	return h.Server
}

var oneClickRoutes = []struct {
	name   string
	path   string
	params url.Values
	want   string
}{
	{"deactivate", model.OneClickDeactivatePath + "alias-1", nil, "deactivate:alias-1"},
	{"delete", model.OneClickDeletePath + "alias-1", nil, "delete:alias-1"},
	{"block email", model.OneClickBlockEmailPath + "alias-1", url.Values{"email": {"news@example.com"}}, "block-email:alias-1:news@example.com"},
	{"block domain", model.OneClickBlockDomainPath + "alias-1", url.Values{"domain": {"example.com"}}, "block-domain:alias-1:example.com"},
}

func signedTarget(t *testing.T, signer *signedurl.Signer, path string, params url.Values, ttl time.Duration) string {
	t.Helper()
	u, err := url.Parse(signer.Sign(path, params, ttl))
	if err != nil {
		t.Fatalf("url.Parse() error = %v", err)
	}
	return u.RequestURI()
}

func oneClickPost(t *testing.T, app *fiber.App, target string) (int, string) {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, target, strings.NewReader("List-Unsubscribe=One-Click"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return doRequest(t, app, req)
}

func doRequest(t *testing.T, app *fiber.App, req *http.Request) (int, string) {
	t.Helper()
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	body, _ := io.ReadAll(res.Body)
	return res.StatusCode, string(body)
}

func TestOneClick_ValidSignedPost(t *testing.T) {
	signer := signedurl.New("test-secret", "https://api.mailx.net")

	for _, route := range oneClickRoutes {
		t.Run(route.name, func(t *testing.T) {
			service := &fakeOneClickService{signer: signer}
			app := newOneClickTestApp(service)
			target := signedTarget(t, signer, route.path, route.params, time.Hour)

			// Repeated calls reach the (idempotent) service action each time
			for range 2 {
				status, body := oneClickPost(t, app, target)
				if status != http.StatusOK || body != "" {
					t.Fatalf("POST = %d %q, want 200 with an empty body", status, body)
				}
			}

			if len(service.calls) != 2 || service.calls[0] != route.want {
				t.Errorf("calls = %v, want [%s %s]", service.calls, route.want, route.want)
			}
		})
	}
}

func TestOneClick_InvalidSignature(t *testing.T) {
	signer := signedurl.New("test-secret", "https://api.mailx.net")
	other := signedurl.New("other-secret", "https://api.mailx.net")

	for _, route := range oneClickRoutes {
		t.Run(route.name, func(t *testing.T) {
			service := &fakeOneClickService{signer: signer}
			app := newOneClickTestApp(service)

			valid := signedTarget(t, signer, route.path, route.params, time.Hour)
			targets := map[string]string{
				"unsigned":     route.path,
				"other secret": signedTarget(t, other, route.path, route.params, time.Hour),
				"expired":      signedTarget(t, signer, route.path, route.params, -time.Minute),
				"other alias":  strings.Replace(valid, "alias-1", "alias-2", 1),
			}
			if route.params != nil {
				targets["tampered param"] = strings.Replace(valid, "example.com", "example.org", 1)
			}

			for name, target := range targets {
				if status, body := oneClickPost(t, app, target); status != http.StatusForbidden || body != "" {
					t.Errorf("%s: POST %s = %d %q, want 403 with an empty body", name, target, status, body)
				}
			}
			if len(service.calls) != 0 {
				t.Errorf("calls = %v, want no action", service.calls)
			}
		})
	}
}

func TestOneClick_SigningNotConfigured(t *testing.T) {
	signer := signedurl.New("test-secret", "https://api.mailx.net")
	service := &fakeOneClickService{}
	app := newOneClickTestApp(service)

	target := signedTarget(t, signer, model.OneClickDeactivatePath+"alias-1", nil, time.Hour)
	if status, body := oneClickPost(t, app, target); status != http.StatusForbidden || body != "" {
		t.Errorf("POST = %d %q, want 403 with an empty body", status, body)
	}
}

func TestOneClick_MethodNotAllowed(t *testing.T) {
	signer := signedurl.New("test-secret", "https://api.mailx.net")
	methods := []string{
		http.MethodGet,
		http.MethodHead,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodOptions,
	}

	for _, route := range oneClickRoutes {
		t.Run(route.name, func(t *testing.T) {
			service := &fakeOneClickService{signer: signer}
			app := newOneClickTestApp(service)
			target := signedTarget(t, signer, route.path, route.params, time.Hour)

			// Non-POST methods must not fall through to the authenticated /v1 group (401)
			for _, method := range methods {
				res, err := app.Test(httptest.NewRequest(method, target, nil))
				if err != nil {
					t.Fatalf("app.Test() error = %v", err)
				}
				body, _ := io.ReadAll(res.Body)
				if res.StatusCode != http.StatusMethodNotAllowed || len(body) != 0 {
					t.Errorf("%s = %d %q, want 405 with an empty body", method, res.StatusCode, body)
				}
				if allow := res.Header.Get(fiber.HeaderAllow); allow != fiber.MethodPost {
					t.Errorf("%s Allow = %q, want %q", method, allow, fiber.MethodPost)
				}
			}
			if len(service.calls) != 0 {
				t.Errorf("calls = %v, want no action on a non-POST method", service.calls)
			}
		})
	}
}

func TestOneClick_ErrorStatuses(t *testing.T) {
	signer := signedurl.New("test-secret", "https://api.mailx.net")

	tests := []struct {
		err  error
		want int
	}{
		{model.ErrOneClickAliasNotFound, http.StatusNotFound},
		{model.ErrOneClickInvalidValue, http.StatusUnprocessableEntity},
		{model.ErrOneClickAction, http.StatusInternalServerError},
	}
	for _, tt := range tests {
		service := &fakeOneClickService{signer: signer, err: tt.err}
		app := newOneClickTestApp(service)

		// An invalid value that was signed still reaches validation
		target := signedTarget(t, signer, model.OneClickBlockEmailPath+"alias-1", url.Values{"email": {"not-an-email"}}, time.Hour)
		if status, body := oneClickPost(t, app, target); status != tt.want || body != "" {
			t.Errorf("%v: POST = %d %q, want %d with an empty body", tt.err, status, body, tt.want)
		}
	}
}

func TestOneClick_RateLimited(t *testing.T) {
	signer := signedurl.New("test-secret", "https://api.mailx.net")
	service := &fakeOneClickService{signer: signer}
	app := newOneClickTestApp(service)

	// The limit is shared across the one-click routes
	for i := range 6 {
		route := oneClickRoutes[i%len(oneClickRoutes)]
		if status, _ := oneClickPost(t, app, signedTarget(t, signer, route.path, route.params, time.Hour)); status != http.StatusOK {
			t.Fatalf("request %d = %d, want 200", i+1, status)
		}
	}

	target := signedTarget(t, signer, model.OneClickDeactivatePath+"alias-1", nil, time.Hour)
	if status, body := oneClickPost(t, app, target); status != http.StatusTooManyRequests || body != "" {
		t.Errorf("7th request = %d %q, want 429 with an empty body", status, body)
	}
}

func TestUpdateListUnsubscribeBehaviour(t *testing.T) {
	service := &fakeOneClickService{behaviour: map[string]model.ListUnsubscribeBehaviour{}}
	h := &Handler{Service: service, Server: fiber.New(), Validator: utils.NewValidator()}
	h.Server.Put("/settings/list-unsubscribe-behaviour", func(c *fiber.Ctx) error {
		c.Locals(auth.USER_ID, "user-1")
		return c.Next()
	}, h.UpdateListUnsubscribeBehaviour)

	put := func(body string) (int, string) {
		req := httptest.NewRequest(http.MethodPut, "/settings/list-unsubscribe-behaviour", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		return doRequest(t, h.Server, req)
	}

	for b := model.ListUnsubscribeOriginalWithFallback; b <= model.ListUnsubscribeBlockDomain; b++ {
		body := `{"list_unsubscribe_behaviour": ` + string(rune('0'+b)) + `}`
		if status, res := put(body); status != http.StatusOK {
			t.Fatalf("PUT %s = %d %s, want 200", body, status, res)
		}
		if got := service.behaviour["user-1"]; got != b {
			t.Errorf("stored behaviour = %d, want %d", got, b)
		}
	}

	for _, body := range []string{
		`{"list_unsubscribe_behaviour": 9}`,
		`{"list_unsubscribe_behaviour": -1}`,
		`{"list_unsubscribe_behaviour": "3"}`,
		`{"list_unsubscribe_behaviour": null}`,
		`{}`,
		`not json`,
	} {
		if status, res := put(body); status != http.StatusUnprocessableEntity || !strings.Contains(res, `"error"`) {
			t.Errorf("PUT %s = %d %s, want 422 with an error", body, status, res)
		}
		if got := service.behaviour["user-1"]; got != model.ListUnsubscribeBlockDomain {
			t.Errorf("PUT %s changed the stored behaviour to %d", body, got)
		}
	}
}
