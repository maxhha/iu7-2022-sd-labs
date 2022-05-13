package repositories

import (
	"iu7-2022-sd-labs/buisness/entities"
)

type ProductOrderField string

const (
	ProductOrderFieldCreationDate ProductOrderField = "CreationDate"
)

type ProductFilter struct {
	IDs       []string
	NameQuery string
}

type ProductOrder struct {
	By   ProductOrderField
	Desc bool
}

type ProductFindParams struct {
	Filter *ProductFilter
	Order  *ProductOrder
	Slice  *ForwardSlice
}

type ProductRepository interface {
	Get(id string) (entities.Product, error)
	ShareLock(id string) (entities.Product, error)
	Find(params *ProductFindParams) ([]entities.Product, error)
	Create(room *entities.Product) error
	Update(id string, updateFn func(room *entities.Product) error) (entities.Product, error)
	Delete(id string) (entities.Product, error)
}
