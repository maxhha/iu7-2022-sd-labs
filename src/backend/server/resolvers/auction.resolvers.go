package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"iu7-2022-sd-labs/server/generated"
	"iu7-2022-sd-labs/server/models"
)

func (r *auctionResolver) Room(ctx context.Context, obj *models.Auction) (*models.Room, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *auctionResolver) Product(ctx context.Context, obj *models.Auction) (*models.Product, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *auctionResolver) BidStepTable(ctx context.Context, obj *models.Auction) (*models.BidStepTable, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *auctionResolver) MinAmount(ctx context.Context, obj *models.Auction) (*float64, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *auctionResolver) StartedAt(ctx context.Context, obj *models.Auction) (string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *auctionResolver) FinishedAt(ctx context.Context, obj *models.Auction) (*string, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *auctionResolver) Offers(ctx context.Context, obj *models.Auction, first *int, after *string, filter *models.OfferFilter) (*models.OfferConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Auctions(ctx context.Context, first *int, after *string, filter *models.AuctionFilter) (*models.AuctionConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

// Auction returns generated.AuctionResolver implementation.
func (r *Resolver) Auction() generated.AuctionResolver { return &auctionResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type auctionResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
