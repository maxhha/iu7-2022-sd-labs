package interactors

import (
	"fmt"
	"iu7-2022-sd-labs/buisness/entities"
	"iu7-2022-sd-labs/buisness/ports/bus"
	"iu7-2022-sd-labs/buisness/ports/interactors"
	"iu7-2022-sd-labs/buisness/ports/repositories"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type OfferSuite struct {
	InteractorSuite
	interactor OfferInteractor
}

func TestOfferSuite(t *testing.T) {
	suite.Run(t, new(OfferSuite))
}

func (s *OfferSuite) SetupTest() {
	s.InteractorSuite.SetupTest()

	s.interactor = NewOfferInteractor(
		s.repo,
		s.eventBus,
		s.payService,
	)
}

func (s *OfferSuite) TestCreate() {
	offerID := "test-offer"
	consumerID := "test-consumer"
	consumer := *s.NewConsumerPtr().SetID(consumerID)
	tableID := "test-table"
	table := *s.NewBidStepTablePtr().SetID(tableID)
	auctionID := "test-auction"
	auction := *s.NewAuctionPtr().
		SetID(auctionID).
		SetBidStepTableID(tableID)

	auctionWithMinAmount := auction
	auctionWithMinAmount.SetMinAmount(decimal.NewFromInt(20))

	greaterMaxOffer := *s.NewOfferPtr().SetAmount(decimal.NewFromInt(20))

	params := interactors.OfferCreateParams{
		ConsumerID: consumerID,
		AuctionID:  auctionID,
		Amount:     decimal.NewFromInt(10),
	}

	type Case struct {
		Name   string
		Mock   func(c *Case)
		Result entities.Offer
		Error  string
	}

	cases := []Case{
		{
			"Case: fail get consumer",
			func(c *Case) {
				s.repo.ConsumerMock.On("Get", consumerID).
					Return(consumer, repositories.ErrNotFound).
					Once()
			},
			*s.NewOfferPtr(),
			"consumer repo get: not found",
		},
		{
			"Case: fail update auction",
			func(c *Case) {
				s.repo.ConsumerMock.On("Get", consumerID).
					Return(consumer, nil).
					Once()

				s.repo.AuctionMock.On("Lock", auctionID).
					Return(auction, repositories.ErrNotFound).
					Once()
			},
			*s.NewOfferPtr(),
			"auction repo lock: not found",
		},
		{
			"Case: fail get table",
			func(c *Case) {
				s.repo.ConsumerMock.On("Get", consumerID).
					Return(consumer, nil).
					Once()

				s.repo.AuctionMock.On("Lock", auctionID).
					Return(auction, nil).
					Once()

				s.repo.BidStepTableMock.On("Get", tableID).Return(
					entities.BidStepTable{},
					repositories.ErrNotFound,
				).Once()
			},
			*s.NewOfferPtr(),
			"table repo get: not found",
		},
		{
			"Case: fail get max offer",
			func(c *Case) {
				s.repo.ConsumerMock.On("Get", consumerID).
					Return(consumer, nil).
					Once()

				s.repo.AuctionMock.On("Lock", auctionID).
					Return(auction, nil).
					Once()

				s.repo.BidStepTableMock.On("Get", tableID).
					Return(table, nil).
					Once()

				s.repo.OfferMock.On("Find", mock.Anything).
					Return([]entities.Offer{}, repositories.ErrNotFound).
					Once()
			},
			*s.NewOfferPtr(),
			"offer repo find max offer: not found",
		},
		{
			"Case: amount is less then auction min bid",
			func(c *Case) {
				s.repo.ConsumerMock.On("Get", consumerID).
					Return(consumer, nil).
					Once()

				s.repo.AuctionMock.On("Lock", auctionID).
					Return(auctionWithMinAmount, nil).
					Once()

				s.repo.BidStepTableMock.On("Get", tableID).Return(table, nil).Once()

				s.repo.OfferMock.On("Find", mock.Anything).
					Return([]entities.Offer{}, nil).
					Once()
			},
			*s.NewOfferPtr(),
			"offered amount is less than min amount",
		},
		{
			"Case: offered amount is not allowed by bid step table",
			func(c *Case) {
				s.repo.ConsumerMock.On("Get", consumerID).
					Return(consumer, nil).
					Once()

				s.repo.AuctionMock.On("Lock", auctionID).
					Return(auction, nil).
					Once()

				s.repo.BidStepTableMock.On("Get", tableID).Return(table, nil).Once()

				s.repo.OfferMock.On("Find", mock.Anything).
					Return([]entities.Offer{greaterMaxOffer}, nil).
					Once()
			},
			*s.NewOfferPtr(),
			"table is not allowed bid:",
		},
		{
			"Case: success",
			func(c *Case) {
				s.repo.ConsumerMock.On("Get", consumerID).
					Return(consumer, nil).
					Once()

				s.repo.AuctionMock.On("Lock", auctionID).
					Return(auction, nil).
					Once()

				s.repo.BidStepTableMock.On("Get", tableID).Return(table, nil).Once()

				s.repo.OfferMock.On("Find", mock.Anything).
					Return([]entities.Offer{}, nil).
					Once()

				s.repo.OfferMock.On("Create", mock.Anything).
					Return(func(offer *entities.Offer) error {
						offer.SetID(offerID)
						return nil
					}).
					Once()

				s.eventBus.On("Notify", &bus.EvtOfferCreated{
					Offer: c.Result,
				}).Once()
			},
			*s.NewOfferPtr().
				SetID(offerID).
				SetConsumerID(consumerID).
				SetAuctionID(auctionID).
				SetAmount(params.Amount),
			"",
		},
	}

	for _, c := range cases {
		c.Mock(&c)
		result, err := s.interactor.Create(&params)

		if len(c.Error) == 0 {
			require.NoError(s.T(), err, c.Name)
			require.Equal(s.T(), c.Result, result, c.Name)
		} else {
			require.ErrorContains(s.T(), err, c.Error, c.Name)
		}
	}
}

func (s *OfferSuite) TestFind() {
	n := 10
	offers := make([]entities.Offer, n)
	for i := 0; i < n; i++ {
		offers[i].SetID(fmt.Sprintf("table%d", i))
	}

	params := repositories.OfferFindParams{}
	s.repo.OfferMock.On("Find", &params).Return(offers, nil)

	result, err := s.interactor.Find(&params)
	require.NoError(s.T(), err)
	require.Equal(s.T(), offers, result)
}

func (s *OfferSuite) TestPay() {
	auctionID := "test-auction"
	offerID := "test-offer"
	offer := *s.NewOfferPtr().SetID(offerID).SetAuctionID(auctionID)

	greaterMaxOffer := *s.NewOfferPtr().
		SetID("other-offer").
		SetAmount(decimal.NewFromInt(20))

	payLink := "test-pay-link"

	type Case struct {
		Name   string
		Mock   func(c *Case)
		Result string
		Error  string
	}

	cases := []Case{
		{
			"Case: fail get offer",
			func(c *Case) {
				s.repo.OfferMock.On("Get", offerID).
					Return(entities.Offer{}, repositories.ErrNotFound).
					Once()
			},
			"",
			"offer repo get: not found",
		},
		{
			"Case: fail find max offer",
			func(c *Case) {
				s.repo.OfferMock.On("Get", offerID).
					Return(offer, nil).
					Once()

				s.repo.OfferMock.On("Find", mock.Anything).
					Return([]entities.Offer{}, repositories.ErrNotFound).
					Once()
			},
			"",
			"offer repo find max offer: not found",
		},
		{
			"Case: pay not max offer",
			func(c *Case) {
				s.repo.OfferMock.On("Get", offerID).
					Return(offer, nil).
					Once()

				s.repo.OfferMock.On("Find", mock.Anything).
					Return([]entities.Offer{greaterMaxOffer}, nil).
					Once()
			},
			"",
			"offer is not max",
		},
		{
			"Case: success",
			func(c *Case) {
				s.repo.OfferMock.On("Get", offerID).
					Return(offer, nil).
					Once()

				s.repo.OfferMock.On("Find", mock.Anything).
					Return([]entities.Offer{offer}, nil).
					Once()

				s.payService.On("PayLink", mock.Anything).
					Return(payLink, nil)
			},
			payLink,
			"",
		},
	}

	for _, c := range cases {
		c.Mock(&c)
		result, err := s.interactor.Pay(offerID)

		if len(c.Error) == 0 {
			require.NoError(s.T(), err, c.Name)
			require.Equal(s.T(), c.Result, result, c.Name)
		} else {
			require.ErrorContains(s.T(), err, c.Error, c.Name)
		}
	}
}
