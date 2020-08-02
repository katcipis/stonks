// Package api is responsible for exporting services
// as an HTTP API.
package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/katcipis/stonks/users"
	"github.com/katcipis/stonks/users/manager"
)

// CreateUserRequestBody is the request body required to create users
type CreateUserRequestBody struct {
	FullName string `json:"fullname"`
	Email    string `json:"email"`
	Password string `json:"password`
}

type CreateUserResponse struct {
	ID string `json:"id"`
}

// Error contains error information used in error responses
type Error struct {
	Message string `json:"message"`
}

// Error response represents the response body
// of all requests that failed.
type ErrorResponse struct {
	Error Error `json:"error"`
}

// Config has all configuration needed by the api, like timeout
// configurations.
type Config struct {
	CreateUserTimeout time.Duration
}

// New creates a new HTTP handler with all the service routes.
func New(usersManager *manager.Manager, cfg Config) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/users", func(res http.ResponseWriter, req *http.Request) {
		// TODO: test unsupported methods
		dec := json.NewDecoder(req.Body)
		parsedReq := CreateUserRequestBody{}
		// TODO: test invalid JSON on request
		dec.Decode(&parsedReq)

		ctx, cancel := context.WithTimeout(context.Background(), cfg.CreateUserTimeout)
		defer cancel()

		// TODO: test created user ID
		userID, err := usersManager.CreateUser(ctx, parsedReq.Email, parsedReq.FullName, parsedReq.Password)
		if err != nil {
			if errors.Is(err, users.InvalidUserParamErr) {
				res.WriteHeader(http.StatusBadRequest)
				// Invalid user param errors are guaranteed
				// to be safe to send to users. If a service is
				// external care must be taken to not leak details
				// that can be a potential security threat.
				// When that is not the case I like the idea of
				// informative error responses as detailed here:
				//
				// - https://commandcenter.blogspot.com/2017/12/error-handling-in-upspin.html
				//
				// I'm specially fond to the idea of a cross service
				// operational trace (instead of stack traces).
				// But I never tried it yet.
				res.Write(errorResponse(err.Error()))
				return
			}
			// TODO: log err
			// Specially when you can't give much detail on errors for
			// security reasons it would be a good idea to have
			// a tracking id for errors to help map the error to
			// the logs, not sure if I'm going to have time to add this.
			res.WriteHeader(http.StatusInternalServerError)
			res.Write(errorResponse("internal error trying to create user"))
			return
		}

		res.WriteHeader(http.StatusCreated)
		res.Write(jsonResponse(CreateUserResponse{ID: userID}))
	})
	return mux
}

func errorResponse(message string) []byte {
	return jsonResponse(ErrorResponse{
		Error: Error{Message: message},
	})
}

func jsonResponse(v interface{}) []byte {
	// TODO: handle and log err
	res, _ := json.Marshal(v)
	return res
}
