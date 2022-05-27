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

func (r *mutationResolver) CreateOffer(ctx context.Context, input models.CreateOfferInput) (*models.OfferResult, error) {
	consumer, err := ports.ForConsumer(ctx)
	if err != nil {
		return nil, err
	}

	offer, err := r.offerInteractor.Create(&interactors.OfferCreateParams{
		ConsumerID: consumer.ID(),
		AuctionID:  input.AuctionID,
		Amount:     input.Amount,
	})

	if err != nil {
		return nil, Wrap(err, "offerInteractor.Create")
	}

	return &models.OfferResult{
		Offer: (&models.Offer{}).From(&offer),
	}, nil
}

func (r *mutationResolver) PayOffer(ctx context.Context, offerID string) (*models.PayOfferResult, error) {
	consumer, err := ports.ForConsumer(ctx)
	if err != nil {
		return nil, err
	}

	offer, err := r.dataloader.LoadOffer(ctx, offerID)
	if err != nil {
		return nil, Wrap(err, "dataloader.LoadOffer")
	}

	if offer.ConsumerID() != consumer.ID() {
		return nil, ErrDenied
	}

	link, err := r.offerInteractor.Pay(offerID)
	if err != nil {
		return nil, Wrap(err, "offerInteractor.Pay")
	}

	return &models.PayOfferResult{
		Link: link,
	}, nil
}

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
