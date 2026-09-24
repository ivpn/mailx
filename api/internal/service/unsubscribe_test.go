package service

import (
	"context"
	"errors"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"ivpn.net/email/api/config"
	"ivpn.net/email/api/internal/model"
	"ivpn.net/email/api/internal/utils/signedurl"
)

const testAliasName = "myalias@mailx.net"

func testSigner() *signedurl.Signer {
	return signedurl.New("test-secret", "https://api.mailx.net")
}

func TestBuildUnsubscribeURLs(t *testing.T) {
	urls := buildUnsubscribeURLs(testSigner(), "alias-1", "News@Example.com")

	for name, got := range map[string]string{
		"Deactivate":  urls.Deactivate,
		"Delete":      urls.Delete,
		"BlockEmail":  urls.BlockEmail,
		"BlockDomain": urls.BlockDomain,
	} {
		if got == "" {
			t.Errorf("%s URL is empty", name)
		}
	}

	if !strings.HasPrefix(urls.Deactivate, "https://api.mailx.net/v1/deactivate-one-click/alias-1?") {
		t.Errorf("Deactivate = %q", urls.Deactivate)
	}
	if !strings.HasPrefix(urls.Delete, "https://api.mailx.net/v1/delete-one-click/alias-1?") {
		t.Errorf("Delete = %q", urls.Delete)
	}

	u, _ := url.Parse(urls.BlockEmail)
	if got := u.Query().Get("email"); got != "news@example.com" {
		t.Errorf("BlockEmail email param = %q, want the normalized sender", got)
	}
	u, _ = url.Parse(urls.BlockDomain)
	if got := u.Query().Get("domain"); got != "example.com" {
		t.Errorf("BlockDomain domain param = %q, want example.com", got)
	}

	for name, raw := range map[string]string{"Deactivate": urls.Deactivate, "BlockEmail": urls.BlockEmail} {
		u, _ := url.Parse(raw)
		if !testSigner().Verify(u.Path, u.RawQuery) {
			t.Errorf("%s URL doesn't verify", name)
		}
	}
}

func TestBuildUnsubscribeURLs_Unavailable(t *testing.T) {
	if urls := buildUnsubscribeURLs(nil, "alias-1", "a@example.com"); urls != (unsubscribeURLs{}) {
		t.Errorf("without a signer, got %+v, want no URLs", urls)
	}
	if urls := buildUnsubscribeURLs(testSigner(), "", "a@example.com"); urls != (unsubscribeURLs{}) {
		t.Errorf("without an alias ID, got %+v, want no URLs", urls)
	}

	urls := buildUnsubscribeURLs(testSigner(), "alias-1", "not-an-address")
	if urls.Deactivate == "" || urls.Delete == "" {
		t.Error("deactivate/delete URLs should not depend on the sender")
	}
	if urls.BlockEmail != "" || urls.BlockDomain != "" {
		t.Errorf("block URLs for an unusable sender = %q / %q, want none", urls.BlockEmail, urls.BlockDomain)
	}
}

func TestListUnsubscribeHeaders_OneClickBehaviours(t *testing.T) {
	urls := unsubscribeURLs{
		Deactivate:  "https://api.mailx.net/deactivate",
		Delete:      "https://api.mailx.net/delete",
		BlockEmail:  "https://api.mailx.net/block-email",
		BlockDomain: "https://api.mailx.net/block-domain",
	}
	original := "<https://list.example.com/unsub>, <mailto:unsub@list.example.com>"

	tests := []struct {
		behaviour model.ListUnsubscribeBehaviour
		want      string
	}{
		{model.ListUnsubscribeDeactivate, urls.Deactivate},
		{model.ListUnsubscribeDelete, urls.Delete},
		{model.ListUnsubscribeBlockEmail, urls.BlockEmail},
		{model.ListUnsubscribeBlockDomain, urls.BlockDomain},
	}
	for _, tt := range tests {
		got := listUnsubscribeHeaders(tt.behaviour, urls, original, "List-Unsubscribe=One-Click", testAliasName)
		if got.Value != "<"+tt.want+">" {
			t.Errorf("behaviour %d: Value = %q, want <%s>", tt.behaviour, got.Value, tt.want)
		}
		if got.Post != "List-Unsubscribe=One-Click" {
			t.Errorf("behaviour %d: Post = %q, want List-Unsubscribe=One-Click", tt.behaviour, got.Post)
		}
	}
}

