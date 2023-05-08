package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"reflect"
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
	Custom struct {
		Code       int    `validate:"min:200|max:500"`
		CodeIn     int    `validate:"in:200,404,500"`
		BodyStrLen string `validate:"len:5"`
		BodyStrIn  string `validate:"in:pupa,lupa"`

		Email string `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr []error
		desc        string
	}{
		{
			in: User{
				ID:     "101",
				Name:   "HW zavupa",
				Age:    20,
				Email:  "pupa@lupa.ru",
				Role:   "admin",
				Phones: []string{"00000000000", "11111111111"},
				meta:   nil,
			},
			expectedErr: nil,
			desc:        "positive",
		},
		{
			in: Token{
				Header:    []byte("nil"),
				Payload:   []byte("nil"),
				Signature: []byte("nil"),
			},
			expectedErr: nil,
			desc:        "positive",
		},
		{
			in: Response{
				Code: 200,
				Body: "{\"Result\" : true}",
			},
			expectedErr: nil,
			desc:        "positive",
		},
		{
			in: Custom{
				Code:       499,
				CodeIn:     200,
				BodyStrLen: "admin",
				BodyStrIn:  "pupa",
				Email:      "ro@super.ru",
			},
			expectedErr: nil,
			desc:        "positive",
		},
		{
			in: Custom{
				Code:       600,
				BodyStrLen: "admin",
				Email:      "ro@super.ru",
			},
			expectedErr: []error{ErrMax},
			desc:        "negative",
		},
	}

	for _, tt := range tests {

		switch tt.desc {
		case "positive":
			t.Run(fmt.Sprintf("%v", reflect.TypeOf(tt.in)), func(t *testing.T) {
				tt := tt
				err := Validate(tt.in)
				require.NoError(t, err)

			})
		case "negative":
			t.Run(fmt.Sprintf("%v", reflect.TypeOf(tt.in)), func(t *testing.T) {
				tt := tt
				err := Validate(tt.in)
				var valErr ValidationErrors

				require.ErrorAs(t, err, &valErr)

				for i, e := range tt.expectedErr {
					require.ErrorIs(t, valErr[i], e)

				}
			})
		}
	}
}
