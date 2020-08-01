package users_test

import (
	"testing"

	"github.com/katcipis/stonks/users"
)

func TestEmailValidation(t *testing.T) {
	type Test struct {
		name      string
		email     string
		wantEmail users.Email
		wantErr   bool
	}

	tests := []Test{
		{
			name:      "ValidEmail",
			email:     "valid@valid.com",
			wantEmail: "valid@valid.com",
		},
		{
			name:      "TrimmesSpacesFromValidEmail",
			email:     "    valid@valid.com  ",
			wantEmail: "valid@valid.com",
		},
		{
			name:    "ValidMailWithNameFails",
			email:   "Test <valid@valid.com>",
			wantErr: true,
		},
		{
			name:    "InvalidMailWithNameFails",
			email:   "Test <invalid>",
			wantErr: true,
		},
		{
			name:    "InvalidMailShouldFail",
			email:   "invalid",
			wantErr: true,
		},
		{
			name:    "InvalidMailWithoutDomainFails",
			email:   "invalid@",
			wantErr: true,
		},
		{
			name:    "InvalidEmptyStringFails",
			email:   "",
			wantErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			email, err := users.ParseEmail(test.email)
			if err != nil {
				if !test.wantErr {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}
			if test.wantErr {
				t.Fatal("expected error, got nil")
			}

			if email != test.wantEmail {
				t.Fatalf("got email %q wanted %q", email, test.wantEmail)
			}
		})
	}
}
