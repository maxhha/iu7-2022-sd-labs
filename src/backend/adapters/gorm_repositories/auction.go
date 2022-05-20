package gorm_repositories

//go:generate go run ../../codegen/gorm_repository/main.go --out auction_gen.go --entity Auction --methods Get,Lock,orderQuery,sliceQuery,Find,Create,Delete

import (
	"database/sql"
	"iu7-2022-sd-labs/buisness/entities"
	"iu7-2022-sd-labs/buisness/ports/repositories"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

var auctionFieldToColumn = map[repositories.AuctionOrderField]string{
	repositories.AuctionOrderFieldCreationDate: "created_at",
}

type Auction struct {
	ID             string `gorm:"<-:false;default:generated()"`
	RoomID         string
	ProductID      string
	BidStepTableID string
	MinAmount      decimal.Decimal
	StartedAt      time.Time
	FinishedAt     sql.NullTime
	CreatedAt      time.Time `gorm:"<-:create"`
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt
}

func (obj *Auction) From(e *entities.Auction) *Auction {
	if e == nil {
		return nil
	}

	obj.ID = e.ID()
	obj.RoomID = e.RoomID()
	obj.ProductID = e.ProductID()
	obj.BidStepTableID = e.BidStepTableID()
	obj.MinAmount = e.MinAmount()
	obj.StartedAt = e.StartedAt()
	obj.FinishedAt = sql.NullTime(e.FinishedAt())

	return obj
}

func (obj *Auction) Into(e *entities.Auction) *entities.Auction {
	if e == nil {
		return nil
	}

	e.SetID(obj.ID)
	e.SetRoomID(obj.RoomID)
	e.SetProductID(obj.ProductID)
	e.SetBidStepTableID(obj.BidStepTableID)
	e.SetStartedAt(obj.StartedAt)
	e.SetMinAmount(obj.MinAmount)
	e.SetFinishedAt(entities.NullTime(obj.FinishedAt))

	return e
}

func (r *AuctionRepository) filterQuery(query *gorm.DB, filter *repositories.AuctionFilter) (*gorm.DB, error) {
	if filter == nil {
		return query, nil
	}

	if len(filter.IDs) > 0 {
		query = query.Where("id in ?", filter.IDs)
	}

	if len(filter.ProductIDs) > 0 {
		query = query.Where("product_id in ?", filter.ProductIDs)
	}

	if len(filter.RoomIDs) > 0 {
		query = query.Where("room_id in ?", filter.RoomIDs)
	}

	return query, nil
}
