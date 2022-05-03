package entities

import "iu7-2022-sd-labs/buisness/errors"

type Room struct {
	id          string
	organizerID string
	name        string
	address     string
	consumerIDs []string
}

func NewRoom() Room {
	return Room{}
}

func (obj *Room) ID() string {
	return obj.id
}

func (obj *Room) SetID(id string) *Room {
	obj.id = id
	return obj
}

func (obj *Room) OrganizerID() string {
	return obj.organizerID
}

func (obj *Room) SetOrganizerID(id string) *Room {
	obj.organizerID = id
	return obj
}

func (obj *Room) Name() string {
	return obj.name
}

func (obj *Room) SetName(name string) *Room {
	obj.name = name
	return obj
}

func (obj *Room) Address() string {
	return obj.address
}

func (obj *Room) SetAddress(address string) *Room {
	obj.address = address
	return obj
}

func (obj *Room) ConsumerIDs() []string {
	return obj.consumerIDs
}

func (obj *Room) SetConsumerIDs(consumerIDs []string) *Room {
	obj.consumerIDs = consumerIDs
	return obj
}

func (obj *Room) AddConsumerID(consumerID string) error {
	for _, existingConsumerID := range obj.ConsumerIDs() {
		if existingConsumerID == consumerID {
			return errors.Wrapf(ErrAlreadyExists, "consumerID='%s'", consumerID)
		}
	}

	obj.SetConsumerIDs(append(obj.ConsumerIDs(), consumerID))
	return nil
}

func (obj *Room) RemoveConsumerID(consumerID string) error {
	newConsumerIDs := make([]string, 0, len(obj.consumerIDs))
	removed := false

	for _, existingConsumerID := range obj.ConsumerIDs() {
		if existingConsumerID == consumerID {
			removed = true
		} else {
			newConsumerIDs = append(newConsumerIDs, existingConsumerID)
		}
	}

	if !removed {
		return errors.Wrapf(ErrNotFound, "consumerID='%s'", consumerID)
	}

	if len(newConsumerIDs) > 0 {
		obj.SetConsumerIDs(newConsumerIDs)
	} else {
		obj.SetConsumerIDs(nil)
	}

	return nil
}
