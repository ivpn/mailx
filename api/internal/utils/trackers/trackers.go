// Package trackers removes open-tracking pixels from the HTML body of forwarded emails.
package trackers

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"slices"
	"strconv"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// MaxHTMLSize is the largest HTML body that gets cleaned; larger bodies are forwarded unchanged.
const MaxHTMLSize = 5 << 20

// MaxReportEntries caps the trackers listed in an encoded report; the total is always kept.
const MaxReportEntries = 50

// maxDomainLength is the longest valid hostname, so a crafted URL can't bloat the report.
const maxDomainLength = 253

var ErrTooLarge = errors.New("html body too large to clean")

// Element is where a tracker was found.
type Element string

const (
	ElementImg        Element = "img"
	ElementSrcset     Element = "srcset"
	ElementBackground Element = "background"
	ElementLink       Element = "link"
)

// Reason is why a URL or element was treated as a tracker.
type Reason string

const (
	ReasonKnownDomain Reason = "domain"
	ReasonTinyImage   Reason = "tiny"
	ReasonHiddenImage Reason = "hidden"
	ReasonURLPattern  Reason = "pattern"
)

// Tracker is one removed tracker. The short JSON keys keep the report link compact.
type Tracker struct {
	Domain  string  `json:"d"`
	Element Element `json:"e"`
	Reason  Reason  `json:"r"`
}

type Result struct {
	HTML    string
	Removed []Tracker
}

// Remove strips tracking pixels from an HTML body. The HTML is only re-serialized when a tracker
// was removed, so a body without trackers comes back unchanged. On error the original HTML is
// returned alongside the error.
func Remove(src string) (res Result, err error) {
	res.HTML = src
	if len(src) > MaxHTMLSize {
		return res, ErrTooLarge
	}

	defer func() {
		if r := recover(); r != nil {
			res, err = Result{HTML: src}, fmt.Errorf("clean html: %v", r)
		}
	}()

	root, err := parse(src)
	if err != nil {
		return res, err
	}

	var c cleaner
	c.walk(root)
	if len(c.removed) == 0 {
		return res, nil
	}

	out, err := render(root)
	if err != nil {
		return Result{HTML: src}, err
	}

	return Result{HTML: out, Removed: c.removed}, nil
}

// parse returns a root node whose children are the parsed input. Only full documents are parsed
// as documents, since that adds <html>, <head> and <body> wrappers to a fragment.
func parse(src string) (*html.Node, error) {
	if isDocument(src) {
		return html.Parse(strings.NewReader(src))
	}

	return parseFragment(src)
}

func parseFragment(src string) (*html.Node, error) {
	root := &html.Node{Type: html.ElementNode, Data: "body", DataAtom: atom.Body}
	nodes, err := html.ParseFragment(strings.NewReader(src), root)
	if err != nil {
		return nil, err
	}
	for _, n := range nodes {
		root.AppendChild(n)
	}

	return root, nil
}

// render serializes root's children, leaving out the root itself.
func render(root *html.Node) (string, error) {
	var b strings.Builder
	for n := range root.ChildNodes() {
		if err := html.Render(&b, n); err != nil {
			return "", err
		}
	}

	return b.String(), nil
}

func isDocument(src string) bool {
	lower := asciiLower(src)
	for _, tag := range []string{"<!doctype", "<html", "<head", "<body"} {
		if strings.Contains(lower, tag) {
			return true
		}
	}

	return false
}

type cleaner struct {
	removed []Tracker
}

func (c *cleaner) add(host string, element Element, reason Reason) {
	if len(host) > maxDomainLength {
		host = host[:maxDomainLength]
	}
	c.removed = append(c.removed, Tracker{Domain: host, Element: element, Reason: reason})
}

func (c *cleaner) walk(n *html.Node) {
	for child := n.FirstChild; child != nil; {
		next := child.NextSibling
		if child.Type == html.ElementNode {
			if c.clean(child) {
				n.RemoveChild(child)
			} else {
				c.walk(child)
			}
		}
		child = next
	}
}

