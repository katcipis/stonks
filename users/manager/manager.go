// Package manager is responsible for managing users.
// Handling operations like creation, listing and deletion safely.
package manager

import (
	"context"
	"fmt"

	"github.com/katcipis/stonks/users"
)

// UsersStore is responsible for storing and retrieving user information
type UsersStore interface {
	// Adds a new user on storage returning its ID in the case of success
	// or a non-nil error in the case of failure.
	// The following errors MUST be returned (possibly wrapped)
	// giving specific conditions:
	//
	// - If the user already exists: users.UserAlreadyExistsErr
	//
	// All other errors are to be considered internal errors.
	AddUser(ctx context.Context, email users.Email, fullname string, hashedPassword string) (string, error)
}

// Authorizer is responsible for authorization and security related operations
type Authorizer interface {

	// PasswordHash creates safe hashes from
	// passwords, suitable for storage and comparison later
	// using IsPasswordMatch.
	PasswordHash(pass string) (string, error)
}

// Manager is responsible for managing users, doing
// operations like creation, listing and deletion safely.
// It does that by the composition of interfaces providing
// storage and authorization.
type Manager struct {
	auth  Authorizer
	store UsersStore
}

// New creates a new users manager
func New(a Authorizer, s UsersStore) *Manager {
	return &Manager{
		auth:  a,
		store: s,
	}
}

// Creates a new user, returning its ID in the case of success
// or a non-nil error in the case of failure.
// The following errors can be expected to be wrapped in the returned
// error giving specific conditions:
//
// - If any of the parameters is invalid: users.InvalidUserParamErr
// - If any the user already exists: users.UserAlreadyExistsErr
//
// All other errors are to be considered internal errors.
func (m *Manager) CreateUser(ctx context.Context, email string, fullname string, password string) (string, error) {
	if fullname == "" {
		return "", fmt.Errorf("%w:empty name", users.InvalidUserParamErr)
	}
	if password == "" {
		return "", fmt.Errorf("%w:empty password", users.InvalidUserParamErr)
	}
	validEmail, err := users.ParseEmail(email)
	if err != nil {
		return "", fmt.Errorf("%w:invalid email:%v", users.InvalidUserParamErr, err)
	}

	hashed, err := m.auth.PasswordHash(password)
	if err != nil {
		return "", fmt.Errorf("error creating password hash:%v", err)
	}
	return m.store.AddUser(ctx, validEmail, fullname, hashed)
}
