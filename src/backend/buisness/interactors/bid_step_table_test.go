package interactors

import (
	"fmt"
	"iu7-2022-sd-labs/buisness/entities"
	"iu7-2022-sd-labs/buisness/ports/interactors"
	"iu7-2022-sd-labs/buisness/ports/repositories"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type BidStepTableSuite struct {
	InteractorSuite
	interactor BidStepTableInteractor
}

func TestBidStepTableSuite(t *testing.T) {
	suite.Run(t, new(BidStepTableSuite))
}

func (s *BidStepTableSuite) SetupTest() {
	s.InteractorSuite.SetupTest()

	s.interactor = NewBidStepTableInteractor(
		s.organizerRepo,
		s.tableRepo,
	)
}

func (s *BidStepTableSuite) TestCreate() {
	name := "table-name"
	organizer := entities.NewOrganizer()
	organizer.SetID("test-organizer")

	type Case struct {
		Name   string
		Params interactors.BidStepTableCreateParams
		Mock   func(c *Case)
		Result entities.BidStepTable
		Error  error
	}

	cases := []Case{
		{
			"Case: fail get organizer",
			interactors.BidStepTableCreateParams{
				OrganizerID: "unknown-organizer",
			},
			func(c *Case) {
				s.organizerRepo.On("Get", "unknown-organizer").
					Return(entities.Organizer{}, repositories.ErrNotFound).
					Once()
			},
			*s.NewBidStepTablePtr(),
			repositories.ErrNotFound,
		},
		{
			"Case: invalid table",
			interactors.BidStepTableCreateParams{
				OrganizerID: organizer.ID(),
				Name:        name,
			},
			func(c *Case) {
				s.organizerRepo.On("Get", organizer.ID()).
					Return(organizer, nil).
					Once()
			},
			*s.NewBidStepTablePtr(),
			entities.ErrIsEmpty,
		},
		{
			"Case: success",
			interactors.BidStepTableCreateParams{
				OrganizerID: organizer.ID(),
				Name:        name,
				Rows: []interactors.BidStepRow{{
					FromAmount: decimal.Zero,
					Step:       decimal.NewFromInt(1),
				}},
			},
			func(c *Case) {
				s.organizerRepo.On("Get", organizer.ID()).
					Return(organizer, nil).
					Once()

				s.tableRepo.On("Create", mock.Anything).
					Run(func(args mock.Arguments) {
						table := args.Get(0).(*entities.BidStepTable)
						table.SetID(c.Result.ID())
					}).
					Return(nil).
					Once()
			},
			*s.NewBidStepTablePtr().
				SetID("test-table").
				SetOrganizerID(organizer.ID()).
				SetName(name).
				SetRows([]entities.BidStepRow{
					*s.NewBidStepRowPtr().
						SetFromAmount(decimal.Zero).
						SetStep(decimal.NewFromInt(1)),
				}),
			nil,
		},
	}

	for _, c := range cases {
		c.Mock(&c)
		result, err := s.interactor.Create(&c.Params)

		if c.Error == nil {
			require.NoError(s.T(), err, c.Name)
			require.Equal(s.T(), c.Result, result, c.Name)
		} else {
			require.ErrorIs(s.T(), err, c.Error, c.Name)
		}
	}
}

func (s *BidStepTableSuite) TestFind() {
	n := 10
	tables := make([]entities.BidStepTable, n)
	for i := 0; i < n; i++ {
		tables[i].SetID(fmt.Sprintf("table%d", i))
	}

	params := repositories.BidStepTableFindParams{}
	s.tableRepo.On("Find", &params).Return(tables, nil)

	result, err := s.interactor.Find(&params)
	require.NoError(s.T(), err)
	require.Equal(s.T(), tables, result)
}

func (s *BidStepTableSuite) TestUpdate() {
	table := entities.NewBidStepTable()
	table.
		SetID("test-id").
		SetOrganizerID("test-organizer").
		SetName("test-table")

	newName := "test-new-table"

	type Case struct {
		Name   string
		Params interactors.BidStepTableUpdateParams
		Mock   func(c *Case)
		Result entities.BidStepTable
		Error  error
	}

	cases := []Case{
		{
			"Case: fail get table",
			interactors.BidStepTableUpdateParams{
				ID: "unknown-table",
			},
			func(c *Case) {
				s.tableRepo.On("Get", "unknown-table").
					Return(entities.BidStepTable{}, repositories.ErrNotFound).
					Once()
			},
			*s.NewBidStepTablePtr(),
			repositories.ErrNotFound,
		},
		{
			"Case: invalid table",
			interactors.BidStepTableUpdateParams{
				ID:   table.ID(),
				Name: newName,
			},
			func(c *Case) {
				s.tableRepo.On("Get", table.ID()).
					Return(table, nil).
					Once()
			},
			*s.NewBidStepTablePtr(),
			entities.ErrIsEmpty,
		},
		{
			"Case: success",
			interactors.BidStepTableUpdateParams{
				ID:   table.ID(),
				Name: newName,
				Rows: []interactors.BidStepRow{{
					FromAmount: decimal.Zero,
					Step:       decimal.NewFromInt(1),
				}},
			},
			func(c *Case) {
				s.tableRepo.On("Get", table.ID()).
					Return(table, nil).
					Once()

				s.tableRepo.On("Update", &c.Result).
					Return(nil).
					Once()
			},
			*s.NewBidStepTablePtr().
				SetID(table.ID()).
				SetOrganizerID(table.OrganizerID()).
				SetName(newName).
				SetRows([]entities.BidStepRow{
					*s.NewBidStepRowPtr().
						SetFromAmount(decimal.Zero).
						SetStep(decimal.NewFromInt(1)),
				}),
			nil,
		},
	}

	for _, c := range cases {
		c.Mock(&c)
		result, err := s.interactor.Update(&c.Params)

		if c.Error == nil {
			require.NoError(s.T(), err, c.Name)
			require.Equal(s.T(), c.Result, result, c.Name)
		} else {
			require.ErrorIs(s.T(), err, c.Error, c.Name)
		}
	}
}
