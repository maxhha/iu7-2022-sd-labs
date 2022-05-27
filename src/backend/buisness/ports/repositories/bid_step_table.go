package repositories

import (
	"iu7-2022-sd-labs/buisness/entities"
)

type BidStepTableOrderField string

const (
	BidStepTableOrderFieldCreationDate BidStepTableOrderField = "CreatoinDate"
)

type BidStepTableFilter struct {
	IDs          []string
	NameQuery    string
	OrganizerIDs []string
}

type BidStepTableOrder struct {
	By   BidStepTableOrderField
	Desc bool
}

type BidStepTableFindParams struct {
	Filter *BidStepTableFilter
	Order  *BidStepTableOrder
	Slice  *ForwardSlice
}

type BidStepTableRepository interface {
	Get(id string) (entities.BidStepTable, error)
	Find(params *BidStepTableFindParams) ([]entities.BidStepTable, error)
	Create(table *entities.BidStepTable) error
	Update(id string, updateFn func(table *entities.BidStepTable) error) (entities.BidStepTable, error)
}
