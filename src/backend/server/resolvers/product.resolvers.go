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
	organizer, err := ports.ForOrganizer(ctx)
	if err != nil {
		return false, err
	}

	product, err := r.dataloader.LoadProduct(ctx, productID)
	if err != nil {
		return false, Wrap(err, "dataloader.LoadProduct")
	}

	if product.OrganizerID() != organizer.ID() {
		return false, ErrDenied
	}

	err = r.productInteractor.Delete(productID)
	return err == nil, Wrap(err, "productInteractor.Delete")
}

func (r *mutationResolver) UpdateProduct(ctx context.Context, input models.UpdateProductInput) (*models.ProductResult, error) {
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

	product, err = r.productInteractor.Update(&interactors.ProductUpdateParams{
		ID:   input.ProductID,
		Name: input.Name,
	})
	if err != nil {
		return nil, Wrap(err, "productInteractor.Update")
	}

	return &models.ProductResult{
		Product: (&models.Product{}).From(&product),
	}, nil
}

func (r *productResolver) Organizer(ctx context.Context, obj *models.Product) (*models.Organizer, error) {
	ent, err := r.dataloader.LoadOrganizer(ctx, obj.OrganizerID)
	if err != nil {
		return nil, Wrap(err, "dataloader.LoadOrganizer")
	}
	return (&models.Organizer{}).From(&ent), nil
}

func (r *productResolver) Auctions(ctx context.Context, obj *models.Product, first *int, after *string, filter *models.AuctionFilter) (*models.AuctionConnection, error) {
	if len(filter.Products) > 0 {
		return nil, fmt.Errorf("filter products must be empty")
	}
	filter.Products = []string{obj.ID}
	return r.generatedPagination__Auctions(ctx, first, after, filter)
}

func (r *queryResolver) Products(ctx context.Context, first *int, after *string, filter *models.ProductFilter) (*models.ProductConnection, error) {
	return r.generatedPagination__Products(ctx, first, after, filter)
}

// Product returns generated.ProductResolver implementation.
func (r *Resolver) Product() generated.ProductResolver { return &productResolver{r} }

type productResolver struct{ *Resolver }
