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

func (r *organizerResolver) BidStepTables(ctx context.Context, obj *models.Organizer) ([]models.BidStepTable, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Organizers(ctx context.Context, first *int, after *string, filter *models.OrganizerFilter) (*models.OrganizerConnection, error) {
	ents, err := r.organizerInteractor.Find(&repositories.OrganizerFindParams{
		Filter: filter.Into(&repositories.OrganizerFilter{}),
		Order: &repositories.OrganizerOrder{
			By:   repositories.OrganizerOrderFieldCreationDate,
			Desc: true,
		},
		Slice: models.FillForwardSlice(first, after, &repositories.ForwardSlice{}),
	})
	if err != nil {
		return nil, Wrap(err, "organizer interactor find")
	}

	if len(ents) == 0 {
		return &models.OrganizerConnection{}, nil
	}

	objs := models.OrganizerEdgesArrayFromEntites(ents)

	hasNextPage := false
	if first != nil {
		hasNextPage = len(objs) > *first
		objs = objs[:len(objs)-1]
	}

	return &models.OrganizerConnection{
		PageInfo: &models.PageInfo{
			HasNextPage: hasNextPage,
			StartCursor: &objs[0].Cursor,
			EndCursor:   &objs[len(objs)-1].Cursor,
		},
		Edges: objs,
	}, nil
}

// Organizer returns generated.OrganizerResolver implementation.
func (r *Resolver) Organizer() generated.OrganizerResolver { return &organizerResolver{r} }

type organizerResolver struct{ *Resolver }
