package api

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"ivpn.net/email/api/internal/middleware/auth"
	"ivpn.net/email/api/internal/utils"
)

const testAliasID = "3f1c2b4a-5d6e-4f70-8a9b-0c1d2e3f4a5b"

var errInactiveSub = errors.New("Your subscription is not active. Please renew to pause aliases.")

// fakePauseService records the pauses the handlers pass on.
type fakePauseService struct {
	Service
	paused   []AliasPauseReq
	unpaused []string
	err      error
}

func (f *fakePauseService) PauseAlias(ctx context.Context, ID string, userID string, value int, unit string) (time.Time, error) {
	if f.err != nil {
		return time.Time{}, f.err
	}
	f.paused = append(f.paused, AliasPauseReq{Value: value, Unit: unit})
	return time.Date(2026, 10, 15, 8, 0, 0, 0, time.UTC), nil
}

func (f *fakePauseService) UnpauseAlias(ctx context.Context, ID string, userID string) error {
	if f.err != nil {
		return f.err
	}
	f.unpaused = append(f.unpaused, ID)
	return nil
}

func newPauseTestApp(service *fakePauseService) *fiber.App {
	h := &Handler{Service: service, Server: fiber.New(), Validator: utils.NewValidator()}
	withUser := func(c *fiber.Ctx) error {
		c.Locals(auth.USER_ID, "user-1")
		return c.Next()
	}
	h.Server.Put("/alias/:id/pause", withUser, h.PauseAlias)
	h.Server.Delete("/alias/:id/pause", withUser, h.UnpauseAlias)

	return h.Server
}

func TestPauseAliasHandler(t *testing.T) {
	service := &fakePauseService{}
	app := newPauseTestApp(service)

	status, data := doRuleRequest(t, app, http.MethodPut, "/alias/"+testAliasID+"/pause", `{"value": 3, "unit": "hours"}`)
	if status != 200 || data["message"] != UpdateAliasSuccess || data["paused_until"] != "2026-10-15T08:00:00Z" {
		t.Fatalf("PUT pause = %d %+v", status, data)
	}
	if len(service.paused) != 1 || service.paused[0] != (AliasPauseReq{Value: 3, Unit: "hours"}) {
		t.Errorf("paused = %+v", service.paused)
	}

	for _, body := range []string{`{"value": 1, "unit": "hours"}`, `{"value": 10, "unit": "days"}`, `{"value": 5, "unit": "weeks"}`} {
		if status, data := doRuleRequest(t, app, http.MethodPut, "/alias/"+testAliasID+"/pause", body); status != 200 {
			t.Errorf("PUT pause %s = %d %+v", body, status, data)
		}
	}

	invalid := map[string]string{
		"value 0":       `{"value": 0, "unit": "days"}`,
		"value 11":      `{"value": 11, "unit": "days"}`,
		"unknown unit":  `{"value": 1, "unit": "months"}`,
		"missing value": `{"unit": "days"}`,
		"missing unit":  `{"value": 1}`,
		"empty body":    `{}`,
		"wrong type":    `{"value": "1", "unit": "days"}`,
	}
	for name, body := range invalid {
		t.Run(name, func(t *testing.T) {
			before := len(service.paused)
			status, data := doRuleRequest(t, app, http.MethodPut, "/alias/"+testAliasID+"/pause", body)
			if status != 400 || data["error"] == nil || data["error"] == "" {
				t.Errorf("PUT pause %s = %d %+v, want 400 with an error", body, status, data)
			}
			if len(service.paused) != before {
				t.Error("an invalid pause reached the service")
			}
		})
	}

	t.Run("service errors are returned", func(t *testing.T) {
		service.err = errInactiveSub
		defer func() { service.err = nil }()

		status, data := doRuleRequest(t, app, http.MethodPut, "/alias/"+testAliasID+"/pause", `{"value": 1, "unit": "days"}`)
		if status != 400 || data["error"] != errInactiveSub.Error() {
			t.Errorf("PUT pause = %d %+v", status, data)
		}
	})
}

func TestUnpauseAliasHandler(t *testing.T) {
	service := &fakePauseService{}
	app := newPauseTestApp(service)

	status, data := doRuleRequest(t, app, http.MethodDelete, "/alias/"+testAliasID+"/pause", "")
	if status != 200 || data["message"] != UpdateAliasSuccess {
		t.Fatalf("DELETE pause = %d %+v", status, data)
	}
	if len(service.unpaused) != 1 || service.unpaused[0] != testAliasID {
		t.Errorf("unpaused = %+v", service.unpaused)
	}

	service.err = errInactiveSub
	status, data = doRuleRequest(t, app, http.MethodDelete, "/alias/"+testAliasID+"/pause", "")
	if status != 400 || data["error"] != errInactiveSub.Error() {
		t.Errorf("DELETE pause = %d %+v", status, data)
	}
}