func TestListUnsubscribeHeaders_BlockWithUnusableSenderFallsBack(t *testing.T) {
	urls := unsubscribeURLs{Deactivate: "https://api.mailx.net/deactivate", Delete: "https://api.mailx.net/delete"}

	for _, behaviour := range []model.ListUnsubscribeBehaviour{model.ListUnsubscribeBlockEmail, model.ListUnsubscribeBlockDomain} {
		got := listUnsubscribeHeaders(behaviour, urls, "<https://list.example.com/unsub>", "List-Unsubscribe=One-Click", testAliasName)
		if got.Value != "<https://list.example.com/unsub>" || got.Post != "List-Unsubscribe=One-Click" {
			t.Errorf("behaviour %d with original header = %+v, want the original header", behaviour, got)
		}

		got = listUnsubscribeHeaders(behaviour, urls, "", "", testAliasName)
		if got.Value != "<"+urls.Deactivate+">" {
			t.Errorf("behaviour %d without original header = %+v, want the deactivate URL", behaviour, got)
		}
	}
}

func TestListUnsubscribeHeaders_Original(t *testing.T) {
	urls := unsubscribeURLs{Deactivate: "https://api.mailx.net/deactivate"}

	tests := []struct {
		name         string
		original     string
		originalPost string
		want         model.ListUnsubscribeHeaders
	}{
		{
			name:     "mailto rewritten to the alias reply address",
			original: "<mailto:unsub@list.example.com>",
			want:     model.ListUnsubscribeHeaders{Value: "<mailto:myalias+unsub=list.example.com@mailx.net>"},
		},
		{
			name:     "mailto subject query kept",
			original: "<mailto:unsub@list.example.com?subject=unsubscribe%20me>",
			want:     model.ListUnsubscribeHeaders{Value: "<mailto:myalias+unsub=list.example.com@mailx.net?subject=unsubscribe%20me>"},
		},
		{
			name:         "https kept and one-click post passed through",
			original:     "<https://list.example.com/unsub?id=1&t=2>",
			originalPost: "List-Unsubscribe=One-Click",
			want: model.ListUnsubscribeHeaders{
				Value: "<https://list.example.com/unsub?id=1&t=2>",
				Post:  "List-Unsubscribe=One-Click",
			},
		},
		{
			name:         "mixed, folded header normalised",
			original:     "<mailto:unsub@list.example.com?subject=stop>,\r\n\t<https://list.example.com/unsub>",
			originalPost: "List-Unsubscribe=One-Click",
			want: model.ListUnsubscribeHeaders{
				Value: "<mailto:myalias+unsub=list.example.com@mailx.net?subject=stop>, <https://list.example.com/unsub>",
				Post:  "List-Unsubscribe=One-Click",
			},
		},
		{
			name:         "one-click post dropped without an https URI",
			original:     "<mailto:unsub@list.example.com>",
			originalPost: "List-Unsubscribe=One-Click",
			want:         model.ListUnsubscribeHeaders{Value: "<mailto:myalias+unsub=list.example.com@mailx.net>"},
		},
		{
			name:     "mailto that can't be rewritten is dropped",
			original: "<mailto:unsub+token@list.example.com>, <https://list.example.com/unsub>",
			want:     model.ListUnsubscribeHeaders{Value: "<https://list.example.com/unsub>"},
		},
		{
			name:     "only unusable URIs falls back to deactivate",
			original: "<mailto:unsub+token@list.example.com>, <ftp://list.example.com>",
			want:     model.ListUnsubscribeHeaders{Value: "<https://api.mailx.net/deactivate>", Post: "List-Unsubscribe=One-Click"},
		},
		{
			name: "no original header falls back to deactivate",
			want: model.ListUnsubscribeHeaders{Value: "<https://api.mailx.net/deactivate>", Post: "List-Unsubscribe=One-Click"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := listUnsubscribeHeaders(model.ListUnsubscribeOriginalWithFallback, urls, tt.original, tt.originalPost, testAliasName)
			if got != tt.want {
				t.Errorf("got %+v, want %+v", got, tt.want)
			}
		})
	}
}

