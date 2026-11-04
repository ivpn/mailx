# gomail.v2 fork

In-house fork of [go-gomail/gomail](https://github.com/go-gomail/gomail) (`gopkg.in/gomail.v2`).

**Base:** upstream commit `81ebce5c23dfd25c6c67194b37d3dd3f338c98b1` (2016-04-11), which is also upstream's last commit. Upstream is unmaintained and its maintained fork [go-mail/mail](https://github.com/go-mail/mail) is archived, so there is nothing to re-sync from. Security fixes must be made here. govulncheck and OSV can't see this code because it isn't a module.

Don't replace this directory with upstream or with `go-mail/mail`. The local changes below would be lost, and some of them would be lost silently.

## Local changes

| File | Change | Lost silently on re-sync? |
|---|---|---|
| `message.go` | `NewRawMessage()`, `SetRawBody()`, `newRawPart()`, used for PGP/MIME (`utils.EncryptWithPGPMIME`) | No: compile error |
| `writeto.go` | `writePart` omits `Content-Type`/`Content-Transfer-Encoding` when the part has no content type (raw PGP/MIME bodies) | **Yes**: covered by `utils/email_test.go` |
| `smtp.go` | `Dialer.StartTLSPolicy`. `NewDialer` sets `MandatoryStartTLS`, so `Dial` fails with `ErrStartTLSUnsupported` when the server doesn't offer STARTTLS, instead of continuing in plaintext | **Yes**: covered by `starttls_policy_test.go` |
| `smtp.go` | Default TLS config sets `MinVersion: tls.VersionTLS12` | **Yes**: covered by `starttls_policy_test.go` |
| `auth.go` | LOGIN auth refuses unencrypted connections except to localhost. Upstream allowed it whenever the server advertised LOGIN | **Yes**: covered by `auth_test.go` |
| `mime_go14.go` | Removed (pre-Go 1.5 fallback; dropped the `gopkg.in/alexcesaro/quotedprintable.v3` dependency) | n/a |
| `auth.go`, `send.go`, `message.go` | `slices` / `maps` modernisation | n/a |

## Tests

Upstream's tests (`auth_test.go`, `example_test.go`, `message_test.go`, `send_test.go`, `smtp_test.go`) are kept. `auth_test.go` is adjusted for the LOGIN change. `starttls_policy_test.go` covers the TLS changes.

```sh
go test ./internal/utils/gomail.v2/
```
