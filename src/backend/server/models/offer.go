package models

import (
	"iu7-2022-sd-labs/buisness/entities"
	"iu7-2022-sd-labs/buisness/ports/repositories"
	"time"

	"github.com/shopspring/decimal"
)

type Offer struct {
	ID         string `json:"id"`
	ConsumerID string
	AuctionID  string
	Amount     decimal.Decimal `json:"amount"`
	CreatedAt  time.Time       `json:"createdAt"`
}

func (obj *Offer) From(ent *entities.Offer) *Offer {
	if ent == nil {
		return nil
	}

	obj.ID = ent.ID()
	obj.ConsumerID = ent.ConsumerID()
	obj.AuctionID = ent.AuctionID()
	obj.Amount = ent.Amount()
	obj.CreatedAt = ent.CreatedAt()

	return obj
}

func (obj *OfferFilter) Into(ent *repositories.OfferFilter) *repositories.OfferFilter {
	if obj == nil {
		return nil
	}

	ent.IDs = obj.Ids
	ent.ConsumerIDs = obj.Consumers
	ent.AuctionIDs = obj.Auctions

	return ent
}

func OfferEdgesArrayFromEntites(ents []entities.Offer) []OfferConnectionEdge {
	objs := make([]OfferConnectionEdge, 0, len(ents))
	for _, ent := range ents {
		obj := Offer{}
		objs = append(objs, OfferConnectionEdge{
			Cursor: ent.ID(),
			Node:   obj.From(&ent),
		})
	}
	return objs
}
