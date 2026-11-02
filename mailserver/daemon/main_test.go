package main

import (
	"bufio"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// fakeAPI stands in for /v1/email/domain/check and counts the requests it gets.
func fakeAPI(t *testing.T, status map[string]int) *int {
	t.Helper()

	calls := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		body, _ := io.ReadAll(r.Body)
		for domain, code := range status {
			if strings.Contains(string(body), `"`+domain+`"`) {
				w.WriteHeader(code)
				return
			}
		}
		w.WriteHeader(http.StatusBadRequest)
	}))
	t.Cleanup(srv.Close)

	apiURL = srv.URL
	cache = map[string]cacheEntry{}
	localDomains = map[string]bool{"irelay.app": true}

	return &calls
}

func TestSocketmapReply(t *testing.T) {
	calls := fakeAPI(t, map[string]int{
		"custom.example": http.StatusOK,
		"broken.example": http.StatusInternalServerError,
	})

	tests := []struct {
		req  string
		want string
	}{
		{"relay irelay.app", "OK irelay.app"},
		{"relay IRELAY.APP", "OK irelay.app"},
		{"relay custom.example", "OK custom.example"},
		{"relay gmail.com", "NOTFOUND "},
		{"relay broken.example", "TEMP domain lookup failed"},
		// Parent-domain lookups never match, so subdomains of our domains are not relayed
		{"relay .irelay.app", "NOTFOUND "},
		{"relay ", "NOTFOUND "},
		{"other irelay.app", "PERM unknown map"},
		{"relay", "PERM unknown map"},
	}

	for _, tt := range tests {
		if got := socketmapReply(tt.req); got != tt.want {
			t.Errorf("socketmapReply(%q) = %q, want %q", tt.req, got, tt.want)
		}
	}

	// Local and cached domains don't call the API; API errors are not cached
	before := *calls
	socketmapReply("relay irelay.app")
	socketmapReply("relay custom.example")
	socketmapReply("relay gmail.com")
	if *calls != before {
		t.Errorf("cached lookups called the API %d times", *calls-before)
	}

	socketmapReply("relay broken.example")
	if *calls != before+1 {
		t.Errorf("failed lookup was cached")
	}
}

func TestHandleSocketmap(t *testing.T) {
	fakeAPI(t, nil)

	client, server := net.Pipe()
	go handleSocketmap(server)
	defer client.Close()

	reader := bufio.NewReader(client)

	// Postfix sends several lookups on one connection
	for _, tt := range []struct{ req, want string }{
		{"relay irelay.app", "OK irelay.app"},
		{"relay gmail.com", "NOTFOUND "},
	} {
		go writeNetstring(client, tt.req)

		got, err := readNetstring(reader)
		if err != nil {
			t.Fatalf("reading reply to %q: %v", tt.req, err)
		}
		if got != tt.want {
			t.Errorf("reply to %q = %q, want %q", tt.req, got, tt.want)
		}
	}
}

func TestReadNetstring(t *testing.T) {
	tests := []struct {
		in      string
		want    string
		wantErr bool
	}{
		{"16:relay irelay.app,", "relay irelay.app", false},
		{"0:,", "", false},
		{"5:relay;", "", true},
		{":relay,", "", true},
		{"x:relay,", "", true},
		{"10001:", "", true},
		{"99999999999999999999:", "", true},
		{"10:short,", "", true},
	}

	for _, tt := range tests {
		got, err := readNetstring(bufio.NewReader(strings.NewReader(tt.in)))
		if (err != nil) != tt.wantErr || got != tt.want {
			t.Errorf("readNetstring(%q) = %q, %v; want %q, err %v", tt.in, got, err, tt.want, tt.wantErr)
		}
	}
}
