package utils

import (
	"bytes"
	"errors"
	"fmt"
	"net/mail"
	"strconv"
	"strings"
)

var (
	ErrAuthServIDNotConfigured = errors.New("authserv-id is not configured")
	ErrNoOwnAuthResults        = errors.New("no Authentication-Results header from own mail server")
	ErrDuplicateOwnAuthResults = errors.New("multiple Authentication-Results headers with own authserv-id")
	ErrMalformedOwnAuthResults = errors.New("malformed Authentication-Results header from own mail server")
	ErrDMARCPolicyReject       = errors.New("dmarc=fail with policy=reject")
	ErrNoAlignedAuthentication = errors.New("no aligned DMARC, DKIM or SPF pass")
	ErrDMARCFail               = errors.New("dmarc=fail")
)

type EmailAuthConfig struct {
	// AuthServIDs are the authserv-ids our own mail servers (Rspamd) write into
	// Authentication-Results, one per server. Headers with any other
	// authserv-id are ignored, since the sender or earlier hops can write them.
	AuthServIDs []string
	// ARCTrustedSealers are domains (e.g. google.com) whose ARC seal is trusted
	// to vouch for a message that passed authentication before reaching us.
	ARCTrustedSealers []string
}

// resInfo is one "method=result property=value ..." entry of an
// Authentication-Results header.
type resInfo struct {
	method   string
	result   string
	props    map[string]string
	comments string
}

// authResults is a parsed Authentication-Results (RFC 8601) or
// ARC-Authentication-Results (RFC 8617) header.
type authResults struct {
	servID   string
	instance int
	results  []resInfo
	// malformed says why the header could not be parsed cleanly, e.g. an
	// unterminated quote or comment. Empty for a well-formed header.
	malformed string
}

// knownMethods are the registered authentication methods (IANA "Email
// Authentication Methods"). Rspamd only writes auth, dkim, dmarc, spf and arc.
var knownMethods = map[string]bool{
	"arc": true, "auth": true, "bimi": true, "dkim": true, "dkim-adsp": true, "dkim-atps": true,
	"dmarc": true, "dnswl": true, "domainkeys": true, "iprev": true, "rrvs": true,
	"sender-id": true, "smime": true, "spf": true, "vbr": true,
}

// VerifyEmailAuth reports whether the message's From address is authenticated.
// It is used for replies and sends, where From is the only proof that the mail
// comes from the alias owner, so it fails closed: only our own mail server's
// Authentication-Results header is read, and the message is rejected if that
// header is missing, duplicated or cfg.AuthServIDs is not configured.
func VerifyEmailAuth(data []byte, cfg EmailAuthConfig) (bool, error) {
	msg, err := mail.ReadMessage(bytes.NewReader(data))
	if err != nil {
		return false, err
	}

	msgID := msg.Header.Get("Message-ID")

	own, err := ownAuthResults(msg.Header, cfg.AuthServIDs)
	if err != nil {
		return false, withMessageID(err, msgID)
	}

	fromDomain, err := headerFromDomain(msg.Header)
	if err != nil {
		return false, withMessageID(err, msgID)
	}

	pass, authErr := senderAuthenticated(own, fromDomain)
	if pass {
		return true, nil
	}

	if arcVouches(msg.Header, own, fromDomain, cfg.ARCTrustedSealers) {
		return true, nil
	}

	return false, withMessageID(authErr, msgID)
}

// VerifyForwardAuth reports whether an inbound message may be forwarded to an
// alias's recipients. The sender is a third party, so alignment is left to the
// mail server's spam filtering; the message is only dropped when our own mail
// server found dmarc=fail and the From domain asks for rejection, unless a
// trusted ARC sealer vouches for it. A missing own header does not block
// forwarding, so a mail server misconfiguration doesn't drop inbound mail, but
// a duplicated or malformed one does, since only a forger produces one.
func VerifyForwardAuth(data []byte, cfg EmailAuthConfig) (bool, error) {
	msg, err := mail.ReadMessage(bytes.NewReader(data))
	if err != nil {
		return false, err
	}

	msgID := msg.Header.Get("Message-ID")

	own, err := ownAuthResults(msg.Header, cfg.AuthServIDs)
	if errors.Is(err, ErrDuplicateOwnAuthResults) || errors.Is(err, ErrMalformedOwnAuthResults) {
		return false, withMessageID(err, msgID)
	}
	if err != nil {
		return true, withMessageID(err, msgID)
	}

	if !dmarcPolicyReject(own) {
		return true, nil
	}

	fromDomain, err := headerFromDomain(msg.Header)
	if err == nil && arcVouches(msg.Header, own, fromDomain, cfg.ARCTrustedSealers) {
		return true, nil
	}

	return false, withMessageID(ErrDMARCPolicyReject, msgID)
}

