package ports

import (
	"context"
	"iu7-2022-sd-labs/buisness/entities"
)

type DataLoader interface {
	WithNewLoader(ctx context.Context) context.Context
	LoadOrganizer(ctx context.Context, id string) (entities.Organizer, error)
}
