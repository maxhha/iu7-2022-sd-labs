package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"iu7-2022-sd-labs/buisness/ports/interactors"
	"iu7-2022-sd-labs/buisness/ports/repositories"
	"iu7-2022-sd-labs/server/generated"
	"iu7-2022-sd-labs/server/models"
	"iu7-2022-sd-labs/server/ports"
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

func (r *mutationResolver) UpdateOrganizer(ctx context.Context, name string) (*models.OrganizerResult, error) {
	organizer, err := ports.ForOrganizer(ctx)
	if err != nil {
		return nil, Wrap(err, "for organizer")
	}

	organizer, err = r.organizerInteractor.Update(&interactors.OrganizerUpdateParams{
		ID:   organizer.ID(),
		Name: name,
	})
	if err != nil {
		return nil, Wrap(err, "orgainzer interactor update")
	}

	return &models.OrganizerResult{
		Orgainzer: (&models.Organizer{}).From(&organizer),
	}, nil
}

func (r *mutationResolver) BlockConsumer(ctx context.Context, consumerID string) (*models.OrganizerResult, error) {
	organizer, err := ports.ForOrganizer(ctx)
	if err != nil {
		return nil, err
	}

	_, err = r.blockListInteractor.AddConsumer(organizer.ID(), consumerID)
	if err != nil {
		return nil, err
	}

	return &models.OrganizerResult{
		Orgainzer: (&models.Organizer{}).From(&organizer),
	}, nil
}

func (r *mutationResolver) UnblockConsumer(ctx context.Context, consumerID string) (*models.OrganizerResult, error) {
	organizer, err := ports.ForOrganizer(ctx)
	if err != nil {
		return nil, err
	}

	_, err = r.blockListInteractor.RemoveConsumer(organizer.ID(), consumerID)
	if err != nil {
		return nil, err
	}

	return &models.OrganizerResult{
		Orgainzer: (&models.Organizer{}).From(&organizer),
	}, nil
}

func (r *organizerResolver) BidStepTables(ctx context.Context, obj *models.Organizer, first *int, after *string, filter *models.BidStepTableFilter) (*models.BidStepTableConnection, error) {
	if len(filter.Organizers) > 0 {
		return nil, fmt.Errorf("filter organizers must be empty")
	}
	filter.Organizers = []string{obj.ID}
	return r.generatedPagination__BidStepTables(ctx, first, after, filter)
}

func (r *organizerResolver) Products(ctx context.Context, obj *models.Organizer, first *int, after *string, filter *models.ProductFilter) (*models.ProductConnection, error) {
	if len(filter.Organizers) > 0 {
		return nil, fmt.Errorf("filter organizers must be empty")
	}
	filter.Organizers = []string{obj.ID}
	return r.generatedPagination__Products(ctx, first, after, filter)
}

func (r *organizerResolver) BlockList(ctx context.Context, obj *models.Organizer) ([]models.Consumer, error) {
	blockLists, err := r.blockListInteractor.Find(&repositories.BlockListFindParams{
		Filter: &repositories.BlockListFilter{
			OrganizerIDs: []string{obj.ID},
		},
	})
	if err != nil {
		return nil, Wrap(err, "blockListInteractor.Find")
	}

	if len(blockLists) == 0 {
		return nil, nil
	}

	ents, err := r.dataloader.LoadManyConsumers(ctx, blockLists[0].ConsumerIDs())
	if err != nil {
		return nil, Wrap(err, "dataloader.LoadManyConsumers")
	}

	return models.ConsumerArrayFromEntites(ents), nil
}

func (r *queryResolver) Organizers(ctx context.Context, first *int, after *string, filter *models.OrganizerFilter) (*models.OrganizerConnection, error) {
	return r.generatedPagination__Organizers(ctx, first, after, filter)
}

// Organizer returns generated.OrganizerResolver implementation.
func (r *Resolver) Organizer() generated.OrganizerResolver { return &organizerResolver{r} }

type organizerResolver struct{ *Resolver }
