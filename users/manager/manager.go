// Package manager is responsible for managing users.
// Handling operations like creation, listing and deletion safely.
package manager

import "github.com/katcipis/stonks/users"

// UsersStorage is responsible for storing and retrieving user information
type UsersStorage interface {
}

// Manager is responsible for handling users related
// operations like creation, listing and deletion safely.
// It does that by the composition of interfaces providing
// storage and authorization.
type Manager struct {
}

// New creates a new users manager
func New(s UsersStorage) *Manager {
	return nil
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
	return "", nil
}
