package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"strings"
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

	TestCase struct {
		Name float64 `validate:"in:20.0,11.1,333"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			User{
				ID:     "123",
				Name:   "Mur",
				Age:    0,
				Email:  "wrongEmail",
				Role:   "cat",
				Phones: []string{"89065555555"},
			},
			ValidationErrors{
				ValidationError{
					Field: "ID",
					Err:   ErrorStrLengthInvalid,
				},
				ValidationError{
					"Age",
					ErrorMinValue,
				},
				ValidationError{
					"Email",
					ErrorIsNotMatch,
				},
				ValidationError{
					"Role",
					ErrorIsNotInSlice,
				},
			},
		}, {
			User{
				ID:     strings.Repeat("1", 36),
				Name:   "Mur",
				Age:    22,
				Email:  "correct@email.com",
				Role:   "stuff",
				Phones: []string{"89065555555"},
			},
			nil,
		}, {
			strings.Repeat("1", 36),
			ErrorIsNotStructType,
		}, {
			Token{
				Header:    nil,
				Payload:   nil,
				Signature: nil,
			},
			nil,
		}, {
			Response{
				Code: 200,
				Body: "{}",
			},
			nil,
		}, {
			TestCase{
				Name: 11.4,
			},
			ErrorUnknownFieldType,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()
			err := Validate(tt.in)
			require.Equal(t, tt.expectedErr, err)
			_ = tt
		})
	}
}
