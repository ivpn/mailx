package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
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

var localDomains = map[string]bool{}

var apiURL string
var psk string

type domainCheckRequest struct {
	Name string `json:"name"`
}

func main() {
	loadConfig()
	loadLocalDomains()

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

	// Skip check for local domains
	if localDomains[domain] {
		fmt.Fprintf(conn, "action=FILTER curl_email:\n\n")
		return
	}

	// Check cache
	if valid, found := checkCache(domain); found {
		if valid {
			fmt.Fprintf(conn, "action=FILTER curl_email:\n\n")
		} else {
			fmt.Fprintf(conn, "action=REJECT domain not configured\n\n")
		}
		return
	}

	// API check
	valid := checkDomain(domain)
	storeCache(domain, valid)

	if valid {
		fmt.Fprintf(conn, "action=FILTER curl_email:\n\n")
	} else {
		fmt.Fprintf(conn, "action=REJECT domain not configured\n\n")
	}
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

func checkDomain(domain string) bool {
	client := http.Client{
		Timeout: 5 * time.Second,
	}

	payload := domainCheckRequest{
		Name: domain,
	}

	jsonData, err := json.Marshal(payload)

	if err != nil {
		return false
	}

	req, err := http.NewRequest(
		"POST",
		apiURL+"/v1/email/domain/check",
		bytes.NewBuffer(jsonData),
	)

	if err != nil {
		return false
	}

	req.Header.Set("Authorization", "Bearer "+psk)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)

	if err != nil {
		return false
	}

	defer resp.Body.Close()

	return resp.StatusCode == 200
}
