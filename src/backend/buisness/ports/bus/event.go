package bus

import "iu7-2022-sd-labs/buisness/entities"

type Event interface {
	IsEvent()
}

type EvtConsumerEnteredRoom struct {
	Consumer entities.Consumer
	Room     entities.Room
}

func (e *EvtConsumerEnteredRoom) IsEvent() {}

type EvtConsumerExitedRoom struct {
	Consumer entities.Consumer
	Room     entities.Room
}

func (e *EvtConsumerExitedRoom) IsEvent() {}

type EvtAuctionCreated struct {
	Auction entities.Auction
}

func (e *EvtAuctionCreated) IsEvent() {}

type EvtAuctionCancelled struct {
	Auction entities.Auction
	Reason  string
}

func (e *EvtAuctionCancelled) IsEvent() {}

type EvtOfferCreated struct {
	Offer entities.Offer
}

func (e *EvtOfferCreated) IsEvent() {}

type EventBus interface {
	Notify(event Event)
	Subscribe() (<-chan Event, int)
	Unsubscribe(id int)
}
