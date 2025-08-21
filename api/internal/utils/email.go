package utils

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"io"
	"math/big"
	"net/mail"
	"regexp"
	"strings"

	"github.com/ProtonMail/gopenpgp/v3/crypto"
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

func EncryptWithPGPMIME(data []byte, fromAddr, fromName, subject, recipientEmail, recipientKey string) (*gomail.Message, error) {
	// --- 1) Serialize the original cleartext email ---
	raw, err := mail.ReadMessage(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("read original email: %w", err)
	}

	// --- 2) Parse recipient public key ---
	publicKey, err := crypto.NewKeyFromArmored(recipientKey)
	if err != nil {
		return nil, fmt.Errorf("parse public key: %w", err)
	}

	// --- 3) Read body and encrypt serialized message ---
	bodyBytes, err := io.ReadAll(raw.Body)
	if err != nil {
		return nil, fmt.Errorf("read message body: %w", err)
	}

	pgp := crypto.PGP()
	encHandle, err := pgp.Encryption().Recipient(publicKey).New()
	if err != nil {
		return nil, fmt.Errorf("create encryption handle: %w", err)
	}

	pgpMessage, err := encHandle.Encrypt(bodyBytes)
	if err != nil {
		return nil, fmt.Errorf("encrypt payload: %w", err)
	}

	armored, err := pgpMessage.Armor()
	if err != nil {
		return nil, fmt.Errorf("armor ciphertext: %w", err)
	}

	// --- 4) Build PGP/MIME multipart body ---
	boundary := "pgp-boundary-" + randomChars(16) // Generate a random boundary string
	var body bytes.Buffer

	// Part 1: version info
	body.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	body.WriteString("Content-Type: application/pgp-encrypted\r\n\r\n")
	body.WriteString("Version: 1\r\n\r\n")

	// Part 2: ciphertext
	body.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	body.WriteString("Content-Type: application/octet-stream; name=\"encrypted.asc\"\r\n")
	body.WriteString("Content-Disposition: inline; filename=\"encrypted.asc\"\r\n\r\n")
	body.WriteString(armored)
	if !strings.HasSuffix(armored, "\n") {
		body.WriteString("\r\n")
	}
	body.WriteString(fmt.Sprintf("\r\n--%s--\r\n", boundary))

	// --- 5) Build final gomail.Message ---
	em := gomail.NewRawMessage()
	em.SetAddressHeader("From", fromAddr, fromName)
	em.SetHeader("To", recipientEmail)
	em.SetHeader("Subject", subject)
	em.SetHeader("MIME-Version", "1.0")
	em.SetHeader("Content-Type", fmt.Sprintf("multipart/encrypted; protocol=\"application/pgp-encrypted\" boundary=\"%s\"", boundary))

	// --- 6) Attach our prebuilt multipart/encrypted body ---
	em.SetRawBody("text/plain", body.String())

	// --- 7) Print the final email message as raw string
	buf := new(bytes.Buffer)
	_, err = em.WriteTo(buf)
	if err != nil {
		return nil, fmt.Errorf("write email to buffer: %w", err)
	}
	fmt.Println("PGP/MIME email message:")
	fmt.Println(buf.String())

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
