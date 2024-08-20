package hw09structvalidator

import (
	"encoding/json"
	"errors"
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
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: User{
				ID:     "7e756c80-a1b9-4be6-b2b2-0d4670fc7f55",
				Name:   "TaoGunner",
				Age:    18,
				Email:  "taogunner@gmail.com",
				Role:   "stuff",
				Phones: []string{},
				meta:   []byte{},
			},
			expectedErr: nil,
		},
		{
			in: User{
				ID:     "7e756c80-a1b9-4be6-b2b2-0d4670fc7f55",
				Name:   "TaoGunner",
				Age:    18,
				Email:  "taogunner@gmail.com",
				Role:   "stuff1",
				Phones: []string{},
				meta:   []byte{},
			},
			expectedErr: append(ValidationErrors{}, ValidationError{Field: "Role", Err: errValidationStringIn}),
		},
		{
			in: User{
				ID:    "7e756c80-a1b9-4be6-b2b2-0d4670fc7f15",
				meta:  []byte{},
				Email: "taogunner@gmail.com",
				Age:   60,
				Role:  "admin",
			},
			expectedErr: append(ValidationErrors{}, ValidationError{Field: "Age", Err: errValidationIntMax}),
		},
		{
			in:          App{Version: "v4.50"},
			expectedErr: nil,
		},
		{
			in:          App{Version: "4.50"},
			expectedErr: ValidationErrors{{Field: "Version", Err: errValidationStringLen}},
		},
		{
			in:          Response{Code: 500, Body: "text"},
			expectedErr: nil,
		},
		{
			in:          Response{Code: 0, Body: "text"},
			expectedErr: ValidationErrors{{Field: "Code", Err: errValidationIntIn}},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			errs := Validate(tt.in)
			if tt.expectedErr == nil {
				require.NoError(t, errs)
				return
			}
			if errors.Is(tt.expectedErr, ValidationErrors{}) {
				var validErr ValidationErrors
				// Сравнение типа ошибки
				require.ErrorAs(t, errs, &validErr)
			}
			// Сравнение текста ошибки
			require.Equal(t, errs.Error(), tt.expectedErr.Error())

			_ = tt
		})
	}
}
