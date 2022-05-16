package resolvers

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"iu7-2022-sd-labs/server/models"
	"iu7-2022-sd-labs/server/ports"

	"github.com/hashicorp/go-multierror"
)

func (r *queryResolver) Viewer(ctx context.Context) (models.Viewer, error) {
	var errors error

	organzier, err := ports.ForOrganizer(ctx)
	if err == nil {
		return (&models.Organizer{}).From(&organzier), nil
	} else {
		errors = multierror.Append(errors, err)
	}

	consumer, err := ports.ForConsumer(ctx)
	if err == nil {
		return (&models.Consumer{}).From(&consumer), nil
	} else {
		errors = multierror.Append(errors, err)
	}

	return nil, errors
}
