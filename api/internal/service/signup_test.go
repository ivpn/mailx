package service

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	nethttp "net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"ivpn.net/email/api/config"
	"ivpn.net/email/api/internal/client/http"
	"ivpn.net/email/api/internal/model"
)

const (
	signupTestSessionID = "sess-1"
	signupTestPreauthID = "pa-1"
	signupTestToken     = "pa-token"
)

// signupTestStore is a minimal in-memory UserStore fake scoped to signup tests; embeds Store (nil)
// so any unimplemented method call panics loudly rather than silently no-oping.
type signupTestStore struct {
	Store
	users        map[string]model.User
	subs         map[string]model.Subscription
	settings     map[string]model.Settings
	failSettings bool
}

func newSignupTestStore() *signupTestStore {
	return &signupTestStore{
		users:    map[string]model.User{},
		subs:     map[string]model.Subscription{},
		settings: map[string]model.Settings{},
	}
}

func (s *signupTestStore) GetUser(ctx context.Context, ID string) (model.User, error) {
	user, ok := s.users[ID]
	if !ok {
		return model.User{}, errNotFound
	}
	return user, nil
}

func (s *signupTestStore) GetUserByEmailUnfinishedSignup(ctx context.Context, email string) (model.User, error) {
	for _, user := range s.users {
		if user.Email == email && user.PasswordHash == "" {
			return user, nil
		}
	}
	return model.User{}, errNotFound
}

func (s *signupTestStore) CheckDuplicateRecipient(ctx context.Context, email string) (bool, error) {
	return false, nil
}

// PostUser mirrors the unique index on users.email
func (s *signupTestStore) PostUser(ctx context.Context, user model.User) (model.User, error) {
	for _, u := range s.users {
		if u.Email == user.Email {
			return model.User{}, &mysql.MySQLError{Number: 1062, Message: "Duplicate entry for key 'uni_users_email'"}
		}
	}
	user.ID = uuid.New().String()
	s.users[user.ID] = user
	return user, nil
}

func (s *signupTestStore) DeleteUser(ctx context.Context, ID string) error {
	delete(s.users, ID)
	return nil
}

func (s *signupTestStore) GetSubscriptionByTokenHash(ctx context.Context, tokenHash string) (model.Subscription, error) {
	return model.Subscription{}, errNotFound
}

func (s *signupTestStore) PostSubscription(ctx context.Context, sub model.Subscription) error {
	s.subs[sub.UserID] = sub
	return nil
}

func (s *signupTestStore) DeleteSubscription(ctx context.Context, userID string) error {
	delete(s.subs, userID)
	return nil
}

func (s *signupTestStore) PostSettings(ctx context.Context, settings model.Settings) error {
	if s.failSettings {
		return errors.New("settings insert failed")
	}
	s.settings[settings.UserID] = settings
	return nil
}

func (s *signupTestStore) DeleteSettings(ctx context.Context, userID string) error {
	delete(s.settings, userID)
	return nil
}

// newSignupTestService wires a Service to the fake store, a pre-auth session in the
// fake cache, and an httptest server acting as both the preauth service and the
// signup webhook. The webhook responds with whatever status webhookStatus holds.
func newSignupTestService(t *testing.T, store *signupTestStore, webhookStatus *atomic.Int32) *Service {
	t.Helper()

	tokenHash := sha256.Sum256([]byte(signupTestToken))
	preauth := model.Preauth{
		ID:          signupTestPreauthID,
		TokenHash:   base64.StdEncoding.EncodeToString(tokenHash[:]),
		ActiveUntil: time.Now().Add(24 * time.Hour),
		Tier:        "Tier 1",
	}

	srv := httptest.NewServer(nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		switch r.URL.Path {
		case "/preauth/" + signupTestPreauthID:
			_ = json.NewEncoder(w).Encode(preauth)
		case "/webhook":
			w.WriteHeader(int(webhookStatus.Load()))
			_, _ = w.Write([]byte(`{"error":"webhook test response"}`))
		default:
			w.WriteHeader(nethttp.StatusNotFound)
		}
	}))
	t.Cleanup(srv.Close)

	cache := newTotpTestCache()
	paSession, err := json.Marshal(model.PASession{ID: signupTestSessionID, Token: signupTestToken, PreauthId: signupTestPreauthID})
	if err != nil {
		t.Fatalf("failed to marshal pre-auth session: %v", err)
	}
	_ = cache.Set(context.Background(), "pasession_"+signupTestSessionID, string(paSession), time.Hour)

	return &Service{
		Store: store,
		Cache: cache,
		Http: http.Http{Cfg: config.APIConfig{
			PreauthURL:       srv.URL + "/preauth",
			SignupWebhookURL: srv.URL + "/webhook",
		}},
	}
}

