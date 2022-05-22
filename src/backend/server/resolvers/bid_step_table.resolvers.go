package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"iu7-2022-sd-labs/server/generated"
	"iu7-2022-sd-labs/server/models"
)

func (r *bidStepRowResolver) FromAmount(ctx context.Context, obj *models.BidStepRow) (float64, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *bidStepRowResolver) Step(ctx context.Context, obj *models.BidStepRow) (float64, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *bidStepTableResolver) Organizer(ctx context.Context, obj *models.BidStepTable) (*models.Organizer, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) BidStepTables(ctx context.Context, first *int, after *string, filter *models.BidStepTableFilter) (*models.BidStepTableConnection, error) {
	panic(fmt.Errorf("not implemented"))
}

// BidStepRow returns generated.BidStepRowResolver implementation.
func (r *Resolver) BidStepRow() generated.BidStepRowResolver { return &bidStepRowResolver{r} }

// BidStepTable returns generated.BidStepTableResolver implementation.
func (r *Resolver) BidStepTable() generated.BidStepTableResolver { return &bidStepTableResolver{r} }

type bidStepRowResolver struct{ *Resolver }
type bidStepTableResolver struct{ *Resolver }
