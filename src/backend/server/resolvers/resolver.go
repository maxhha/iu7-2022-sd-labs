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
	auth                ports.Auth
}

func New(
	organizerInteractor interactors.OrganizerInteractor,
	auth ports.Auth,
) Resolver {
	return Resolver{
		organizerInteractor,
		auth,
	}
}
