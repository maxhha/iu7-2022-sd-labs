package entities

import (
	"time"

	"github.com/shopspring/decimal"
)

type Offer struct {
	id         string
	consumerID string
	auctionID  string
	amount     decimal.Decimal
	createdAt  time.Time
}

func NewOffer() Offer {
	return Offer{}
}

func (obj *Offer) ID() string {
	return obj.id
}

func (obj *Offer) SetID(id string) *Offer {
	obj.id = id
	return obj
}

func (obj *Offer) ConsumerID() string {
	return obj.consumerID
}

func (obj *Offer) SetConsumerID(consumerID string) *Offer {
	obj.consumerID = consumerID
	return obj
}

func (obj *Offer) AuctionID() string {
	return obj.auctionID
}

func (obj *Offer) SetAuctionID(auctionID string) *Offer {
	obj.auctionID = auctionID
	return obj
}

func (obj *Offer) Amount() decimal.Decimal {
	return obj.amount
}

func (obj *Offer) SetAmount(amount decimal.Decimal) *Offer {
	obj.amount = amount
	return obj
}

func (obj *Offer) CreatedAt() time.Time {
	return obj.createdAt
}

func (obj *Offer) SetCreatedAt(createdAt time.Time) *Offer {
	obj.createdAt = createdAt
	return obj
}
