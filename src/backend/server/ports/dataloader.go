package ports

import (
	"context"
	"iu7-2022-sd-labs/buisness/entities"
)

type DataLoader interface {
	WithNewLoader(ctx context.Context) context.Context
	LoadOrganizer(ctx context.Context, id string) (entities.Organizer, error)
	LoadManyOrganizers(ctx context.Context, ids []string) ([]entities.Organizer, error)
	LoadConsumer(ctx context.Context, id string) (entities.Consumer, error)
	LoadManyConsumers(ctx context.Context, ids []string) ([]entities.Consumer, error)
	LoadRoom(ctx context.Context, id string) (entities.Room, error)
	LoadManyRooms(ctx context.Context, ids []string) ([]entities.Room, error)
}
