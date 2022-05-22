package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"iu7-2022-sd-labs/server/generated"
	"iu7-2022-sd-labs/server/models"
	"iu7-2022-sd-labs/server/ports"
)

func (r *mutationResolver) CreateProduct(ctx context.Context, name string) (*models.ProductResult, error) {
	organizer, err := ports.ForOrganizer(ctx)
	if err != nil {
		return nil, err
	}

	product, err := r.productInteractor.Create(organizer.ID(), name)
	if err != nil {
		return nil, Wrap(err, "product interactor create")
	}

	return &models.ProductResult{
		Product: (&models.Product{}).From(&product),
	}, nil
}

func (r *mutationResolver) DeleteProduct(ctx context.Context, productID string) (bool, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) UpdateProduct(ctx context.Context, input models.UpdateProductInput) (*models.ProductResult, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *productResolver) Organizer(ctx context.Context, obj *models.Product) (*models.Organizer, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *productResolver) Auctions(ctx context.Context, obj *models.Product, first *int, after *string, filter *models.AuctionFilter) (*models.AuctionConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Products(ctx context.Context, first *int, after *string, filter *models.ProductFilter) (*models.ProductConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

// Product returns generated.ProductResolver implementation.
func (r *Resolver) Product() generated.ProductResolver { return &productResolver{r} }

type productResolver struct{ *Resolver }
