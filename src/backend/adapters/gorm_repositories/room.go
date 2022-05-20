package gorm_repositories

//go:generate go run ../../codegen/gorm_repository/main.go --out room_gen.go --entity Room --methods orderQuery,sliceQuery,Delete

import (
	"fmt"
	"iu7-2022-sd-labs/buisness/entities"
	"iu7-2022-sd-labs/buisness/ports/repositories"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const consumersInRoomTable = "room_consumer_rel"

var roomFieldToColumn = map[repositories.RoomOrderField]string{
	repositories.RoomOrderFieldCreationDate: "created_at",
}

type Room struct {
	ID          string `gorm:"<-:false;default:generated()"`
	OrganizerID string
	Name        string
	Address     string
	CreatedAt   time.Time `gorm:"<-:create"`
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt
}

func (obj *Room) From(e *entities.Room) *Room {
	if e == nil {
		return nil
	}

	obj.ID = e.ID()
	obj.OrganizerID = e.OrganizerID()
	obj.Name = e.Name()
	obj.Address = e.Address()

	return obj
}

func (obj *Room) Into(e *entities.Room) *entities.Room {
	if e == nil {
		return nil
	}

	e.SetID(obj.ID)
	e.SetOrganizerID(obj.OrganizerID)
	e.SetName(obj.Name)
	e.SetAddress(obj.Address)

	return e
}

func (r *RoomRepository) setConsumerIDs(db *gorm.DB, ent *entities.Room) ([]string, error) {
	var consumerIDs []string
	err := r.db.Table(consumersInRoomTable).
		Select("consumer_id").
		Where("table_id = ?", ent.ID()).
		Scan(&consumerIDs).
		Error

	if err != nil {
		return nil, Wrap(err, "db select consumers")
	}

	ent.SetConsumerIDs(consumerIDs)
	return consumerIDs, nil
}

func (r *RoomRepository) Get(id string) (entities.Room, error) {
	obj := Room{}
	ent := entities.NewRoom()

	if err := r.db.Take(&obj, "id = ?", id).Error; err != nil {
		return ent, Wrap(err, "db take")
	}

	obj.Into(&ent)
	_, err := r.setConsumerIDs(r.db, &ent)

	return ent, Wrap(err, "set consumer ids")
}

func (r *RoomRepository) filterQuery(query *gorm.DB, filter *repositories.RoomFilter) (*gorm.DB, error) {
	if filter == nil {
		return query, nil
	}

	if len(filter.IDs) > 0 {
		query = query.Where("id in ?", filter.IDs)
	}

	if len(filter.OrganizerIDs) > 0 {
		query = query.Where("organizer_id in ?", filter.OrganizerIDs)
	}

	if len(filter.NameQuery) > 0 {
		query = query.Where("name ~* ?", filter.NameQuery)
	}

	if len(filter.AddressQuery) > 0 {
		query = query.Where("address ~* ?", filter.AddressQuery)
	}

	if len(filter.ConsumerEnteredIDs) > 0 {
		query = query.Joins(
			fmt.Sprintf(
				"JOIN %s cr on cr.room_id = id AND cr.consumer_id in ?",
				consumersInRoomTable,
			),
			filter.ConsumerEnteredIDs,
		)
	}

	return query, nil
}

func (r *RoomRepository) Find(params *repositories.RoomFindParams) ([]entities.Room, error) {
	query := r.db.Model(&Room{})

	if params != nil {
		var err error
		if query, err = r.filterQuery(query, params.Filter); err != nil {
			return nil, Wrap(err, "filter query")
		}

		if query, err = r.orderQuery(query, params.Order); err != nil {
			return nil, Wrap(err, "order query")
		}

		if query, err = r.sliceQuery(query, params.Order, params.Slice); err != nil {
			return nil, Wrap(err, "slice query")
		}
	}

	var objs []Room
	if err := query.Find(&objs).Error; err != nil {
		return nil, Wrap(err, "db find")
	}

	ents := make([]entities.Room, 0, len(objs))

	for _, obj := range objs {
		ent := entities.NewRoom()
		obj.Into(&ent)

		if _, err := r.setConsumerIDs(r.db, &ent); err != nil {
			return nil, Wrapf(err, `set consumer ids for id="%s"`, obj.ID)
		}

		ents = append(ents, ent)
	}

	return ents, nil
}

func (r *RoomRepository) Create(ent *entities.Room) error {
	if len(ent.ConsumerIDs()) > 0 {
		return fmt.Errorf("consumer ids is not empty")
	}

	obj := Room{}
	obj.From(ent)

	if err := r.db.Create(&obj).Error; err != nil {
		return Wrap(err, "db create")
	}

	obj.Into(ent)
	return nil
}

func (r *RoomRepository) Update(id string, updateFn func(ent *entities.Room) error) (entities.Room, error) {
	obj := Room{}
	ent := entities.NewRoom()

	err := r.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Clauses(clause.Locking{
			Strength: "UPDATE",
			Table:    clause.Table{Name: clause.CurrentTable},
		}).
			Take(&obj, "id = ?", id).
			Error

		if err != nil {
			return Wrap(err, "db select lock for update")
		}

		obj.Into(&ent)
		consumerIDs, err := r.setConsumerIDs(tx, &ent)
		if err != nil {
			return Wrap(err, "set consumer ids")
		}

		if err = updateFn(&ent); err != nil {
			return Wrap(err, "updateFn")
		}

		obj.From(&ent)
		if err = tx.Save(obj).Error; err != nil {
			return Wrap(err, "db save")
		}

		roomID := id
		objRelsCreated := make(map[string]struct{}, len(ent.ConsumerIDs()))
		objRelsDeleted := make(map[string]struct{}, len(consumerIDs))

		for _, id := range consumerIDs {
			objRelsDeleted[id] = struct{}{}
		}

		for _, id := range ent.ConsumerIDs() {
			if _, exists := objRelsDeleted[id]; exists {
				delete(objRelsDeleted, id)
			} else {
				objRelsCreated[id] = struct{}{}
			}
		}

		type Rel struct {
			ConsumerID string
			RoomID     string
		}

		objs := make([]Rel, 0, len(objRelsDeleted))
		for id := range objRelsDeleted {
			objs = append(objs, Rel{
				RoomID:     roomID,
				ConsumerID: id,
			})
		}

		if err = tx.Table(consumersInRoomTable).Delete(&objs).Error; err != nil {
			return Wrap(err, `tx delete rels`)
		}

		objs = make([]Rel, 0, len(objRelsCreated))
		for id := range objRelsCreated {
			objs = append(objs, Rel{
				RoomID:     roomID,
				ConsumerID: id,
			})
		}

		if err = tx.Table(consumersInRoomTable).Create(&objs).Error; err != nil {
			return Wrap(err, `tx create rels`)
		}

		return nil
	})

	if err != nil {
		return ent, Wrap(err, "db transaction")
	}

	obj.Into(&ent)
	return ent, nil
}
