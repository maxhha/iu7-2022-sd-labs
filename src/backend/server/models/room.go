package models

import (
	"iu7-2022-sd-labs/buisness/entities"
	"iu7-2022-sd-labs/buisness/ports/repositories"
)

type Room struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Address     string `json:"address"`
	OrganizerID string
	ConsumerIDs []string
}

func (obj *Room) From(ent *entities.Room) *Room {
	if ent == nil {
		return nil
	}

	obj.ID = ent.ID()
	obj.Name = ent.Name()
	obj.Address = ent.Address()
	obj.OrganizerID = ent.OrganizerID()
	obj.ConsumerIDs = ent.ConsumerIDs()

	return obj
}

func (obj *RoomFilter) Into(ent *repositories.RoomFilter) *repositories.RoomFilter {
	if obj == nil {
		return nil
	}

	ent.IDs = obj.Ids
	if obj.Name == nil {
		ent.NameQuery = ""
	} else {
		ent.NameQuery = *obj.Name
	}
	if obj.Address == nil {
		ent.AddressQuery = ""
	} else {
		ent.AddressQuery = *obj.Address
	}
	ent.OrganizerIDs = obj.Orgainzers
	ent.ConsumerEnteredIDs = obj.Consumers

	return ent
}

func RoomEdgesArrayFromEntites(ents []entities.Room) []RoomConnectionEdge {
	objs := make([]RoomConnectionEdge, 0, len(ents))
	for _, ent := range ents {
		obj := Room{}
		objs = append(objs, RoomConnectionEdge{
			Cursor: ent.ID(),
			Node:   obj.From(&ent),
		})
	}
	return objs
}