// clean strips tracker URLs from n's attributes and reports whether n itself should be removed.
func (c *cleaner) clean(n *html.Node) bool {
	switch n.DataAtom {
	case atom.Img:
		if c.isTrackingImage(n) {
			return true
		}
		c.cleanSrcset(n)
	case atom.Source:
		c.cleanSrcset(n)
	case atom.Link:
		if host, reason, ok := match(getAttr(n, "href")); ok {
			c.add(host, ElementLink, reason)
			return true
		}
	case atom.Style:
		for t := range n.ChildNodes() {
			if t.Type == html.TextNode {
				t.Data = c.cleanCSS(t.Data)
			}
		}
	case atom.Noscript:
		// <noscript> content is raw text to the parser, but mail clients render it as HTML
		for t := range n.ChildNodes() {
			if t.Type == html.TextNode {
				t.Data = c.cleanNoscript(t.Data)
			}
		}
	}

	if host, reason, ok := match(getAttr(n, "background")); ok {
		c.add(host, ElementBackground, reason)
		removeAttr(n, "background")
	}
	if i := attrIndex(n, "style"); i >= 0 {
		n.Attr[i].Val = c.cleanCSS(n.Attr[i].Val)
	}

	return false
}

func (c *cleaner) cleanNoscript(text string) string {
	before := len(c.removed)
	root, err := parseFragment(text)
	if err != nil {
		return text
	}
	c.walk(root)
	if len(c.removed) == before {
		return text
	}

	out, err := render(root)
	if err != nil {
		c.removed = c.removed[:before]
		return text
	}

	return out
}

func (c *cleaner) isTrackingImage(n *html.Node) bool {
	src := getAttr(n, "src")
	if host, reason, ok := match(src); ok {
		c.add(host, ElementImg, reason)
		return true
	}

	// Tiny and hidden images only count when they load something remote
	host, remote := remoteHost(src)
	if !remote {
		for _, cand := range splitSrcset(getAttr(n, "srcset")) {
			if host, remote = remoteHost(cand.url); remote {
				break
			}
		}
	}
	if !remote {
		return false
	}

	style := parseStyle(getAttr(n, "style"))
	switch {
	case isHidden(n, style):
		c.add(host, ElementImg, ReasonHiddenImage)
	case isTiny(n, style):
		c.add(host, ElementImg, ReasonTinyImage)
	default:
		return false
	}

	return true
}

func (c *cleaner) cleanSrcset(n *html.Node) {
	i := attrIndex(n, "srcset")
	if i < 0 {
		return
	}

	candidates := splitSrcset(n.Attr[i].Val)
	kept := make([]string, 0, len(candidates))
	for _, cand := range candidates {
		if host, reason, ok := match(cand.url); ok {
			c.add(host, ElementSrcset, reason)
			continue
		}
		kept = append(kept, cand.text)
	}

	switch {
	case len(kept) == len(candidates):
	case len(kept) == 0:
		removeAttr(n, "srcset")
	default:
		n.Attr[i].Val = strings.Join(kept, ", ")
	}
}

// cleanCSS replaces every url(...) pointing at a tracker with none, which is valid wherever
// url(...) is and loads nothing.
func (c *cleaner) cleanCSS(css string) string {
	lower := asciiLower(css)
	var b strings.Builder
	last := 0
	for i := 0; ; {
		j := strings.Index(lower[i:], "url(")
		if j < 0 {
			break
		}
		start := i + j
		k := strings.IndexByte(lower[start:], ')')
		if k < 0 {
			break
		}
		end := start + k + 1

		raw := strings.Trim(strings.TrimSpace(css[start+len("url("):end-1]), `"'`)
		if host, reason, ok := match(raw); ok {
			c.add(host, ElementBackground, reason)
			b.WriteString(css[last:start])
			b.WriteString("none")
			last = end
		}
		i = end
	}

	if last == 0 {
		return css
	}
	b.WriteString(css[last:])

	return b.String()
}

