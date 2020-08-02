// Package api is responsible for exporting services
// as an HTTP API.
package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/katcipis/stonks/users"
	"github.com/katcipis/stonks/users/manager"
)

// CreateUserRequestBody is the request body required to create users
type CreateUserRequestBody struct {
	FullName string `json:"fullname"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// CreateUserResponse is the response body when a user is created with success
type CreateUserResponse struct {
	ID string `json:"id"`
}

// Error contains error information used in error responses
type Error struct {
	Message string `json:"message"`
}

// ErrorResponse represents the response body
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
	const usersPath = "/v1/users"

	mux := http.NewServeMux()
	userslog := log.WithFields(log.Fields{"path": usersPath})

	mux.HandleFunc(usersPath, func(res http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			res.WriteHeader(http.StatusMethodNotAllowed)
			msg := fmt.Sprintf("method %q is not allowed", req.Method)
			logResponseBodyWrite(userslog, res, errorResponse(msg))
			userslog.WithFields(log.Fields{"error": msg}).Warning("method not allowed")
			return
		}
		dec := json.NewDecoder(req.Body)
		parsedReq := CreateUserRequestBody{}

		err := dec.Decode(&parsedReq)
		if err != nil {
			msg := fmt.Sprintf("error parsing JSON request body: %v", err)
			res.WriteHeader(http.StatusBadRequest)
			logResponseBodyWrite(userslog, res, errorResponse(msg))
			userslog.WithFields(log.Fields{"error": msg}).Warning("invalid request body")
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), cfg.CreateUserTimeout)
		defer cancel()

		userID, err := usersManager.CreateUser(ctx, parsedReq.Email, parsedReq.FullName, parsedReq.Password)
		if err != nil {
			if errors.Is(err, users.InvalidUserParamErr) ||
				errors.Is(err, users.UserAlreadyExistsErr) {
				res.WriteHeader(http.StatusBadRequest)
				// Invalid and duplicated user param errors are guaranteed
				// to be safe to send to users (not much info added on the error context).
				// If a service is external care must be taken to not leak details
				// that can be a potential security threat.
				// When that is not the case I like the idea of
				// informative error responses as detailed here:
				//
				// - https://commandcenter.blogspot.com/2017/12/error-handling-in-upspin.html
				//
				// I'm specially fond to the idea of a cross service
				// operational trace (instead of stack traces).
				// But I never tried it yet.
				logResponseBodyWrite(userslog, res, errorResponse(err.Error()))
				userslog.WithFields(log.Fields{"error": err.Error()}).Warning("bad request error")
				return
			}
			// Specially when you can't give much detail on errors for
			// security reasons it would be a good idea to have
			// a tracking id for errors to help map the error to
			// the logs, not sure if I'm going to have time to add this.
			res.WriteHeader(http.StatusInternalServerError)
			logResponseBodyWrite(userslog, res, errorResponse("internal server error"))
			userslog.WithFields(log.Fields{"error": err.Error()}).Error("internal server error")
			return
		}

		res.WriteHeader(http.StatusCreated)
		logResponseBodyWrite(userslog, res, jsonResponse(CreateUserResponse{ID: userID}))
	})
	return mux
}

func logResponseBodyWrite(logger *log.Entry, w io.Writer, data []byte) {
	_, err := w.Write(data)
	if err != nil {
		logger.WithFields(log.Fields{"error": err}).Warning("writing response body")
	}
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
