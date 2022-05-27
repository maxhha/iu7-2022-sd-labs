package interactors

import (
	"iu7-2022-sd-labs/buisness/ports/bus"
	"iu7-2022-sd-labs/buisness/ports/interactors"
	"iu7-2022-sd-labs/buisness/ports/repositories"
	"iu7-2022-sd-labs/buisness/ports/services"
)

type Interactor struct {
	auction      AuctionInteractor
	bidStepTable BidStepTableInteractor
	consumer     ConsumerInteractor
	offer        OfferInteractor
	organizer    OrganizerInteractor
	product      ProductInteractor
	room         RoomInteractor
}

func New(
	repo repositories.Repository,
	eventBus bus.EventBus,
	validatorService services.ConsumerFormValidatorService,
	payService services.OfferPayService,
) Interactor {
	return Interactor{
		auction:      NewAuctionInteractor(repo, eventBus),
		bidStepTable: NewBidStepTableInteractor(repo),
		consumer:     NewConsumerInteractor(repo, eventBus, validatorService),
		offer:        NewOfferInteractor(repo, eventBus, payService),
		organizer:    NewOrganizerInteractor(repo.Organizer()),
		product:      NewProductInteractor(repo.Organizer(), repo.Product()),
		room:         NewRoomInteractor(repo.Organizer(), repo.Room()),
	}
}

func (i *Interactor) Auction() interactors.AuctionInteractor {
	return &i.auction
}

func (i *Interactor) BidStepTable() interactors.BidStepTableInteractor {
	return &i.bidStepTable
}

func (i *Interactor) Consumer() interactors.ConsumerInteractor {
	return &i.consumer
}

func (i *Interactor) Offer() interactors.OfferInteractor {
	return &i.offer
}

func (i *Interactor) Organizer() interactors.OrganizerInteractor {
	return &i.organizer
}

func (i *Interactor) Product() interactors.ProductInteractor {
	return &i.product
}

func (i *Interactor) Room() interactors.RoomInteractor {
	return &i.room
}
