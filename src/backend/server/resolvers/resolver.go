package resolvers

import (
	"iu7-2022-sd-labs/buisness/ports/bus"
	"iu7-2022-sd-labs/buisness/ports/interactors"
	"iu7-2022-sd-labs/server/ports"
	"log"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	auctionInteractor      interactors.AuctionInteractor
	bidStepTableInteractor interactors.BidStepTableInteractor
	blockListInteractor    interactors.BlockListInteractor
	consumerInteractor     interactors.ConsumerInteractor
	offerInteractor        interactors.OfferInteractor
	organizerInteractor    interactors.OrganizerInteractor
	productInteractor      interactors.ProductInteractor
	roomInteractor         interactors.RoomInteractor
	auth                   ports.Auth
	dataloader             ports.DataLoader
	eventBus               bus.EventBus
	logger                 *log.Logger
}

func New(
	auctionInteractor interactors.AuctionInteractor,
	bidStepTableInteractor interactors.BidStepTableInteractor,
	blockListInteractor interactors.BlockListInteractor,
	consumerInteractor interactors.ConsumerInteractor,
	offerInteractor interactors.OfferInteractor,
	organizerInteractor interactors.OrganizerInteractor,
	productInteractor interactors.ProductInteractor,
	roomInteractor interactors.RoomInteractor,
	auth ports.Auth,
	dataloader ports.DataLoader,
	eventBus bus.EventBus,
	logger *log.Logger,
) Resolver {
	return Resolver{
		auctionInteractor,
		bidStepTableInteractor,
		blockListInteractor,
		consumerInteractor,
		offerInteractor,
		organizerInteractor,
		productInteractor,
		roomInteractor,
		auth,
		dataloader,
		eventBus,
		logger,
	}
}
