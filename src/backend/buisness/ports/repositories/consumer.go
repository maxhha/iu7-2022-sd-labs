package repositories

import "iu7-2022-sd-labs/buisness/entities"

type ConsumerOrderField string

const (
	ConsumerOrderFieldCreationDate ConsumerOrderField = "CreatoinDate"
)

type ConsumerFormFieldQuery struct {
	Field string
	Query string
}

type ConsumerFilter struct {
	IDs              []string
	NickameQuery     string
	FormFieldQueries []ConsumerFormFieldQuery
}

type ConsumerOrder struct {
	By   ConsumerOrderField
	Desc bool
}

type ConsumerFindParams struct {
	Filter *ConsumerFilter
	Order  *ConsumerOrder
	Slice  *ForwardSlice
}

type ConsumerRepository interface {
	Get(id string) (entities.Consumer, error)
	Find(params *ConsumerFindParams) ([]entities.Consumer, error)
	Create(organizer *entities.Consumer) error
	Update(organizer *entities.Consumer) error
}