func isHidden(n *html.Node, style map[string]string) bool {
	if attrIndex(n, "hidden") >= 0 {
		return true
	}
	if style["display"] == "none" || style["visibility"] == "hidden" {
		return true
	}
	if v, _, ok := parseLength(style["opacity"]); ok && v <= 0 {
		return true
	}
	for _, dim := range []string{"width", "height"} {
		if v, _, ok := parseLength(style[dim]); ok && v == 0 {
			return true
		}
	}

	return false
}

func isTiny(n *html.Node, style map[string]string) bool {
	for _, dim := range []string{"width", "height"} {
		for _, val := range []string{getAttr(n, dim), style[dim]} {
			if v, unit, ok := parseLength(val); ok && v <= 1 && (unit == "" || unit == "px") {
				return true
			}
		}
	}

	return false
}

// parseLength parses a leading number and returns it with the unit that follows it.
func parseLength(s string) (float64, string, bool) {
	s = strings.TrimSpace(s)
	i := 0
	for i < len(s) && (s[i] >= '0' && s[i] <= '9' || s[i] == '.' || i == 0 && (s[i] == '-' || s[i] == '+')) {
		i++
	}
	v, err := strconv.ParseFloat(s[:i], 64)
	if err != nil {
		return 0, "", false
	}

	return v, asciiLower(strings.TrimSpace(s[i:])), true
}

// parseStyle parses an inline style into lowercased property/value pairs; later declarations
// win, as they do in CSS.
func parseStyle(s string) map[string]string {
	style := map[string]string{}
	for decl := range strings.SplitSeq(s, ";") {
		prop, val, ok := strings.Cut(decl, ":")
		if !ok {
			continue
		}
		val = strings.TrimSuffix(asciiLower(strings.TrimSpace(val)), "!important")
		style[asciiLower(strings.TrimSpace(prop))] = strings.TrimSpace(val)
	}

	return style
}

type srcsetCandidate struct {
	url  string
	text string
}

// splitSrcset splits a srcset into candidates. URLs may contain commas (e.g. data: URIs or image
// CDN parameters), so a URL runs to the next whitespace and descriptors to the next comma.
func splitSrcset(s string) []srcsetCandidate {
	var out []srcsetCandidate
	for i := 0; i < len(s); {
		for i < len(s) && (isSpace(s[i]) || s[i] == ',') {
			i++
		}
		if i >= len(s) {
			break
		}

		start := i
		for i < len(s) && !isSpace(s[i]) {
			i++
		}
		u := s[start:i]
		if strings.HasSuffix(u, ",") {
			u = strings.TrimRight(u, ",")
			out = append(out, srcsetCandidate{url: u, text: u})
			continue
		}

		depth := 0
	descriptors:
		for ; i < len(s); i++ {
			switch s[i] {
			case '(':
				depth++
			case ')':
				depth = max(depth-1, 0)
			case ',':
				if depth == 0 {
					break descriptors
				}
			}
		}
		out = append(out, srcsetCandidate{url: u, text: strings.TrimSpace(s[start:i])})
	}

	return out
}

// match reports whether raw points at a known tracker domain or tracking URL pattern, along
// with the host it points at.
func match(raw string) (string, Reason, bool) {
	host, path, ok := splitURL(raw)
	if !ok {
		return "", "", false
	}
	if isKnownDomain(host) {
		return host, ReasonKnownDomain, true
	}
	for _, p := range urlPatterns {
		if p.MatchString(path) {
			return host, ReasonURLPattern, true
		}
	}

	return "", "", false
}

func remoteHost(raw string) (string, bool) {
	host, _, ok := splitURL(raw)
	return host, ok
}

