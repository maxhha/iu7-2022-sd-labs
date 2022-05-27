package jwt_auth

import (
	"errors"
	"fmt"
)

var ErrWrongSigningMethod = errors.New("unexpected signing method")

func Wrap(err error, s string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf(s+": %w", err)
}

func Wrapf(err error, format string, a ...interface{}) error {
	if err == nil {
		return nil
	}

	return fmt.Errorf(format+": %w", append(a, err)...)
}

func Is(err error, target error) bool {
	return errors.Is(err, target)
}
