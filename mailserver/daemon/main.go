package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

type cacheEntry struct {
	valid bool
	exp   time.Time
}

var cache = map[string]cacheEntry{}
var mu sync.RWMutex

var positiveTTL = 10 * time.Minute
var negativeTTL = 2 * time.Minute

const curlEmailLogPath = "/var/log/mail/curl-email.log"

var logClearInterval = 24 * time.Hour

var localDomains = map[string]bool{}

// Postfix socketmap table that answers relay_domains lookups:
// relay_domains = socketmap:inet:daemon:10026:relay
const socketmapName = "relay"
const socketmapMaxRequest = 10000
const socketmapIdleTimeout = 2 * time.Minute

var apiURL string
var psk string

type domainCheckRequest struct {
	Name string `json:"name"`
}

func main() {
	loadConfig()
	loadLocalDomains()

	go runLogCleaner()
	go serveSocketmap()

	ln, err := net.Listen("tcp", ":10025")

	if err != nil {
		panic(err)
	}

	fmt.Println("Domain Policy Service (daemon) listening on 10025")

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}

		go handle(conn)
	}
}

func loadConfig() {
	apiURL = os.Getenv("API_URL")
	psk = os.Getenv("PSK")

	if ttl := os.Getenv("CACHE_POSITIVE_TTL"); ttl != "" {
		if d, err := time.ParseDuration(ttl); err == nil {
			positiveTTL = d
		}
	}

	if ttl := os.Getenv("CACHE_NEGATIVE_TTL"); ttl != "" {
		if d, err := time.ParseDuration(ttl); err == nil {
			negativeTTL = d
		}
	}

	if interval := os.Getenv("CURL_EMAIL_LOG_CLEAR_INTERVAL"); interval != "" {
		if d, err := time.ParseDuration(interval); err == nil {
			logClearInterval = d
		}
	}
}

func loadLocalDomains() {
	domains := os.Getenv("LOCAL_DOMAINS")

	if domains == "" {
		return
	}

	for _, d := range strings.Split(domains, ",") {
		domain := strings.TrimSpace(d)
		if domain != "" {
			localDomains[domain] = true
		}
	}
}

func runLogCleaner() {
	ticker := time.NewTicker(logClearInterval)
	defer ticker.Stop()

	for range ticker.C {
		clearCurlEmailLog()
	}
}

func clearCurlEmailLog() {
	// O_CREATE so this is resilient if the daemon starts before the mailserver has created the file
	f, err := os.OpenFile(curlEmailLogPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("failed to clear curl-email log:", err)
		return
	}

	f.Close()
}

func handle(conn net.Conn) {
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(10 * time.Second))
	scanner := bufio.NewScanner(conn)
	var recipient string

	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			break
		}

		if strings.HasPrefix(line, "recipient=") {
			recipient = strings.TrimPrefix(line, "recipient=")
		}
	}

	if recipient == "" {
		fmt.Fprintf(conn, "action=DUNNO\n\n")
		return
	}

	parts := strings.Split(recipient, "@")

	if len(parts) != 2 {
		fmt.Fprintf(conn, "action=REJECT invalid recipient\n\n")
		return
	}

	domain := strings.ToLower(parts[1])

	valid, err := lookupDomain(domain)
	if err != nil {
		fmt.Println("domain check failed:", err)
	}

	if valid {
		fmt.Fprintf(conn, "action=FILTER curl_email:\n\n")
	} else {
		fmt.Fprintf(conn, "action=REJECT domain not configured\n\n")
	}
}

// lookupDomain reports whether mail for domain is handled by the API: a local
// domain or a verified custom domain. API errors are returned and not cached.
func lookupDomain(domain string) (bool, error) {
	// Skip check for local domains
	if localDomains[domain] {
		return true, nil
	}

	// Check cache
	if valid, found := checkCache(domain); found {
		return valid, nil
	}

	// API check
	valid, err := checkDomain(domain)
	if err != nil {
		return false, err
	}

	storeCache(domain, valid)

	return valid, nil
}

