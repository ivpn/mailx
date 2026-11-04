package service

import (
	"context"
	"errors"
	"log"
	"net/url"
	"regexp"
	"strings"
	"time"

	"ivpn.net/email/api/config"
	"ivpn.net/email/api/internal/client/mailer"
	"ivpn.net/email/api/internal/model"
	"ivpn.net/email/api/internal/utils"
	"ivpn.net/email/api/internal/utils/signedurl"
)

const (
	// oneClickTTL is how long the one-click URLs in a forwarded message stay valid.
	oneClickTTL = 30 * 24 * time.Hour
	// oneClickNotifyTTL dedupes the notification sent when a one-click request deactivates or
	// deletes an alias, so repeated requests don't flood the user's inbox.
	oneClickNotifyTTL = time.Hour

	listUnsubscribeOneClick = "List-Unsubscribe=One-Click"
)

var (
	ErrUpdateListUnsubscribeBehaviour = errors.New("Unable to update List-Unsubscribe behaviour.")

	listUnsubscribeFoldRE = regexp.MustCompile(`\s*[\r\n]+\s*`)
	listUnsubscribeURIRE  = regexp.MustCompile(`<([^<>]*)>`)

	blocklistValidator = utils.NewValidator()

	// sendNotification sends a one-click unsubscribe notification; tests replace it.
	sendNotification = func(cfg config.SMTPClientConfig, to string, subject string, templateFile string, data any) error {
		return mailer.New(cfg).SendTemplate(to, subject, templateFile, data)
	}
)

func (s *Service) UpdateListUnsubscribeBehaviour(ctx context.Context, userID string, behaviour model.ListUnsubscribeBehaviour) error {
	err := s.Store.UpdateListUnsubscribeBehaviour(ctx, userID, behaviour)
	if err != nil {
		log.Printf("error updating list-unsubscribe behaviour: %s", err.Error())
		return ErrUpdateListUnsubscribeBehaviour
	}

	return nil
}

// OneClickSigner returns the signer for one-click URLs, or nil when it isn't configured.
func (s *Service) OneClickSigner() *signedurl.Signer {
	return signedurl.New(s.Cfg.API.UnsubscribeSecret, s.Cfg.API.PublicURL)
}

// ListUnsubscribeHeaders returns the List-Unsubscribe headers for a message forwarded
// through alias, per the user's List-Unsubscribe behaviour setting.
func (s *Service) ListUnsubscribeHeaders(settings model.Settings, alias model.Alias, msg model.Msg) model.ListUnsubscribeHeaders {
	urls := buildUnsubscribeURLs(s.OneClickSigner(), alias.ID, msg.From)
	return listUnsubscribeHeaders(settings.ListUnsubscribeBehaviour, urls, msg.ListUnsubscribe, msg.ListUnsubscribePost, alias.Name)
}

// unsubscribeURLs are the signed one-click URLs for a forwarded message; empty means unavailable.
type unsubscribeURLs struct {
	Deactivate  string
	Delete      string
	BlockEmail  string
	BlockDomain string
}

func buildUnsubscribeURLs(signer *signedurl.Signer, aliasID string, sender string) unsubscribeURLs {
	// Without an alias ID (a catch-all alias that couldn't be saved) there's nothing to act on.
	if signer == nil || aliasID == "" {
		return unsubscribeURLs{}
	}

	urls := unsubscribeURLs{
		Deactivate: signer.Sign(model.OneClickDeactivatePath+aliasID, nil, oneClickTTL),
		Delete:     signer.Sign(model.OneClickDeletePath+aliasID, nil, oneClickTTL),
	}

	email := model.NormalizeBlocklistValue(model.BlocklistEmail, sender)
	if blocklistValidator.Var(email, model.BlocklistValueRules(model.BlocklistEmail)) == nil {
		urls.BlockEmail = signer.Sign(model.OneClickBlockEmailPath+aliasID, url.Values{"email": {email}}, oneClickTTL)
	}

	if at := strings.LastIndex(email, "@"); at >= 0 {
		domain := model.NormalizeBlocklistValue(model.BlocklistDomain, email[at+1:])
		if blocklistValidator.Var(domain, model.BlocklistValueRules(model.BlocklistDomain)) == nil {
			urls.BlockDomain = signer.Sign(model.OneClickBlockDomainPath+aliasID, url.Values{"domain": {domain}}, oneClickTTL)
		}
	}

	return urls
}

