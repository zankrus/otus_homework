package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

type ValidatorTestCases struct {
	in          interface{}
	expectedErr []error
	flag        string
	desc        string
}
type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
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
		CodeMin    int    `validate:"min:200"`
		CodeMax    int    `validate:"max:500"`
		CodeMix    int    `validate:"min:200|max:500"`
		CodeIn     int    `validate:"in:200,404,500"`
		BodyStrLen string `validate:"len:5"`
		BodyStrIn  string `validate:"in:pupa,lupa"`
		BodyStrMix string `validate:"in:pupa,lupa|len:4"`
		Email      string `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
	}
)

func TestValidate(t *testing.T) {
	tests := []ValidatorTestCases{
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
			flag:        "positive",
		},
		{
			in: Token{
				Header:    []byte("nil"),
				Payload:   []byte("nil"),
				Signature: []byte("nil"),
			},
			expectedErr: nil,
			flag:        "positive",
		},
		{
			in: Response{
				Code: 200,
				Body: "{\"Result\" : true}",
			},
			expectedErr: nil,
			flag:        "positive",
		},
		{
			in: Custom{
				CodeMin:    499,
				CodeIn:     200,
				CodeMix:    400,
				BodyStrLen: "admin",
				BodyStrIn:  "pupa",
				BodyStrMix: "lupa",
				Email:      "ro@super.ru",
			},
			expectedErr: nil,
			flag:        "positive",
		},
		{
			in: Custom{
				CodeMin:    199,
				CodeMax:    700,
				CodeIn:     201,
				CodeMix:    900,
				BodyStrLen: "admin1",
				BodyStrIn:  "pupa1",
				BodyStrMix: "lupa2",
				Email:      "rosuper.ru",
			},
			expectedErr: []error{ErrMin, ErrMax, ErrMax, ErrIn, ErrLen, ErrIn, ErrIn, ErrLen, ErrRegex},
			flag:        "negative_validation",
			desc:        "Все ошибки валидации и комбинации",
		},
		{
			in:          "Not a Struct",
			expectedErr: []error{ErrNotAStruct},
			flag:        "negative_software",
			desc:        "На входе не структура",
		},
		{
			in: struct {
				Code int `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
			}{
				Code: 1,
			},
			expectedErr: []error{ErrMissmatchTagAndType},
			flag:        "negative_software",
			desc:        "Тэг не подходит для типа",
		},
		{
			in: struct {
				Code int `validate:"regexp:"`
			}{
				Code: 1,
			},
			expectedErr: []error{ErrBrokenTag},
			flag:        "negative_software",
			desc:        "Сломанный тэг",
		},
	}

	for _, tt := range tests {
		tt := tt

		switch tt.flag {
		case "positive":
			t.Run(fmt.Sprintf("%v", reflect.TypeOf(tt.in)), func(t *testing.T) {
				tt := tt
				err := Validate(tt.in)
				require.NoError(t, err)
			})
		case "negative_validation":
			t.Run(fmt.Sprintf("%v", reflect.TypeOf(tt.in)), func(t *testing.T) {
				tt := tt
				err := Validate(tt.in)
				var valErr ValidationErrors

				require.ErrorAs(t, err, &valErr)

				for i, e := range tt.expectedErr {
					require.ErrorIs(t, valErr[i].Err, e)
				}
			})
		case "negative_software":
			t.Run(fmt.Sprintf("%v", reflect.TypeOf(tt.in)), func(t *testing.T) {
				tt := tt
				err := Validate(tt.in)
				require.ErrorContains(t, err, tt.expectedErr[0].Error())
			})
		}
	}
}