// serveSocketmap answers Postfix socketmap lookups for relay_domains, so mail
// Postfix generates itself (bounces) is also routed to the API by relay_transport.
func serveSocketmap() {
	ln, err := net.Listen("tcp", ":10026")

	if err != nil {
		panic(err)
	}

	fmt.Println("Relay Domains Socketmap (daemon) listening on 10026")

	for {
		conn, err := ln.Accept()
		if err != nil {
			continue
		}

		go handleSocketmap(conn)
	}
}

// handleSocketmap serves one Postfix connection, which stays open for many
// lookups. Each request is a netstring "<name> <key>".
func handleSocketmap(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {
		conn.SetDeadline(time.Now().Add(socketmapIdleTimeout))

		req, err := readNetstring(reader)
		if err != nil {
			return
		}

		if err := writeNetstring(conn, socketmapReply(req)); err != nil {
			return
		}
	}
}

func socketmapReply(req string) string {
	name, key, ok := strings.Cut(req, " ")
	if !ok || name != socketmapName {
		return "PERM unknown map"
	}

	domain := strings.ToLower(key)

	// Postfix also looks up parent domains as ".example.com".
	// Only exact domains are ours, as in the policy check.
	if domain == "" || strings.HasPrefix(domain, ".") {
		return "NOTFOUND "
	}

	valid, err := lookupDomain(domain)
	if err != nil {
		// Postfix defers the message and retries, instead of routing it elsewhere
		fmt.Println("domain check failed:", err)
		return "TEMP domain lookup failed"
	}

	if !valid {
		return "NOTFOUND "
	}

	return "OK " + domain
}

func readNetstring(r *bufio.Reader) (string, error) {
	length := 0
	digits := 0

	for {
		b, err := r.ReadByte()
		if err != nil {
			return "", err
		}

		if b == ':' && digits > 0 {
			break
		}

		if b < '0' || b > '9' {
			return "", errors.New("invalid netstring length")
		}

		length = length*10 + int(b-'0')
		digits++

		if length > socketmapMaxRequest {
			return "", errors.New("netstring too long")
		}
	}

	buf := make([]byte, length+1)
	if _, err := io.ReadFull(r, buf); err != nil {
		return "", err
	}

	if buf[length] != ',' {
		return "", errors.New("missing netstring terminator")
	}

	return string(buf[:length]), nil
}

func writeNetstring(w io.Writer, s string) error {
	_, err := fmt.Fprintf(w, "%d:%s,", len(s), s)
	return err
}

func checkCache(domain string) (bool, bool) {
	mu.RLock()
	entry, ok := cache[domain]
	mu.RUnlock()

	if !ok {
		return false, false
	}

	if time.Now().After(entry.exp) {
		mu.Lock()
		delete(cache, domain)
		mu.Unlock()
		return false, false
	}

	return entry.valid, true
}

func storeCache(domain string, valid bool) {
	var ttl time.Duration

	if valid {
		ttl = positiveTTL
	} else {
		ttl = negativeTTL
	}

	mu.Lock()

	cache[domain] = cacheEntry{
		valid: valid,
		exp:   time.Now().Add(ttl),
	}

	mu.Unlock()
}

// checkDomain asks the API whether domain is a verified custom domain.
// The API answers 200 for a verified domain and 400 for anything else;
// any other result is an error, not a "no".
func checkDomain(domain string) (bool, error) {
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	payload := domainCheckRequest{
		Name: domain,
	}

	jsonData, err := json.Marshal(payload)

	if err != nil {
		return false, err
	}

	req, err := http.NewRequest(
		"POST",
		apiURL+"/v1/email/domain/check",
		bytes.NewBuffer(jsonData),
	)

	if err != nil {
		return false, err
	}

	req.Header.Set("Authorization", "Bearer "+psk)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)

	if err != nil {
		return false, err
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		return true, nil
	case http.StatusBadRequest:
		return false, nil
	default:
		return false, fmt.Errorf("unexpected status %d for %s", resp.StatusCode, domain)
	}
}
