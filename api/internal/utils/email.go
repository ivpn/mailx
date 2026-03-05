package utils

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io"
	"math/big"
	"mime"
	"net/mail"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/ProtonMail/gopenpgp/v3/crypto"
	xcharset "golang.org/x/net/html/charset"
	"ivpn.net/email/api/internal/utils/gomail.v2"
)

func RemoveHeader(text string) string {
	re := regexp.MustCompile(`(?m)^This email was sent to .+? from .+\n?`)
	return re.ReplaceAllString(text, "")
}

func RemoveHtmlHeader(html string) string {
	// Relaxed regex: match any <table> containing "This email was sent to" and ending at </table>
	re := regexp.MustCompile(`(?is)<table[^>]*>.*?This email was sent to.*?</table>`)
	cleaned := re.ReplaceAllString(html, "")

	// Optionally clean up one or more immediate trailing <br> tags or empty <div><br></div>
	cleaned = regexp.MustCompile(`(?i)(\s*<br\s*/?>\s*|<div[^>]*>\s*(<br\s*/?>)?\s*</div>)+`).ReplaceAllString(cleaned, "")

	return cleaned
}

func EncryptWithPGPInline(plainText string, recipientKey string) (string, error) {
	publicKey, err := crypto.NewKeyFromArmored(recipientKey)
	if err != nil {
		return "", fmt.Errorf("parse public key: %w", err)
	}

	pgp := crypto.PGP()
	encHandle, err := pgp.Encryption().Recipient(publicKey).New()
	if err != nil {
		return "", fmt.Errorf("create encryption handle: %w", err)
	}

	pgpMessage, err := encHandle.Encrypt([]byte(plainText))
	if err != nil {
		return "", fmt.Errorf("encrypt text: %w", err)
	}

	armored, err := pgpMessage.ArmorBytes()
	if err != nil {
		return "", fmt.Errorf("armor ciphertext: %w", err)
	}

	return string(armored), nil
}

func EncryptWithPGPMIME(orig *gomail.Message, fromAddr, fromName, subject, recipientEmail, recipientKey string) (*gomail.Message, error) {
	// --- 1) Serialize the original email ---
	var buf bytes.Buffer
	if _, err := orig.WriteTo(&buf); err != nil {
		return nil, fmt.Errorf("serialize original email: %w", err)
	}

	// --- 2) Parse recipient public key ---
	publicKey, err := crypto.NewKeyFromArmored(recipientKey)
	if err != nil {
		return nil, fmt.Errorf("parse public key: %w", err)
	}

	// --- 3) Encrypt body ---
	pgp := crypto.PGP()
	encHandle, err := pgp.Encryption().Recipient(publicKey).New()
	if err != nil {
		return nil, fmt.Errorf("create encryption handle: %w", err)
	}

	pgpMessage, err := encHandle.Encrypt(buf.Bytes())
	if err != nil {
		return nil, fmt.Errorf("encrypt payload: %w", err)
	}

	armored, err := pgpMessage.ArmorBytes()
	if err != nil {
		return nil, fmt.Errorf("armor ciphertext: %w", err)
	}

	// Normalize line endings to CRLF
	armoredStr := strings.ReplaceAll(string(armored), "\n", "\r\n")

	// --- 4) Build PGP/MIME multipart body ---
	boundary := "boundary-" + randomChars(16)
	var body bytes.Buffer

	// Part 1: version
	body.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	body.WriteString("Content-Type: application/pgp-encrypted\r\n")
	body.WriteString("Content-Description: PGP/MIME version identification\r\n")
	body.WriteString("Content-Transfer-Encoding: 7bit\r\n\r\n")
	body.WriteString("Version: 1\r\n\r\n")

	// Part 2: encrypted content
	body.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	body.WriteString("Content-Type: application/octet-stream; name=\"encrypted.asc\"\r\n")
	body.WriteString("Content-Description: OpenPGP encrypted message\r\n")
	body.WriteString("Content-Disposition: inline; filename=\"encrypted.asc\"\r\n")
	body.WriteString("Content-Transfer-Encoding: 7bit\r\n\r\n")
	body.WriteString(armoredStr)
	if !strings.HasSuffix(armoredStr, "\r\n") {
		body.WriteString("\r\n")
	}

	// End boundary
	body.WriteString(fmt.Sprintf("\r\n--%s--\r\n", boundary))

	// --- 5) Build final raw email ---
	em := gomail.NewRawMessage()
	em.SetAddressHeader("From", fromAddr, fromName)
	em.SetHeader("To", recipientEmail)
	em.SetHeader("Subject", subject)
	em.SetHeader("Date", time.Now().UTC().Format(time.RFC1123Z))
	em.SetHeader("Content-Type", fmt.Sprintf("multipart/encrypted; protocol=\"application/pgp-encrypted\"; boundary=\"%s\"", boundary))

	// --- 6) Attach fully prebuilt multipart body ---
	em.SetRawBody(body.String())

	return em, nil
}

func randomChars(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz0123456789")
	b := make([]rune, n)
	for i := range b {
		index, err := cryptoRandInt(len(letterRunes))
		if err != nil {
			// Handle error, return empty string or fallback
			return ""
		}
		b[i] = letterRunes[index]
	}
	return string(b)
}

func cryptoRandInt(max int) (int, error) {
	nBig, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		return 0, err
	}
	return int(nBig.Int64()), nil
}

