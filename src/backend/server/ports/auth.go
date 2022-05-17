package ports

import (
	"context"
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

type contextKey struct {
	name string
}

var organizerContextKey = contextKey{"organizer"}
var consumerContextKey = contextKey{"consumer"}
var ErrUnauthorized = errors.New("unauthorized")

func ForOrganizer(ctx context.Context) (entities.Organizer, error) {
	ent, ok := ctx.Value(organizerContextKey).(entities.Organizer)
	if !ok {
		return entities.Organizer{}, ErrUnauthorized
	}
	return ent, nil
}

func WithOrganizer(c context.Context, ent entities.Organizer) context.Context {
	return context.WithValue(c, organizerContextKey, ent)
}

func ForConsumer(ctx context.Context) (entities.Consumer, error) {
	ent, ok := ctx.Value(consumerContextKey).(entities.Consumer)
	if !ok {
		return entities.Consumer{}, ErrUnauthorized
	}
	return ent, nil
}

func WithConsumer(c context.Context, ent entities.Consumer) context.Context {
	return context.WithValue(c, consumerContextKey, ent)
}
