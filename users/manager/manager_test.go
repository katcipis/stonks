package manager_test

import (
	"errors"
	"testing"

	"github.com/katcipis/stonks/users/manager"
)

func TestUserCreation(t *testing.T) {

	type Test struct {
		name         string
		userEmail    manager.Email
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
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			storage := newUsersStorage()
			usersManager := manager.New(storage)
			userID, err := usersManager.CreateUser(test.userEmail, test.userName, test.userPassword)
			if !errors.Is(test.wantErr, err) {
				t.Fatalf("got err [%v] but want err[%v]", err, test.wantErr)
			}

			gotUser, ok := storage.userByID(userID)
			if !ok {
				t.Fatalf("unable to find user id %q on storage", userID)
			}

			wantUser := User{
				id:       userID,
				name:     test.userName,
				password: test.userPassword,
				email:    test.userEmail,
			}

			if gotUser != wantUser {
				t.Fatalf("got user[%+v] != wanted [%+v]", gotUser, wantUser)
			}
		})
	}
}

// UsersStorage is a simple in memory user storage implementation used in tests
type UsersStorage struct {
}

// User is a user representation specific for test purposes
type User struct {
	id       string
	name     string
	password string
	email    manager.Email
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
	return &UsersStorage{}
}

func (s *UsersStorage) userByID(id string) (User, bool) {
	return User{}, false
}
