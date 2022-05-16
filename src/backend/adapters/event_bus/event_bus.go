package event_bus

import (
	"iu7-2022-sd-labs/buisness/ports/bus"
	"sync"
)

type EventBus struct {
	mu          sync.Mutex
	subsCounter int
	subs        map[int]chan bus.Event
}

func NewEventBus() EventBus {
	return EventBus{}
}

func (b *EventBus) Subscribe() (<-chan bus.Event, int) {
	b.mu.Lock()
	defer b.mu.Unlock()

	ch := make(chan bus.Event, 1)
	id := b.subsCounter
	b.subsCounter++
	b.subs[id] = ch

	return ch, id
}

func (b *EventBus) Unsubscribe(id int) {
	b.mu.Lock()
	defer b.mu.Unlock()

	delete(b.subs, id)
}

func (b *EventBus) Notify(event bus.Event) {
	b.mu.Lock()

	for _, ch := range b.subs {
		go func(ch chan<- bus.Event, event bus.Event) {
			ch <- event
		}(ch, event)
	}

	b.mu.Unlock()
}
