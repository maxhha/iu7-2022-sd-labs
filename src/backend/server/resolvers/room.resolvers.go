package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"iu7-2022-sd-labs/buisness/ports/repositories"
	"iu7-2022-sd-labs/server/generated"
	"iu7-2022-sd-labs/server/models"
	"iu7-2022-sd-labs/server/ports"
)

func (r *mutationResolver) CreateRoom(ctx context.Context, name string, address string) (*models.RoomResult, error) {
	organizer, err := ports.ForOrganizer(ctx)
	if err != nil {
		return nil, err
	}

	room, err := r.roomInteractor.Create(organizer.ID(), name, address)
	if err != nil {
		return nil, Wrap(err, "room interactor create")
	}

	return &models.RoomResult{
		Room: (&models.Room{}).From(&room),
	}, nil
}

func (r *queryResolver) Rooms(ctx context.Context, first *int, after *string, filter *models.RoomFilter) (*models.RoomConnection, error) {
	return r.generatedPagination__Rooms(ctx, first, after, filter)
}

func (r *roomResolver) Organizer(ctx context.Context, obj *models.Room) (*models.Organizer, error) {
	organizer, err := r.dataloader.LoadOrganizer(ctx, obj.OrganizerID)
	if err != nil {
		return nil, Wrap(err, "dataloader load organizer")
	}

	return (&models.Organizer{}).From(&organizer), nil
}

func (r *roomResolver) Consumers(ctx context.Context, obj *models.Room) ([]models.Consumer, error) {
	consumers, err := r.consumerInteractor.Find(&repositories.ConsumerFindParams{
		Filter: &repositories.ConsumerFilter{
			IDs: obj.ConsumerIDs,
		},
	})
	if err != nil {
		return nil, Wrap(err, "consumer interactor find")
	}

	return models.ConsumerArrayFromEntites(consumers), nil
}

// Room returns generated.RoomResolver implementation.
func (r *Resolver) Room() generated.RoomResolver { return &roomResolver{r} }

type roomResolver struct{ *Resolver }
