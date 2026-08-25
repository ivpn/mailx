package service

import (
	"context"
	"errors"
	"strings"
	"testing"

	"ivpn.net/email/api/config"
	"ivpn.net/email/api/internal/model"
)

var errNotFound = errors.New("not found")

// fakeStore implements the Store interface. It embeds Store (as a nil
// interface) so tests only need to override the methods FindRecipients and
// its helpers actually call; anything else panics if invoked unexpectedly.
type fakeStore struct {
	Store

	aliases            map[string]model.Alias
	domains            map[string]model.Domain
	settingsByUser     map[string]model.Settings
	verifiedRecipients map[string][]model.Recipient
	recipients         map[string][]model.Recipient
	postedMessages     []model.Message
}

func newFakeStore() *fakeStore {
	return &fakeStore{
		aliases:            map[string]model.Alias{},
		domains:            map[string]model.Domain{},
		settingsByUser:     map[string]model.Settings{},
		verifiedRecipients: map[string][]model.Recipient{},
		recipients:         map[string][]model.Recipient{},
	}
}

func (f *fakeStore) GetAliasByName(name string) (model.Alias, error) {
	alias, ok := f.aliases[name]
	if !ok {
		return model.Alias{}, errNotFound
	}
	return alias, nil
}

func (f *fakeStore) GetVerifiedDomainByName(ctx context.Context, name string) (model.Domain, error) {
	domain, ok := f.domains[name]
	if !ok {
		return model.Domain{}, errNotFound
	}
	return domain, nil
}

func (f *fakeStore) GetSettings(ctx context.Context, userID string) (model.Settings, error) {
	return f.settingsByUser[userID], nil
}

func (f *fakeStore) GetVerifiedRecipients(ctx context.Context, emails string, userID string) ([]model.Recipient, error) {
	wanted := strings.Split(emails, ",")
	var matches []model.Recipient
	for _, rcp := range f.verifiedRecipients[userID] {
		for _, email := range wanted {
			if rcp.Email == email {
				matches = append(matches, rcp)
			}
		}
	}
	return matches, nil
}

func (f *fakeStore) GetRecipients(ctx context.Context, userID string) ([]model.Recipient, error) {
	return f.recipients[userID], nil
}

func (f *fakeStore) PostMessage(ctx context.Context, message model.Message) error {
	f.postedMessages = append(f.postedMessages, message)
	return nil
}

func newTestService(store *fakeStore) *Service {
	return &Service{
		Cfg: config.Config{
			API: config.APIConfig{
				Domains: "mailx.net",
			},
		},
		Store: store,
	}
}

func TestFindRecipients_PlusTagResolvesExistingAlias(t *testing.T) {
	store := newFakeStore()
	store.aliases["myalias@mailx.net"] = model.Alias{
		BaseModel:  model.BaseModel{ID: "alias-1"},
		Name:       "myalias@mailx.net",
		UserID:     "user-1",
		Enabled:    true,
		Recipients: "rcpt@example.com",
	}
	store.recipients["user-1"] = []model.Recipient{{Email: "rcpt@example.com"}}
	s := newTestService(store)

	rcps, alias, msgType, err := s.FindRecipients("sender@somewhere.com", "myalias+shop@mailx.net", model.Send)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if alias.Name != "myalias@mailx.net" {
		t.Errorf("expected alias name myalias@mailx.net, got %s", alias.Name)
	}
	if msgType != model.Forward {
		t.Errorf("expected msgType Forward, got %v", msgType)
	}
	if len(rcps) != 1 || rcps[0].Email != "rcpt@example.com" {
		t.Errorf("expected recipient rcpt@example.com, got %+v", rcps)
	}
}

func TestFindRecipients_WildcardAliasFallbackWhenBaseAliasMissing(t *testing.T) {
	store := newFakeStore()
	store.aliases["*+news@customdomain.com"] = model.Alias{
		BaseModel:  model.BaseModel{ID: "alias-2"},
		Name:       "*+news@customdomain.com",
		UserID:     "user-2",
		Enabled:    true,
		CatchAll:   true,
		Recipients: "rcpt@example.com",
	}
	store.domains["customdomain.com"] = model.Domain{Name: "customdomain.com", UserID: "user-2", Enabled: true}
	store.recipients["user-2"] = []model.Recipient{{Email: "rcpt@example.com"}}
	s := newTestService(store)

	rcps, alias, msgType, err := s.FindRecipients("sender@somewhere.com", "anything+news@customdomain.com", model.Send)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if alias.Name != "*+news@customdomain.com" {
		t.Errorf("expected wildcard alias match, got %s", alias.Name)
	}
	if msgType != model.Forward {
		t.Errorf("expected msgType Forward, got %v", msgType)
	}
	if len(rcps) != 1 || rcps[0].Email != "rcpt@example.com" {
		t.Errorf("expected recipient rcpt@example.com, got %+v", rcps)
	}
}