func TestListUnsubscribeHeaders_NoSigner(t *testing.T) {
	// One-click URLs unavailable (signing not configured): keep the original, else no header.
	got := listUnsubscribeHeaders(model.ListUnsubscribeDeactivate, unsubscribeURLs{}, "<https://list.example.com/unsub>", "", testAliasName)
	if got.Value != "<https://list.example.com/unsub>" {
		t.Errorf("got %+v, want the original header", got)
	}

	got = listUnsubscribeHeaders(model.ListUnsubscribeDeactivate, unsubscribeURLs{}, "", "", testAliasName)
	if got != (model.ListUnsubscribeHeaders{}) {
		t.Errorf("got %+v, want no header", got)
	}
}

func TestRewriteListUnsubscribe_MailtoRoundTripsToReply(t *testing.T) {
	// The rewritten address must resolve back to the alias and list address in the reply flow.
	for _, aliasName := range []string{"myalias@mailx.net", "*+news@custom.com"} {
		value, _ := rewriteListUnsubscribe("<mailto:unsub@list.example.com>", aliasName)
		address := strings.TrimSuffix(strings.TrimPrefix(value, "<mailto:"), ">")

		gotAlias, gotTo := model.ParseReplyTo(address)
		if gotAlias != aliasName || gotTo != "unsub@list.example.com" {
			t.Errorf("alias %q: ParseReplyTo(%q) = %q, %q", aliasName, address, gotAlias, gotTo)
		}
	}
}

// oneClickStore extends fakeStore with what the one-click actions use.
type oneClickStore struct {
	*fakeStore

	blocklist []model.BlocklistItem
	users     map[string]model.User
}

func (o *oneClickStore) GetAliasByIDUnscoped(ctx context.Context, ID string) (model.Alias, error) {
	for _, a := range o.aliases {
		if a.ID == ID {
			return a, nil
		}
	}
	return model.Alias{}, errNotFound
}

func (o *oneClickStore) DeleteAlias(ctx context.Context, ID string, userID string) error {
	for name, a := range o.aliases {
		if a.ID == ID && a.UserID == userID {
			a.DeletedAt = gorm.DeletedAt{Time: time.Now(), Valid: true}
			o.aliases[name] = a
		}
	}
	return nil
}

func (o *oneClickStore) GetBlocklistCount(ctx context.Context, userID string) (int, error) {
	return len(o.blocklist), nil
}

func (o *oneClickStore) PostBlocklistItem(ctx context.Context, item model.BlocklistItem) (model.BlocklistItem, error) {
	for _, existing := range o.blocklist {
		if existing.UserID == item.UserID && existing.Type == item.Type && existing.Value == item.Value {
			return model.BlocklistItem{}, &mysql.MySQLError{Number: 1062}
		}
	}
	o.blocklist = append(o.blocklist, item)
	return item, nil
}

func (o *oneClickStore) GetUser(ctx context.Context, ID string) (model.User, error) {
	user, ok := o.users[ID]
	if !ok {
		return model.User{}, errNotFound
	}
	return user, nil
}

type sentNotification struct {
	to, subject, template string
	data                  map[string]any
}

