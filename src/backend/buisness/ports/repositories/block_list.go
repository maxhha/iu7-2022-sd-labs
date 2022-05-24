package repositories

import (
	"iu7-2022-sd-labs/buisness/entities"
)

type BlockListOrderField string

const (
	BlockListOrderFieldCreationDate BlockListOrderField = "CreationDate"
)

type BlockListFilter struct {
	OrganizerIDs []string
	ConsumerIDs  []string
}

type BlockListOrder struct {
	By   BlockListOrderField
	Desc bool
}

type BlockListFindParams struct {
	Filter *BlockListFilter
	Order  *BlockListOrder
	Slice  *ForwardSlice
}

type BlockListRepository interface {
	Find(params *BlockListFindParams) ([]entities.BlockList, error)
	UpdateOrCreateByOrganizer(organizerID string, updateFn func(ent *entities.BlockList) error) (entities.BlockList, error)
}
