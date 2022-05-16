package resolvers

import (
	"iu7-2022-sd-labs/buisness/ports/interactors"
	"iu7-2022-sd-labs/server/ports"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	organizerInteractor interactors.OrganizerInteractor
	consumerInteractor  interactors.ConsumerInteractor
	roomInteractor      interactors.RoomInteractor
	auth                ports.Auth
	dataloader          ports.DataLoader
}

func New(
	organizerInteractor interactors.OrganizerInteractor,
	consumerInteractor interactors.ConsumerInteractor,
	roomInteractor interactors.RoomInteractor,
	auth ports.Auth,
	dataloader ports.DataLoader,
) Resolver {
	return Resolver{
		organizerInteractor,
		consumerInteractor,
		roomInteractor,
		auth,
		dataloader,
	}
}
