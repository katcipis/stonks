package users

import (
	"fmt"
	"net/mail"
)

// Email represents an email address on the form: "name@domain"
type Email string

// ParseEmail will parse the email address, returning the
// valid email address with trimmed spaces on success or a non-nil
// error in case of failure.
func ParseEmail(email string) (Email, error) {
	// I usually would prefer to guarantee invariants on the
	// type constructor, but in Go you can't enforce types
	// to be built with a constructor so it is always possible
	// to create a variable with zero values that are invalid.
	// So for these scenarios I end up using validation methods.
	p, err := mail.ParseAddress(string(email))
	if err != nil {
		return "", err
	}
	if p.Name != "" {
		return "", fmt.Errorf("email should be on the form user@domain not %q", email)
	}
	return Email(p.Address), nil
}
