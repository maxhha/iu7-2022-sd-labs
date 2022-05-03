package interactors

import (
	"iu7-2022-sd-labs/buisness/entities"
	"iu7-2022-sd-labs/buisness/ports/repositories"
)

type ProductUpdateParams struct {
	ID   string
	Name string
}

type ProductInteractor interface {
	Create(organizerID string, name string) (entities.Product, error)
	Find(params *repositories.ProductFindParams) ([]entities.Product, error)
	Update(params *ProductUpdateParams) (entities.Product, error)
	Delete(id string) error
}
