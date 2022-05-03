package entities

import (
	"time"

	"github.com/shopspring/decimal"
)

type Auction struct {
	id             string
	roomID         string
	productID      string
	bidStepTableID string
	minAmount      decimal.Decimal
	startedAt      time.Time
	finishedAt     NullTime
}

func NewAuction() Auction {
	return Auction{}
}

func (obj *Auction) ID() string {
	return obj.id
}

func (obj *Auction) SetID(id string) *Auction {
	obj.id = id
	return obj
}

func (obj *Auction) RoomID() string {
	return obj.roomID
}

func (obj *Auction) SetRoomID(roomID string) *Auction {
	obj.roomID = roomID
	return obj
}

func (obj *Auction) ProductID() string {
	return obj.productID
}

func (obj *Auction) SetProductID(productID string) *Auction {
	obj.productID = productID
	return obj
}

func (obj *Auction) BidStepTableID() string {
	return obj.bidStepTableID
}

func (obj *Auction) SetBidStepTableID(bidStepTableID string) *Auction {
	obj.bidStepTableID = bidStepTableID
	return obj
}

func (obj *Auction) MinAmount() decimal.Decimal {
	return obj.minAmount
}

func (obj *Auction) SetMinAmount(minAmount decimal.Decimal) *Auction {
	obj.minAmount = minAmount
	return obj
}

func (obj *Auction) StartedAt() time.Time {
	return obj.startedAt
}

func (obj *Auction) SetStartedAt(startedAt time.Time) *Auction {
	obj.startedAt = startedAt
	return obj
}

func (obj *Auction) FinishedAt() NullTime {
	return obj.finishedAt
}

func (obj *Auction) SetFinishedAt(finishedAt NullTime) *Auction {
	obj.finishedAt = finishedAt
	return obj
}
