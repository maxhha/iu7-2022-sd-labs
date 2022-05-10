package interactors

import (
	"iu7-2022-sd-labs/buisness/entities"
	"iu7-2022-sd-labs/buisness/errors"
	"iu7-2022-sd-labs/buisness/ports/interactors"
	"iu7-2022-sd-labs/buisness/ports/repositories"
)

type ProductInteractor struct {
	organizerRepo repositories.OrganizerRepository
	productRepo   repositories.ProductRepository
}

func NewProductInteractor(
	organizerRepo repositories.OrganizerRepository,
	productRepo repositories.ProductRepository,
) ProductInteractor {
	return ProductInteractor{organizerRepo, productRepo}
}

func (interactor *ProductInteractor) Create(
	organizerID string,
	name string,
) (entities.Product, error) {
	org, err := interactor.organizerRepo.Get(organizerID)
	if err != nil {
		return entities.Product{}, errors.Wrap(err, "organizer repo get")
	}

	product := entities.NewProduct()
	product.
		SetOrganizerID(org.ID()).
		SetName(name)

	err = interactor.productRepo.Create(&product)
	return product, errors.Wrap(err, "product repo create")
}

func (interactor *ProductInteractor) Find(
	params *repositories.ProductFindParams,
) ([]entities.Product, error) {
	products, err := interactor.productRepo.Find(params)
	return products, errors.Wrap(err, "product repo find")
}

func (interactor *ProductInteractor) Update(
	params *interactors.ProductUpdateParams,
) (entities.Product, error) {
	product, err := interactor.productRepo.Update(params.ID, func(product *entities.Product) error {
		product.SetName(params.Name)
		return nil
	})

	return product, errors.Wrap(err, "product repo update")
}

func (interactor *ProductInteractor) Delete(id string) error {
	_, err := interactor.productRepo.Delete(id)
	return errors.Wrap(err, "product repo delete")
}
