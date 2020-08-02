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

func (s *Storage) AddUser(
	ctx context.Context,
	email users.Email,
	fullname string,
	hashedPassword string,
) (string, error) {
	return "", nil
}
