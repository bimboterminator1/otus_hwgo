package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		name        string
		in          interface{}
		expectedErr error
	}{
		{
			name: "Valid User",
			in: User{
				ID:     "123456789012345678901234567890123456",
				Name:   "John Doe",
				Age:    25,
				Email:  "johndoe@example.com",
				Role:   "admin",
				Phones: []string{"12345678901"},
			},
			expectedErr: nil,
		},
		{
			name: "Invalid User - ID length mismatch",
			in: User{
				ID:     "123",
				Name:   "John Doe",
				Age:    25,
				Email:  "johndoe@example.com",
				Role:   "admin",
				Phones: []string{"12345678901"},
			},
			expectedErr: ValidationErrors{
				{Field: "ID", Err: ErrStringLengthMismatch},
			},
		},
		{
			name: "Invalid User - Age below minimum",
			in: User{
				ID:     "123456789012345678901234567890123456",
				Name:   "John Doe",
				Age:    17,
				Email:  "johndoe@example.com",
				Role:   "admin",
				Phones: []string{"12345678901"},
			},
			expectedErr: ValidationErrors{
				{Field: "Age", Err: ErrValueIsLessThanMinValue},
			},
		},
		{
			name: "Invalid User - Email regex mismatch",
			in: User{
				ID:     "123456789012345678901234567890123456",
				Name:   "John Doe",
				Age:    25,
				Email:  "invalid-email",
				Role:   "admin",
				Phones: []string{"12345678901"},
			},
			expectedErr: ValidationErrors{
				{Field: "Email", Err: ErrRegexpMatchFailed},
			},
		},
		{
			name: "Invalid User - Role not in list",
			in: User{
				ID:     "123456789012345678901234567890123456",
				Name:   "John Doe",
				Age:    25,
				Email:  "johndoe@example.com",
				Role:   "guest",
				Phones: []string{"12345678901"},
			},
			expectedErr: ValidationErrors{
				{Field: "Role", Err: ErrValueNotInList},
			},
		},
		{
			name: "Invalid User - Phone length mismatch",
			in: User{
				ID:     "123456789012345678901234567890123456",
				Name:   "John Doe",
				Age:    25,
				Email:  "johndoe@example.com",
				Role:   "admin",
				Phones: []string{"123"},
			},
			expectedErr: ValidationErrors{
				{Field: "Phones", Err: ErrStringLengthMismatch},
			},
		},
		{
			name: "Valid App",
			in: App{
				Version: "1.2.3",
			},
			expectedErr: nil,
		},
		{
			name: "Invalid App - Version length mismatch",
			in: App{
				Version: "1.2",
			},
			expectedErr: ValidationErrors{
				{Field: "Version", Err: ErrStringLengthMismatch},
			},
		},
		{
			name: "Valid Token",
			in: Token{
				Header:    []byte("header"),
				Payload:   []byte("payload"),
				Signature: []byte("signature"),
			},
			expectedErr: nil,
		},
		{
			name: "Valid Response",
			in: Response{
				Code: 200,
				Body: "OK",
			},
			expectedErr: nil,
		},
		{
			name: "Invalid Response - Code not in list",
			in: Response{
				Code: 400,
				Body: "Bad Request",
			},
			expectedErr: ValidationErrors{
				{Field: "Code", Err: ErrValueNotInList},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			err := Validate(tt.in)
			require.ErrorIs(t, err, tt.expectedErr)
			_ = tt
		})
	}
}