// emptyCharsetEncodedWordRe matches RFC 2047 encoded words with an empty
// charset field, e.g. =??q?text?= or =??B?text?=.
var emptyCharsetEncodedWordRe = regexp.MustCompile(`=\?\?([bqBQ]\?[^?]*)\?=`)

// fixEmptyCharsetEncodedWords replaces malformed RFC 2047 encoded words that
// have an empty charset (=??encoding?content?=) with UTF-8 as the charset.
// Some mailers omit the charset field while still encoding content as UTF-8.
// Without this fix the letters library (and other strict parsers) return
// "cannot lookup encoding" and refuse to process the message.
func fixEmptyCharsetEncodedWords(s string) string {
	return emptyCharsetEncodedWordRe.ReplaceAllString(s, "=?UTF-8?$1?=")
}

// PreprocessEmailData normalises raw email bytes so that the standard
// net/mail parser and the letters library can handle them reliably.
// It deliberately does NOT attempt to decode RFC 2047 encoded-words in
// address headers (From, To, Cc, …): Go's net/mail.ParseAddress already
// handles RFC 2047 display names natively, and running a second decode
// pass with a no-op charset reader produces invalid UTF-8 bytes that
// cause parsing to fail.  Use SafeDecodeAddressName to decode display
// names after parsing.
//
// It does fix malformed encoded words with an empty charset field
// (=??q?…?=) by substituting UTF-8, so that strict third-party parsers
// such as letters do not abort with "cannot lookup encoding".
func PreprocessEmailData(data []byte) ([]byte, error) {
	msg, err := mail.ReadMessage(bytes.NewReader(data))
	if err != nil {
		return data, nil // Return original data if it can't be parsed
	}

	var buf bytes.Buffer

	// Write headers, fixing any malformed encoded words in each value.
	for key := range msg.Header {
		for _, value := range msg.Header[key] {
			buf.WriteString(key)
			buf.WriteString(": ")
			buf.WriteString(fixEmptyCharsetEncodedWords(value))
			buf.WriteString("\r\n")
		}
	}

	// Blank line between headers and body
	buf.WriteString("\r\n")

	// Copy body
	_, err = io.Copy(&buf, msg.Body)
	if err != nil {
		return data, nil // Return original data on error
	}

	return buf.Bytes(), nil
}

// SafeDecodeAddressName decodes an RFC 2047 encoded display name returned
// by net/mail.ParseAddress.  If the decoded string is not valid UTF-8
// (e.g. the encoded word used a non-UTF-8 charset and the bytes were
// returned verbatim), it falls back to stripping encoded-word tokens and
// returning only the plain-text fragments, producing a best-effort
// human-readable name rather than a garbled or empty string.
func SafeDecodeAddressName(name string) string {
	if !strings.Contains(name, "=?") {
		return name
	}

	decoder := mime.WordDecoder{}
	decoded, err := decoder.DecodeHeader(name)
	if err == nil && utf8.ValidString(decoded) {
		return decoded
	}

	// Fallback: strip encoded-word tokens (including those with an empty
	// charset, e.g. =??q?…?=) and keep only the plain-text fragments.
	re := regexp.MustCompile(`=\?[^?]*\?[bqBQ]\?[^?]*\?=`)
	plain := re.ReplaceAllString(name, "")
	return strings.TrimSpace(plain)
}

// DecodeHeaderWithCharset decodes an RFC 2047 encoded header value using a
// charset-aware decoder backed by golang.org/x/net/html/charset. Unlike the
// default mime.WordDecoder, this handles legacy charsets such as windows-1251,
// koi8-r, and iso-8859-* by transcoding them to UTF-8. If decoding fails the
// original string is returned as-is so parsing can continue.
func DecodeHeaderWithCharset(s string) string {
	decoder := mime.WordDecoder{
		CharsetReader: func(charset string, input io.Reader) (io.Reader, error) {
			return xcharset.NewReaderLabel(charset, input)
		},
	}
	decoded, err := decoder.DecodeHeader(s)
	if err != nil {
		return s
	}
	return decoded
}

// NormalizeAddressSeparators replaces semicolon address separators with commas
// so that Go's mail.ParseAddressList can handle Outlook-style To/Cc headers
// such as "a@b.com; c@d.com" or a trailing semicolon "a@b.com;". The returned
// string has any leading/trailing commas and whitespace trimmed.
func NormalizeAddressSeparators(s string) string {
	s = strings.ReplaceAll(s, ";", ",")
	return strings.Trim(s, ", \t")
}

// CleanupMalformedEncodedAddress attempts to extract a valid email address
// from a malformed RFC 2047 encoded string
func CleanupMalformedEncodedAddress(addr string) string {
	// Look for email address in angle brackets
	if idx := strings.Index(addr, "<"); idx != -1 {
		if endIdx := strings.Index(addr[idx:], ">"); endIdx != -1 {
			email := addr[idx : idx+endIdx+1]
			// Try to decode any remaining encoded-words before the email
			before := addr[:idx]

			// Check if there's an encoded-word
			if strings.Contains(before, "=?") {
				// Remove the malformed encoded-word entirely
				// and just use the email address
				return strings.TrimSpace(email[1 : len(email)-1])
			}

			return strings.TrimSpace(before) + " " + email
		}
	}

	// If no angle brackets, return as-is
	return addr
}
