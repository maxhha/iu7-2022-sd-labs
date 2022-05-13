package repositories

import (
	"iu7-2022-sd-labs/buisness/entities"
)

type RoomOrderField string

const (
	RoomOrderFieldCreationDate RoomOrderField = "CreationDate"
)

type RoomFilter struct {
	IDs          []string
	NameQuery    string
	AddressQuery string
}

type RoomOrder struct {
	By   RoomOrderField
	Desc bool
}

type RoomFindParams struct {
	Filter *RoomFilter
	Order  *RoomOrder
	Slice  *ForwardSlice
}

type RoomRepository interface {
	Get(id string) (entities.Room, error)
	Find(params *RoomFindParams) ([]entities.Room, error)
	Create(room *entities.Room) error
	Update(id string, updateFn func(room *entities.Room) error) (entities.Room, error)
	Delete(id string) (entities.Room, error)
}
