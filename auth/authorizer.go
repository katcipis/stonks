package auth

import "golang.org/x/crypto/bcrypt"

// Authorizer is responsible for authorization and security related operations
type Authorizer struct {
}

// New creates a new Authorizer
func New() *Authorizer {
	return &Authorizer{}
}

// PasswordHash is responsible for creating safe hashes from
// passwords, suitable for storage and comparison later
// using IsPasswordMatch.
func (*Authorizer) PasswordHash(pass string) (string, error) {
	v, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	return string(v), err
}

// HashMatchesPassword checks if the salted hash matches the
// given plain text password, returning true if there is a match,
// false otherwise.
func (*Authorizer) HashMatchesPassword(saltedhash string, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(saltedhash), []byte(password)) == nil
}
