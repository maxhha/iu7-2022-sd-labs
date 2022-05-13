package gorm_repositories

import (
	"iu7-2022-sd-labs/buisness/entities"
	"iu7-2022-sd-labs/buisness/ports/repositories"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type BidStepTableRepositorySuite struct {
	RepositorySuite
}

func TestBidStepTableRepositorySuite(t *testing.T) {
	suite.Run(t, new(BidStepTableRepositorySuite))
}

func (s *BidStepTableRepositorySuite) TestGet() {
	id := "test-table"
	table := *entities.NewBidStepTablePtr().
		SetID(id).
		SetRows([]entities.BidStepRow{
			*entities.NewBidStepRowPtr(),
		})

	type Case struct {
		Name   string
		Mock   func(*Case)
		Result entities.BidStepTable
		Error  error
	}

	cases := []Case{{
		"Case: not found",
		func(c *Case) {
			s.SqlMock.ExpectQuery(s.SQL(`
				SELECT \* FROM "bid_step_tables"
				WHERE id = \$1
				AND "bid_step_tables"\."deleted_at" IS NULL
			`)).
				WithArgs(id).
				WillReturnError(gorm.ErrRecordNotFound)
		},
		entities.NewBidStepTable(),
		repositories.ErrNotFound,
	}, {
		"Case: fail find rows",
		func(c *Case) {
			obj := BidStepTable{}
			obj.From(&table)

			s.SqlMock.ExpectQuery(s.SQL(`
				SELECT \* FROM "bid_step_tables"
				WHERE id = \$1
				AND "bid_step_tables"\."deleted_at" IS NULL
				LIMIT 1
			`)).
				WithArgs(id).
				WillReturnRows(MockRows(obj))

			s.SqlMock.ExpectQuery(s.SQL(`
				SELECT \* FROM "bid_step_rows"
				WHERE table_id = \$1
			`)).
				WithArgs(id).
				WillReturnError(gorm.ErrRecordNotFound)
		},
		entities.NewBidStepTable(),
		repositories.ErrNotFound,
	}, {
		"Case: success",
		func(c *Case) {
			obj := BidStepTable{}
			obj.From(&table)

			s.SqlMock.ExpectQuery(s.SQL(`
				SELECT \* FROM "bid_step_tables"
				WHERE id = \$1
				AND "bid_step_tables"\."deleted_at" IS NULL
				LIMIT 1
			`)).
				WithArgs(id).
				WillReturnRows(MockRows(obj))

			s.SqlMock.ExpectQuery(s.SQL(`
				SELECT \* FROM "bid_step_rows"
				WHERE table_id = \$1
			`)).
				WithArgs(id).
				WillReturnRows(MockRows(bidStepRowsFromEntity(table)))
		},
		table,
		nil,
	}}

	for _, c := range cases {
		c.Mock(&c)
		result, err := s.repo.BidStepTable().Get(id)

		if c.Error == nil {
			require.NoError(s.T(), err, c.Name)
			require.Equal(s.T(), c.Result, result, c.Name)
		} else {
			require.ErrorIs(s.T(), err, c.Error, c.Name)
		}
	}
}
