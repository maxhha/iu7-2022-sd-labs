package interactors

import (
	"iu7-2022-sd-labs/buisness/entities"
	"iu7-2022-sd-labs/buisness/ports/mocks"

	"github.com/stretchr/testify/suite"
)

type InteractorSuite struct {
	suite.Suite
	roomRepo         *mocks.RoomRepository
	organizerRepo    *mocks.OrganizerRepository
	offerRepo        *mocks.OfferRepository
	tableRepo        *mocks.BidStepTableRepository
	productRepo      *mocks.ProductRepository
	auctionRepo      *mocks.AuctionRepository
	consumerRepo     *mocks.ConsumerRepository
	validatorService *mocks.ConsumerFormValidatorService
	payService       *mocks.OfferPayService
	eventBus         *mocks.EventBus
}

func (s *InteractorSuite) SetupTest() {
	s.roomRepo = mocks.NewRoomRepository(s.T())
	s.organizerRepo = mocks.NewOrganizerRepository(s.T())
	s.offerRepo = mocks.NewOfferRepository(s.T())
	s.tableRepo = mocks.NewBidStepTableRepository(s.T())
	s.productRepo = mocks.NewProductRepository(s.T())
	s.auctionRepo = mocks.NewAuctionRepository(s.T())
	s.eventBus = mocks.NewEventBus(s.T())
	s.consumerRepo = mocks.NewConsumerRepository(s.T())
	s.validatorService = mocks.NewConsumerFormValidatorService(s.T())
	s.payService = mocks.NewOfferPayService(s.T())
}

func (s *InteractorSuite) NewOrganizerPtr() *entities.Organizer {
	obj := entities.NewOrganizer()
	return &obj
}

func (s *InteractorSuite) NewConsumerPtr() *entities.Consumer {
	obj := entities.NewConsumer()
	return &obj
}

func (s *InteractorSuite) NewAuctionPtr() *entities.Auction {
	obj := entities.NewAuction()
	return &obj
}

func (s *InteractorSuite) NewRoomPtr() *entities.Room {
	obj := entities.NewRoom()
	return &obj
}

func (s *InteractorSuite) NewOfferPtr() *entities.Offer {
	obj := entities.NewOffer()
	return &obj
}

func (s *InteractorSuite) NewBidStepTablePtr() *entities.BidStepTable {
	obj := entities.NewBidStepTable()
	return &obj
}

func (s *InteractorSuite) NewBidStepRowPtr() *entities.BidStepRow {
	obj := entities.NewBidStepRow()
	return &obj
}

func (s *InteractorSuite) NewProductPtr() *entities.Product {
	obj := entities.NewProduct()
	return &obj
}
