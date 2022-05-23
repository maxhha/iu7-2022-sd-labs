package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"iu7-2022-sd-labs/server/generated"
	"iu7-2022-sd-labs/server/models"
)

func (r *bidStepTableResolver) Organizer(ctx context.Context, obj *models.BidStepTable) (*models.Organizer, error) {
	ent, err := r.dataloader.LoadOrganizer(ctx, obj.OrganizerID)
	if err != nil {
		return nil, Wrap(err, "dataloader LoadOrganizer")
	}
	return (&models.Organizer{}).From(&ent), nil
}

func (r *queryResolver) BidStepTables(ctx context.Context, first *int, after *string, filter *models.BidStepTableFilter) (*models.BidStepTableConnection, error) {
	return r.generatedPagination__BidStepTables(ctx, first, after, filter)
}

// BidStepTable returns generated.BidStepTableResolver implementation.
func (r *Resolver) BidStepTable() generated.BidStepTableResolver { return &bidStepTableResolver{r} }

type bidStepTableResolver struct{ *Resolver }
