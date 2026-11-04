package api

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"ivpn.net/email/api/internal/middleware/auth"
	"ivpn.net/email/api/internal/model"
	"ivpn.net/email/api/internal/utils"
)

const testRuleID = "8d6c0f43-3a5b-4b8e-9f3e-1c2d3e4f5a6b"

// fakeRuleService records what the handlers pass on. Rules of users other than user-1 are
// "not found", like the real service.
type fakeRuleService struct {
	Service
	rules     []model.Rule
	posted    []model.Rule
	updated   []model.Rule
	active    map[string]bool
	reordered [][]string
	deleted   []string
	err       error
}

func (f *fakeRuleService) GetRules(ctx context.Context, userID string) ([]model.Rule, error) {
	return f.rules, f.err
}

func (f *fakeRuleService) PostRule(ctx context.Context, rule model.Rule) (model.Rule, error) {
	if f.err != nil {
		return model.Rule{}, f.err
	}
	f.posted = append(f.posted, rule)
	rule.ID = testRuleID
	return rule, nil
}

func (f *fakeRuleService) UpdateRule(ctx context.Context, rule model.Rule) (model.Rule, error) {
	if f.err != nil {
		return model.Rule{}, f.err
	}
	if rule.ID != testRuleID {
		return model.Rule{}, model.ErrRuleNotFound
	}
	f.updated = append(f.updated, rule)
	return rule, nil
}

func (f *fakeRuleService) UpdateRuleActive(ctx context.Context, ID string, userID string, active bool) error {
	if f.err != nil {
		return f.err
	}
	if ID != testRuleID {
		return model.ErrRuleNotFound
	}
	f.active[ID] = active
	return nil
}

func (f *fakeRuleService) ReorderRules(ctx context.Context, userID string, ids []string) error {
	if f.err != nil {
		return f.err
	}
	f.reordered = append(f.reordered, ids)
	return nil
}

func (f *fakeRuleService) DeleteRule(ctx context.Context, ID string, userID string) error {
	if f.err != nil {
		return f.err
	}
	if ID != testRuleID {
		return model.ErrRuleNotFound
	}
	f.deleted = append(f.deleted, ID)
	return nil
}

func newRuleTestApp(service *fakeRuleService) *fiber.App {
	h := &Handler{Service: service, Server: fiber.New(), Validator: utils.NewValidator()}
	withUser := func(c *fiber.Ctx) error {
		c.Locals(auth.USER_ID, "user-1")
		return c.Next()
	}
	h.Server.Get("/rules", withUser, h.GetRules)
	h.Server.Put("/rules/order", withUser, h.ReorderRules)
	h.Server.Post("/rule", withUser, h.PostRule)
	h.Server.Put("/rule/:id", withUser, h.UpdateRule)
	h.Server.Put("/rule/:id/active", withUser, h.UpdateRuleActive)
	h.Server.Delete("/rule/:id", withUser, h.DeleteRule)

	return h.Server
}

func doRuleRequest(t *testing.T, app *fiber.App, method string, path string, body string) (int, map[string]any) {
	t.Helper()

	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	defer res.Body.Close()

	raw, _ := io.ReadAll(res.Body)
	var data map[string]any
	_ = json.Unmarshal(raw, &data)

	return res.StatusCode, data
}

const testRuleBody = `{
	"name": "Bills",
	"operator": "OR",
	"conditions": [
		{"type": "subject", "match": "contains", "values": ["invoice", "receipt"]},
		{"type": "is_mailing_list"}
	],
	"actions": [
		{"type": "forward_to", "value": "rcp-1,rcp-2"},
		{"type": "stop"}
	],
	"forwards": true,
	"replies": false,
	"sends": true
}`

