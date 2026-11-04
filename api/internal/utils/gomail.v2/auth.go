package gomail

import (
	"bytes"
	"errors"
	"fmt"
	"net/smtp"
)

// loginAuth is an smtp.Auth that implements the LOGIN authentication mechanism.
type loginAuth struct {
	username string
	password string
	host     string
}

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	// Like net/smtp's PlainAuth, never send the password over an unencrypted
	// connection except to localhost. Upstream gomail allowed it whenever the
	// server advertised LOGIN, which an on-path attacker can fake.
	if !server.TLS && !isLocalhost(server.Name) {
		return "", nil, errors.New("gomail: unencrypted connection")
	}
	if server.Name != a.host {
		return "", nil, errors.New("gomail: wrong host name")
	}
	return "LOGIN", nil, nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if !more {
		return nil, nil
	}

	switch {
	case bytes.Equal(fromServer, []byte("Username:")):
		return []byte(a.username), nil
	case bytes.Equal(fromServer, []byte("Password:")):
		return []byte(a.password), nil
	default:
		return nil, fmt.Errorf("gomail: unexpected server challenge: %s", fromServer)
	}
}

func isLocalhost(name string) bool {
	return name == "localhost" || name == "127.0.0.1" || name == "::1"
}
