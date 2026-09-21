package service

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"gorm.io/gorm"
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
	subscription       model.Subscription
	postAliasErr       error
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

func (f *fakeStore) GetDomains(ctx context.Context, userID string) ([]model.Domain, error) {
	var result []model.Domain
	for _, d := range f.domains {
		if d.UserID == userID {
			result = append(result, d)
		}
	}
	return result, nil
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

func (f *fakeStore) GetSubscription(ctx context.Context, userID string) (model.Subscription, error) {
	return f.subscription, nil
}

func (f *fakeStore) PostAlias(ctx context.Context, alias model.Alias, maxDaily int, maxInboundHourly int) (model.Alias, error) {
	if f.postAliasErr != nil {
		return model.Alias{}, f.postAliasErr
	}
	return alias, nil
}

// GetAliases is a minimal filter (userID + wildcard flag only) sufficient for the
// wildcard-related PostAlias/GetWildcardDomainInfo tests that use it.
func (f *fakeStore) GetAliases(ctx context.Context, userID string, limit int, offset int, sortBy string, sortOrder string, wildcard string, search string, status string) ([]model.Alias, error) {
	var result []model.Alias
	for _, a := range f.aliases {
		if a.UserID != userID {
			continue
		}
		if wildcard == "true" && !a.Wildcard {
			continue
		}
		if wildcard == "false" && a.Wildcard {
			continue
		}
		result = append(result, a)
	}
	return result, nil
}

func (f *fakeStore) GetAliasesNoStats(ctx context.Context, userID string, limit int, offset int, sortBy string, sortOrder string, wildcard string, search string, status string) ([]model.Alias, error) {
	return f.GetAliases(ctx, userID, limit, offset, sortBy, sortOrder, wildcard, search, status)
}

func (f *fakeStore) GetAliasUnscoped(ctx context.Context, ID string, userID string) (model.Alias, error) {
	for _, a := range f.aliases {
		if a.ID == ID && a.UserID == userID {
			return a, nil
		}
	}
	return model.Alias{}, errNotFound
}

func (f *fakeStore) ForgetAlias(ctx context.Context, ID string, userID string) error {
	for name, a := range f.aliases {
		if a.ID == ID && a.UserID == userID {
			delete(f.aliases, name)
			return nil
		}
	}
	return errNotFound
}

func (f *fakeStore) BulkUpdateAliasEnabled(ctx context.Context, ids []string, userID string, enabled bool) error {
	idSet := aliasIDSet(ids)
	for name, a := range f.aliases {
		if idSet[a.ID] && a.UserID == userID {
			a.Enabled = enabled
			f.aliases[name] = a
		}
	}
	return nil
}

func (f *fakeStore) BulkUpdateAliasPinned(ctx context.Context, ids []string, userID string, pinned bool) error {
	idSet := aliasIDSet(ids)
	for name, a := range f.aliases {
		if idSet[a.ID] && a.UserID == userID {
			a.Pinned = pinned
			f.aliases[name] = a
		}
	}
	return nil
}

func (f *fakeStore) BulkDeleteAlias(ctx context.Context, ids []string, userID string) error {
	idSet := aliasIDSet(ids)
	matched := 0
	for _, a := range f.aliases {
		if idSet[a.ID] && a.UserID == userID && !a.DeletedAt.Valid {
			matched++
		}
	}
	if matched != len(idSet) {
		return model.ErrBulkAliasNotEligible
	}

	for name, a := range f.aliases {
		if idSet[a.ID] && a.UserID == userID {
			a.DeletedAt = gorm.DeletedAt{Time: time.Now(), Valid: true}
			f.aliases[name] = a
		}
	}
	return nil
}

func (f *fakeStore) BulkRestoreAlias(ctx context.Context, ids []string, userID string) error {
	idSet := aliasIDSet(ids)
	matched := 0
	for _, a := range f.aliases {
		if idSet[a.ID] && a.UserID == userID && a.DeletedAt.Valid {
			matched++
		}
	}
	if matched != len(idSet) {
		return model.ErrBulkAliasNotEligible
	}

	for name, a := range f.aliases {
		if idSet[a.ID] && a.UserID == userID {
			a.DeletedAt = gorm.DeletedAt{}
			f.aliases[name] = a
		}
	}
	return nil
}

func (f *fakeStore) GetAliasesUnscopedByIDs(ctx context.Context, ids []string, userID string) ([]model.Alias, error) {
	idSet := aliasIDSet(ids)
	var result []model.Alias
	for _, a := range f.aliases {
		if idSet[a.ID] && a.UserID == userID {
			result = append(result, a)
		}
	}
	return result, nil
}

func (f *fakeStore) BulkForgetAlias(ctx context.Context, ids []string, userID string) error {
	idSet := aliasIDSet(ids)
	for name, a := range f.aliases {
		if idSet[a.ID] && a.UserID == userID {
			delete(f.aliases, name)
		}
	}
	return nil
}

func aliasIDSet(ids []string) map[string]bool {
	set := make(map[string]bool, len(ids))
	for _, id := range ids {
		set[id] = true
	}
	return set
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
		Wildcard:   true,
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

func TestFindRecipients_DotWildcardAliasFallbackWhenBaseAliasMissing(t *testing.T) {
	store := newFakeStore()
	store.aliases["*.news@customdomain.com"] = model.Alias{
		BaseModel:  model.BaseModel{ID: "alias-2b"},
		Name:       "*.news@customdomain.com",
		UserID:     "user-2b",
		Enabled:    true,
		Wildcard:   true,
		Recipients: "rcpt@example.com",
	}
	store.domains["customdomain.com"] = model.Domain{Name: "customdomain.com", UserID: "user-2b", Enabled: true}
	store.recipients["user-2b"] = []model.Recipient{{Email: "rcpt@example.com"}}
	s := newTestService(store)

	rcps, alias, msgType, err := s.FindRecipients("sender@somewhere.com", "anything.news@customdomain.com", model.Send)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if alias.Name != "*.news@customdomain.com" {
		t.Errorf("expected dot wildcard alias match, got %s", alias.Name)
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

// A "+" tag that doesn't match any existing Wildcard Alias is just an ordinary address:
// it must still be eligible for catch-all auto-creation, using the full address (including
// the tag) as the new alias name, not a tag-stripped base.
func TestFindRecipients_UnmatchedPlusTagIsAutoCreatedWithFullAddress(t *testing.T) {
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
	if alias.Name != "random+tag@customdomain.com" {
		t.Errorf("expected the full address random+tag@customdomain.com to become the new alias, got %s", alias.Name)
	}
	if alias.Origin != model.Inbound {
		t.Errorf("expected Origin == Inbound so PostInboundAlias auto-creates the alias, got %v", alias.Origin)
	}
	if msgType != model.Forward {
		t.Errorf("expected msgType Forward, got %v", msgType)
	}
	if len(rcps) != 1 || rcps[0].Email != "catchall@example.com" {
		t.Errorf("expected recipient catchall@example.com, got %+v", rcps)
	}
}

// When a "+" Wildcard Alias actually exists for the suffix, the address must resolve to
// that alias (via the fallback earlier in FindRecipients) rather than the domain catch-all.
func TestFindRecipients_PlusTagWithMatchingWildcardRidesWildcardNotCatchAll(t *testing.T) {
	store := newFakeStore()
	store.aliases["*+shop@customdomain.com"] = model.Alias{
		BaseModel:  model.BaseModel{ID: "alias-6"},
		Name:       "*+shop@customdomain.com",
		UserID:     "user-6",
		Enabled:    true,
		Wildcard:   true,
		Recipients: "rcpt@example.com",
	}
	store.domains["customdomain.com"] = model.Domain{
		Name:      "customdomain.com",
		UserID:    "user-6",
		Enabled:   true,
		CatchAll:  true,
		Recipient: "catchall@example.com",
	}
	store.recipients["user-6"] = []model.Recipient{{Email: "rcpt@example.com"}}
	s := newTestService(store)

	_, alias, _, err := s.FindRecipients("sender@somewhere.com", "newalias+shop@customdomain.com", model.Send)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if alias.Name != "*+shop@customdomain.com" {
		t.Errorf("expected the plus Wildcard Alias to be used, got %s", alias.Name)
	}
}

// A reply-encoded address that doesn't resolve to a real alias must never auto-provision
// one - unlike a plain tag, there's no sensible "full address" to create an alias from.
func TestFindRecipients_ReplyEncodedAddressWithoutMatchingAliasNotAutoCreated(t *testing.T) {
	store := newFakeStore()
	store.domains["customdomain.com"] = model.Domain{
		Name:      "customdomain.com",
		UserID:    "user-6e",
		Enabled:   true,
		CatchAll:  true,
		Recipient: "catchall@example.com",
	}
	store.verifiedRecipients["user-6e"] = []model.Recipient{
		{Email: "catchall@example.com", IsActive: true},
		{Email: "sender@somewhere.com", IsActive: true},
	}
	s := newTestService(store)

	_, alias, _, err := s.FindRecipients("sender@somewhere.com", "noalias+contact=external.com@customdomain.com", model.Send)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if alias.Origin == model.Inbound {
		t.Errorf("expected Origin != Inbound so PostInboundAlias is never invoked, got %v", alias.Origin)
	}
}

// A "." only blocks catch-all auto-creation when it actually resolves to an existing
// Wildcard Alias (handled by the fallback above, before this point is ever reached);
// otherwise a dotted address is a normal address and may still be auto-created.
func TestFindRecipients_DottedAddressWithoutMatchingWildcardIsAutoCreated(t *testing.T) {
	store := newFakeStore()
	store.domains["customdomain.com"] = model.Domain{
		Name:      "customdomain.com",
		UserID:    "user-6b",
		Enabled:   true,
		CatchAll:  true,
		Recipient: "catchall@example.com",
	}
	store.verifiedRecipients["user-6b"] = []model.Recipient{{Email: "catchall@example.com", IsActive: true}}
	s := newTestService(store)

	_, alias, _, err := s.FindRecipients("sender@somewhere.com", "newalias.shop@customdomain.com", model.Send)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if alias.Origin != model.Inbound {
		t.Errorf("expected Origin == Inbound so PostInboundAlias auto-creates the alias, got %v", alias.Origin)
	}
}

// When a "." Wildcard Alias actually exists for the suffix, the address must resolve to
// that alias (via the fallback earlier in FindRecipients) rather than the domain catch-all.
func TestFindRecipients_DottedAddressWithMatchingWildcardRidesWildcardNotCatchAll(t *testing.T) {
	store := newFakeStore()
	store.aliases["*.shop@customdomain.com"] = model.Alias{
		BaseModel:  model.BaseModel{ID: "alias-7"},
		Name:       "*.shop@customdomain.com",
		UserID:     "user-6d",
		Enabled:    true,
		Wildcard:   true,
		Recipients: "rcpt@example.com",
	}
	store.domains["customdomain.com"] = model.Domain{
		Name:      "customdomain.com",
		UserID:    "user-6d",
		Enabled:   true,
		CatchAll:  true,
		Recipient: "catchall@example.com",
	}
	store.recipients["user-6d"] = []model.Recipient{{Email: "rcpt@example.com"}}
	s := newTestService(store)

	_, alias, _, err := s.FindRecipients("sender@somewhere.com", "newalias.shop@customdomain.com", model.Send)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if alias.Name != "*.shop@customdomain.com" {
		t.Errorf("expected the dot Wildcard Alias to be used, got %s", alias.Name)
	}
}

// Regression: hasTag must only look at the local part. The domain itself always contains a
// ".", so a plain address (no "+" or "." in the local part) must still be eligible for
// catch-all auto-creation, even though the full address contains dots from the domain.
func TestFindRecipients_PlainAddressOnCatchAllDomainIsAutoCreated(t *testing.T) {
	store := newFakeStore()
	store.domains["domain.net"] = model.Domain{
		Name:      "domain.net",
		UserID:    "user-6c",
		Enabled:   true,
		CatchAll:  true,
		Recipient: "catchall@example.com",
	}
	store.verifiedRecipients["user-6c"] = []model.Recipient{{Email: "catchall@example.com", IsActive: true}}
	s := newTestService(store)

	_, alias, _, err := s.FindRecipients("sender@somewhere.com", "catchalldomainnew@domain.net", model.Send)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if alias.Origin != model.Inbound {
		t.Errorf("expected Origin == Inbound so PostInboundAlias auto-creates the alias, got %v", alias.Origin)
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
