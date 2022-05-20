package gorm_repositories

import (
	"database/sql"
	"fmt"
	"iu7-2022-sd-labs/buisness/entities"
	"iu7-2022-sd-labs/buisness/ports/repositories"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type AuctionRepositorySuite struct {
	RepositorySuite
}


func TestAuctionRepositorySuite(t *testing.T) {
	suite.Run(t, new(AuctionRepositorySuite))
}

func (s *AuctionRepositorySuite) TestGet() {
	id := "test-auction"

	type Case struct {
		Name   string
		Mock   func(*Case)
		Result entities.Auction
		Error  error
	}

	cases := []Case{{
		"Case: not found",
		func(c *Case) {
			s.SqlMock.ExpectQuery(s.SQL(`
				SELECT \* FROM "auctions" WHERE id = \$1
				AND "auctions"\."deleted_at" IS NULL
				LIMIT 1
			`)).
				WithArgs(id).
				WillReturnError(gorm.ErrRecordNotFound)
		},
		entities.NewAuction(),
		repositories.ErrNotFound,
	}, {
		"Case: success",
		func(c *Case) {
			s.SqlMock.ExpectQuery(s.SQL(`
				SELECT \* FROM "auctions" WHERE id = \$1
				AND "auctions"\."deleted_at" IS NULL
				LIMIT 1
			`)).
				WithArgs(id).
				WillReturnRows(MockRows(Auction{ID: id}))
		},
		*entities.NewAuctionPtr().SetID(id),
		nil,
	}}

	for _, c := range cases {
		c.Mock(&c)
		result, err := s.repo.Auction().Get(id)

		if c.Error == nil {
			require.NoError(s.T(), err, c.Name)
			require.Equal(s.T(), c.Result, result, c.Name)
		} else {
			require.ErrorIs(s.T(), err, c.Error, c.Name)
		}
	}
}

func (s *AuctionRepositorySuite) TestLock() {
	id := "test-auction"

	type Case struct {
		Name   string
		Mock   func(*Case)
		Result entities.Auction
		Error  error
	}

	cases := []Case{{
		"Case: not found",
		func(c *Case) {
			s.SqlMock.ExpectQuery(s.SQL(`
				SELECT \* FROM "auctions" WHERE id = \$1
				AND "auctions"\."deleted_at" IS NULL
				LIMIT 1 
				FOR UPDATE
			`)).
				WithArgs(id).
				WillReturnError(gorm.ErrRecordNotFound)
		},
		entities.NewAuction(),
		repositories.ErrNotFound,
	}, {
		"Case: success",
		func(c *Case) {
			s.SqlMock.ExpectQuery(s.SQL(`
				SELECT \* FROM "auctions" WHERE id = \$1
				AND "auctions"\."deleted_at" IS NULL
				LIMIT 1
				FOR UPDATE
			`)).
				WithArgs(id).
				WillReturnRows(MockRows(Auction{ID: id}))
		},
		*entities.NewAuctionPtr().SetID(id),
		nil,
	}}

	for _, c := range cases {
		c.Mock(&c)
		result, err := s.repo.Auction().Lock(id)

		if c.Error == nil {
			require.NoError(s.T(), err, c.Name)
			require.Equal(s.T(), c.Result, result, c.Name)
		} else {
			require.ErrorIs(s.T(), err, c.Error, c.Name)
		}
	}
}

func (s *AuctionRepositorySuite) TestFind() {
	n := 10
	objs := make([]Auction, 0, n)
	ents := make([]entities.Auction, 0, n)
	for i := 0; i < n; i++ {
		id := fmt.Sprintf("test-auction%d", i)
		objs = append(objs, Auction{ID: id})
		ents = append(ents, *entities.NewAuctionPtr().SetID(id))
	}

	type Case struct {
		Name   string
		Mock   func(*Case)
		Result []entities.Auction
		Error  error
	}

	cases := []Case{{
		"Case: not found",
		func(c *Case) {
			s.SqlMock.ExpectQuery(s.SQL(`
				SELECT \* FROM "auctions" WHERE "auctions"\."deleted_at" IS NULL
			`)).
				WillReturnError(gorm.ErrRecordNotFound)
		},
		nil,
		repositories.ErrNotFound,
	}, {
		"Case: success",
		func(c *Case) {
			s.SqlMock.ExpectQuery(s.SQL(`
				SELECT \* FROM "auctions" WHERE "auctions"\."deleted_at" IS NULL
			`)).
				WillReturnRows(MockRows(objs))
		},
		ents,
		nil,
	}}

	for _, c := range cases {
		c.Mock(&c)
		result, err := s.repo.Auction().Find(&repositories.AuctionFindParams{})

		if c.Error == nil {
			require.NoError(s.T(), err, c.Name)
			require.Equal(s.T(), c.Result, result, c.Name)
		} else {
			require.ErrorIs(s.T(), err, c.Error, c.Name)
		}
	}
}

func (s *AuctionRepositorySuite) TestCreate() {
	id := "test-auction"
	roomID := "test-room-id"
	productID := "test-product-id"
	tableID := "test-table-id"
	minAmount := decimal.NewFromFloatWithExponent(123.45, -2)
	startedAt := time.Now()
	finishedAt := time.Now().Add(time.Duration(1) * time.Hour)

	ent := *entities.NewAuctionPtr().
		SetRoomID(roomID).
		SetProductID(productID).
		SetBidStepTableID(tableID).
		SetMinAmount(minAmount).
		SetStartedAt(startedAt).
		SetFinishedAt(entities.NullTime{
			Valid: true,
			Time:  finishedAt,
		})

	createdEnt := ent
	createdEnt.SetID(id)

	type Case struct {
		Name   string
		Mock   func(*Case)
		Result entities.Auction
		Error  error
	}

	cases := []Case{{
		"Case: fail create",
		func(c *Case) {
			s.SqlMock.ExpectQuery(s.SQL(`
				INSERT INTO "auctions" \("room_id","product_id","bid_step_table_id","min_amount","started_at","finished_at","created_at","updated_at","deleted_at"\)
				VALUES \(\$1,\$2,\$3,\$4,\$5,\$6,\$7,\$8,\$9\) 
				RETURNING "id"
			`)).
				WithArgs(roomID, productID, tableID, minAmount, startedAt, finishedAt, sqlmock.AnyArg(), sqlmock.AnyArg(), nil).
				WillReturnError(gorm.ErrRecordNotFound)
		},
		entities.NewAuction(),
		repositories.ErrNotFound,
	}, {
		"Case: success",
		func(c *Case) {
			obj := Auction{}
			obj.From(&ent)
			obj.ID = id

			s.SqlMock.ExpectQuery(s.SQL(`
				INSERT INTO "auctions" \("room_id","product_id","bid_step_table_id","min_amount","started_at","finished_at","created_at","updated_at","deleted_at"\)
				VALUES \(\$1,\$2,\$3,\$4,\$5,\$6,\$7,\$8,\$9\) 
				RETURNING "id"
			`)).
				WithArgs(roomID, productID, tableID, minAmount, startedAt, finishedAt, sqlmock.AnyArg(), sqlmock.AnyArg(), nil).
				WillReturnRows(MockRows(obj))
		},
		createdEnt,
		nil,
	}}

	for _, c := range cases {
		c.Mock(&c)
		result := ent
		err := s.repo.Auction().Create(&result)

		if c.Error == nil {
			require.NoError(s.T(), err, c.Name)
			require.Equal(s.T(), c.Result, result, c.Name)
		} else {
			require.ErrorIs(s.T(), err, c.Error, c.Name)
		}
	}
}

func (s *AuctionRepositorySuite) TestDelete() {
	id := "test-auction"
	auction := *entities.NewAuctionPtr().SetID(id)

	type Case struct {
		Name   string
		Mock   func(*Case)
		Result entities.Auction
		Error  error
	}

	cases := []Case{{
		"Case: not found",
		func(c *Case) {
			s.SqlMock.ExpectBegin()

			s.SqlMock.ExpectQuery(s.SQL(`
				SELECT \* FROM "auctions" WHERE id = \$1
				AND "auctions"\."deleted_at" IS NULL
				LIMIT 1 
				FOR UPDATE
			`)).
				WithArgs(id).
				WillReturnError(gorm.ErrRecordNotFound)

			s.SqlMock.ExpectRollback()
		},
		entities.NewAuction(),
		repositories.ErrNotFound,
	}, {
		"Case: fail delete",
		func(c *Case) {
			s.SqlMock.ExpectBegin()

			s.SqlMock.ExpectQuery(s.SQL(`
				SELECT \* FROM "auctions" WHERE id = \$1
				AND "auctions"\."deleted_at" IS NULL
				LIMIT 1 
				FOR UPDATE
			`)).
				WithArgs(id).
				WillReturnRows(MockRows(Auction{ID: id}))

			s.SqlMock.ExpectExec(s.SQL(`
				UPDATE "auctions" SET "deleted_at"=\$1
				WHERE "auctions"\."id" = \$2
				AND "auctions"\."deleted_at" IS NULL
			`)).
				WithArgs(sqlmock.AnyArg(), id).
				WillReturnError(sql.ErrConnDone)

			s.SqlMock.ExpectRollback()
		},
		entities.NewAuction(),
		sql.ErrConnDone,
	}, {
		"Case: success",
		func(c *Case) {
			s.SqlMock.ExpectBegin()

			s.SqlMock.ExpectQuery(s.SQL(`
				SELECT \* FROM "auctions" WHERE id = \$1
				AND "auctions"\."deleted_at" IS NULL
				LIMIT 1 
				FOR UPDATE
			`)).
				WithArgs(id).
				WillReturnRows(MockRows(Auction{ID: id}))

			s.SqlMock.ExpectExec(s.SQL(`
				UPDATE "auctions" SET "deleted_at"=\$1
				WHERE "auctions"\."id" = \$2
				AND "auctions"\."deleted_at" IS NULL
			`)).
				WithArgs(sqlmock.AnyArg(), id).
				WillReturnResult(sqlmock.NewResult(0, 1))

			s.SqlMock.ExpectCommit()
		},
		auction,
		nil,
	}}

	for _, c := range cases {
		c.Mock(&c)
		result, err := s.repo.Auction().Delete(id)

		if c.Error == nil {
			require.NoError(s.T(), err, c.Name)
			require.Equal(s.T(), c.Result, result, c.Name)
		} else {
			require.ErrorIs(s.T(), err, c.Error, c.Name)
		}
	}
}
