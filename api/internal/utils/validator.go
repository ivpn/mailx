package utils

import (
	"log"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	*validator.Validate
}

func NewValidator() Validator {
	v := Validator{validator.New()}

	err := v.RegisterValidation("password", passwordValidation)
	if err != nil {
		log.Println("error registering password validation:", err)
	}

	err = v.RegisterValidation("pgp", pgpKeyValidation)
	if err != nil {
		log.Println("error registering pgp key validation:", err)
	}

	err = v.RegisterValidation("emailx", sqlEmailValidation)
	if err != nil {
		log.Println("error registering sql email validation:", err)
	}

	err = v.RegisterValidation("search", searchValidation)
	if err != nil {
		log.Println("error registering search validation:", err)
	}

	err = v.RegisterValidation("emaillocalpart", emailLocalPartValidation)
	if err != nil {
		log.Println("error registering email local part validation:", err)
	}

	return v
}

func ValidateEmail(email string) error {
	validator := NewValidator()
	return validator.Var(email, "required,emailx")
}

func passwordValidation(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	if len(password) < 12 || len(password) > 64 {
		return false
	}

	var uppercase = regexp.MustCompile(`[A-Z]`).MatchString
	var lowercase = regexp.MustCompile(`[a-z]`).MatchString
	var number = regexp.MustCompile(`[0-9]`).MatchString
	var specialChar = regexp.MustCompile(`[-_+=~!@#$%^&*(),;.?":{}|<>]`).MatchString

	if !uppercase(password) {
		return false
	}

	if !lowercase(password) {
		return false
	}

	if !number(password) {
		return false
	}

	if !specialChar(password) {
		return false
	}

	return true
}

func sqlEmailValidation(fl validator.FieldLevel) bool {
	email := fl.Field().String()

	// “omitempty” double check
	if email == "" {
		return true
	}

	// Check if the email is valid
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

func pgpKeyValidation(fl validator.FieldLevel) bool {
	key := fl.Field().String()

	// “omitempty” double check
	if key == "" {
		return true
	}

	// Ignore hash
	if len(key) == 64 {
		return true
	}

	// Check that the key starts with a valid PGP header
	return strings.HasPrefix(key, "-----BEGIN PGP PUBLIC KEY BLOCK-----") && strings.HasSuffix(key, "-----END PGP PUBLIC KEY BLOCK-----")
}

func searchValidation(fl validator.FieldLevel) bool {
	value := fl.Field().String()

	// “omitempty” double check
	if value == "" {
		return true
	}

	re := regexp.MustCompile(`^[-a-zA-Z0-9 ._+@]+$`)
	return re.MatchString(value)
}

func emailLocalPartValidation(fl validator.FieldLevel) bool {
	localPart := fl.Field().String()

	if len(localPart) < 1 || len(localPart) > 64 {
		return false
	}

	// RFC 5321/5322 allow a broad set of characters in the local part.
	// This covers the common, practical subset: letters, digits, and
	// the typical unquoted special characters.
	var validChars = regexp.MustCompile(`^[a-zA-Z0-9.!#$%&'*+/=?^_` + "`" + `{|}~-]+$`).MatchString

	if !validChars(localPart) {
		return false
	}

	// Local part must not start or end with a dot
	if strings.HasPrefix(localPart, ".") || strings.HasSuffix(localPart, ".") {
		return false
	}

	// No consecutive dots allowed
	if strings.Contains(localPart, "..") {
		return false
	}

	return true
}
