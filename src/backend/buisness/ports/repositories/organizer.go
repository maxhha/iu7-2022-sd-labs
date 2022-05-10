package repositories

import (
	"iu7-2022-sd-labs/buisness/entities"
)

type OrganizerOrderField string

const (
	OrganizerOrderFieldCreationDate OrganizerOrderField = "CreatoinDate"
)

type OrganizerFilter struct {
	IDs       []string
	NameQuery string
}

type OrganizerOrder struct {
	By   OrganizerOrderField
	Desc bool
}

type OrganizerFindParams struct {
	Filter *OrganizerFilter
	Order  *OrganizerOrder
	Slice  *ForwardSlice
}

type OrganizerRepository interface {
	Get(id string) (entities.Organizer, error)
	Find(params *OrganizerFindParams) ([]entities.Organizer, error)
	Create(organizer *entities.Organizer) error
	Update(organizer *entities.Organizer) error
}
