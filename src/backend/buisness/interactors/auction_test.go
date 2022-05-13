package interactors

import (
	"fmt"
	"iu7-2022-sd-labs/buisness/entities"
	"iu7-2022-sd-labs/buisness/ports/bus"
	"iu7-2022-sd-labs/buisness/ports/interactors"
	"iu7-2022-sd-labs/buisness/ports/repositories"
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type AuctionSuite struct {
	InteractorSuite
	interactor AuctionInteractor
}

func TestAuctionSuite(t *testing.T) {
	suite.Run(t, new(AuctionSuite))
}

func (s *AuctionSuite) SetupTest() {
	s.InteractorSuite.SetupTest()

	s.interactor = NewAuctionInteractor(s.repo, s.eventBus)
}

func (s *AuctionSuite) TestCreate() {
	room := entities.NewRoom()
	room.SetID("test-room")

	table := entities.NewBidStepTable()
	table.SetID("test-table")

	expectedProduct := entities.NewProduct()
	expectedProduct.SetID("test-product")

	id := "test-auction"
	minAmount := decimal.NewFromFloat(123.45)
	startedAt := time.Now().UTC().Add(time.Duration(2) * time.Hour)

	type Case struct {
		Name   string
		Params interactors.AuctionCreateParams
		Mock   func(c *Case)
		Result entities.Auction
		Error  error
	}

	cases := []Case{
		{
			"Case: Get room error",
			interactors.AuctionCreateParams{
				RoomID: "unknown-room",
			},
			func(c *Case) {
				s.repo.RoomMock.On("Get", "unknown-room").
					Return(entities.Room{}, repositories.ErrNotFound).
					Once()
			},
			entities.Auction{},
			repositories.ErrNotFound,
		},
		{
			"Case: Get table error",
			interactors.AuctionCreateParams{
				RoomID:         room.ID(),
				BidStepTableID: "unknown-table",
			},
			func(c *Case) {
				s.repo.RoomMock.On("Get", room.ID()).
					Return(room, nil).
					Once()

				s.repo.BidStepTableMock.On("Get", "unknown-table").
					Return(entities.BidStepTable{}, repositories.ErrNotFound).
					Once()
			},
			entities.Auction{},
			repositories.ErrNotFound,
		},
		{
			"Case: Create error",
			interactors.AuctionCreateParams{
				RoomID:         room.ID(),
				BidStepTableID: table.ID(),
				ProductID:      expectedProduct.ID(),
			},
			func(c *Case) {
				s.repo.RoomMock.On("Get", room.ID()).
					Return(room, nil).
					Once()

				s.repo.BidStepTableMock.On("Get", table.ID()).
					Return(table, nil).
					Once()

				s.repo.ProductMock.On("ShareLock", expectedProduct.ID()).
					Return(expectedProduct, nil).
					Once()

				s.repo.AuctionMock.On("Create", mock.Anything).
					Return(repositories.ErrNotFound).
					Once()
			},
			entities.Auction{},
			repositories.ErrNotFound,
		},
		{
			"Case: Success",
			interactors.AuctionCreateParams{
				RoomID:         room.ID(),
				BidStepTableID: table.ID(),
				ProductID:      expectedProduct.ID(),
				StartedAt:      startedAt,
				MinAmount:      minAmount,
			},
			func(c *Case) {
				s.repo.RoomMock.On("Get", room.ID()).
					Return(room, nil).
					Once()

				s.repo.BidStepTableMock.On("Get", table.ID()).
					Return(table, nil).
					Once()

				s.repo.ProductMock.On("ShareLock", expectedProduct.ID()).
					Return(expectedProduct, nil).
					Once()

				s.repo.AuctionMock.On("Create", mock.Anything).
					Run(func(args mock.Arguments) {
						auction := args.Get(0).(*entities.Auction)
						auction.SetID(id)
					}).
					Return(nil).
					Once()

				s.eventBus.On("Notify", &bus.EvtAuctionCreated{
					Auction: c.Result,
				}).Once()
			},
			*s.NewAuctionPtr().
				SetID(id).
				SetStartedAt(startedAt).
				SetMinAmount(minAmount).
				SetProductID(expectedProduct.ID()).
				SetRoomID(room.ID()).
				SetBidStepTableID(table.ID()),
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

func (s *AuctionSuite) TestFind() {
	n := 10
	orgs := make([]entities.Auction, n)
	for i := 0; i < n; i++ {
		orgs[i].SetID(fmt.Sprintf("org%d", i))
	}

	params := repositories.AuctionFindParams{}
	s.repo.AuctionMock.On("Find", &params).Return(orgs, nil)

	result, err := s.interactor.Find(&params)
	require.NoError(s.T(), err)
	require.Equal(s.T(), orgs, result)
}

func (s *AuctionSuite) TestCancel() {
	id := "test-auction"
	reason := "test-reason"

	type Case struct {
		Name  string
		Mock  func(c *Case)
		Error error
	}

	cases := []Case{
		{
			"Case: fail delete",
			func(c *Case) {
				s.repo.AuctionMock.On("Delete", id).
					Return(entities.Auction{}, repositories.ErrNotFound).
					Once()
			},
			repositories.ErrNotFound,
		},
		{
			"Case: success",
			func(c *Case) {
				auction := entities.NewAuction()
				auction.SetID(id)
				s.repo.AuctionMock.On("Delete", id).
					Return(auction, nil).
					Once()

				s.eventBus.On("Notify", &bus.EvtAuctionCancelled{
					Auction: auction,
					Reason:  reason,
				})
			},
			nil,
		},
	}

	for _, c := range cases {
		c.Mock(&c)
		err := s.interactor.Cancel(id, reason)

		if c.Error == nil {
			require.NoError(s.T(), err, c.Name)
		} else {
			require.ErrorIs(s.T(), err, c.Error, c.Name)
		}
	}
}
