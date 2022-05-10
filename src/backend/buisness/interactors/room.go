package interactors

import (
	"iu7-2022-sd-labs/buisness/entities"
	"iu7-2022-sd-labs/buisness/errors"
	"iu7-2022-sd-labs/buisness/ports/repositories"
)

type RoomInteractor struct {
	organizerRepo repositories.OrganizerRepository
	roomRepo      repositories.RoomRepository
}

func NewRoomInteractor(
	organizerRepo repositories.OrganizerRepository,
	roomRepo repositories.RoomRepository,
) RoomInteractor {
	return RoomInteractor{organizerRepo, roomRepo}
}

func (interactor *RoomInteractor) Create(organizerID string, name string, address string) (entities.Room, error) {
	org, err := interactor.organizerRepo.Get(organizerID)
	if err != nil {
		return entities.Room{}, errors.Wrap(err, "organizer repo get")
	}

	room := entities.NewRoom()
	room.
		SetOrganizerID(org.ID()).
		SetName(name).
		SetAddress(address)

	err = interactor.roomRepo.Create(&room)
	return room, errors.Wrap(err, "room repo create")
}

func (interactor *RoomInteractor) Find(params *repositories.RoomFindParams) ([]entities.Room, error) {
	rooms, err := interactor.roomRepo.Find(params)
	return rooms, errors.Wrap(err, "room repo find")
}

func (interactor *RoomInteractor) Delete(id string) error {
	_, err := interactor.roomRepo.Delete(id)
	return errors.Wrap(err, "room repo delete")
}
