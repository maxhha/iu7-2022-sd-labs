package interactors

import (
	"iu7-2022-sd-labs/buisness/entities"
	"iu7-2022-sd-labs/buisness/ports/repositories"
)

type ConsumerForm map[string]interface{}

type ConsumerUpdateParams struct {
	ID       string
	Nickname string
	Form     ConsumerForm
}

type ConsumerInteractor interface {
	Create(nickname string, form ConsumerForm) (entities.Consumer, error)
	Find(params *repositories.ConsumerFindParams) ([]entities.Consumer, error)
	Update(params *ConsumerUpdateParams) (entities.Consumer, error)
	EnterRoom(consumerID string, roomID string) error
	ExitRoom(consumerID string, roomID string) error
}
