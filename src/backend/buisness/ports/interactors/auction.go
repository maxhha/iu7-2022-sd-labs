package interactors

import (
	"iu7-2022-sd-labs/buisness/entities"
	"iu7-2022-sd-labs/buisness/ports/repositories"
	"time"

	"github.com/shopspring/decimal"
)

type AuctionCreateParams struct {
	RoomID         string
	BidStepTableID string
	ProductID      string
	StartedAt      time.Time
	MinAmount      decimal.Decimal
}

type AuctionInteractor interface {
	Create(params *AuctionCreateParams) (entities.Auction, error)
	Find(params *repositories.AuctionFindParams) ([]entities.Auction, error)
	Cancel(id string, reason string) error
}
