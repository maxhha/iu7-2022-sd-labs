package interactors

import (
	"iu7-2022-sd-labs/buisness/entities"
	"iu7-2022-sd-labs/buisness/errors"
	"iu7-2022-sd-labs/buisness/ports/interactors"
	"iu7-2022-sd-labs/buisness/ports/repositories"
)

type OrganizerInteractor struct {
	repo repositories.OrganizerRepository
}

func NewOrganizerInteractor(
	repo repositories.OrganizerRepository,
) OrganizerInteractor {
	return OrganizerInteractor{repo: repo}
}

func (interactor *OrganizerInteractor) Create(
	name string,
) (entities.Organizer, error) {
	org := entities.NewOrganizer()
	org.SetName(name)

	err := interactor.repo.Create(&org)
	return org, errors.Wrap(err, "repo create")
}

func (interactor *OrganizerInteractor) Find(
	params *repositories.OrganizerFindParams,
) ([]entities.Organizer, error) {
	orgs, err := interactor.repo.Find(params)
	return orgs, errors.Wrap(err, "repo find")
}

func (interactor *OrganizerInteractor) Update(
	params *interactors.OrganizerUpdateParams,
) (entities.Organizer, error) {
	org, err := interactor.repo.Get(params.ID)
	if err != nil {
		return org, errors.Wrap(err, "repo get")
	}
	org.SetName(params.Name)

	err = interactor.repo.Update(&org)

	return org, errors.Wrap(err, "repo update")
}