// splitURL returns the lowercased host and path of an http(s) or protocol-relative URL, read the
// way browsers read it. Relative URLs and cid: or data: URIs load nothing remote, so they're
// never reported.
func splitURL(raw string) (host, path string, ok bool) {
	s := strings.Map(func(r rune) rune {
		switch r {
		case '\t', '\n', '\r':
			return -1
		case '\\':
			return '/'
		}
		return r
	}, strings.TrimSpace(raw))

	lower := asciiLower(s)
	switch {
	case strings.HasPrefix(lower, "https:"):
		s = s[len("https:"):]
	case strings.HasPrefix(lower, "http:"):
		s = s[len("http:"):]
	case strings.HasPrefix(s, "//"):
	default:
		return "", "", false
	}
	s = strings.TrimLeft(s, "/")

	authority := s
	if i := strings.IndexAny(s, "/?#"); i >= 0 {
		authority, path = s[:i], s[i:]
	}
	if i := strings.LastIndexByte(authority, '@'); i >= 0 {
		authority = authority[i+1:]
	}
	host = authority
	if strings.HasPrefix(host, "[") {
		if i := strings.IndexByte(host, ']'); i >= 0 {
			host = host[:i+1]
		}
	} else if i := strings.LastIndexByte(host, ':'); i >= 0 {
		host = host[:i]
	}
	if h, err := url.PathUnescape(host); err == nil {
		host = h
	}
	host = strings.TrimSuffix(asciiLower(host), ".")
	if host == "" {
		return "", "", false
	}

	if i := strings.IndexAny(path, "?#"); i >= 0 {
		path = path[:i]
	}
	if p, err := url.PathUnescape(path); err == nil {
		path = p
	}

	return host, asciiLower(path), true
}

func isKnownDomain(host string) bool {
	for {
		if _, ok := knownDomains[host]; ok {
			return true
		}
		i := strings.IndexByte(host, '.')
		if i < 0 {
			return false
		}
		host = host[i+1:]
	}
}

func attrIndex(n *html.Node, key string) int {
	return slices.IndexFunc(n.Attr, func(a html.Attribute) bool {
		return a.Namespace == "" && a.Key == key
	})
}

func getAttr(n *html.Node, key string) string {
	if i := attrIndex(n, key); i >= 0 {
		return n.Attr[i].Val
	}

	return ""
}

func removeAttr(n *html.Node, key string) {
	n.Attr = slices.DeleteFunc(n.Attr, func(a html.Attribute) bool {
		return a.Namespace == "" && a.Key == key
	})
}

func isSpace(b byte) bool {
	return b == ' ' || b == '\t' || b == '\n' || b == '\r' || b == '\f'
}

// asciiLower lowercases ASCII letters only, so byte offsets match the input. strings.ToLower
// can change the length of some non-ASCII characters.
func asciiLower(s string) string {
	i := strings.IndexFunc(s, func(r rune) bool { return r >= 'A' && r <= 'Z' })
	if i < 0 {
		return s
	}

	b := []byte(s)
	for j := i; j < len(b); j++ {
		if b[j] >= 'A' && b[j] <= 'Z' {
			b[j] += 'a' - 'A'
		}
	}

	return string(b)
}

type report struct {
	Total    int       `json:"n"`
	Trackers []Tracker `json:"t"`
}

// EncodeReport encodes removed trackers as base64url compact JSON for the report link. Only the
// first MaxReportEntries are listed, alongside the total.
func EncodeReport(removed []Tracker) string {
	r := report{
		Total:    len(removed),
		Trackers: removed[:min(len(removed), MaxReportEntries)],
	}
	b, _ := json.Marshal(r) // #nosec G104 -- only strings and ints, can't fail

	return base64.RawURLEncoding.EncodeToString(b)
}

// DecodeReport decodes a report produced by EncodeReport.
func DecodeReport(s string) (int, []Tracker, error) {
	b, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		return 0, nil, err
	}

	var r report
	if err := json.Unmarshal(b, &r); err != nil {
		return 0, nil, err
	}

	return r.Total, r.Trackers, nil
}
