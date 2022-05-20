package mocks

import (
	repositories "iu7-2022-sd-labs/buisness/ports/repositories"
	testing "testing"
)

type Repository struct {
	AuctionMock      *AuctionRepository
	BidStepTableMock *BidStepTableRepository
	ConsumerMock     *ConsumerRepository
	OfferMock        *OfferRepository
	OrganizerMock    *OrganizerRepository
	ProductMock      *ProductRepository
	RoomMock         *RoomRepository
}

func NewRepository(t testing.TB) *Repository {
	return &Repository{
		NewAuctionRepository(t),
		NewBidStepTableRepository(t),
		NewConsumerRepository(t),
		NewOfferRepository(t),
		NewOrganizerRepository(t),
		NewProductRepository(t),
		NewRoomRepository(t),
	}
}

func (r *Repository) Auction() repositories.AuctionRepository {
	return r.AuctionMock
}
func (r *Repository) BidStepTable() repositories.BidStepTableRepository {
	return r.BidStepTableMock
}
func (r *Repository) Consumer() repositories.ConsumerRepository {
	return r.ConsumerMock
}
func (r *Repository) Offer() repositories.OfferRepository {
	return r.OfferMock
}
func (r *Repository) Organizer() repositories.OrganizerRepository {
	return r.OrganizerMock
}
func (r *Repository) Product() repositories.ProductRepository {
	return r.ProductMock
}
func (r *Repository) Room() repositories.RoomRepository {
	return r.RoomMock
}
func (r *Repository) Atomic(fn func(r repositories.Repository) error) error {
	return fn(r)
}
