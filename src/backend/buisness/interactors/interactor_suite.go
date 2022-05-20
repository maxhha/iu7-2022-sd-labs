package interactors

import (
	"iu7-2022-sd-labs/buisness/ports/mocks"

	"github.com/stretchr/testify/suite"
)

type InteractorSuite struct {
	suite.Suite
	repo             *mocks.Repository
	validatorService *mocks.ConsumerFormValidatorService
	payService       *mocks.OfferPayService
	eventBus         *mocks.EventBus
}

func (s *InteractorSuite) SetupTest() {
	s.repo = mocks.NewRepository(s.T())
	s.eventBus = mocks.NewEventBus(s.T())
	s.validatorService = mocks.NewConsumerFormValidatorService(s.T())
	s.payService = mocks.NewOfferPayService(s.T())
}
