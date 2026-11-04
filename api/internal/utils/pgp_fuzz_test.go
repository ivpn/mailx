package utils

import (
	"strings"
	"testing"

	"github.com/ProtonMail/gopenpgp/v3/crypto"
)

// FuzzEncryptWithPGPInline fuzzes the user-supplied recipient public key:
// users paste an armored key in settings, and every forwarded message is
// encrypted to it via crypto.NewKeyFromArmored and Encrypt.
func FuzzEncryptWithPGPInline(f *testing.F) {
	key, err := crypto.PGP().KeyGeneration().AddUserId("fuzz", "fuzz@example.com").New().GenerateKey()
	if err != nil {
		f.Fatal(err)
	}
	pub, err := key.GetArmoredPublicKey()
	if err != nil {
		f.Fatal(err)
	}

	seeds := []string{
		pub,
		pub[:len(pub)/2],
		strings.Replace(pub, "PUBLIC KEY BLOCK", "PRIVATE KEY BLOCK", 2),
		"-----BEGIN PGP PUBLIC KEY BLOCK-----\n\n-----END PGP PUBLIC KEY BLOCK-----",
		"-----BEGIN PGP PUBLIC KEY BLOCK-----\n\nAAAA\n-----END PGP PUBLIC KEY BLOCK-----",
		"",
		"not a key",
	}
	for _, s := range seeds {
		f.Add(s)
	}

	f.Fuzz(func(t *testing.T, recipientKey string) {
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("EncryptWithPGPInline panicked on key (len=%d): %v", len(recipientKey), r)
			}
		}()
		_, _ = EncryptWithPGPInline("Hello", recipientKey)
	})
}
