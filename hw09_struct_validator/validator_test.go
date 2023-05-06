package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"testing"
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
		Code    int    `validate:"min:300"`
		BodyStr string `validate:"len:36|regexp:^\\w+@\\w+\\.\\w+$"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
		desc        string
	}{
		{
			in: Custom{
				Code:    200,
				BodyStr: "{\"Result\" : true}",
			},
			expectedErr: nil,
			desc:        "Пустая структура",
		},
		// ...
		// Place your code here.
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt.desc), func(t *testing.T) {
			tt := tt
			t.Parallel()

			// Place your code here.

			err := Validate(tt.in)
			require.Equal(t, tt.expectedErr, err)

			_ = tt
		})
	}
}
