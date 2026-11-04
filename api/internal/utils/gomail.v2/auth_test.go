package gomail

import (
	"net/smtp"
	"testing"
)

const (
	testUser = "user"
	testPwd  = "pwd"
	testHost = "smtp.example.com"
)

type authTest struct {
	auths      []string
	challenges []string
	tls        bool
	wantData   []string
	wantError  bool
}

func TestNoAdvertisement(t *testing.T) {
	testLoginAuth(t, &authTest{
		auths:     []string{},
		tls:       false,
		wantError: true,
	})
}

func TestNoAdvertisementTLS(t *testing.T) {
	testLoginAuth(t, &authTest{
		auths:      []string{},
		challenges: []string{"Username:", "Password:"},
		tls:        true,
		wantData:   []string{"", testUser, testPwd},
	})
}

// Upstream allowed LOGIN over plaintext whenever the server advertised it;
// the fork refuses, since a stripped-STARTTLS MITM can advertise anything.
func TestLoginNoTLS(t *testing.T) {
	testLoginAuth(t, &authTest{
		auths:     []string{"PLAIN", "LOGIN"},
		tls:       false,
		wantError: true,
	})
}

func TestLoginNoTLSLocalhost(t *testing.T) {
	testLoginAuthHost(t, "localhost", &authTest{
		auths:      []string{"LOGIN"},
		challenges: []string{"Username:", "Password:"},
		tls:        false,
		wantData:   []string{"", testUser, testPwd},
	})
}

func TestLoginTLS(t *testing.T) {
	testLoginAuth(t, &authTest{
		auths:      []string{"LOGIN"},
		challenges: []string{"Username:", "Password:"},
		tls:        true,
		wantData:   []string{"", testUser, testPwd},
	})
}

func testLoginAuth(t *testing.T, test *authTest) {
	testLoginAuthHost(t, testHost, test)
}

func testLoginAuthHost(t *testing.T, host string, test *authTest) {
	auth := &loginAuth{
		username: testUser,
		password: testPwd,
		host:     host,
	}
	server := &smtp.ServerInfo{
		Name: host,
		TLS:  test.tls,
		Auth: test.auths,
	}
	proto, toServer, err := auth.Start(server)
	if err != nil && !test.wantError {
		t.Fatalf("loginAuth.Start(): %v", err)
	}
	if err != nil && test.wantError {
		return
	}
	if test.wantError {
		t.Fatal("loginAuth.Start(): expected error, got nil")
	}
	if proto != "LOGIN" {
		t.Errorf("invalid protocol, got %q, want LOGIN", proto)
	}

	i := 0
	got := string(toServer)
	if got != test.wantData[i] {
		t.Errorf("Invalid response, got %q, want %q", got, test.wantData[i])
	}

	for _, challenge := range test.challenges {
		i++
		if i >= len(test.wantData) {
			t.Fatalf("unexpected challenge: %q", challenge)
		}

		toServer, err = auth.Next([]byte(challenge), true)
		if err != nil {
			t.Fatalf("loginAuth.Auth(): %v", err)
		}
		got = string(toServer)
		if got != test.wantData[i] {
			t.Errorf("Invalid response, got %q, want %q", got, test.wantData[i])
		}
	}
}
