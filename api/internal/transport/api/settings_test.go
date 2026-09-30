package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"ivpn.net/email/api/internal/middleware/auth"
	"ivpn.net/email/api/internal/model"
	"ivpn.net/email/api/internal/utils"
)

// fakeSettingsService embeds Service (as a nil interface) so only the methods under test need
// implementing.
type fakeSettingsService struct {
	Service
	blockTrackers map[string]bool
}

func (f *fakeSettingsService) UpdateBlockTrackers(ctx context.Context, userID string, enabled bool) error {
	f.blockTrackers[userID] = enabled
	return nil
}

func (f *fakeSettingsService) GetDefaults(ctx context.Context, userID string) (model.Settings, []string, error) {
	return model.Settings{BlockTrackers: f.blockTrackers[userID]}, nil, nil
}

func (f *fakeSettingsService) GetVerifiedDomains(ctx context.Context, userID string) ([]model.Domain, error) {
	return nil, nil
}

func newSettingsTestApp(service *fakeSettingsService) *fiber.App {
	h := &Handler{Service: service, Server: fiber.New(), Validator: utils.NewValidator()}
	withUser := func(c *fiber.Ctx) error {
		c.Locals(auth.USER_ID, "user-1")
		return c.Next()
	}
	h.Server.Put("/settings/block-trackers", withUser, h.UpdateBlockTrackers)
	h.Server.Get("/defaults", withUser, h.GetDefaults)

	return h.Server
}

func TestUpdateBlockTrackers(t *testing.T) {
	service := &fakeSettingsService{blockTrackers: map[string]bool{"user-1": true}}
	app := newSettingsTestApp(service)

	put := func(body string) (int, string) {
		req := httptest.NewRequest(http.MethodPut, "/settings/block-trackers", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		return doRequest(t, app, req)
	}

	for _, want := range []bool{false, true, false} {
		body := `{"block_trackers": false}`
		if want {
			body = `{"block_trackers": true}`
		}
		if status, res := put(body); status != http.StatusOK || !strings.Contains(res, UpdateBlockTrackersSuccess) {
			t.Fatalf("PUT %s = %d %s, want 200", body, status, res)
		}
		if got := service.blockTrackers["user-1"]; got != want {
			t.Errorf("stored block_trackers = %v, want %v", got, want)
		}
	}

	for _, tt := range []struct {
		body string
		want string
	}{
		{`{"block_trackers": "true"}`, ErrBlockTrackersInvalid},
		{`{"block_trackers": 1}`, ErrBlockTrackersInvalid},
		{`not json`, ErrBlockTrackersInvalid},
		{`{"block_trackers": null}`, ErrBlockTrackersRequired},
		{`{}`, ErrBlockTrackersRequired},
	} {
		if status, res := put(tt.body); status != http.StatusUnprocessableEntity || !strings.Contains(res, tt.want) {
			t.Errorf("PUT %s = %d %s, want 422 %q", tt.body, status, res, tt.want)
		}
		if service.blockTrackers["user-1"] {
			t.Errorf("PUT %s changed the stored setting", tt.body)
		}
	}
}

func TestGetDefaultsIncludesBlockTrackers(t *testing.T) {
	service := &fakeSettingsService{blockTrackers: map[string]bool{"user-1": true}}
	app := newSettingsTestApp(service)

	status, res := doRequest(t, app, httptest.NewRequest(http.MethodGet, "/defaults", nil))
	if status != http.StatusOK || !strings.Contains(res, `"block_trackers":true`) {
		t.Errorf("GET /defaults = %d %s, want block_trackers true", status, res)
	}
}
