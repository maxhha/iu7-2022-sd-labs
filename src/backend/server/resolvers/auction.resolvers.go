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

func (r *auctionResolver) Room(ctx context.Context, obj *models.Auction) (*models.Room, error) {
	ent, err := r.dataloader.LoadRoom(ctx, obj.RoomID)
	if err != nil {
		return nil, Wrap(err, "dataloader loadRoom")
	}

	return (&models.Room{}).From(&ent), nil
}

func (r *auctionResolver) Product(ctx context.Context, obj *models.Auction) (*models.Product, error) {
	ent, err := r.dataloader.LoadProduct(ctx, obj.ProductID)
	if err != nil {
		return nil, Wrap(err, "dataloader loadProduct")
	}

	return (&models.Product{}).From(&ent), nil
}

func (r *auctionResolver) BidStepTable(ctx context.Context, obj *models.Auction) (*models.BidStepTable, error) {
	ent, err := r.dataloader.LoadBidStepTable(ctx, obj.BidStepTableID)
	if err != nil {
		return nil, Wrap(err, "dataloader loadBidStepTable")
	}

	return (&models.BidStepTable{}).From(&ent), nil
}

func (r *auctionResolver) Offers(ctx context.Context, obj *models.Auction, first *int, after *string, filter *models.OfferFilter) (*models.OfferConnection, error) {
	if len(filter.Auctions) > 0 {
		return nil, fmt.Errorf("filter auctions must be empty")
	}
	filter.Auctions = []string{obj.ID}
	return r.generatedPagination__Offers(ctx, first, after, filter)
}

func (r *mutationResolver) CreateAuction(ctx context.Context, input models.CreateAuctionInput) (*models.AuctionResult, error) {
	organizer, err := ports.ForOrganizer(ctx)
	if err != nil {
		return nil, err
	}

	product, err := r.dataloader.LoadProduct(ctx, input.ProductID)
	if err != nil {
		return nil, Wrap(err, "dataloader.LoadProduct")
	}

	if product.OrganizerID() != organizer.ID() {
		return nil, ErrDenied
	}

	ent, err := r.auctionInteractor.Create(&interactors.AuctionCreateParams{
		RoomID:         input.RoomID,
		BidStepTableID: input.BidStepTableID,
		ProductID:      product.ID(),
		StartedAt:      input.StartedAt,
		MinAmount:      input.MinAmount,
	})

	if err != nil {
		return nil, Wrap(err, "auctionInteractor.Create")
	}

	return &models.AuctionResult{
		Auction: (&models.Auction{}).From(&ent),
	}, nil
}

func (r *mutationResolver) CancelAuction(ctx context.Context, input models.CancelAuctionInput) (bool, error) {
	organizer, err := ports.ForOrganizer(ctx)
	if err != nil {
		return false, err
	}

	auction, err := r.dataloader.LoadAuction(ctx, input.AuctionID)
	if err != nil {
		return false, Wrap(err, "dataloader.LoadAuction")
	}

	product, err := r.dataloader.LoadProduct(ctx, auction.ProductID())
	if err != nil {
		return false, Wrap(err, "dataloader.LoadProduct")
	}

	if product.OrganizerID() != organizer.ID() {
		return false, ErrDenied
	}

	err = r.auctionInteractor.Cancel(input.AuctionID, input.Reason)
	return err == nil, Wrap(err, "auctionInteractor.Cancel")
}

func (r *queryResolver) Auctions(ctx context.Context, first *int, after *string, filter *models.AuctionFilter) (*models.AuctionConnection, error) {
	return r.generatedPagination__Auctions(ctx, first, after, filter)
}

// Auction returns generated.AuctionResolver implementation.
func (r *Resolver) Auction() generated.AuctionResolver { return &auctionResolver{r} }

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type auctionResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