func listUnsubscribeHeaders(behaviour model.ListUnsubscribeBehaviour, urls unsubscribeURLs, originalValue string, originalPost string, aliasName string) model.ListUnsubscribeHeaders {
	oneClick := func(target string) model.ListUnsubscribeHeaders {
		return model.ListUnsubscribeHeaders{Value: "<" + target + ">", Post: listUnsubscribeOneClick}
	}

	switch {
	case behaviour == model.ListUnsubscribeDeactivate && urls.Deactivate != "":
		return oneClick(urls.Deactivate)
	case behaviour == model.ListUnsubscribeDelete && urls.Delete != "":
		return oneClick(urls.Delete)
	case behaviour == model.ListUnsubscribeBlockEmail && urls.BlockEmail != "":
		return oneClick(urls.BlockEmail)
	case behaviour == model.ListUnsubscribeBlockDomain && urls.BlockDomain != "":
		return oneClick(urls.BlockDomain)
	}

	// OriginalWithFallback, or the selected one-click URL is unavailable
	if value, hasHTTP := rewriteListUnsubscribe(originalValue, aliasName); value != "" {
		headers := model.ListUnsubscribeHeaders{Value: value}
		// RFC 8058 one-click only applies to an HTTP(S) URI
		if hasHTTP && strings.EqualFold(strings.TrimSpace(originalPost), listUnsubscribeOneClick) {
			headers.Post = listUnsubscribeOneClick
		}
		return headers
	}

	if urls.Deactivate != "" {
		return oneClick(urls.Deactivate)
	}

	return model.ListUnsubscribeHeaders{}
}

// rewriteListUnsubscribe normalises an original List-Unsubscribe value and rewrites its mailto:
// URIs to the alias's reply address, so an unsubscribe email goes out from the alias instead of
// the user's real address. HTTP(S) URIs are kept as they are; any other URI, and any mailto:
// that can't be rewritten to a reply address that resolves back to the same list address, is
// dropped rather than risk exposing the real address. hasHTTP reports whether an HTTP(S) URI
// was kept.
func rewriteListUnsubscribe(value string, aliasName string) (rewritten string, hasHTTP bool) {
	value = listUnsubscribeFoldRE.ReplaceAllString(strings.TrimSpace(value), " ")

	var uris []string
	for _, m := range listUnsubscribeURIRE.FindAllStringSubmatch(value, -1) {
		uri := strings.TrimSpace(m[1])
		lower := strings.ToLower(uri)

		switch {
		case strings.HasPrefix(lower, "https://") || strings.HasPrefix(lower, "http://"):
			uris = append(uris, "<"+uri+">")
			hasHTTP = true
		case strings.HasPrefix(lower, "mailto:"):
			if mailto, ok := rewriteMailto(uri[len("mailto:"):], aliasName); ok {
				uris = append(uris, "<"+mailto+">")
			}
		}
	}

	return strings.Join(uris, ", "), hasHTTP
}

func rewriteMailto(target string, aliasName string) (string, bool) {
	address, query, _ := strings.Cut(target, "?")
	address, err := url.PathUnescape(strings.TrimSpace(address))
	if err != nil || strings.Contains(address, ",") || utils.ValidateEmail(address) != nil {
		return "", false
	}

	replyTo := model.GenerateReplyTo(aliasName, address)
	if parsedAlias, parsedTo := model.ParseReplyTo(replyTo); parsedAlias != aliasName || parsedTo != address {
		return "", false
	}

	mailto := "mailto:" + replyTo
	if query = mailtoKeptQuery(query); query != "" {
		mailto += "?" + query
	}

	return mailto, true
}

// mailtoKeptQuery keeps only the subject and body fields of a mailto: query. Any other field
// (to, cc, bcc, …) is dropped, as it could make the mail client send a copy straight from the
// user's real address. Kept fields are left as the sender encoded them, since re-encoding would
// turn spaces into "+", which mailto: (RFC 6068) reads as a literal plus.
func mailtoKeptQuery(query string) string {
	var kept []string
	for field := range strings.SplitSeq(query, "&") {
		name, _, _ := strings.Cut(field, "=")
		name, err := url.PathUnescape(name)
		if err != nil {
			continue
		}

		if strings.EqualFold(name, "subject") || strings.EqualFold(name, "body") {
			kept = append(kept, field)
		}
	}

	return strings.Join(kept, "&")
}

// OneClickDeactivateAlias deactivates an alias on a one-click unsubscribe request. It is
// idempotent: an alias that's already inactive or deleted is left as is, with no notification.
func (s *Service) OneClickDeactivateAlias(ctx context.Context, aliasID string, req model.OneClickRequest) error {
	alias, err := s.oneClickAlias(ctx, aliasID)
	if err != nil {
		return err
	}

	if alias.DeletedAt.Valid || !alias.Enabled {
		return nil
	}

	err = s.Store.BulkUpdateAliasEnabled(ctx, []string{alias.ID}, alias.UserID, false)
	if err != nil {
		log.Printf("error deactivating alias via one-click unsubscribe: %s", err.Error())
		return model.ErrOneClickAction
	}

	log.Printf("One-Click Unsubscribe deactivated alias %s", alias.ID)
	s.notifyOneClick(ctx, alias, req, "unsubscribe-deactivate-notify:", "Alias deactivated via one-click unsubscribe", "alias_deactivated_unsubscribe.tmpl")

	return nil
}

