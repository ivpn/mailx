package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"ivpn.net/email/api/config"
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
	values        map[string]string
}

func (f *fakeSettingsService) UpdateDefaultDomain(ctx context.Context, userID string, domain string) error {
	f.values["domain"] = domain
	return nil
}

func (f *fakeSettingsService) UpdateDefaultRecipient(ctx context.Context, userID string, recipient string) error {
	f.values["recipient"] = recipient
	return nil
}

func (f *fakeSettingsService) UpdateAliasFormat(ctx context.Context, userID string, format string) error {
	f.values["alias_format"] = format
	return nil
}

func (f *fakeSettingsService) UpdateFromName(ctx context.Context, userID string, fromName string) error {
	f.values["from_name"] = fromName
	return nil
}

func (f *fakeSettingsService) GetVerifiedRecipients(ctx context.Context, emails string, userID string) ([]model.Recipient, error) {
	var rcps []model.Recipient
	for email := range strings.SplitSeq(emails, ",") {
		if email == "me@example.com" || email == "work@example.com" {
			rcps = append(rcps, model.Recipient{Email: email})
		}
	}
	return rcps, nil
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
	return []model.Domain{{Name: "custom.example"}}, nil
}

func newFakeSettingsService() *fakeSettingsService {
	return &fakeSettingsService{blockTrackers: map[string]bool{"user-1": true}, toggles: map[string]bool{}, values: map[string]string{}}
}

func newSettingsTestApp(service *fakeSettingsService) *fiber.App {
	h := &Handler{Cfg: config.APIConfig{Domains: "mailx.net,mailx.io"}, Service: service, Server: fiber.New(), Validator: utils.NewValidator()}
	withUser := func(c *fiber.Ctx) error {
		c.Locals(auth.USER_ID, "user-1")
		return c.Next()
	}
	h.Server.Put("/settings/block-trackers", withUser, h.UpdateBlockTrackers)
	h.Server.Put("/settings/remove-header", withUser, h.UpdateRemoveHeader)
	h.Server.Put("/settings/log-issues", withUser, h.UpdateLogIssues)
	h.Server.Put("/settings/domain", withUser, h.UpdateDefaultDomain)
	h.Server.Put("/settings/recipient", withUser, h.UpdateDefaultRecipient)
	h.Server.Put("/settings/alias-format", withUser, h.UpdateAliasFormat)
	h.Server.Put("/settings/from-name", withUser, h.UpdateFromName)
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

func TestUpdateStringSettings(t *testing.T) {
	tests := []struct {
		path        string
		field       string
		success     string
		errRequired string
		errInvalid  string
		valid       map[string]string // request value -> stored value
		invalid     []string
		required    []string
	}{
		{
			path: "/settings/domain", field: "domain", success: UpdateDefaultDomainSuccess,
			errRequired: ErrDefaultDomainRequired, errInvalid: ErrDefaultDomainInvalid,
			valid:    map[string]string{"mailx.net": "mailx.net", "mailx.io": "mailx.io", " custom.example ": "custom.example"},
			invalid:  []string{"mailx", "other.example", "mailx.net,mailx.io"},
			required: []string{"", "   "},
		},
		{
			path: "/settings/recipient", field: "recipient", success: UpdateDefaultRcpSuccess,
			errRequired: ErrDefaultRcpRequired, errInvalid: ErrDefaultRcpInvalid,
			valid:    map[string]string{"me@example.com": "me@example.com", "work@example.com": "work@example.com"},
			invalid:  []string{"stranger@example.com", "me@example.com,work@example.com"},
			required: []string{"", "   "},
		},
		{
			path: "/settings/alias-format", field: "alias_format", success: UpdateAliasFormatSuccess,
			errRequired: ErrAliasFormatRequired, errInvalid: ErrAliasFormatInvalid,
			valid:    map[string]string{"words": "words", "random": "random", "uuid": "uuid"},
			invalid:  []string{"custom", "wildcard", "Words"},
			required: []string{"", "   "},
		},
		{
			path: "/settings/from-name", field: "from_name", success: UpdateFromNameSuccess,
			errRequired: ErrFromNameRequired, errInvalid: ErrFromNameInvalid,
			valid: map[string]string{
				"Jane Doe":                             "Jane Doe",
				" Jane ":                               "Jane",
				"":                                     "",
				"Łukasz Żółć":                          "Łukasz Żółć",
				strings.Repeat("ž", MaxFromNameLength): strings.Repeat("ž", MaxFromNameLength),
			},
			invalid: []string{"Jane\r\nBcc: x@example.com", "Jane\u0000", strings.Repeat("a", MaxFromNameLength+1)},
		},
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
			body := func(v string) string {
				return `{"` + tt.field + `": "` + v + `"}`
			}

			for v, want := range tt.valid {
				if status, res := put(body(v)); status != http.StatusOK || !strings.Contains(res, tt.success) {
					t.Fatalf("PUT %s = %d %s, want 200 %q", body(v), status, res, tt.success)
				}
				if got := service.values[tt.field]; got != want {
					t.Errorf("PUT %s stored %q, want %q", body(v), got, want)
				}
			}

			service.values[tt.field] = "unchanged"
			bad := []struct {
				body string
				want string
			}{
				{`{"` + tt.field + `": 1}`, tt.errInvalid},
				{`not json`, tt.errInvalid},
				{`{"` + tt.field + `": null}`, tt.errRequired},
				{`{}`, tt.errRequired},
			}
			for _, v := range tt.invalid {
				bad = append(bad, struct {
					body string
					want string
				}{body(v), tt.errInvalid})
			}
			for _, v := range tt.required {
				bad = append(bad, struct {
					body string
					want string
				}{body(v), tt.errRequired})
			}
			for _, b := range bad {
				if status, res := put(b.body); status != http.StatusUnprocessableEntity || !strings.Contains(res, b.want) {
					t.Errorf("PUT %s = %d %s, want 422 %q", b.body, status, res, b.want)
				}
				if got := service.values[tt.field]; got != "unchanged" {
					t.Errorf("PUT %s changed the stored value to %q", b.body, got)
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
