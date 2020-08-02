// Package storage is responsible for storing and retrieving users
package storage

import (
	"context"

	"github.com/katcipis/stonks/users"
)

// Storage is responsible for storing and retrieving users
type Storage struct {
}

// New creates a new Storage given the database address
// and the password to securely connect to it.
func New(addr string, password string) *Storage {
	return nil
}

// AddUser adds a user with the given parameters, returning its ID in the case
// of success or an error otherwise.
// If an user with the given email already exists it returns users.UserAlreadyExistsErr
func (s *Storage) AddUser(
	ctx context.Context,
	email users.Email,
	fullname string,
	hashedPassword string,
) (string, error) {
	return "TODO:FORNOW:-)", nil
}
