package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"ivpn.net/email/api/internal/middleware/auth"
	"ivpn.net/email/api/internal/model"
	"ivpn.net/email/api/internal/utils"
)

// fakeSettingsService embeds Service (as a nil interface) so only the methods under test need
// implementing. It stores each toggle by setting name.
type fakeSettingsService struct {
	Service
	blockTrackers map[string]bool
	toggles       map[string]bool
}

func (f *fakeSettingsService) UpdateBlockTrackers(ctx context.Context, userID string, enabled bool) error {
	f.blockTrackers[userID] = enabled
	f.toggles["block_trackers"] = enabled
	return nil
}

func (f *fakeSettingsService) UpdateRemoveHeader(ctx context.Context, userID string, removeHeader bool) error {
	f.toggles["remove_header"] = removeHeader
	return nil
}

func (f *fakeSettingsService) UpdateLogIssues(ctx context.Context, userID string, enabled bool) error {
	f.toggles["log_issues"] = enabled
	return nil
}

func (f *fakeSettingsService) GetDefaults(ctx context.Context, userID string) (model.Settings, []string, error) {
	return model.Settings{BlockTrackers: f.blockTrackers[userID]}, nil, nil
}

func (f *fakeSettingsService) GetVerifiedDomains(ctx context.Context, userID string) ([]model.Domain, error) {
	return nil, nil
}

func newFakeSettingsService() *fakeSettingsService {
	return &fakeSettingsService{blockTrackers: map[string]bool{"user-1": true}, toggles: map[string]bool{}}
}

func newSettingsTestApp(service *fakeSettingsService) *fiber.App {
	h := &Handler{Service: service, Server: fiber.New(), Validator: utils.NewValidator()}
	withUser := func(c *fiber.Ctx) error {
		c.Locals(auth.USER_ID, "user-1")
		return c.Next()
	}
	h.Server.Put("/settings/block-trackers", withUser, h.UpdateBlockTrackers)
	h.Server.Put("/settings/remove-header", withUser, h.UpdateRemoveHeader)
	h.Server.Put("/settings/log-issues", withUser, h.UpdateLogIssues)
	h.Server.Get("/defaults", withUser, h.GetDefaults)

	return h.Server
}

func TestUpdateBoolSettings(t *testing.T) {
	tests := []struct {
		path        string
		field       string
		success     string
		errRequired string
		errInvalid  string
	}{
		{"/settings/block-trackers", "block_trackers", UpdateBlockTrackersSuccess, ErrBlockTrackersRequired, ErrBlockTrackersInvalid},
		{"/settings/remove-header", "remove_header", UpdateRemoveHeaderSuccess, ErrRemoveHeaderRequired, ErrRemoveHeaderInvalid},
		{"/settings/log-issues", "log_issues", UpdateLogIssuesSuccess, ErrLogIssuesRequired, ErrLogIssuesInvalid},
	}

	for _, tt := range tests {
		t.Run(tt.field, func(t *testing.T) {
			service := newFakeSettingsService()
			app := newSettingsTestApp(service)

			put := func(body string) (int, string) {
				req := httptest.NewRequest(http.MethodPut, tt.path, strings.NewReader(body))
				req.Header.Set("Content-Type", "application/json")
				return doRequest(t, app, req)
			}

			for _, want := range []bool{true, false, true} {
				body := `{"` + tt.field + `": ` + strconv.FormatBool(want) + `}`
				if status, res := put(body); status != http.StatusOK || !strings.Contains(res, tt.success) {
					t.Fatalf("PUT %s = %d %s, want 200 %q", body, status, res, tt.success)
				}
				if got := service.toggles[tt.field]; got != want {
					t.Errorf("stored %s = %v, want %v", tt.field, got, want)
				}
			}

			for _, bad := range []struct {
				body string
				want string
			}{
				{`{"` + tt.field + `": "false"}`, tt.errInvalid},
				{`{"` + tt.field + `": 0}`, tt.errInvalid},
				{`not json`, tt.errInvalid},
				{`{"` + tt.field + `": null}`, tt.errRequired},
				{`{}`, tt.errRequired},
			} {
				if status, res := put(bad.body); status != http.StatusUnprocessableEntity || !strings.Contains(res, bad.want) {
					t.Errorf("PUT %s = %d %s, want 422 %q", bad.body, status, res, bad.want)
				}
				if !service.toggles[tt.field] {
					t.Errorf("PUT %s changed the stored setting", bad.body)
				}
			}
		})
	}
}

func TestGetDefaultsIncludesBlockTrackers(t *testing.T) {
	app := newSettingsTestApp(newFakeSettingsService())

	status, res := doRequest(t, app, httptest.NewRequest(http.MethodGet, "/defaults", nil))
	if status != http.StatusOK || !strings.Contains(res, `"block_trackers":true`) {
		t.Errorf("GET /defaults = %d %s, want block_trackers true", status, res)
	}
}
