package interactors

import (
	"iu7-2022-sd-labs/buisness/entities"
	"iu7-2022-sd-labs/buisness/errors"
	"iu7-2022-sd-labs/buisness/ports/bus"
	"iu7-2022-sd-labs/buisness/ports/interactors"
	"iu7-2022-sd-labs/buisness/ports/repositories"
	"iu7-2022-sd-labs/buisness/ports/services"
)

type OfferInteractor struct {
	consumerRepo repositories.ConsumerRepository
	tableRepo    repositories.BidStepTableRepository
	offerRepo    repositories.OfferRepository
	auctionRepo  repositories.AuctionRepository
	eventBus     bus.EventBus
	payService   services.OfferPayService
}

func NewOfferInteractor(
	consumerRepo repositories.ConsumerRepository,
	tableRepo repositories.BidStepTableRepository,
	offerRepo repositories.OfferRepository,
	auctionRepo repositories.AuctionRepository,
	eventBus bus.EventBus,
	payService services.OfferPayService,
) OfferInteractor {
	return OfferInteractor{
		consumerRepo,
		tableRepo,
		offerRepo,
		auctionRepo,
		eventBus,
		payService,
	}
}

func (interactor *OfferInteractor) getMaxOffer(auctionID string) (entities.Offer, error) {
	maxOffer, err := interactor.offerRepo.Find(&repositories.OfferFindParams{
		Filter: &repositories.OfferFilter{
			AuctionIDs: []string{auctionID},
		},
		Order: &repositories.OfferOrder{
			By:   repositories.OfferOrderFieldAmount,
			Desc: true,
		},
		Slice: &repositories.ForwardSlice{
			Limit: 1,
		},
	})
	if err != nil {
		return entities.Offer{}, errors.Wrap(err, "offer repo find max offer")
	}

	if len(maxOffer) == 0 {
		return entities.Offer{}, entities.ErrNotFound
	}

	return maxOffer[0], nil
}

func (interactor *OfferInteractor) Create(
	params *interactors.OfferCreateParams,
) (entities.Offer, error) {
	var offer entities.Offer
	consumer, err := interactor.consumerRepo.Get(params.ConsumerID)
	if err != nil {
		return offer, errors.Wrap(err, "consumer repo get")
	}

	_, err = interactor.auctionRepo.Update(params.AuctionID, func(auction *entities.Auction) error {
		table, err := interactor.tableRepo.Get(auction.BidStepTableID())
		if err != nil {
			return errors.Wrap(err, "table repo get")
		}

		maxOffer, err := interactor.getMaxOffer(auction.ID())
		if err == nil {
			err := table.IsAllowedBid(maxOffer.Amount(), params.Amount)
			if err != nil {
				return errors.Wrap(err, "table is not allowed bid")
			}
		} else if errors.Is(err, entities.ErrNotFound) {
			if params.Amount.LessThan(auction.MinAmount()) {
				return errors.Wrapf(
					interactors.ErrOfferedAmountIsLessThanMinAmount,
					"amount=%s min=%s",
					params.Amount,
					auction.MinAmount(),
				)
			}
		} else {
			return errors.Wrap(err, "get max offer")
		}

		offer = entities.NewOffer()
		offer.
			SetConsumerID(consumer.ID()).
			SetAuctionID(auction.ID()).
			SetAmount(params.Amount)

		err = interactor.offerRepo.Create(&offer)

		return errors.Wrap(err, "offer repo create")
	})

	if err != nil {
		return offer, errors.Wrap(err, "auction repo update")
	}

	interactor.eventBus.Notify(&bus.EvtOfferCreated{
		Offer: offer,
	})

	return offer, nil
}

func (interactor *OfferInteractor) Find(
	params *repositories.OfferFindParams,
) ([]entities.Offer, error) {
	offers, err := interactor.offerRepo.Find(params)
	return offers, errors.Wrap(err, "offer repo find")
}

func (interactor *OfferInteractor) Pay(id string) (string, error) {
	offer, err := interactor.offerRepo.Get(id)
	if err != nil {
		return "", errors.Wrap(err, "offer repo get")
	}

	maxOffer, err := interactor.getMaxOffer(offer.AuctionID())
	if err != nil {
		return "", errors.Wrap(err, "get max offer")
	}

	if maxOffer.ID() != offer.ID() {
		return "", interactors.ErrOfferIsNotMax
	}

	link, err := interactor.payService.PayLink(&offer)

	return link, errors.Wrap(err, "pay service pay link")
}
