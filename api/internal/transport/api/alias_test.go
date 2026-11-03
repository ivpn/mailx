package api

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"ivpn.net/email/api/internal/middleware/auth"
	"ivpn.net/email/api/internal/model"
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

// fakeImportService records the import rows the handler passes on and skips the ones on
// skipDomain, like the service does for a domain that isn't the user's.
type fakeImportService struct {
	Service
	rows       []model.AliasImportReq
	skipDomain string
}

func (f *fakeImportService) ImportAliases(ctx context.Context, reqs []model.AliasImportReq, userID string) ([]model.Alias, []model.AliasImportSkip, error) {
	f.rows = reqs
	var aliases []model.Alias
	var skipped []model.AliasImportSkip
	for _, req := range reqs {
		name := req.LocalPart + "@" + req.Domain
		if req.Domain == f.skipDomain {
			skipped = append(skipped, model.AliasImportSkip{Row: req.Row, Alias: name, Reason: "domain"})
			continue
		}
		aliases = append(aliases, model.Alias{Name: name})
	}
	return aliases, skipped, nil
}

func doImportRequest(t *testing.T, service *fakeImportService, csv string) (int, map[string]any) {
	t.Helper()

	h := &Handler{Service: service, Server: fiber.New(), Validator: utils.NewValidator()}
	h.Server.Post("/aliases/import", func(c *fiber.Ctx) error {
		c.Locals(auth.USER_ID, "user-1")
		return c.Next()
	}, h.ImportAliases)

	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	part, err := w.CreateFormFile("file", "aliases.csv")
	if err != nil {
		t.Fatalf("CreateFormFile() error = %v", err)
	}
	_, _ = part.Write([]byte(csv))
	_ = w.Close()

	req := httptest.NewRequest(http.MethodPost, "/aliases/import", &body)
	req.Header.Set("Content-Type", w.FormDataContentType())
	res, err := h.Server.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	defer res.Body.Close()

	raw, _ := io.ReadAll(res.Body)
	var data map[string]any
	_ = json.Unmarshal(raw, &data)

	return res.StatusCode, data
}

func TestImportAliasesHandler(t *testing.T) {
	t.Run("accepts the same names as alias creation", func(t *testing.T) {
		service := &fakeImportService{}
		csv := "alias,description,enabled,recipients\n" +
			"averyveryverylongaliasname.29@example.net,,true,rcpt@example.net\n" +
			"short@example.net,,true,rcpt@example.net\n" +
			"a@example.net,,true,rcpt@example.net\n" +
			strings.Repeat("x", 64) + "@example.net,,true,rcpt@example.net\n"

		status, data := doImportRequest(t, service, csv)
		if status != 200 || data["count"] != float64(4) || data["message"] != AliasImportSuccess {
			t.Fatalf("POST import = %d %+v", status, data)
		}
		if skipped, _ := data["skipped"].([]any); skipped == nil || len(skipped) != 0 {
			t.Errorf("skipped = %+v, want an empty list", data["skipped"])
		}
		if len(service.rows) != 4 || service.rows[0].Row != 2 || service.rows[3].Row != 5 {
			t.Errorf("rows passed to the service = %+v", service.rows)
		}
	})

	t.Run("reports skipped rows with their reason in row order", func(t *testing.T) {
		service := &fakeImportService{skipDomain: "other.net"}
		csv := "alias,description,enabled,recipients\n" +
			"valid@example.net,,true,rcpt@example.net\n" +
			strings.Repeat("x", 65) + "@example.net,,true,rcpt@example.net\n" +
			"no-at-sign,,true,rcpt@example.net\n" +
			"skipped@other.net,,true,rcpt@example.net\n" +
			"bad..name@example.net,,true,rcpt@example.net\n" +
			"norcpt@example.net,,true,\n"

		status, data := doImportRequest(t, service, csv)
		if status != 200 || data["count"] != float64(1) || data["message"] != AliasImportPartial {
			t.Fatalf("POST import = %d %+v", status, data)
		}

		want := []map[string]any{
			{"row": float64(3), "alias": strings.Repeat("x", 65) + "@example.net", "reason": ErrImportLocalPart},
			{"row": float64(4), "alias": "no-at-sign", "reason": ErrImportAddress},
			{"row": float64(5), "alias": "skipped@other.net", "reason": "domain"},
			{"row": float64(6), "alias": "bad..name@example.net", "reason": ErrImportLocalPart},
			{"row": float64(7), "alias": "norcpt@example.net", "reason": ErrImportNoRcp},
		}
		skipped, _ := data["skipped"].([]any)
		if len(skipped) != len(want) {
			t.Fatalf("skipped = %+v, want %+v", skipped, want)
		}
		for i, w := range want {
			got, _ := skipped[i].(map[string]any)
			for k, v := range w {
				if got[k] != v {
					t.Errorf("skipped[%d][%s] = %v, want %v", i, k, got[k], v)
				}
			}
		}
	})
}
