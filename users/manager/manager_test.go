package manager_test

import (
	"errors"
	"strconv"
	"testing"

	"github.com/katcipis/stonks/auth"
	"github.com/katcipis/stonks/users"
	"github.com/katcipis/stonks/users/manager"
)

func TestUserCreation(t *testing.T) {

	type Test struct {
		name         string
		userEmail    users.Email
		userName     string
		userPassword string
		wantErr      error
	}

	tests := []Test{
		{
			name:         "Success",
			userEmail:    "hi@notadmin.com",
			userName:     "User Name",
			userPassword: "Some Password",
		},
		{
			name:         "SuccessForAdminUser",
			userEmail:    "user@test.com",
			userName:     "Admin User Has Domain test.com",
			userPassword: "admin password",
		},
		{
			name:         "FailureOnEmptyUserName",
			userName:     "",
			userEmail:    "user@test.com",
			userPassword: "admin password",
			wantErr:      users.InvalidUserParamErr,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			storage := newUsersStorage()
			authorizer := auth.New()
			usersManager := manager.New(authorizer, storage)
			userID, err := usersManager.CreateUser(test.userEmail, test.userName, test.userPassword)

			if test.wantErr != nil {
				if !errors.Is(err, test.wantErr) {
					t.Fatalf("got err [%v] but want err[%v]", err, test.wantErr)
				}
				return
			}

			gotUser, ok := storage.userByID(userID)
			if !ok {
				t.Fatalf("unable to find user id %q on storage", userID)
			}

			if gotUser.fullname != test.userName {
				t.Errorf("got name %q want %q", gotUser.fullname, test.userName)
			}

			if gotUser.email != test.userEmail {
				t.Errorf("got email %q want %q", gotUser.email, test.userEmail)
			}

			if !authorizer.HashMatchesPassword(gotUser.hashedPassword, test.userPassword) {
				t.Errorf("got hashed password %q that doesn't match password %q", gotUser.hashedPassword, test.userPassword)
			}
		})
	}
}

// UsersStorage is a simple in memory user storage implementation used in tests
type UsersStorage struct {
	idCount int
	users   map[string]User
}

// User is a user representation specific for test purposes
type User struct {
	id             string
	fullname       string
	hashedPassword string
	email          users.Email
}

func newUsersStorage() *UsersStorage {
	// WHY: after having a lot of brittle tests and problems with
	// mocks I went on the direction of using fake implementations,
	// and sometimes valid ones but in memory useful just for local
	// testing and experimentation. So far I have been happier in the
	// sense that the tests are easier to understand and to debug.
	//
	// In one of his posts Kent Beck also talks a little about
	// coupling tests too much on code structure (mocks can be a source of that):
	//
	// - https://medium.com/@kentbeck_7670/programmer-test-principles-d01c064d7934
	return &UsersStorage{
		users: map[string]User{},
	}
}

func (s *UsersStorage) AddUser(email users.Email, fullname string, pass string) (string, error) {
	s.idCount++
	id := strconv.Itoa(s.idCount)
	s.users[id] = User{
		id:             id,
		fullname:       fullname,
		hashedPassword: pass,
		email:          email,
	}
	return id, nil
}

func (s *UsersStorage) userByID(id string) (User, bool) {
	v, ok := s.users[id]
	return v, ok
}

func hashedPassword(t *testing.T, a *auth.Authorizer, pass string) string {
	t.Helper()

	hashed, err := a.PasswordHash(pass)
	if err != nil {
		t.Fatal(err)
	}
	return hashed
}
