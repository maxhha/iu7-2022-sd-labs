package interactors

import (
	"iu7-2022-sd-labs/buisness/entities"
	"iu7-2022-sd-labs/buisness/ports/repositories"
)

type OrganizerUpdateParams struct {
	ID   string
	Name string
}

type OrganizerInteractor interface {
	Create(name string) (entities.Organizer, error)
	Find(params *repositories.OrganizerFindParams) ([]entities.Organizer, error)
	Update(params *OrganizerUpdateParams) (entities.Organizer, error)
}