func newOneClickTestService(t *testing.T, aliases ...model.Alias) (*Service, *oneClickStore, chan sentNotification) {
	t.Helper()

	store := &oneClickStore{fakeStore: newFakeStore(), users: map[string]model.User{"user-1": {Email: "account@example.com"}}}
	for _, a := range aliases {
		store.aliases[a.Name] = a
	}
	s := &Service{
		Cfg:   config.Config{API: config.APIConfig{ApiAllowOrigin: "https://app.mailx.net"}},
		Store: store,
		Cache: newTotpTestCache(),
	}

	sent := make(chan sentNotification, 10)
	orig := sendNotification
	sendNotification = func(cfg config.SMTPClientConfig, to string, subject string, templateFile string, data any) error {
		sent <- sentNotification{to: to, subject: subject, template: templateFile, data: data.(map[string]any)}
		return nil
	}
	t.Cleanup(func() { sendNotification = orig })

	return s, store, sent
}

func expectNotifications(t *testing.T, sent chan sentNotification, want int) []sentNotification {
	t.Helper()

	var got []sentNotification
	for range want {
		select {
		case n := <-sent:
			got = append(got, n)
		case <-time.After(time.Second):
			t.Fatalf("got %d notifications, want %d", len(got), want)
		}
	}

	select {
	case n := <-sent:
		t.Fatalf("unexpected extra notification: %+v", n)
	case <-time.After(50 * time.Millisecond):
	}

	return got
}

func testOneClickAlias() model.Alias {
	return model.Alias{BaseModel: model.BaseModel{ID: "alias-1"}, Name: testAliasName, UserID: "user-1", Enabled: true, Recipients: "rcpt@example.com"}
}

func TestOneClickDeactivateAlias(t *testing.T) {
	s, store, sent := newOneClickTestService(t, testOneClickAlias())
	store.settingsByUser["user-1"] = model.Settings{Recipient: "default@example.com"}
	req := model.OneClickRequest{IP: "203.0.113.7", UserAgent: "Gmail"}

	for range 3 {
		if err := s.OneClickDeactivateAlias(context.Background(), "alias-1", req); err != nil {
			t.Fatalf("OneClickDeactivateAlias() error = %v", err)
		}
	}

	if store.aliases[testAliasName].Enabled {
		t.Error("alias is still enabled")
	}

	got := expectNotifications(t, sent, 1)
	n := got[0]
	if n.to != "default@example.com" || n.template != "alias_deactivated_unsubscribe.tmpl" {
		t.Errorf("notification = %+v, want deactivate template to the default recipient", n)
	}
	if n.data["alias"] != testAliasName || n.data["ip"] != "203.0.113.7" || n.data["userAgent"] != "Gmail" || n.data["origin"] != "https://app.mailx.net" {
		t.Errorf("notification data = %+v", n.data)
	}
}

func TestOneClickDeactivateAlias_AlreadyInactiveOrDeleted(t *testing.T) {
	inactive := testOneClickAlias()
	inactive.Enabled = false
	deleted := testOneClickAlias()
	deleted.ID, deleted.Name = "alias-2", "deleted@mailx.net"
	deleted.DeletedAt = gorm.DeletedAt{Time: time.Now(), Valid: true}

	s, store, sent := newOneClickTestService(t, inactive, deleted)
	for _, id := range []string{"alias-1", "alias-2"} {
		if err := s.OneClickDeactivateAlias(context.Background(), id, model.OneClickRequest{}); err != nil {
			t.Fatalf("OneClickDeactivateAlias(%s) error = %v", id, err)
		}
	}
	if store.aliases["deleted@mailx.net"].Enabled != true {
		t.Error("a deleted alias shouldn't be touched")
	}

	expectNotifications(t, sent, 0)
}

