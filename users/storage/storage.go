// Package storage is responsible for storing and retrieving users
package storage

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/katcipis/stonks/users"
)

// Storage is responsible for storing and retrieving users
type Storage struct {
	connPool *pgxpool.Pool
}

// New creates a new Storage given the database address,user
// and the password to securely connect to it.
func New(host string, database string, user string, password string) (*Storage, error) {
	const timeout = 10 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	pgconfig := fmt.Sprintf("postgres://%s:%s@%s:5432/%s", user, password, host, database)
	pool, err := pgxpool.Connect(ctx, pgconfig)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}
	return &Storage{connPool: pool}, nil
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
	sqlStatement := `INSERT INTO users.users (email, fullname, password_hash) VALUES ($1, $2, $3) RETURNING id`
	rows, err := s.connPool.Query(ctx, sqlStatement, email, fullname, hashedPassword)
	if err != nil {
		return "", fmt.Errorf("error inserting new user:%v", err)
	}
	defer rows.Close()

	var userID int64
	rows.Next()
	err = rows.Err()
	if err != nil {
		pgerr, ok := err.(*pgconn.PgError)
		err := fmt.Errorf("error scanning new user result:%v", err)
		if !ok {
			return "", err
		}
		// From: https://www.postgresql.org/docs/11/errcodes-appendix.html
		const uniqueViolationErrorCode = "23505"
		if pgerr.Code == uniqueViolationErrorCode {
			return "", fmt.Errorf("%w:%s", users.UserAlreadyExistsErr, email)
		}
		return "", err
	}
	err = rows.Scan(&userID)
	return strconv.FormatInt(userID, 10), nil
}
