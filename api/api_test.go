package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/katcipis/stonks/api"
	"github.com/katcipis/stonks/auth"
	"github.com/katcipis/stonks/users/manager"
	"github.com/katcipis/stonks/users/storage"
)

// WHY: Usually I would do more testing on the isolated level and just validate
// on integration tests that everything fits together and works for common
// usage scenarios (error handling is easier tested isolated with injected failures).
// But I ended up doing more integration test given the time limitation, if I was
// in a situation that had a tight deadline to put in production I would feel
// safer having some isolated tests + more integration tests than just a lot
// of isolated ones.

func TestUserCreation(t *testing.T) {
	type Test struct {
		name           string
		requestBody    []byte
		wantStatusCode int
	}

	tests := []Test{
		{
			name: "Success",
			requestBody: toJSON(t, api.CreateUserRequestBody{
				FullName: "Success",
				Email:    "stonks@corp.com",
				Password: "weakpass",
			}),
			wantStatusCode: http.StatusCreated,
		},
		{
			name: "FailsIfNameIsMissing",
			requestBody: toJSON(t, api.CreateUserRequestBody{
				Email:    "stonks2@corp.com",
				Password: "weakpass",
			}),
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "FailsIfEmailIsMissing",
			requestBody: toJSON(t, api.CreateUserRequestBody{
				FullName: "Success",
				Password: "weakpass",
			}),
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name: "FailsIfPasswordIsMissing",
			requestBody: toJSON(t, api.CreateUserRequestBody{
				FullName: "Success",
				Email:    "stonks3@corp.com",
			}),
			wantStatusCode: http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			const userDBAddr = "todo"
			const userDBPass = "testing"

			usersStorage := storage.New(userDBAddr, userDBPass)
			authorizer := auth.New()
			usersManager := manager.New(authorizer, usersStorage)

			service := api.New(usersManager, api.Config{
				CreateUserTimeout: 10 * time.Second,
			})

			server := httptest.NewServer(service)
			defer server.Close()

			createUserURL := server.URL + "/v1/users"
			request := newRequest(t, http.MethodPost, createUserURL, test.requestBody)
			client := server.Client()

			res, err := client.Do(request)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()

			if res.StatusCode != test.wantStatusCode {
				t.Fatalf("got response %d want %d", res.StatusCode, test.wantStatusCode)
			}

			if test.wantStatusCode != http.StatusCreated {
				dec := json.NewDecoder(res.Body)
				wantErr := api.ErrorResponse{}
				err := dec.Decode(&wantErr)
				if err != nil {
					t.Fatal(err)
				}
				// Validate that a message is sent, but not its contents
				// since the message is for human inspection only and
				// should be handled opaquely by code.
				// If necessary we can introduce error codes (strings or ints),
				// but it does not seem necessary for now.
				// If we add some tracking ID for errors this would also
				// be the place to check for them.
				if wantErr.Error.Message == "" {
					t.Fatalf("expected an error message on status code %d", test.wantStatusCode)
				}
				return
			}

			// TODO: check on storage if the user has been created
		})
	}
}

func toJSON(t *testing.T, v interface{}) []byte {
	t.Helper()

	j, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	return j
}

func newRequest(t *testing.T, method string, url string, body []byte) *http.Request {
	t.Helper()

	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}

	return req
}
