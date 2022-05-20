package interactors

import (
	"iu7-2022-sd-labs/buisness/entities"
	"iu7-2022-sd-labs/buisness/errors"
	"iu7-2022-sd-labs/buisness/ports/bus"
	"iu7-2022-sd-labs/buisness/ports/interactors"
	"iu7-2022-sd-labs/buisness/ports/repositories"
	"iu7-2022-sd-labs/buisness/ports/services"
)

type ConsumerInteractor struct {
	repo             repositories.Repository
	eventBus         bus.EventBus
	validatorService services.ConsumerFormValidatorService
}

func NewConsumerInteractor(
	repo repositories.Repository,
	eventBus bus.EventBus,
	validatorService services.ConsumerFormValidatorService,
) ConsumerInteractor {
	return ConsumerInteractor{
		repo,
		eventBus,
		validatorService,
	}
}

func (interactor *ConsumerInteractor) Create(
	nickname string,
	form interactors.ConsumerForm,
) (entities.Consumer, error) {
	if err := interactor.validatorService.Validate(form); err != nil {
		err = errors.Wrap(err, "form validator service validate")
		return entities.Consumer{}, err
	}

	consumer := entities.NewConsumer()
	consumer.
		SetNickname(nickname).
		SetForm(form)

	err := interactor.repo.Consumer().Create(&consumer)
	return consumer, errors.Wrap(err, "consumer repo create")
}

func (interactor *ConsumerInteractor) Find(
	params *repositories.ConsumerFindParams,
) ([]entities.Consumer, error) {
	consumers, err := interactor.repo.Consumer().Find(params)
	return consumers, errors.Wrap(err, "consumer repo find")
}

func (interactor *ConsumerInteractor) Update(
	params *interactors.ConsumerUpdateParams,
) (entities.Consumer, error) {
	consumer, err := interactor.repo.Consumer().Get(params.ID)
	if err != nil {
		return consumer, errors.Wrap(err, "consumer repo get")
	}

	if err = interactor.validatorService.Validate(params.Form); err != nil {
		return consumer, errors.Wrap(err, "form validator service validate")
	}

	consumer.
		SetNickname(params.Nickname).
		SetForm(params.Form)

	err = interactor.repo.Consumer().Update(&consumer)

	return consumer, errors.Wrap(err, "consumer repo update")
}

func (interactor *ConsumerInteractor) EnterRoom(
	consumerID string,
	roomID string,
) error {
	consumer, err := interactor.repo.Consumer().Get(consumerID)
	if err != nil {
		return errors.Wrap(err, "consumer repo get")
	}

	room, err := interactor.repo.Room().Update(roomID, func(room *entities.Room) error {
		err := room.AddConsumerID(consumer.ID())
		return errors.Wrap(err, "room add consumer id")
	})

	if err != nil {
		return errors.Wrap(err, "room repo update")
	}

	event := bus.EvtConsumerEnteredRoom{
		Consumer: consumer,
		Room:     room,
	}
	interactor.eventBus.Notify(&event)

	return nil
}

func (interactor *ConsumerInteractor) ExitRoom(
	consumerID string,
	roomID string,
) error {
	consumer, err := interactor.repo.Consumer().Get(consumerID)
	if err != nil {
		return errors.Wrap(err, "consumer repo get")
	}

	room, err := interactor.repo.Room().Update(roomID, func(room *entities.Room) error {
		err := room.RemoveConsumerID(consumer.ID())
		return errors.Wrap(err, "room remove consumer id")
	})

	if err != nil {
		return errors.Wrap(err, "room repo update")
	}

	event := bus.EvtConsumerExitedRoom{
		Consumer: consumer,
		Room:     room,
	}
	interactor.eventBus.Notify(&event)

	return nil
}
