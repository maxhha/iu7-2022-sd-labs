package gorm_repositories

//go:generate go run ../../codegen/gorm_repository/main.go --out offer_gen.go --entity Offer --methods Get,orderQuery,sliceQuery,Find,Create

import (
	"iu7-2022-sd-labs/buisness/entities"
	"iu7-2022-sd-labs/buisness/ports/repositories"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

var offerFieldToColumn = map[repositories.OfferOrderField]string{
	repositories.OfferOrderFieldCreationDate: "created_at",
	repositories.OfferOrderFieldAmount:       "amount",
}

type Offer struct {
	ID         string `gorm:"<-:false;default:generated()"`
	ConsumerID string
	AuctionID  string
	Amount     decimal.Decimal
	CreatedAt  time.Time `gorm:"<-:create"`
	DeletedAt  gorm.DeletedAt
}

func (obj *Offer) From(e *entities.Offer) *Offer {
	if e == nil {
		return nil
	}

	obj.ID = e.ID()
	obj.ConsumerID = e.ConsumerID()
	obj.AuctionID = e.AuctionID()
	obj.Amount = e.Amount()
	obj.CreatedAt = e.CreatedAt()

	return obj
}

func (obj *Offer) Into(e *entities.Offer) *entities.Offer {
	if e == nil {
		return nil
	}

	e.SetID(obj.ID)
	e.SetConsumerID(obj.ConsumerID)
	e.SetAuctionID(obj.AuctionID)
	e.SetAmount(obj.Amount)
	e.SetCreatedAt(obj.CreatedAt)

	return e
}

func (r *OfferRepository) filterQuery(query *gorm.DB, filter *repositories.OfferFilter) (*gorm.DB, error) {
	if filter == nil {
		return query, nil
	}

	if len(filter.IDs) > 0 {
		query = query.Where("id in ?", filter.IDs)
	}

	if len(filter.ConsumerIDs) > 0 {
		query = query.Where("consumer_id in ?", filter.ConsumerIDs)
	}

	if len(filter.AuctionIDs) > 0 {
		query = query.Where("auction_id in ?", filter.AuctionIDs)
	}

	return query, nil
}