func TestGetRules(t *testing.T) {
	service := &fakeRuleService{rules: []model.Rule{{BaseModel: model.BaseModel{ID: testRuleID}, Name: "Bills", InvalidRecipients: true, Applied: 3}}}
	app := newRuleTestApp(service)

	req := httptest.NewRequest(http.MethodGet, "/rules", nil)
	res, err := app.Test(req)
	if err != nil {
		t.Fatalf("app.Test() error = %v", err)
	}
	defer res.Body.Close()

	var rules []map[string]any
	if err := json.NewDecoder(res.Body).Decode(&rules); err != nil {
		t.Fatalf("decode error = %v", err)
	}
	if res.StatusCode != 200 || len(rules) != 1 || rules[0]["name"] != "Bills" || rules[0]["invalid_recipients"] != true || rules[0]["applied"] != float64(3) {
		t.Errorf("GET /rules = %d %+v", res.StatusCode, rules)
	}
	if _, ok := rules[0]["user_id"]; ok {
		t.Error("user_id exposed")
	}

	service.err = errors.New("Unable to retrieve rules.")
	status, data := doRuleRequest(t, app, http.MethodGet, "/rules", "")
	if status != 400 || data["error"] != "Unable to retrieve rules." {
		t.Errorf("GET /rules on error = %d %+v", status, data)
	}
}

func TestPostRule(t *testing.T) {
	service := &fakeRuleService{}
	app := newRuleTestApp(service)

	status, data := doRuleRequest(t, app, http.MethodPost, "/rule", testRuleBody)
	if status != 201 || data["message"] != PostRuleSuccess {
		t.Fatalf("POST /rule = %d %+v", status, data)
	}
	if rule, _ := data["rule"].(map[string]any); rule["id"] != testRuleID {
		t.Errorf("response rule = %+v", data["rule"])
	}

	got := service.posted[0]
	if got.UserID != "user-1" || got.Name != "Bills" || got.Operator != model.RuleOr || !got.Forwards || got.Replies || !got.Sends {
		t.Errorf("posted rule = %+v", got)
	}
	if len(got.Conditions) != 2 || got.Conditions[0].Values[1] != "receipt" || got.Conditions[1].Type != model.RuleConditionIsMailingList {
		t.Errorf("posted conditions = %+v", got.Conditions)
	}
	if len(got.Actions) != 2 || got.Actions[0].Value != "rcp-1,rcp-2" || got.Actions[1].Type != model.RuleActionStop {
		t.Errorf("posted actions = %+v", got.Actions)
	}

	t.Run("service errors are returned", func(t *testing.T) {
		service.err = errors.New("You’ve reached the maximum number of rules.")
		defer func() { service.err = nil }()

		status, data := doRuleRequest(t, app, http.MethodPost, "/rule", testRuleBody)
		if status != 400 || data["error"] != "You’ve reached the maximum number of rules." {
			t.Errorf("POST /rule = %d %+v", status, data)
		}
	})

	invalid := map[string]string{
		"malformed JSON":         `{"name": `,
		"missing name":           `{"operator": "AND", "conditions": [{"type": "subject"}], "actions": [{"type": "stop"}], "forwards": true}`,
		"bad operator":           strings.Replace(testRuleBody, `"OR"`, `"XOR"`, 1),
		"no conditions":          `{"name": "x", "operator": "AND", "actions": [{"type": "stop"}], "forwards": true}`,
		"condition without type": strings.Replace(testRuleBody, `"type": "is_mailing_list"`, `"type": ""`, 1),
		"action without type":    strings.Replace(testRuleBody, `"type": "stop"`, `"type": ""`, 1),
		"values not strings":     strings.Replace(testRuleBody, `["invoice", "receipt"]`, `[1, 2]`, 1),
	}
	for name, body := range invalid {
		t.Run(name, func(t *testing.T) {
			before := len(service.posted)
			status, data := doRuleRequest(t, app, http.MethodPost, "/rule", body)
			if status != 400 || data["error"] != ErrInvalidRequest {
				t.Errorf("POST /rule = %d %+v", status, data)
			}
			if len(service.posted) != before {
				t.Error("an invalid request reached the service")
			}
		})
	}
}

func TestUpdateRule(t *testing.T) {
	service := &fakeRuleService{}
	app := newRuleTestApp(service)

	status, data := doRuleRequest(t, app, http.MethodPut, "/rule/"+testRuleID, testRuleBody)
	if status != 200 || data["message"] != UpdateRuleSuccess {
		t.Fatalf("PUT /rule/:id = %d %+v", status, data)
	}
	if got := service.updated[0]; got.ID != testRuleID || got.UserID != "user-1" || got.Name != "Bills" {
		t.Errorf("updated rule = %+v", got)
	}

	status, data = doRuleRequest(t, app, http.MethodPut, "/rule/another-users-rule", testRuleBody)
	if status != 404 || data["error"] != model.ErrRuleNotFound.Error() {
		t.Errorf("PUT another user's rule = %d %+v, want 404", status, data)
	}

	status, _ = doRuleRequest(t, app, http.MethodPut, "/rule/"+testRuleID, `{"name": "x"}`)
	if status != 400 {
		t.Errorf("PUT with an invalid body = %d, want 400", status)
	}
}

