package domain

import (
	"regexp"
	"strings"
)

var (
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]{3,20}$`)
	e164Regex     = regexp.MustCompile(`^\+[1-9][0-9]{7,14}$`)
	letterRegex   = regexp.MustCompile(`[A-Za-z]`)
	numberRegex   = regexp.MustCompile(`[0-9]`)
)

type FieldErrors map[string]string

func NormalizeEmail(email string) string       { return strings.ToLower(strings.TrimSpace(email)) }
func NormalizeUsername(username string) string { return strings.ToLower(strings.TrimSpace(username)) }
func NormalizePhone(phone string) string       { return strings.TrimSpace(phone) }
func ValidateUsername(username string) bool    { return usernameRegex.MatchString(username) }
func ValidatePhoneE164(phone string) bool      { return e164Regex.MatchString(phone) }
func ValidatePassword(password string) bool {
	if len(password) < 8 {
		return false
	}
	return letterRegex.MatchString(password) && numberRegex.MatchString(password)
}
