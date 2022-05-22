package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"iu7-2022-sd-labs/server/generated"
	"iu7-2022-sd-labs/server/models"
)

func (r *offerResolver) Consumer(ctx context.Context, obj *models.Offer) (*models.Consumer, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *offerResolver) Auction(ctx context.Context, obj *models.Offer) (*models.Auction, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *offerResolver) Amount(ctx context.Context, obj *models.Offer) (float64, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *offerResolver) CreatedAt(ctx context.Context, obj *models.Offer) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Offers(ctx context.Context, first *int, after *string, filter *models.OfferFilter) (*models.OfferConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

// Offer returns generated.OfferResolver implementation.
func (r *Resolver) Offer() generated.OfferResolver { return &offerResolver{r} }

type offerResolver struct{ *Resolver }
