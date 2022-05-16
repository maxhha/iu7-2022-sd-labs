package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"iu7-2022-sd-labs/server/generated"
	"iu7-2022-sd-labs/server/models"
)

func (r *mutationResolver) CreateConsumer(ctx context.Context, nickname string, form map[string]interface{}) (*models.TokenResult, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Consumers(ctx context.Context, first *int, after *string, filter *models.ConsumerFilter) (*models.ConsumerConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
