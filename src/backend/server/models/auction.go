package models

import (
	"database/sql"
	"iu7-2022-sd-labs/buisness/entities"
	"iu7-2022-sd-labs/buisness/ports/repositories"
	"time"

	"github.com/shopspring/decimal"
)

type Auction struct {
	ID             string `json:"id"`
	RoomID         string
	ProductID      string
	BidStepTableID string
	MinAmount      decimal.NullDecimal `json:"minAmount"`
	StartedAt      time.Time           `json:"startedAt"`
	FinishedAt     sql.NullTime        `json:"finishedAt"`
}

func (obj *Auction) From(ent *entities.Auction) *Auction {
	if ent == nil {
		return nil
	}

	obj.ID = ent.ID()
	obj.RoomID = ent.RoomID()
	obj.ProductID = ent.ProductID()
	obj.BidStepTableID = ent.BidStepTableID()
	if ent.MinAmount().IsZero() {
		obj.MinAmount.Valid = false
	} else {
		obj.MinAmount.Valid = true
		obj.MinAmount.Decimal = ent.MinAmount()
	}
	obj.StartedAt = ent.StartedAt()
	obj.FinishedAt = sql.NullTime(ent.FinishedAt())

	return obj
}

func (obj *AuctionFilter) Into(ent *repositories.AuctionFilter) *repositories.AuctionFilter {
	if obj == nil {
		return nil
	}

	ent.IDs = obj.Ids
	ent.RoomIDs = obj.Rooms
	ent.ProductIDs = obj.Products

	return ent
}

func AuctionEdgesArrayFromEntites(ents []entities.Auction) []AuctionConnectionEdge {
	objs := make([]AuctionConnectionEdge, 0, len(ents))
	for _, ent := range ents {
		obj := Auction{}
		objs = append(objs, AuctionConnectionEdge{
			Cursor: ent.ID(),
			Node:   obj.From(&ent),
		})
	}
	return objs
}
