package interactors

import (
	"iu7-2022-sd-labs/buisness/entities"
	"iu7-2022-sd-labs/buisness/errors"
	"iu7-2022-sd-labs/buisness/ports/repositories"
)

type BlockListInteractor struct {
	repo repositories.Repository
}

func NewBlockListInteractor(repo repositories.Repository) BlockListInteractor {
	return BlockListInteractor{repo}
}

func (interactor *BlockListInteractor) Find(params *repositories.BlockListFindParams) ([]entities.BlockList, error) {
	ents, err := interactor.repo.BlockList().Find(params)
	return ents, errors.Wrap(err, "repo block list find")
}

func (interactor *BlockListInteractor) AddConsumer(organizerID string, consumerID string) (entities.BlockList, error) {
	blockList, err := interactor.repo.BlockList().UpdateOrCreateByOrganizer(organizerID, func(blockList *entities.BlockList) error {
		err := blockList.AddConsumerID(consumerID)
		return errors.Wrap(err, "block list add consumer id")
	})

	return blockList, errors.Wrap(err, "block list repo update")
}

func (interactor *BlockListInteractor) RemoveConsumer(organizerID string, consumerID string) (entities.BlockList, error) {
	blockList, err := interactor.repo.BlockList().UpdateOrCreateByOrganizer(organizerID, func(blockList *entities.BlockList) error {
		err := blockList.RemoveConsumerID(consumerID)
		return errors.Wrap(err, "block list remove consumer id")
	})

	return blockList, errors.Wrap(err, "block list repo update")
}
