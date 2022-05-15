package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"iu7-2022-sd-labs/server/generated"
	"iu7-2022-sd-labs/server/models"
)

func (r *mutationResolver) CreateOrganizer(ctx context.Context, name string) (*models.TokenResult, error) {
	organizer, err := r.organizerInteractor.Create(name)
	if err != nil {
		return nil, Wrap(err, "organizer interactor")
	}

	token, err := r.auth.NewOrganizerToken(&organizer)
	if err != nil {
		return nil, Wrap(err, "new organizer token")
	}

	return &models.TokenResult{
		Token: token,
	}, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

type mutationResolver struct{ *Resolver }