// ownAuthResults returns the single Authentication-Results header written by
// one of our own mail servers, identified by servIDs. A message is received by
// one of them, so headers from more than one count as duplicates. The header
// must also have the shape Rspamd writes (see checkOwnAuthResults).
func ownAuthResults(header mail.Header, servIDs []string) (authResults, error) {
	configured := false
	for _, servID := range servIDs {
		if strings.TrimSpace(servID) != "" {
			configured = true
			break
		}
	}
	if !configured {
		return authResults{}, ErrAuthServIDNotConfigured
	}

	var own []authResults
	var ownValue string
	for _, value := range header["Authentication-Results"] {
		parsed := parseAuthResults(value)
		if isOwnServID(parsed.servID, servIDs) {
			own = append(own, parsed)
			ownValue = value
		}
	}

	switch len(own) {
	case 0:
		return authResults{}, ErrNoOwnAuthResults
	case 1:
		if reason := checkOwnAuthResults(own[0], ownValue); reason != "" {
			return authResults{}, fmt.Errorf("%w: %s", ErrMalformedOwnAuthResults, reason)
		}
		return own[0], nil
	default:
		return authResults{}, ErrDuplicateOwnAuthResults
	}
}

// checkOwnAuthResults returns why r, parsed from value, doesn't look like a
// header written by Rspamd, or "" if it does. Rspamd copies sender-controlled
// values (e.g. the SMTP MAIL FROM or DKIM d=) into its header, so a value that
// broke out of its quotes would show up as an extra result. Rspamd quotes such
// values but escapes only '"', not '\', so a value with a backslash can't be
// parsed unambiguously and is rejected. Rspamd writes all DKIM results first,
// then at most one dmarc and one spf result, and only registered methods.
func checkOwnAuthResults(r authResults, value string) string {
	if strings.Contains(value, `\`) {
		return "backslash in header"
	}
	if r.malformed != "" {
		return r.malformed
	}

	counts := map[string]int{}
	afterDKIM := false
	for _, ri := range r.results {
		if !knownMethods[ri.method] {
			return "unknown method " + strconv.Quote(ri.method)
		}
		if ri.method != "dkim" {
			afterDKIM = true
		} else if afterDKIM {
			return "dkim result after other results"
		}
		counts[ri.method]++
	}

	for _, method := range []string{"dmarc", "spf"} {
		if counts[method] > 1 {
			return "multiple " + method + " results"
		}
	}

	return ""
}

func isOwnServID(servID string, servIDs []string) bool {
	if servID == "" {
		return false
	}

	for _, own := range servIDs {
		if strings.EqualFold(servID, strings.TrimSpace(own)) {
			return true
		}
	}

	return false
}

// senderAuthenticated applies DMARC to r: dmarc=fail fails, and otherwise
// dmarc=pass for fromDomain passes, whatever order the results are in. Without
// a DMARC result, any DKIM signature or SPF pass whose domain is aligned with
// fromDomain is enough.
func senderAuthenticated(r authResults, fromDomain string) (bool, error) {
	dmarcPass := false
	for _, ri := range r.find("dmarc") {
		switch ri.result {
		case "pass":
			if strings.EqualFold(ri.props["header.from"], fromDomain) {
				dmarcPass = true
			}
		case "fail":
			return false, fmt.Errorf("%w, fromDomain: %s", ErrDMARCFail, fromDomain)
		}
	}
	if dmarcPass {
		return true, nil
	}

	for _, ri := range r.find("dkim") {
		if ri.result == "pass" && relaxedMatch(fromDomain, dkimDomain(ri)) {
			return true, nil
		}
	}

	for _, ri := range r.find("spf") {
		if ri.result == "pass" && relaxedMatch(fromDomain, extractDomain(ri.props["smtp.mailfrom"])) {
			return true, nil
		}
	}

	return false, fmt.Errorf("%w, fromDomain: %s", ErrNoAlignedAuthentication, fromDomain)
}

// dmarcPolicyReject reports whether r has dmarc=fail for a domain that
// publishes p=reject. Rspamd only reports the policy in a comment.
func dmarcPolicyReject(r authResults) bool {
	for _, ri := range r.find("dmarc") {
		if ri.result != "fail" {
			continue
		}
		if ri.props["policy.dmarc"] == "reject" || strings.Contains(ri.comments, "policy=reject") {
			return true
		}
	}
	return false
}

// arcVouches reports whether a trusted ARC sealer recorded an aligned pass for
// fromDomain before the message reached us. Our own mail server must report
// arc=pass, meaning it verified the whole ARC chain, so the latest ARC-Seal and
// the ARC-Authentication-Results header it covers are genuine. Rspamd may write
// one arc result per verified seal, so every one of them must be a pass.
func arcVouches(header mail.Header, own authResults, fromDomain string, trustedSealers []string) bool {
	arc := own.find("arc")
	if len(arc) == 0 {
		return false
	}
	for _, ri := range arc {
		if ri.result != "pass" {
			return false
		}
	}

	instance, sealer := latestARCSeal(header["Arc-Seal"])
	if instance == 0 || !isTrustedSealer(sealer, trustedSealers) {
		return false
	}

	var sealed []authResults
	for _, value := range header["Arc-Authentication-Results"] {
		parsed := parseAuthResults(value)
		if parsed.instance == instance {
			sealed = append(sealed, parsed)
		}
	}
	if len(sealed) != 1 {
		return false
	}

	pass, _ := senderAuthenticated(sealed[0], fromDomain)
	return pass
}

// latestARCSeal returns the instance (i=) and signing domain (d=) of the
// ARC-Seal with the highest instance, or 0 if there is none or it is ambiguous.
func latestARCSeal(seals []string) (int, string) {
	var instance, count int
	var domain string

	for _, seal := range seals {
		tags := parseTagList(seal)
		i, err := strconv.Atoi(tags["i"])
		if err != nil || i < 1 {
			continue
		}

		switch {
		case i > instance:
			instance, domain, count = i, tags["d"], 1
		case i == instance:
			count++
		}
	}

	if count != 1 {
		return 0, ""
	}
	return instance, domain
}

// parseTagList parses a DKIM-style "tag=value; tag=value" list, as used by
// ARC-Seal, lowercasing tags and values.
func parseTagList(value string) map[string]string {
	tags := map[string]string{}
	for part := range strings.SplitSeq(value, ";") {
		tag, val, ok := strings.Cut(part, "=")
		if !ok {
			continue
		}
		tags[strings.ToLower(strings.TrimSpace(tag))] = strings.ToLower(strings.TrimSpace(val))
	}
	return tags
}

// parseAuthResults parses an Authentication-Results header value, or an
// ARC-Authentication-Results value with its leading "i=N;" instance tag.
// Methods, results and property names and values are lowercased; comments are
// kept per entry because some results (e.g. Rspamd's DMARC policy) are only
// reported there.
func parseAuthResults(value string) authResults {
	segments, terminated := splitAuthResults(value)
	parsed := authResults{}
	if !terminated {
		parsed.malformed = "unterminated quoted string or comment"
	}

	if len(segments) > 0 {
		if tag, val, ok := strings.Cut(strings.TrimSpace(segments[0].text), "="); ok && strings.EqualFold(strings.TrimSpace(tag), "i") {
			parsed.instance, _ = strconv.Atoi(strings.TrimSpace(val))
			segments = segments[1:]
		}
	}

	if len(segments) == 0 {
		return parsed
	}

	if fields := splitFields(segments[0].text); len(fields) > 0 {
		parsed.servID = strings.ToLower(fields[0])
	}

	for _, seg := range segments[1:] {
		fields := splitFields(seg.text)
		if len(fields) == 0 {
			continue
		}

		method, result, ok := strings.Cut(fields[0], "=")
		if !ok {
			// "none" stands alone when there are no results (RFC 8601, section 2.2).
			if (len(segments) != 2 || !strings.EqualFold(fields[0], "none")) && parsed.malformed == "" {
				parsed.malformed = "result without method=value: " + strconv.Quote(fields[0])
			}
			continue
		}
		method, _, _ = strings.Cut(method, "/")

		ri := resInfo{
			method:   strings.ToLower(method),
			result:   strings.ToLower(unquote(result)),
			props:    map[string]string{},
			comments: strings.ToLower(seg.comments),
		}
		for _, field := range fields[1:] {
			if key, val, ok := strings.Cut(field, "="); ok {
				ri.props[strings.ToLower(key)] = strings.ToLower(unquote(val))
			}
		}
		parsed.results = append(parsed.results, ri)
	}

	return parsed
}

type authSegment struct {
	text     string
	comments string
}

// splitAuthResults splits a header value on ";" into segments, keeping
// parenthesized comments (which may nest) and quoted strings intact and
// separating comment text from the rest. Quoted strings are also honoured
// inside comments, because Rspamd quotes sender-controlled values there (e.g.
// "domain of "a)b"@example.com designates ..."). It reports whether every quoted
// string and comment was closed.
func splitAuthResults(value string) ([]authSegment, bool) {
	var segments []authSegment
	var text, comments strings.Builder
	depth := 0
	inQuote, escaped := false, false

	for _, r := range value {
		switch {
		case escaped:
			escaped = false
		case r == '\\' && (inQuote || depth > 0):
			escaped = true
			continue
		case r == '"':
			inQuote = !inQuote
		case inQuote:
		case r == '(':
			depth++
			if depth == 1 {
				comments.WriteByte(' ')
				continue
			}
		case r == ')' && depth > 0:
			depth--
			if depth == 0 {
				text.WriteByte(' ')
				continue
			}
		case depth == 0 && r == ';':
			segments = append(segments, authSegment{text: text.String(), comments: comments.String()})
			text.Reset()
			comments.Reset()
			continue
		}

		if depth > 0 {
			comments.WriteRune(r)
		} else {
			text.WriteRune(r)
		}
	}

	segments = append(segments, authSegment{text: text.String(), comments: comments.String()})
	return segments, !inQuote && depth == 0 && !escaped
}

// splitFields splits on whitespace outside quoted strings.
func splitFields(s string) []string {
	var fields []string
	var field strings.Builder
	inQuote := false

	for _, r := range s {
		switch {
		case r == '"':
			inQuote = !inQuote
		case !inQuote && (r == ' ' || r == '\t' || r == '\r' || r == '\n'):
			if field.Len() > 0 {
				fields = append(fields, field.String())
				field.Reset()
			}
			continue
		}
		field.WriteRune(r)
	}

	if field.Len() > 0 {
		fields = append(fields, field.String())
	}
	return fields
}

func unquote(s string) string {
	if len(s) > 1 && s[0] == '"' && s[len(s)-1] == '"' {
		return s[1 : len(s)-1]
	}
	return s
}

func (r authResults) find(method string) []resInfo {
	var found []resInfo
	for _, ri := range r.results {
		if ri.method == method {
			found = append(found, ri)
		}
	}
	return found
}

// dkimDomain returns the signing domain of a DKIM result: header.d, or the
// domain of header.i (Gmail reports only header.i=@domain).
func dkimDomain(ri resInfo) string {
	if d := ri.props["header.d"]; d != "" {
		return d
	}
	if i := ri.props["header.i"]; i != "" {
		return extractDomain(i)
	}
	return ""
}

func headerFromDomain(header mail.Header) (string, error) {
	fromAddr, err := mail.ParseAddress(header.Get("From"))
	if err != nil {
		return "", err
	}
	return extractDomain(fromAddr.Address), nil
}

func withMessageID(err error, msgID string) error {
	return fmt.Errorf("%w, Message-ID: %s", err, msgID)
}

func extractDomain(email string) string {
	// Trim whitespace
	email = strings.TrimSpace(email)

	// Trim surrounding quotes if present
	if len(email) > 1 && email[0] == '"' && email[len(email)-1] == '"' {
		email = email[1 : len(email)-1]
	}

	// Find the last '@' character and return the domain part
	if at := strings.LastIndex(email, "@"); at != -1 {
		return email[at+1:]
	}

	return email
}

func relaxedMatch(fromDomain, authDomain string) bool {
	if fromDomain == "" || authDomain == "" {
		return false
	}

	fromDomain = strings.ToLower(fromDomain)
	authDomain = strings.ToLower(authDomain)

	if fromDomain == authDomain {
		return true
	}

	// Require a label boundary so e.g. "notexample.com" does not match "example.com".
	return strings.HasSuffix(fromDomain, "."+authDomain) || strings.HasSuffix(authDomain, "."+fromDomain)
}

// isTrustedSealer reports whether domain is (or is a subdomain of) one of the
// operator-configured ARC sealers.
func isTrustedSealer(domain string, trustedSealers []string) bool {
	if domain == "" {
		return false
	}

	domain = strings.ToLower(domain)
	for _, trusted := range trustedSealers {
		trusted = strings.ToLower(strings.TrimSpace(trusted))
		if trusted != "" && (domain == trusted || strings.HasSuffix(domain, "."+trusted)) {
			return true
		}
	}

	return false
}