func TestOneClickDeleteAlias(t *testing.T) {
	s, store, sent := newOneClickTestService(t, testOneClickAlias())

	for range 3 {
		if err := s.OneClickDeleteAlias(context.Background(), "alias-1", model.OneClickRequest{}); err != nil {
			t.Fatalf("OneClickDeleteAlias() error = %v", err)
		}
	}

	if !store.aliases[testAliasName].DeletedAt.Valid {
		t.Error("alias is not deleted")
	}

	got := expectNotifications(t, sent, 1)
	// No default recipient in settings: falls back to the account email.
	if got[0].to != "account@example.com" || got[0].template != "alias_deleted_unsubscribe.tmpl" {
		t.Errorf("notification = %+v, want delete template to the account email", got[0])
	}
}

func TestOneClick_AliasNotFound(t *testing.T) {
	s, _, _ := newOneClickTestService(t)
	ctx := context.Background()

	if err := s.OneClickDeactivateAlias(ctx, "missing", model.OneClickRequest{}); !errors.Is(err, model.ErrOneClickAliasNotFound) {
		t.Errorf("OneClickDeactivateAlias() error = %v, want ErrOneClickAliasNotFound", err)
	}
	if err := s.OneClickDeleteAlias(ctx, "missing", model.OneClickRequest{}); !errors.Is(err, model.ErrOneClickAliasNotFound) {
		t.Errorf("OneClickDeleteAlias() error = %v, want ErrOneClickAliasNotFound", err)
	}
	if err := s.OneClickBlock(ctx, "missing", model.BlocklistEmail, "a@example.com"); !errors.Is(err, model.ErrOneClickAliasNotFound) {
		t.Errorf("OneClickBlock() error = %v, want ErrOneClickAliasNotFound", err)
	}
}

func TestOneClickBlock(t *testing.T) {
	s, store, sent := newOneClickTestService(t, testOneClickAlias())
	ctx := context.Background()

	for range 2 {
		if err := s.OneClickBlock(ctx, "alias-1", model.BlocklistEmail, "News@Example.com"); err != nil {
			t.Fatalf("OneClickBlock(email) error = %v", err)
		}
		if err := s.OneClickBlock(ctx, "alias-1", model.BlocklistDomain, "example.com"); err != nil {
			t.Fatalf("OneClickBlock(domain) error = %v", err)
		}
	}

	want := []model.BlocklistItem{
		{UserID: "user-1", Type: model.BlocklistEmail, Value: "news@example.com"},
		{UserID: "user-1", Type: model.BlocklistDomain, Value: "example.com"},
	}
	if len(store.blocklist) != len(want) {
		t.Fatalf("blocklist = %+v, want %+v", store.blocklist, want)
	}
	for i := range want {
		if store.blocklist[i] != want[i] {
			t.Errorf("blocklist[%d] = %+v, want %+v", i, store.blocklist[i], want[i])
		}
	}

	expectNotifications(t, sent, 0)
}

func TestOneClickBlock_InvalidValue(t *testing.T) {
	s, store, _ := newOneClickTestService(t, testOneClickAlias())
	ctx := context.Background()

	tests := []struct {
		itemType model.BlocklistType
		value    string
	}{
		{model.BlocklistEmail, ""},
		{model.BlocklistEmail, "not-an-email"},
		{model.BlocklistEmail, strings.Repeat("a", 250) + "@example.com"},
		{model.BlocklistDomain, ""},
		{model.BlocklistDomain, "not a domain"},
		{model.BlocklistDomain, "-bad-.com"},
	}
	for _, tt := range tests {
		if err := s.OneClickBlock(ctx, "alias-1", tt.itemType, tt.value); !errors.Is(err, model.ErrOneClickInvalidValue) {
			t.Errorf("OneClickBlock(%s, %q) error = %v, want ErrOneClickInvalidValue", tt.itemType, tt.value, err)
		}
	}
	if len(store.blocklist) != 0 {
		t.Errorf("blocklist = %+v, want nothing added", store.blocklist)
	}
}
