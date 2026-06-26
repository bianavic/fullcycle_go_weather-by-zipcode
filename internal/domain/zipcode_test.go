package domain_test

import (
	"errors"
	"testing"

	"github.com/bianavic/fullcycle_go_weather-by-zipcode/internal/domain"
)

func TestNewZipcode(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    domain.Zipcode
		wantErr error
	}{
		{name: "accepts 8 digits", input: "01001000", want: domain.Zipcode("01001000")},
		{name: "rejects fewer than 8 digits", input: "1234567", wantErr: domain.ErrInvalidZipcode},
		{name: "rejects more than 8 digits", input: "123456789", wantErr: domain.ErrInvalidZipcode},
		{name: "rejects non-digit characters", input: "0100100a", wantErr: domain.ErrInvalidZipcode},
		{name: "rejects empty string", input: "", wantErr: domain.ErrInvalidZipcode},
		{name: "rejects hyphenated format", input: "01001-000", wantErr: domain.ErrInvalidZipcode},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := domain.NewZipcode(tc.input)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("err = %v, want %v", err, tc.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}

func TestZipcodeValidate(t *testing.T) {
	t.Run("valid zipcode passes", func(t *testing.T) {
		z := domain.Zipcode("01001000")
		if err := z.Validate(); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("invalid zipcode returns ErrInvalidZipcode", func(t *testing.T) {
		z := domain.Zipcode("abc")
		if err := z.Validate(); !errors.Is(err, domain.ErrInvalidZipcode) {
			t.Errorf("err = %v, want %v", err, domain.ErrInvalidZipcode)
		}
	})
}

func TestZipcodeString(t *testing.T) {
	z := domain.Zipcode("01001000")
	if z.String() != "01001000" {
		t.Errorf("got %q, want %q", z.String(), "01001000")
	}
}
