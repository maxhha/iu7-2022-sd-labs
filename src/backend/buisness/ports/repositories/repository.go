package repositories

type Repository interface {
	Auction() AuctionRepository
	BidStepTable() BidStepTableRepository
	Consumer() ConsumerRepository
	Offer() OfferRepository
	Organizer() OrganizerRepository
	Product() ProductRepository
	Room() RoomRepository
	BlockList() BlockListRepository
	Atomic(fn func(tx Repository) error) error
}
