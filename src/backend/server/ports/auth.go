package ports

import (
	"errors"
	"iu7-2022-sd-labs/buisness/entities"
)

var ErrWrongSubject = errors.New("wrong subject")

type Auth interface {
	NewOrganizerToken(organizer *entities.Organizer) (string, error)
	ParseOrganizerToken(token string) (entities.Organizer, error)
	NewConsumerToken(organizer *entities.Consumer) (string, error)
	ParseConsumerToken(token string) (entities.Consumer, error)
}
