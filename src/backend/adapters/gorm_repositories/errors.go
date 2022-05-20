package gorm_repositories

import (
	"errors"
	"fmt"
	"iu7-2022-sd-labs/buisness/ports/repositories"

	"gorm.io/gorm"
)

var errorsMap = map[error]error{
	gorm.ErrRecordNotFound: repositories.ErrNotFound,
}

func mapError(err error) error {
	for check, converted := range errorsMap {
		if errors.Is(err, check) {
			return fmt.Errorf("%w (converted from: %v)", converted, err)
		}
	}

	return err
}

func Wrap(err error, s string) error {
	if err == nil {
		return nil
	}

	return fmt.Errorf(s+": %w", mapError(err))
}

func Wrapf(err error, format string, a ...interface{}) error {
	if err == nil {
		return nil
	}

	return fmt.Errorf(format+": %w", append(a, mapError(err))...)
}

func Is(err error, target error) bool {
	return errors.Is(err, target)
}
