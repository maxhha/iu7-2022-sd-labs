package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"iu7-2022-sd-labs/buisness/ports/bus"
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
	if len(obj.ConsumerIDs) == 0 {
		return nil, nil
	}

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

func (r *subscriptionResolver) ConsumersInRoomUpdated(ctx context.Context, roomID string) (<-chan *models.Room, error) {
	eventChan, subID := r.eventBus.Subscribe()
	go func() {
		<-ctx.Done()
		r.eventBus.Unsubscribe(subID)
	}()

	ch := make(chan *models.Room, 1)

	go func() {
		for event := range eventChan {
			switch event := event.(type) {
			case *bus.EvtConsumerEnteredRoom:
				if roomID == event.Room.ID() {
					ch <- (&models.Room{}).From(&event.Room)
				}
			case *bus.EvtConsumerExitedRoom:
				if roomID == event.Room.ID() {
					ch <- (&models.Room{}).From(&event.Room)
				}
			default:
			}
		}
	}()

	return ch, nil
}

// Room returns generated.RoomResolver implementation.
func (r *Resolver) Room() generated.RoomResolver { return &roomResolver{r} }

// Subscription returns generated.SubscriptionResolver implementation.
func (r *Resolver) Subscription() generated.SubscriptionResolver { return &subscriptionResolver{r} }

type roomResolver struct{ *Resolver }
type subscriptionResolver struct{ *Resolver }
