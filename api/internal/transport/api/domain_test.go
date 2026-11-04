package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"ivpn.net/email/api/config"
	"ivpn.net/email/api/internal/middleware/auth"
	"ivpn.net/email/api/internal/model"
	"ivpn.net/email/api/internal/utils"
)

const testDomainID = "00e60d58-7565-430c-8bc5-d14e76b6f3a0"

// fakeDomainService embeds Service (as a nil interface) so only the methods under test need
// implementing. It keeps the last saved domain.
type fakeDomainService struct {
	Service
	saved *model.Domain
}

func (f *fakeDomainService) GetDomain(ctx context.Context, domainID string, userID string) (model.Domain, error) {
	return model.Domain{BaseModel: model.BaseModel{ID: domainID}, UserID: userID, Name: "custom.example"}, nil
}

func (f *fakeDomainService) UpdateDomain(ctx context.Context, domain model.Domain) error {
	f.saved = &domain
	return nil
}

func (f *fakeDomainService) GetVerifiedRecipients(ctx context.Context, emails string, userID string) ([]model.Recipient, error) {
	var rcps []model.Recipient
	for email := range strings.SplitSeq(emails, ",") {
		if email == "me@example.com" || email == "work@example.com" {
			rcps = append(rcps, model.Recipient{Email: email})
		}
	}
	return rcps, nil
}

func newDomainTestApp(service *fakeDomainService) *fiber.App {
	h := &Handler{Cfg: config.APIConfig{Domains: "mailx.net,mailx.io"}, Service: service, Server: fiber.New(), Validator: utils.NewValidator()}
	withUser := func(c *fiber.Ctx) error {
		c.Locals(auth.USER_ID, "user-1")
		return c.Next()
	}
	h.Server.Put("/domain/:id", withUser, h.UpdateDomain)

	return h.Server
}

func TestUpdateDomainRecipient(t *testing.T) {
	tests := []struct {
		recipient string
		wantSaved bool
	}{
		{"", true},
		{"me@example.com", true},
		{"stranger@example.com", false},
		{"me@example.com,work@example.com", false},
		{"me@example.com,stranger@example.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.recipient, func(t *testing.T) {
			service := &fakeDomainService{}
			app := newDomainTestApp(service)

			body := `{"id": "` + testDomainID + `", "recipient": "` + tt.recipient + `", "enabled": true, "catch_all": true}`
			req := httptest.NewRequest(http.MethodPut, "/domain/"+testDomainID, strings.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			status, res := doRequest(t, app, req)

			if tt.wantSaved {
				if status != http.StatusOK || !strings.Contains(res, UpdateDomainSuccess) {
					t.Fatalf("PUT %s = %d %s, want 200 %q", body, status, res, UpdateDomainSuccess)
				}
				if service.saved == nil || service.saved.Recipient != tt.recipient {
					t.Errorf("saved domain = %+v, want recipient %q", service.saved, tt.recipient)
				}
				return
			}

			if status != http.StatusBadRequest || !strings.Contains(res, ErrDomainRecipientInvalid) {
				t.Errorf("PUT %s = %d %s, want 400 %q", body, status, res, ErrDomainRecipientInvalid)
			}
			if service.saved != nil {
				t.Errorf("PUT %s saved the domain", body)
			}
		})
	}
}