func TestUpdateRuleActive(t *testing.T) {
	service := &fakeRuleService{active: map[string]bool{}}
	app := newRuleTestApp(service)

	status, data := doRuleRequest(t, app, http.MethodPut, "/rule/"+testRuleID+"/active", `{"active": false}`)
	if status != 200 || data["message"] != UpdateRuleActiveOff || service.active[testRuleID] {
		t.Errorf("deactivate = %d %+v %v", status, data, service.active)
	}

	status, data = doRuleRequest(t, app, http.MethodPut, "/rule/"+testRuleID+"/active", `{"active": true}`)
	if status != 200 || data["message"] != UpdateRuleActiveOn || !service.active[testRuleID] {
		t.Errorf("activate = %d %+v %v", status, data, service.active)
	}

	tests := []struct {
		path, body, err string
		status          int
	}{
		{"/rule/" + testRuleID + "/active", `{}`, ErrRuleActiveRequired, 400},
		{"/rule/" + testRuleID + "/active", `{"active": "yes"}`, ErrRuleActiveInvalid, 400},
		{"/rule/other/active", `{"active": true}`, model.ErrRuleNotFound.Error(), 404},
	}
	for _, tt := range tests {
		status, data := doRuleRequest(t, app, http.MethodPut, tt.path, tt.body)
		if status != tt.status || data["error"] != tt.err {
			t.Errorf("PUT %s %s = %d %+v, want %d %q", tt.path, tt.body, status, data, tt.status, tt.err)
		}
	}

	service.err = errors.New("Your subscription is not active. Please renew to update rules.")
	status, data = doRuleRequest(t, app, http.MethodPut, "/rule/"+testRuleID+"/active", `{"active": true}`)
	if status != 400 || data["error"] != service.err.Error() {
		t.Errorf("toggle with an inactive subscription = %d %+v", status, data)
	}
}

func TestReorderRules(t *testing.T) {
	service := &fakeRuleService{}
	app := newRuleTestApp(service)
	other := "1b4e28ba-2fa1-41d2-883f-0016d3cca427"

	status, data := doRuleRequest(t, app, http.MethodPut, "/rules/order", `{"ids": ["`+other+`", "`+testRuleID+`"]}`)
	if status != 200 || data["message"] != ReorderRulesSuccess {
		t.Fatalf("PUT /rules/order = %d %+v", status, data)
	}
	if got := service.reordered[0]; len(got) != 2 || got[0] != other || got[1] != testRuleID {
		t.Errorf("reordered = %v", got)
	}

	for _, body := range []string{`{}`, `{"ids": ["not-a-uuid"]}`, `{"ids": "x"}`} {
		status, data := doRuleRequest(t, app, http.MethodPut, "/rules/order", body)
		if status != 400 || data["error"] != ErrInvalidRequest {
			t.Errorf("PUT /rules/order %s = %d %+v", body, status, data)
		}
	}

	service.err = errors.New("Your rules have changed since this page was loaded. Please reload and try again.")
	status, data = doRuleRequest(t, app, http.MethodPut, "/rules/order", `{"ids": ["`+testRuleID+`"]}`)
	if status != 400 || data["error"] != service.err.Error() {
		t.Errorf("stale reorder = %d %+v", status, data)
	}
}

func TestDeleteRule(t *testing.T) {
	service := &fakeRuleService{}
	app := newRuleTestApp(service)

	status, data := doRuleRequest(t, app, http.MethodDelete, "/rule/"+testRuleID, "")
	if status != 200 || data["message"] != DeleteRuleSuccess || len(service.deleted) != 1 {
		t.Errorf("DELETE /rule/:id = %d %+v", status, data)
	}

	status, data = doRuleRequest(t, app, http.MethodDelete, "/rule/other", "")
	if status != 404 || data["error"] != model.ErrRuleNotFound.Error() {
		t.Errorf("DELETE another user's rule = %d %+v", status, data)
	}
}