// OneClickDeleteAlias soft-deletes an alias on a one-click unsubscribe request. It is
// idempotent: an alias that's already deleted is left as is, with no notification.
func (s *Service) OneClickDeleteAlias(ctx context.Context, aliasID string, req model.OneClickRequest) error {
	alias, err := s.oneClickAlias(ctx, aliasID)
	if err != nil {
		return err
	}

	if alias.DeletedAt.Valid {
		return nil
	}

	err = s.Store.DeleteAlias(ctx, alias.ID, alias.UserID)
	if err != nil {
		log.Printf("error deleting alias via one-click unsubscribe: %s", err.Error())
		return model.ErrOneClickAction
	}

	log.Printf("One-Click Unsubscribe deleted alias %s", alias.ID)
	s.notifyOneClick(ctx, alias, req, "unsubscribe-delete-notify:", "Alias deleted via one-click unsubscribe", "alias_deleted_unsubscribe.tmpl")

	return nil
}

// OneClickBlock adds a sender address or domain to the alias owner's blocklist on a one-click
// unsubscribe request. It is idempotent: an entry that's already blocklisted is a success.
func (s *Service) OneClickBlock(ctx context.Context, aliasID string, itemType model.BlocklistType, value string) error {
	value = model.NormalizeBlocklistValue(itemType, value)
	if blocklistValidator.Var(value, model.BlocklistValueRules(itemType)) != nil {
		return model.ErrOneClickInvalidValue
	}

	alias, err := s.oneClickAlias(ctx, aliasID)
	if err != nil {
		return err
	}

	_, err = s.PostBlocklistItem(ctx, model.BlocklistItem{
		UserID: alias.UserID,
		Type:   itemType,
		Value:  value,
	})
	if errors.Is(err, ErrDuplicateBlocklistItem) {
		return nil
	}
	if err != nil {
		log.Printf("error blocking %s via one-click unsubscribe: %s", itemType, err.Error())
		return model.ErrOneClickAction
	}

	log.Printf("One-Click Unsubscribe blocked %s for alias %s", itemType, alias.ID)

	return nil
}

func (s *Service) oneClickAlias(ctx context.Context, aliasID string) (model.Alias, error) {
	alias, err := s.Store.GetAliasByIDUnscoped(ctx, aliasID)
	if err != nil || alias.ID == "" {
		return model.Alias{}, model.ErrOneClickAliasNotFound
	}

	return alias, nil
}

// notifyOneClick tells the user an alias was deactivated/deleted by a one-click request, at
// most once per alias and action within oneClickNotifyTTL, so it can be undone if unintended.
func (s *Service) notifyOneClick(ctx context.Context, alias model.Alias, req model.OneClickRequest, keyPrefix string, subject string, templateFile string) {
	first, err := s.Cache.SetNX(ctx, keyPrefix+alias.ID, 1, oneClickNotifyTTL)
	if err != nil {
		log.Printf("error deduping one-click unsubscribe notification: %s", err.Error())
		return
	}
	if !first {
		return
	}

	to, err := s.notificationRecipient(ctx, alias.UserID)
	if err != nil {
		log.Printf("error resolving one-click unsubscribe notification recipient: %s", err.Error())
		return
	}

	origin, _, _ := strings.Cut(s.Cfg.API.ApiAllowOrigin, ",")
	data := map[string]any{
		"alias":     alias.Name,
		"ip":        req.IP,
		"userAgent": req.UserAgent,
		"timestamp": time.Now().UTC().Format("January 2, 3:04 PM (MST)"),
		"origin":    strings.TrimSpace(origin),
		"from":      s.Cfg.SMTPClient.SenderName,
	}

	utils.Background(func() {
		err := sendNotification(s.Cfg.SMTPClient, to, subject, templateFile, data)
		if err != nil {
			log.Printf("error sending one-click unsubscribe notification: %s", err.Error())
		}
	})
}

// notificationRecipient returns the user's default recipient, falling back to the account email.
func (s *Service) notificationRecipient(ctx context.Context, userID string) (string, error) {
	settings, err := s.Store.GetSettings(ctx, userID)
	if err == nil && settings.Recipient != "" {
		return settings.Recipient, nil
	}

	user, err := s.Store.GetUser(ctx, userID)
	if err != nil {
		return "", err
	}

	return user.Email, nil
}
