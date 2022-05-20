package repositories

import (
	"iu7-2022-sd-labs/buisness/entities"
)

type AuctionOrderField string

const (
	AuctionOrderFieldCreationDate AuctionOrderField = "CreationDate"
)

type AuctionFilter struct {
	IDs        []string
	ProductIDs []string
	RoomIDs    []string
}

type AuctionOrder struct {
	By   AuctionOrderField
	Desc bool
}

type AuctionFindParams struct {
	Filter *AuctionFilter
	Order  *AuctionOrder
	Slice  *ForwardSlice
}

type AuctionRepository interface {
	Get(id string) (entities.Auction, error)
	Lock(id string) (entities.Auction, error)
	Find(params *AuctionFindParams) ([]entities.Auction, error)
	Create(auction *entities.Auction) error
	Delete(id string) (entities.Auction, error)
}
