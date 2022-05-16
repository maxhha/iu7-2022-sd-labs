package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"iu7-2022-sd-labs/buisness/ports/interactors"
	"iu7-2022-sd-labs/server/generated"
	"iu7-2022-sd-labs/server/models"
	"iu7-2022-sd-labs/server/ports"
)

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

func (r *queryResolver) Consumers(ctx context.Context, first *int, after *string, filter *models.ConsumerFilter) (*models.ConsumerConnection, error) {
	return r.generatedPagination__Consumers(ctx, first, after, filter)
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
