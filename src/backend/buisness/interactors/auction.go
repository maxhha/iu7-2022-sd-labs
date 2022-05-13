package interactors

import (
	"iu7-2022-sd-labs/buisness/entities"
	"iu7-2022-sd-labs/buisness/errors"
	"iu7-2022-sd-labs/buisness/ports/bus"
	"iu7-2022-sd-labs/buisness/ports/interactors"
	"iu7-2022-sd-labs/buisness/ports/repositories"
)

type AuctionInteractor struct {
	repo     repositories.Repository
	eventBus bus.EventBus
}

func NewAuctionInteractor(
	repo repositories.Repository,
	eventBus bus.EventBus,
) AuctionInteractor {
	return AuctionInteractor{
		repo,
		eventBus,
	}
}

func (interactor *AuctionInteractor) Create(
	params *interactors.AuctionCreateParams,
) (entities.Auction, error) {
	var auction entities.Auction
	room, err := interactor.repo.Room().Get(params.RoomID)
	if err != nil {
		return auction, errors.Wrap(err, "room repo get")
	}

	table, err := interactor.repo.BidStepTable().Get(params.BidStepTableID)
	if err != nil {
		return auction, errors.Wrap(err, "table repo get")
	}

	err = interactor.repo.Atomic(func(tx repositories.Repository) error {
		product, err := tx.Product().ShareLock(params.ProductID)

		if err != nil {
			return errors.Wrap(err, "product share lock")
		}

		auction = entities.NewAuction()
		auction.
			SetRoomID(room.ID()).
			SetBidStepTableID(table.ID()).
			SetProductID(product.ID()).
			SetMinAmount(params.MinAmount).
			SetStartedAt(params.StartedAt)

		err = tx.Auction().Create(&auction)
		return errors.Wrap(err, "auction repo create")
	})

	if err != nil {
		return auction, errors.Wrap(err, "repo atomic")
	}

	interactor.eventBus.Notify(&bus.EvtAuctionCreated{
		Auction: auction,
	})

	return auction, nil
}

func (interactor *AuctionInteractor) Find(
	params *repositories.AuctionFindParams,
) ([]entities.Auction, error) {
	auctions, err := interactor.repo.Auction().Find(params)
	return auctions, errors.Wrap(err, "auction repo find")
}

func (interactor *AuctionInteractor) Cancel(id string, reason string) error {
	auction, err := interactor.repo.Auction().Delete(id)
	if err != nil {
		return errors.Wrap(err, "auction repo delete")
	}

	interactor.eventBus.Notify(&bus.EvtAuctionCancelled{
		Auction: auction,
		Reason:  reason,
	})

	return nil
}
