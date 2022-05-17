package interactors

type Interactor interface {
	Auction() AuctionInteractor
	BidStepTable() BidStepTableInteractor
	Consumer() ConsumerInteractor
	Offer() OfferInteractor
	Organizer() OrganizerInteractor
	Product() ProductInteractor
	Room() RoomInteractor
}
