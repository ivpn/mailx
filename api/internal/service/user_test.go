package service

import (
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	"ivpn.net/email/api/internal/model"
	"ivpn.net/email/api/internal/utils"
)

// totpTestStore is a minimal UserStore fake scoped to TOTP tests; embeds Store (nil)
// so any unimplemented method call panics loudly rather than silently no-oping.
type totpTestStore struct {
	Store
}

func (s *totpTestStore) GetUser(ctx context.Context, ID string) (model.User, error) {
	return model.User{Email: "user@example.com"}, nil
}

func (s *totpTestStore) TotpEnable(ctx context.Context, ID string, secret string, backupCodes string) error {
	return nil
}

// totpTestCache is a minimal in-memory Cache fake scoped to TOTP tests.
type totpTestCache struct {
	mu   sync.Mutex
	data map[string]string
}

func newTotpTestCache() *totpTestCache {
	return &totpTestCache{data: map[string]string{}}
}

func (c *totpTestCache) Set(ctx context.Context, key string, value any, exp time.Duration) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = fmt.Sprintf("%v", value)
	return nil
}

func (c *totpTestCache) Get(ctx context.Context, key string) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	v, ok := c.data[key]
	if !ok {
		return "", errNotFound
	}
	return v, nil
}

func (c *totpTestCache) Del(ctx context.Context, key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
	return nil
}

func (c *totpTestCache) Incr(ctx context.Context, key string, exp time.Duration) error {
	return nil
}

func TestTotpEnable_SecretEntropy(t *testing.T) {
	s := &Service{Store: &totpTestStore{}, Cache: newTotpTestCache()}

	res, err := s.TotpEnable(context.Background(), "user-1")
	if err != nil {
		t.Fatalf("TotpEnable returned an error: %v", err)
	}

	decoded, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(res.Secret)
	if err != nil {
		t.Fatalf("failed to decode secret: %v", err)
	}
	if len(decoded) != 20 {
		t.Errorf("Expected a 160-bit (20 byte) secret, got %d bytes", len(decoded))
	}
}

func TestTotpEnableConfirm_BackupCodes(t *testing.T) {
	s := &Service{Store: &totpTestStore{}, Cache: newTotpTestCache()}

	enableRes, err := s.TotpEnable(context.Background(), "user-1")
	if err != nil {
		t.Fatalf("TotpEnable returned an error: %v", err)
	}

	code := computeTotpCodeForTest(t, enableRes.Secret)

	res, err := s.TotpEnableConfirm(context.Background(), "user-1", code)
	if err != nil {
		t.Fatalf("TotpEnableConfirm returned an error: %v", err)
	}

	codes := strings.Fields(res.Backup)
	if len(codes) != 8 {
		t.Fatalf("Expected 8 backup codes, got %d", len(codes))
	}
	for _, c := range codes {
		if len(c) != 10 {
			t.Errorf("Expected backup code length 10, got %d (%s)", len(c), c)
		}
		for _, ch := range c {
			if !strings.ContainsRune(utils.AlphaNumericUserFriendly, ch) {
				t.Errorf("Unexpected character %c in backup code %s", ch, c)
			}
		}
	}
}

// computeTotpCodeForTest duplicates internal/utils/otp_2fa.go's unexported computeCode,
// since it isn't reachable across the package boundary from these service-layer tests.
func computeTotpCodeForTest(t *testing.T, secret string) string {
	t.Helper()

	key, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(secret)
	if err != nil {
		t.Fatalf("failed to decode secret: %v", err)
	}

	value := time.Now().UTC().Unix() / 30
	h := hmac.New(sha1.New, key)
	if err := binary.Write(h, binary.BigEndian, value); err != nil {
		t.Fatalf("failed to compute HMAC: %v", err)
	}
	sum := h.Sum(nil)

	offset := sum[19] & 0x0f
	truncated := binary.BigEndian.Uint32(sum[offset:offset+4]) & 0x7fffffff
	code := truncated % 1000000

	return fmt.Sprintf("%06d", code)
}
