package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/ProtonMail/gopenpgp/v3/crypto"
	"github.com/gofiber/fiber/v2"
	"ivpn.net/email/api/internal/middleware/auth"
	"ivpn.net/email/api/internal/model"
	"ivpn.net/email/api/internal/utils"
)

const testRecipientID = "3f1c2b4a-5d6e-4f70-8a9b-0c1d2e3f4a5b"

// fakeRecipientService records the recipients the handler saves.
type fakeRecipientService struct {
	Service
	recipient model.Recipient
	updated   []model.Recipient
}

func (f *fakeRecipientService) GetRecipient(ctx context.Context, ID string, userID string) (model.Recipient, error) {
	return f.recipient, nil
}

func (f *fakeRecipientService) UpdateRecipient(ctx context.Context, rcp model.Recipient) error {
	f.updated = append(f.updated, rcp)
	return nil
}

func newRecipientTestApp(service *fakeRecipientService) *fiber.App {
	h := &Handler{Service: service, Server: fiber.New(), Validator: utils.NewValidator()}
	withUser := func(c *fiber.Ctx) error {
		c.Locals(auth.USER_ID, "user-1")
		return c.Next()
	}
	h.Server.Put("/recipient", withUser, h.UpdateRecipient)

	return h.Server
}

func testArmoredPublicKey(t *testing.T) string {
	t.Helper()

	key, err := crypto.PGP().KeyGeneration().AddUserId("test", "test@example.com").New().GenerateKey()
	if err != nil {
		t.Fatalf("GenerateKey() error = %v", err)
	}
	public, err := key.GetArmoredPublicKey()
	if err != nil {
		t.Fatalf("GetArmoredPublicKey() error = %v", err)
	}

	return strings.TrimSpace(public)
}

func recipientUpdateBody(t *testing.T, pgpKey string, pgpEnabled bool) string {
	t.Helper()

	body, err := json.Marshal(RecipientReq{ID: testRecipientID, PGPKey: pgpKey, PGPEnabled: pgpEnabled})
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	return string(body)
}

func TestUpdateRecipientRejectsInvalidPGPKey(t *testing.T) {
	tests := []struct {
		desc string
		key  string
	}{
		{"armor with text instead of a key (issue #517)", "-----BEGIN PGP PUBLIC KEY BLOCK-----\nnot a key\n-----END PGP PUBLIC KEY BLOCK-----"},
		{"armor with no key data", "-----BEGIN PGP PUBLIC KEY BLOCK-----\n\n-----END PGP PUBLIC KEY BLOCK-----"},
		{"armor without footer", "-----BEGIN PGP PUBLIC KEY BLOCK-----\nnot a key"},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			service := &fakeRecipientService{recipient: model.Recipient{UserID: "user-1"}}
			app := newRecipientTestApp(service)

			status, data := doRuleRequest(t, app, http.MethodPut, "/recipient", recipientUpdateBody(t, tt.key, true))

			if status != 400 || data["error"] != ErrInvalidPGPKey {
				t.Fatalf("UpdateRecipient() = %d %v, want 400 %q", status, data, ErrInvalidPGPKey)
			}
			if len(service.updated) != 0 {
				t.Errorf("UpdateRecipient() saved %+v, want nothing saved", service.updated)
			}
		})
	}
}

func TestUpdateRecipientAcceptsValidPGPKey(t *testing.T) {
	key := testArmoredPublicKey(t)
	service := &fakeRecipientService{recipient: model.Recipient{UserID: "user-1"}}
	app := newRecipientTestApp(service)

	status, data := doRuleRequest(t, app, http.MethodPut, "/recipient", recipientUpdateBody(t, key, true))

	if status != 200 || data["message"] != UpdateRecipientSuccess {
		t.Fatalf("UpdateRecipient() = %d %v, want 200", status, data)
	}
	if len(service.updated) != 1 || service.updated[0].PGPKey != key || !service.updated[0].PGPEnabled {
		t.Errorf("UpdateRecipient() saved %+v, want key saved with encryption on", service.updated)
	}
}

func TestUpdateRecipientKeepsKeyWhenHashSent(t *testing.T) {
	key := testArmoredPublicKey(t)
	service := &fakeRecipientService{recipient: model.Recipient{UserID: "user-1", PGPKey: key, PGPEnabled: true}}
	app := newRecipientTestApp(service)

	// The app sends the hashed key back when toggling encryption
	status, data := doRuleRequest(t, app, http.MethodPut, "/recipient", recipientUpdateBody(t, utils.HashPGPKey(key), false))

	if status != 200 {
		t.Fatalf("UpdateRecipient() = %d %v, want 200", status, data)
	}
	if len(service.updated) != 1 || service.updated[0].PGPKey != key || service.updated[0].PGPEnabled {
		t.Errorf("UpdateRecipient() saved %+v, want stored key kept with encryption off", service.updated)
	}
}
