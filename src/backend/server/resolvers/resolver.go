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
	organizerInteractor interactors.OrganizerInteractor
	consumerInteractor  interactors.ConsumerInteractor
	roomInteractor      interactors.RoomInteractor
	productInteractor   interactors.ProductInteractor
	auth                ports.Auth
	dataloader          ports.DataLoader
	eventBus            bus.EventBus
	logger              *log.Logger
}

func New(
	organizerInteractor interactors.OrganizerInteractor,
	consumerInteractor interactors.ConsumerInteractor,
	roomInteractor interactors.RoomInteractor,
	productInteractor interactors.ProductInteractor,
	auth ports.Auth,
	dataloader ports.DataLoader,
	eventBus bus.EventBus,
	logger *log.Logger,
) Resolver {
	return Resolver{
		organizerInteractor,
		consumerInteractor,
		roomInteractor,
		productInteractor,
		auth,
		dataloader,
		eventBus,
		logger,
	}
}
