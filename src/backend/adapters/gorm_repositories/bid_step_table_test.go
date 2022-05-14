package gorm_repositories

import (
	"fmt"
	"iu7-2022-sd-labs/buisness/entities"
	"iu7-2022-sd-labs/buisness/ports/repositories"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/shopspring/decimal"
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

func (s *BidStepTableRepositorySuite) TestFind() {
	n := 10
	tables := make([]entities.BidStepTable, 0, n)
	for i := 0; i < n; i++ {
		table := entities.NewBidStepTable()
		table.SetID(fmt.Sprintf("test-table%d", i)).
			SetRows([]entities.BidStepRow{
				*entities.NewBidStepRowPtr().
					SetFromAmount(decimal.NewFromInt(int64(i))),
			})
		tables = append(tables, table)
	}

	type Case struct {
		Name   string
		Mock   func(*Case)
		Result []entities.BidStepTable
		Error  error
	}

	cases := []Case{{
		"Case: find tables error",
		func(c *Case) {
			s.SqlMock.ExpectQuery(s.SQL(`
				SELECT \* FROM "bid_step_tables"
				WHERE "bid_step_tables"\."deleted_at" IS NULL
			`)).
				WillReturnError(gorm.ErrRecordNotFound)
		},
		nil,
		repositories.ErrNotFound,
	}, {
		"Case: success",
		func(c *Case) {
			tableObjs := make([]BidStepTable, 0, n)
			for _, ent := range tables {
				obj := BidStepTable{}
				obj.From(&ent)
				tableObjs = append(tableObjs, obj)
			}

			s.SqlMock.ExpectQuery(s.SQL(`
				SELECT \* FROM "bid_step_tables"
				WHERE "bid_step_tables"\."deleted_at" IS NULL
			`)).
				WillReturnRows(MockRows(tableObjs))

			for _, ent := range tables {
				rowObjs := make([]BidStepRow, 0, len(ent.Rows()))
				for _, rowEnt := range ent.Rows() {
					rowObj := BidStepRow{TableID: ent.ID()}
					rowObj.From(&rowEnt)
					rowObjs = append(rowObjs, rowObj)
				}
				s.SqlMock.ExpectQuery(s.SQL(`
					SELECT \* FROM "bid_step_rows"
					WHERE table_id = \$1
				`)).
					WithArgs(ent.ID()).
					WillReturnRows(MockRows(rowObjs))
			}
		},
		tables,
		nil,
	}}

	for _, c := range cases {
		c.Mock(&c)
		result, err := s.repo.BidStepTable().Find(&repositories.BidStepTableFindParams{})
		if c.Error == nil {
			require.NoError(s.T(), err, c.Name)
			require.Equal(s.T(), c.Result, result, c.Name)
		} else {
			require.ErrorIs(s.T(), err, c.Error, c.Name)
		}
	}
}

func (s *BidStepTableRepositorySuite) TestCreate() {
	id := "test-table"
	organizerID := "test-organizer"
	name := "test-name"
	table := *entities.NewBidStepTablePtr().
		SetName(name).
		SetOrganizerID(organizerID).
		SetRows([]entities.BidStepRow{
			*entities.NewBidStepRowPtr(),
		})

	result := table
	result.SetID(id)

	type Case struct {
		Name   string
		Mock   func(*Case)
		Result entities.BidStepTable
		Error  error
	}

	cases := []Case{{
		"Case: success",
		func(c *Case) {
			s.SqlMock.ExpectBegin()
			s.SqlMock.ExpectQuery(s.SQL(`
				INSERT INTO "bid_step_tables" \("name","organizer_id","created_at","updated_at","deleted_at"\)
				VALUES \(\$1,\$2,\$3,\$4,\$5\)
				RETURNING "id"
			`)).
				WithArgs(name, organizerID, sqlmock.AnyArg(), sqlmock.AnyArg(), nil).
				WillReturnRows(MockRows(*(&BidStepTable{}).From(&result)))

			for _, rowEnt := range table.Rows() {
				s.SqlMock.ExpectExec(s.SQL(`
					INSERT INTO "bid_step_rows" \("table_id","from_amount","step"\)
					VALUES \(\$1,\$2,\$3\)
				`)).
					WithArgs(id, rowEnt.FromAmount().String(), rowEnt.Step().String()).
					WillReturnResult(sqlmock.NewResult(0, 1))
			}

			s.SqlMock.ExpectCommit()
		},
		result,
		nil,
	}}

	for _, c := range cases {
		c.Mock(&c)
		ent := table
		err := s.repo.BidStepTable().Create(&ent)
		if c.Error == nil {
			require.NoError(s.T(), err, c.Name)
			require.Equal(s.T(), c.Result, ent, c.Name)
		} else {
			require.ErrorIs(s.T(), err, c.Error, c.Name)
		}
	}
}

func (s *BidStepTableRepositorySuite) TestUpdate() {
	id := "test-table"
	organizerID := "test-organizer"
	name := "test-name"

	rowUnchanged := *entities.NewBidStepRowPtr().
		SetFromAmount(decimal.Zero)

	rowToDelete := *entities.NewBidStepRowPtr().
		SetFromAmount(decimal.NewFromInt(1))

	rowToUpdate := *entities.NewBidStepRowPtr().
		SetFromAmount(decimal.NewFromInt(2))

	rowUpdated := *entities.NewBidStepRowPtr().
		SetFromAmount(decimal.NewFromInt(2)).
		SetStep(decimal.NewFromInt(1))

	rowCreated := *entities.NewBidStepRowPtr().
		SetFromAmount(decimal.NewFromInt(3)).
		SetStep(decimal.NewFromInt(2))

	table := *entities.NewBidStepTablePtr().
		SetID(id).
		SetName(name).
		SetOrganizerID(organizerID).
		SetRows([]entities.BidStepRow{
			rowUnchanged,
			rowToDelete,
			rowToUpdate,
		})

	result := table
	result.SetRows([]entities.BidStepRow{
		rowUnchanged,
		rowCreated,
		rowUpdated,
	})

	type Case struct {
		Name     string
		Mock     func(*Case)
		UpdateFn func(table *entities.BidStepTable) error
		Result   entities.BidStepTable
		Error    error
	}

	cases := []Case{{
		"Case: success",
		func(c *Case) {
			s.SqlMock.ExpectBegin()
			s.SqlMock.ExpectQuery(s.SQL(`
				SELECT \* FROM "bid_step_tables"
				WHERE id = \$1
				AND "bid_step_tables"\."deleted_at" IS NULL
				LIMIT 1
				FOR UPDATE
			`)).
				WithArgs(id).
				WillReturnRows(MockRows(*(&BidStepTable{}).From(&table)))

			s.SqlMock.ExpectQuery(s.SQL(`
				SELECT \* FROM "bid_step_rows"
				WHERE table_id = \$1
			`)).
				WithArgs(id).
				WillReturnRows(MockRows(bidStepRowsFromEntity(table)))

			s.SqlMock.ExpectExec(s.SQL(`
				UPDATE "bid_step_tables"
				SET "name"=\$1,"organizer_id"=\$2,"updated_at"=\$3,"deleted_at"=\$4
				WHERE "bid_step_tables"\."deleted_at" IS NULL
				AND "id" = \$5
			`)).
				WithArgs(name, organizerID, sqlmock.AnyArg(), nil, id).
				WillReturnResult(sqlmock.NewResult(0, 1))

			s.SqlMock.ExpectExec(s.SQL(`
				DELETE FROM "bid_step_rows"
				WHERE \("bid_step_rows"\."table_id","bid_step_rows"\."from_amount"\) IN \(\(\$1,\$2\)\)
			`)).
				WithArgs(id, rowToDelete.FromAmount().String()).
				WillReturnResult(sqlmock.NewResult(0, 1))

			s.SqlMock.ExpectExec(s.SQL(`
				INSERT INTO "bid_step_rows" \("table_id","from_amount","step"\)
				VALUES \(\$1,\$2,\$3\)
			`)).
				WithArgs(id, rowCreated.FromAmount().String(), rowCreated.Step().String()).
				WillReturnResult(sqlmock.NewResult(1, 0))

			s.SqlMock.ExpectExec(s.SQL(`
				INSERT INTO "bid_step_rows" \("table_id","from_amount","step"\)
				VALUES \(\$1,\$2,\$3\)
				ON CONFLICT \("table_id","from_amount"\)
				DO UPDATE SET "step"="excluded"\."step"
			`)).
				WithArgs(id, rowUpdated.FromAmount().String(), rowUpdated.Step().String()).
				WillReturnResult(sqlmock.NewResult(1, 0))

			s.SqlMock.ExpectCommit()
		},
		func(table *entities.BidStepTable) error {
			table.SetRows(result.Rows())

			return nil
		},
		result,
		nil,
	}}

	for _, c := range cases {
		c.Mock(&c)
		result, err := s.repo.BidStepTable().Update(id, c.UpdateFn)
		if c.Error == nil {
			require.NoError(s.T(), err, c.Name)
			require.Equal(s.T(), c.Result, result, c.Name)
		} else {
			require.ErrorIs(s.T(), err, c.Error, c.Name)
		}
	}
}
