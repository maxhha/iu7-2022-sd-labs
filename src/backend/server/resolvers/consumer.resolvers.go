package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"iu7-2022-sd-labs/buisness/ports/interactors"
	"iu7-2022-sd-labs/server/generated"
	"iu7-2022-sd-labs/server/models"
	"iu7-2022-sd-labs/server/ports"
)

func (r *consumerResolver) Rooms(ctx context.Context, obj *models.Consumer, first *int, after *string, filter *models.RoomFilter) (*models.RoomConnection, error) {
	if len(filter.Consumers) > 0 {
		return nil, fmt.Errorf("consumers must be empty")
	}
	filter.Consumers = []string{obj.ID}
	return r.generatedPagination__Rooms(ctx, first, after, filter)
}

func (r *consumerResolver) Offers(ctx context.Context, obj *models.Consumer, first *int, after *string, filter *models.OfferFilter) (*models.OfferConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) CreateConsumer(ctx context.Context, nickname string, form map[string]interface{}) (*models.TokenResult, error) {
	ent, err := r.consumerInteractor.Create(nickname, form)
	if err != nil {
		return nil, Wrap(err, "consumer interactor create")
	}

	token, err := r.auth.NewConsumerToken(&ent)
	if err != nil {
		return nil, Wrap(err, "new consumer token")
	}

	return &models.TokenResult{
		Token: token,
	}, nil
}

func (r *mutationResolver) UpdateConsumer(ctx context.Context, nickname string, form map[string]interface{}) (*models.ConsumerResult, error) {
	ent, err := ports.ForConsumer(ctx)
	if err != nil {
		return nil, err
	}
	ent, err = r.consumerInteractor.Update(&interactors.ConsumerUpdateParams{
		ID:       ent.ID(),
		Nickname: nickname,
		Form:     form,
	})
	if err != nil {
		return nil, Wrap(err, "consumer interactor update")
	}

	return &models.ConsumerResult{
		Consumer: (&models.Consumer{}).From(&ent),
	}, nil
}

func (r *mutationResolver) EnterRoom(ctx context.Context, roomID string) (bool, error) {
	consumer, err := ports.ForConsumer(ctx)
	if err != nil {
		return false, err
	}

	err = r.consumerInteractor.EnterRoom(consumer.ID(), roomID)
	return err == nil, Wrap(err, "consumer interactor enter room")
}

func (r *mutationResolver) ExitRoom(ctx context.Context, roomID string) (bool, error) {
	consumer, err := ports.ForConsumer(ctx)
	if err != nil {
		return false, err
	}

	err = r.consumerInteractor.ExitRoom(consumer.ID(), roomID)
	return err == nil, Wrap(err, "consumer interactor exit room")
}

func (r *queryResolver) Consumers(ctx context.Context, first *int, after *string, filter *models.ConsumerFilter) (*models.ConsumerConnection, error) {
	return r.generatedPagination__Consumers(ctx, first, after, filter)
}

// Consumer returns generated.ConsumerResolver implementation.
func (r *Resolver) Consumer() generated.ConsumerResolver { return &consumerResolver{r} }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type consumerResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
