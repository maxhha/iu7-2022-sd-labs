package interactors

import (
	"iu7-2022-sd-labs/buisness/entities"
	"iu7-2022-sd-labs/buisness/ports/repositories"
)

type RoomInteractor interface {
	Create(organizerID string, name string, address string) (entities.Room, error)
	Find(params *repositories.RoomFindParams) ([]entities.Room, error)
	Delete(id string) error
}
