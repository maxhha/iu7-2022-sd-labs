package interactors

import (
	"fmt"
	"iu7-2022-sd-labs/buisness/entities"
	"iu7-2022-sd-labs/buisness/ports/interactors"
	"iu7-2022-sd-labs/buisness/ports/repositories"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ProductSuite struct {
	InteractorSuite
	interactor ProductInteractor
}

func TestProductSuite(t *testing.T) {
	suite.Run(t, new(ProductSuite))
}

func (s *ProductSuite) SetupTest() {
	s.InteractorSuite.SetupTest()

	s.interactor = NewProductInteractor(
		s.repo.OrganizerMock,
		s.repo.ProductMock,
	)
}

func (s *ProductSuite) TestCreate() {
	organizerID := "test-organizer"
	name := "product-name"
	productID := "test-product"

	cases := []struct {
		Name    string
		Mock    func()
		Product entities.Product
		Error   error
	}{
		{
			"No organizer",
			func() {
				s.repo.OrganizerMock.
					On("Get", organizerID).
					Return(entities.Organizer{}, repositories.ErrNotFound).
					Once()
			},
			*s.NewProductPtr(),
			repositories.ErrNotFound,
		},
		{
			"Succes",
			func() {
				org := entities.NewOrganizer()
				org.SetID(organizerID)
				s.repo.OrganizerMock.
					On("Get", organizerID).
					Return(org, nil).
					Once()

				expectedProduct := s.
					NewProductPtr().
					SetOrganizerID(organizerID).
					SetName(name)

				s.repo.ProductMock.
					On("Create", expectedProduct).
					Run(func(args mock.Arguments) {
						p, ok := args.Get(0).(*entities.Product)
						require.True(s.T(), ok)
						p.SetID(productID)
					}).
					Return(nil).
					Once()
			},
			*s.NewProductPtr().
				SetOrganizerID(organizerID).
				SetName(name).
				SetID(productID),
			nil,
		},
	}

	for _, c := range cases {
		c.Mock()
		result, err := s.interactor.Create(organizerID, name)

		if c.Error != nil {
			require.ErrorIs(s.T(), err, repositories.ErrNotFound, c.Name)
		} else {
			require.NoError(s.T(), err, c.Name)
			require.Equal(s.T(), c.Product, result)
		}
	}
}

func (s *ProductSuite) TestFind() {
	n := 10
	products := make([]entities.Product, n)
	for i := 0; i < n; i++ {
		products[i].SetID(fmt.Sprintf("product%d", i))
	}

	params := repositories.ProductFindParams{}
	s.repo.ProductMock.On("Find", &params).Return(products, nil)

	result, err := s.interactor.Find(&params)
	require.NoError(s.T(), err)
	require.Equal(s.T(), products, result)
}

func (s *ProductSuite) TestUpdate() {
	id := "test-product"
	name := "product-name"
	expectedProduct := *s.NewProductPtr().SetID(id).SetName(name)

	s.repo.ProductMock.
		On("Update", id, mock.Anything).
		Run(func(args mock.Arguments) {
			updateFn, ok := args.Get(1).(func(product *entities.Product) error)
			require.True(s.T(), ok)

			product := entities.NewProduct()
			err := updateFn(&product)
			require.NoError(s.T(), err)
			require.Equal(s.T(), name, product.Name())
		}).
		Return(expectedProduct, nil)

	result, err := s.interactor.Update(&interactors.ProductUpdateParams{
		ID:   id,
		Name: name,
	})
	require.NoError(s.T(), err)
	require.Equal(s.T(), result, expectedProduct)
}

func (s *ProductSuite) TestDelete() {
	id := "test-product"
	s.repo.ProductMock.On("Delete", id).Return(entities.Product{}, nil)

	err := s.interactor.Delete(id)
	require.NoError(s.T(), err)
}
