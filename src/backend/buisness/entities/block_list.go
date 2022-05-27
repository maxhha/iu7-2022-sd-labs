package entities

import "iu7-2022-sd-labs/buisness/errors"

type BlockList struct {
	organizerID string
	consumerIDs []string
}

func NewBlockList() BlockList {
	return BlockList{}
}

func NewBlockListPtr() *BlockList {
	obj := NewBlockList()
	return &obj
}

func (obj *BlockList) OrganizerID() string {
	return obj.organizerID
}

func (obj *BlockList) SetOrganizerID(id string) *BlockList {
	obj.organizerID = id
	return obj
}

func (obj *BlockList) ConsumerIDs() []string {
	return obj.consumerIDs
}

func (obj *BlockList) SetConsumerIDs(consumerIDs []string) *BlockList {
	obj.consumerIDs = consumerIDs
	return obj
}

func (obj *BlockList) AddConsumerID(consumerID string) error {
	for _, existingConsumerID := range obj.ConsumerIDs() {
		if existingConsumerID == consumerID {
			return errors.Wrapf(ErrAlreadyExists, "consumerID='%s'", consumerID)
		}
	}

	obj.SetConsumerIDs(append(obj.ConsumerIDs(), consumerID))
	return nil
}

func (obj *BlockList) RemoveConsumerID(consumerID string) error {
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
		return errors.Wrapf(ErrConsumerNotFound, "consumerID='%s'", consumerID)
	}

	if len(newConsumerIDs) > 0 {
		obj.SetConsumerIDs(newConsumerIDs)
	} else {
		obj.SetConsumerIDs(nil)
	}

	return nil
}