func signupWithPassword(s *Service, email string) error {
	pass := "correct horse battery staple"
	_, err := s.GetUnfinishedSignupOrPostUser(context.Background(), model.User{Email: email, PasswordPlain: &pass}, uuid.New().String(), signupTestSessionID)
	return err
}

func assertSignupRolledBack(t *testing.T, store *signupTestStore) {
	t.Helper()
	if len(store.users) != 0 || len(store.subs) != 0 || len(store.settings) != 0 {
		t.Errorf("Expected failed signup to be rolled back, got %d users, %d subscriptions, %d settings",
			len(store.users), len(store.subs), len(store.settings))
	}
}

func TestSignupWebhookFailureRollsBack(t *testing.T) {
	store := newSignupTestStore()
	var webhookStatus atomic.Int32
	webhookStatus.Store(nethttp.StatusInternalServerError)
	s := newSignupTestService(t, store, &webhookStatus)

	err := signupWithPassword(s, "user@example.com")
	if !errors.Is(err, ErrSignupWebhook) {
		t.Fatalf("Expected ErrSignupWebhook, got %v", err)
	}

	assertSignupRolledBack(t, store)
}

// A password signup that failed on the webhook used to leave a user row with a
// password hash behind, so every retry hit the unique email index.
func TestSignupRetryAfterWebhookFailure(t *testing.T) {
	store := newSignupTestStore()
	var webhookStatus atomic.Int32
	webhookStatus.Store(nethttp.StatusBadGateway)
	s := newSignupTestService(t, store, &webhookStatus)

	err := signupWithPassword(s, "user@example.com")
	if !errors.Is(err, ErrSignupWebhook) {
		t.Fatalf("Expected ErrSignupWebhook on first attempt, got %v", err)
	}

	webhookStatus.Store(nethttp.StatusOK)

	err = signupWithPassword(s, "user@example.com")
	if err != nil {
		t.Fatalf("Expected retry to succeed, got %v", err)
	}

	if len(store.users) != 1 || len(store.subs) != 1 || len(store.settings) != 1 {
		t.Errorf("Expected exactly one user, subscription and settings, got %d, %d, %d",
			len(store.users), len(store.subs), len(store.settings))
	}
}

func TestSignupSettingsFailureRollsBack(t *testing.T) {
	store := newSignupTestStore()
	store.failSettings = true
	var webhookStatus atomic.Int32
	webhookStatus.Store(nethttp.StatusOK)
	s := newSignupTestService(t, store, &webhookStatus)

	err := signupWithPassword(s, "user@example.com")
	if !errors.Is(err, ErrPostUser) {
		t.Fatalf("Expected ErrPostUser, got %v", err)
	}

	assertSignupRolledBack(t, store)
}

// The API response must not reveal that an email is already registered, and the
// existing account must not be touched by the failed signup.
func TestSignupDuplicateEmailIsGeneric(t *testing.T) {
	store := newSignupTestStore()
	existing := model.User{Email: "user@example.com", PasswordHash: "existing-hash", IsActive: true}
	existing.ID = "existing-user"
	store.users[existing.ID] = existing
	store.subs[existing.ID] = model.Subscription{UserID: existing.ID}
	store.settings[existing.ID] = model.Settings{UserID: existing.ID}

	var webhookStatus atomic.Int32
	webhookStatus.Store(nethttp.StatusOK)
	s := newSignupTestService(t, store, &webhookStatus)

	err := signupWithPassword(s, "user@example.com")
	if !errors.Is(err, ErrPostUser) {
		t.Fatalf("Expected ErrPostUser, got %v", err)
	}
	if strings.Contains(strings.ToLower(err.Error()), "email") {
		t.Errorf("Expected a generic error, got %q", err.Error())
	}

	if _, ok := store.users[existing.ID]; !ok {
		t.Error("Expected existing user to be kept")
	}
	if _, ok := store.subs[existing.ID]; !ok {
		t.Error("Expected existing subscription to be kept")
	}
	if _, ok := store.settings[existing.ID]; !ok {
		t.Error("Expected existing settings to be kept")
	}
}