func TestFindRecipients_ReplyToRoundTripUnaffectedByFix(t *testing.T) {
	store := newFakeStore()
	store.aliases["myalias@mailx.net"] = model.Alias{
		BaseModel: model.BaseModel{ID: "alias-3"},
		Name:      "myalias@mailx.net",
		UserID:    "user-3",
		Enabled:   true,
	}
	store.verifiedRecipients["user-3"] = []model.Recipient{{Email: "sender@somewhere.com", IsActive: true}}
	s := newTestService(store)

	rcps, alias, msgType, err := s.FindRecipients("sender@somewhere.com", "myalias+contact=external.com@mailx.net", model.Reply)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if alias.Name != "myalias@mailx.net" {
		t.Errorf("expected alias myalias@mailx.net, got %s", alias.Name)
	}
	if msgType != model.Reply {
		t.Errorf("expected msgType Reply (passthrough), got %v", msgType)
	}
	if len(rcps) != 1 || rcps[0].Email != "contact@external.com" {
		t.Errorf("expected reply target contact@external.com, got %+v", rcps)
	}
}

func TestFindRecipients_DisabledAliasStillBlockedAfterPlusTagStripped(t *testing.T) {
	store := newFakeStore()
	store.aliases["disabled@mailx.net"] = model.Alias{
		BaseModel: model.BaseModel{ID: "alias-4"},
		Name:      "disabled@mailx.net",
		UserID:    "user-4",
		Enabled:   false,
	}
	s := newTestService(store)

	_, alias, _, err := s.FindRecipients("sender@somewhere.com", "disabled+tag@mailx.net", model.Send)
	if err != ErrDisabledAlias {
		t.Fatalf("expected ErrDisabledAlias, got %v", err)
	}
	if alias.Name != "disabled@mailx.net" {
		t.Errorf("expected resolved alias name disabled@mailx.net, got %s", alias.Name)
	}
	if len(store.postedMessages) != 1 || store.postedMessages[0].Type != model.Block {
		t.Errorf("expected a Block message to be recorded, got %+v", store.postedMessages)
	}
}

func TestFindRecipients_UnmatchedPlusTagFallsThroughToDomainCatchAll(t *testing.T) {
	store := newFakeStore()
	store.domains["customdomain.com"] = model.Domain{
		Name:      "customdomain.com",
		UserID:    "user-5",
		Enabled:   true,
		CatchAll:  true,
		Recipient: "catchall@example.com",
	}
	store.verifiedRecipients["user-5"] = []model.Recipient{{Email: "catchall@example.com", IsActive: true}}
	s := newTestService(store)

	rcps, alias, msgType, err := s.FindRecipients("sender@somewhere.com", "random+tag@customdomain.com", model.Send)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	// No real alias or Wildcard Alias matches, so the base (tag-stripped) name is
	// used purely as a label on the domain-wide catch-all result.
	if alias.Name != "random@customdomain.com" {
		t.Errorf("expected catch-all alias label random@customdomain.com, got %s", alias.Name)
	}
	if alias.Origin == model.Inbound {
		t.Errorf("expected a tagged address not to be marked for auto-creation, got Origin %v", alias.Origin)
	}
	if msgType != model.Forward {
		t.Errorf("expected msgType Forward, got %v", msgType)
	}
	if len(rcps) != 1 || rcps[0].Email != "catchall@example.com" {
		t.Errorf("expected recipient catchall@example.com, got %+v", rcps)
	}
}

// Reproduces the QA report: a plus-tagged address on a catch-all domain must
// never be auto-created as a new alias under its tag-stripped base name.
func TestFindRecipients_TaggedAddressOnCatchAllDomainNotAutoCreated(t *testing.T) {
	store := newFakeStore()
	store.domains["customdomain.com"] = model.Domain{
		Name:      "customdomain.com",
		UserID:    "user-6",
		Enabled:   true,
		CatchAll:  true,
		Recipient: "catchall@example.com",
	}
	store.verifiedRecipients["user-6"] = []model.Recipient{{Email: "catchall@example.com", IsActive: true}}
	s := newTestService(store)

	_, alias, _, err := s.FindRecipients("sender@somewhere.com", "newalias+shop@customdomain.com", model.Send)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if alias.Origin == model.Inbound {
		t.Errorf("expected Origin != Inbound so PostInboundAlias is never invoked, got %v", alias.Origin)
	}
}

func TestFindRecipients_NoAliasNoCatchAllReturnsError(t *testing.T) {
	store := newFakeStore()
	s := newTestService(store)

	_, alias, _, err := s.FindRecipients("sender@somewhere.com", "randomjunk+tag@mailx.net", model.Send)
	if err != ErrGetAliasByName {
		t.Fatalf("expected ErrGetAliasByName, got %v", err)
	}
	if alias.Name != "randomjunk@mailx.net" {
		t.Errorf("expected resolved alias name randomjunk@mailx.net, got %s", alias.Name)
	}
}

func TestResolveForward(t *testing.T) {
	store := newFakeStore()
	store.recipients["user-1"] = []model.Recipient{
		{Email: "a@example.com"},
		{Email: "b@example.com"},
	}
	s := newTestService(store)

	rcps, err := s.resolveForward(model.Alias{UserID: "user-1", Recipients: "a@example.com"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(rcps) != 1 || rcps[0].Email != "a@example.com" {
		t.Errorf("expected only a@example.com, got %+v", rcps)
	}
}

func TestResolveForward_NoRecipientsConfigured(t *testing.T) {
	store := newFakeStore()
	s := newTestService(store)

	_, err := s.resolveForward(model.Alias{UserID: "user-1"})
	if err != ErrNoRecipients {
		t.Fatalf("expected ErrNoRecipients, got %v", err)
	}
}
