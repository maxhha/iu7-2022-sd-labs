package repositories

import (
	"iu7-2022-sd-labs/buisness/entities"
)

type OfferOrderField string

const (
	OfferOrderFieldCreationDate OfferOrderField = "CreationDate"
	OfferOrderFieldAmount       OfferOrderField = "Amount"
)

type OfferFilter struct {
	IDs         []string
	ConsumerIDs []string
	AuctionIDs  []string
}

type OfferOrder struct {
	By   OfferOrderField
	Desc bool
}

type OfferFindParams struct {
	Filter *OfferFilter
	Order  *OfferOrder
	Slice  *ForwardSlice
}

type OfferRepository interface {
	Get(id string) (entities.Offer, error)
	Find(params *OfferFindParams) ([]entities.Offer, error)
	Create(auction *entities.Offer) error
	Update(id string, updateFn func(offer *entities.Offer) error) (entities.Offer, error)
}
