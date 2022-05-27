package interactors

import (
	"iu7-2022-sd-labs/buisness/entities"
	"iu7-2022-sd-labs/buisness/ports/repositories"
)

type BlockListInteractor interface {
	Find(params *repositories.BlockListFindParams) ([]entities.BlockList, error)
	AddConsumer(organizerID string, consumerID string) (entities.BlockList, error)
	RemoveConsumer(organizerID string, consumerID string) (entities.BlockList, error)
}
