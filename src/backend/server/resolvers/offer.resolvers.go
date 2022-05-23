package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"iu7-2022-sd-labs/server/generated"
	"iu7-2022-sd-labs/server/models"
)

func (r *offerResolver) Consumer(ctx context.Context, obj *models.Offer) (*models.Consumer, error) {
	ent, err := r.dataloader.LoadConsumer(ctx, obj.ConsumerID)
	if err != nil {
		return nil, Wrap(err, "dataloader.LoadConsumer")
	}
	return (&models.Consumer{}).From(&ent), nil
}

func (r *offerResolver) Auction(ctx context.Context, obj *models.Offer) (*models.Auction, error) {
	ent, err := r.dataloader.LoadAuction(ctx, obj.AuctionID)
	if err != nil {
		return nil, Wrap(err, "dataloader.LoadAuction")
	}
	return (&models.Auction{}).From(&ent), nil
}

func (r *queryResolver) Offers(ctx context.Context, first *int, after *string, filter *models.OfferFilter) (*models.OfferConnection, error) {
	return r.generatedPagination__Offers(ctx, first, after, filter)
}

// Offer returns generated.OfferResolver implementation.
func (r *Resolver) Offer() generated.OfferResolver { return &offerResolver{r} }

type offerResolver struct{ *Resolver }
