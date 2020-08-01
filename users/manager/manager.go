// Package manager is responsible for managing users.
// Handling operations like creation, listing and deletion safely.
package manager

import "github.com/katcipis/stonks/users"

// UsersStore is responsible for storing and retrieving user information
type UsersStore interface {
	// Adds a new user on storage returning its ID in the case of success
	// or a non-nil error in the case of failure.
	// The following errors MUST be returned (possibly wrapped)
	// giving specific conditions:
	//
	// - If any the user already exists: users.UserAlreadyExistsErr
	//
	// All other errors are to be considered internal errors.
	AddUser(email users.Email, fullname string, hashedPassword string) (string, error)
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
func (m *Manager) CreateUser(email users.Email, fullname string, password string) (string, error) {
	hashed, _ := m.auth.PasswordHash(password)
	// TODO handle password hash generation failures
	return m.store.AddUser(email, fullname, hashed)
}
