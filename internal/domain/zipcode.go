package domain

import "regexp"

type Zipcode string

var zipcodePattern = regexp.MustCompile(`^[0-9]{8}$`)

func NewZipcode(s string) (Zipcode, error) {
	z := Zipcode(s)
	if err := z.Validate(); err != nil {
		return "", err
	}
	return z, nil
}

func (z Zipcode) Validate() error {
	if !zipcodePattern.MatchString(string(z)) {
		return ErrInvalidZipcode
	}
	return nil
}

func (z Zipcode) String() string {
	return string(z)
}
