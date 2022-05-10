package interactors

import (
	"iu7-2022-sd-labs/buisness/entities"
	"iu7-2022-sd-labs/buisness/ports/repositories"

	"github.com/shopspring/decimal"
)

type OfferCreateParams struct {
	ConsumerID string
	AuctionID  string
	Amount     decimal.Decimal
}

type OfferInteractor interface {
	Create(params *OfferCreateParams) (entities.Offer, error)
	Find(params *repositories.OfferFindParams) ([]entities.Offer, error)
	Pay(id string) (string, error)
}
