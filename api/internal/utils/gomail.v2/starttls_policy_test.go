package gomail

import (
	"crypto/tls"
	"errors"
	"io"
	"net"
	"net/smtp"
	"testing"
	"time"
)

// noStartTLSClient is an smtpClient whose server does not advertise STARTTLS
// (as after an on-path attacker strips it). It records the commands it gets.
type noStartTLSClient struct {
	cmds []string
}

func (c *noStartTLSClient) Hello(string) error { return nil }
func (c *noStartTLSClient) Extension(ext string) (bool, string) {
	c.cmds = append(c.cmds, "Extension "+ext)
	if ext == "AUTH" {
		return true, "LOGIN"
	}
	return false, ""
}
func (c *noStartTLSClient) StartTLS(*tls.Config) error {
	c.cmds = append(c.cmds, "StartTLS")
	return nil
}
func (c *noStartTLSClient) Auth(smtp.Auth) error {
	c.cmds = append(c.cmds, "Auth")
	return nil
}
func (c *noStartTLSClient) Mail(string) error             { return nil }
func (c *noStartTLSClient) Rcpt(string) error             { return nil }
func (c *noStartTLSClient) Data() (io.WriteCloser, error) { return nil, nil }
func (c *noStartTLSClient) Quit() error                   { return nil }
func (c *noStartTLSClient) Close() error {
	c.cmds = append(c.cmds, "Close")
	return nil
}

func stubNoStartTLS(t *testing.T) *noStartTLSClient {
	t.Helper()
	client := &noStartTLSClient{}

	origDial, origNew := netDialTimeout, smtpNewClient
	t.Cleanup(func() { netDialTimeout, smtpNewClient = origDial, origNew })

	netDialTimeout = func(string, string, time.Duration) (net.Conn, error) {
		return testConn, nil
	}
	smtpNewClient = func(net.Conn, string) (smtpClient, error) {
		return client, nil
	}

	return client
}

func TestNewDialerRequiresStartTLS(t *testing.T) {
	d := NewDialer(testHost, testPort, testUser, testPwd)
	if d.StartTLSPolicy != MandatoryStartTLS {
		t.Fatalf("NewDialer StartTLSPolicy = %v, want MandatoryStartTLS", d.StartTLSPolicy)
	}
}

func TestDialMandatoryStartTLSFailsClosed(t *testing.T) {
	client := stubNoStartTLS(t)

	d := NewDialer(testHost, testPort, testUser, testPwd)
	_, err := d.Dial()
	if !errors.Is(err, ErrStartTLSUnsupported) {
		t.Fatalf("Dial error = %v, want ErrStartTLSUnsupported", err)
	}

	for _, cmd := range client.cmds {
		if cmd == "Auth" {
			t.Fatal("credentials were sent over a connection without STARTTLS")
		}
	}
	if last := client.cmds[len(client.cmds)-1]; last != "Close" {
		t.Errorf("connection not closed, last command %q", last)
	}
}

func TestDialOpportunisticStartTLSContinues(t *testing.T) {
	stubNoStartTLS(t)

	d := &Dialer{Host: testHost, Port: testPort}
	if _, err := d.Dial(); err != nil {
		t.Fatalf("opportunistic Dial without credentials failed: %v", err)
	}
}

func TestDialSSLIgnoresStartTLSPolicy(t *testing.T) {
	stubNoStartTLS(t)
	origTLS := tlsClient
	t.Cleanup(func() { tlsClient = origTLS })
	tlsClient = func(net.Conn, *tls.Config) *tls.Conn { return testTLSConn }

	d := NewDialer(testHost, testSSLPort, testUser, testPwd)
	d.Auth = testAuth
	if _, err := d.Dial(); err != nil {
		t.Fatalf("implicit-TLS Dial failed: %v", err)
	}
}

func TestDefaultTLSConfigMinVersion(t *testing.T) {
	d := NewDialer(testHost, testPort, testUser, testPwd)
	if v := d.tlsConfig().MinVersion; v != tls.VersionTLS12 {
		t.Fatalf("default MinVersion = %x, want TLS 1.2", v)
	}
}
