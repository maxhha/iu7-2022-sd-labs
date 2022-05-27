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

func (r *bidStepTableResolver) Organizer(ctx context.Context, obj *models.BidStepTable) (*models.Organizer, error) {
	ent, err := r.dataloader.LoadOrganizer(ctx, obj.OrganizerID)
	if err != nil {
		return nil, Wrap(err, "dataloader LoadOrganizer")
	}
	return (&models.Organizer{}).From(&ent), nil
}

func (r *mutationResolver) CreateBidStepTable(ctx context.Context, input models.CreateBidStepTableInput) (*models.BidStepTableResult, error) {
	organizer, err := ports.ForOrganizer(ctx)
	if err != nil {
		return nil, err
	}

	table, err := r.bidStepTableInteractor.Create(&interactors.BidStepTableCreateParams{
		OrganizerID: organizer.ID(),
		Name:        input.Name,
		Rows:        models.BidStepRowInputsArrayIntoInteractorRows(input.Rows),
	})
	if err != nil {
		return nil, Wrap(err, "bidStepTableInteractor.Create")
	}

	return &models.BidStepTableResult{
		BidStepTable: (&models.BidStepTable{}).From(&table),
	}, nil
}

func (r *mutationResolver) UpdateBidStepTable(ctx context.Context, input models.UpdateBidStepTableInput) (*models.BidStepTableResult, error) {
	organizer, err := ports.ForOrganizer(ctx)
	if err != nil {
		return nil, err
	}

	table, err := r.dataloader.LoadBidStepTable(ctx, input.BidStepTableID)
	if err != nil {
		return nil, Wrap(err, "dataloader.LoadBidStepTable")
	}

	if table.OrganizerID() != organizer.ID() {
		return nil, ErrDenied
	}

	table, err = r.bidStepTableInteractor.Update(&interactors.BidStepTableUpdateParams{
		ID:   table.ID(),
		Name: input.Name,
		Rows: models.BidStepRowInputsArrayIntoInteractorRows(input.Rows),
	})
	if err != nil {
		return nil, Wrap(err, "bidStepTableInteractor.Update")
	}

	return &models.BidStepTableResult{
		BidStepTable: (&models.BidStepTable{}).From(&table),
	}, nil
}

func (r *queryResolver) BidStepTables(ctx context.Context, first *int, after *string, filter *models.BidStepTableFilter) (*models.BidStepTableConnection, error) {
	return r.generatedPagination__BidStepTables(ctx, first, after, filter)
}

// BidStepTable returns generated.BidStepTableResolver implementation.
func (r *Resolver) BidStepTable() generated.BidStepTableResolver { return &bidStepTableResolver{r} }

type bidStepTableResolver struct{ *Resolver }
