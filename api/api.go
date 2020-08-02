// Package api is responsible for exporting services
// as an HTTP API.
package api

import (
	"net/http"
	"time"

	"github.com/katcipis/stonks/users/manager"
)

// CreateUserRequestBody is the request body required to create users
type CreateUserRequestBody struct {
	FullName string `json:"fullname"`
	Email    string `json:"email"`
	Password string `json:"password`
}

// Error response represents the response body
// of all requests that failed.
type ErrorResponse struct {
	Error struct {
		Message string `json:"message"`
	}
}

// Config has all configuration needed by the api, like timeout
// configurations.
type Config struct {
	CreateUserTimeout time.Duration
}

// New creates a new HTTP handler with all the service routes.
func New(usersManager *manager.Manager, cfg Config) http.Handler {
	return nil
}
