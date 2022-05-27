package gorm_repositories

//go:generate go run ../../codegen/gorm_repository/main.go --out block_list_gen.go --entity BlockList --methods orderQuery

import (
	"fmt"
	"iu7-2022-sd-labs/buisness/entities"
	"iu7-2022-sd-labs/buisness/ports/repositories"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type BlockedConsumer struct {
	OrganizerID string    `gorm:"<-:create"`
	ConsumerID  string    `gorm:"<-:create"`
	CreatedAt   time.Time `gorm:"<-:create"`
}

func (c *BlockedConsumer) key() string {
	return c.OrganizerID
}

var blockListFieldToColumn = map[repositories.BlockListOrderField]string{
	repositories.BlockListOrderFieldCreationDate: "created_at",
}

func BlockedConsumersToEntities(objs []BlockedConsumer) []entities.BlockList {
	entIndexes := make(map[string]int)
	ents := make([]entities.BlockList, 0)
	for _, obj := range objs {
		var ent entities.BlockList
		entIndex, exists := entIndexes[obj.OrganizerID]

		if exists {
			ent = ents[entIndex]
		} else {
			entIndex = len(ents)
			entIndexes[obj.OrganizerID] = entIndex
			ent = entities.NewBlockList()
			ents = append(ents, ent)
		}

		ent.SetConsumerIDs(append(ent.ConsumerIDs(), obj.ConsumerID))
		ents[entIndex] = ent
	}

	return ents
}

func BlockedConsumersFromEntities(ents []entities.BlockList) []BlockedConsumer {
	objs := make([]BlockedConsumer, 0)
	for _, ent := range ents {
		organizerID := ent.OrganizerID()
		for _, consumerID := range ent.ConsumerIDs() {
			objs = append(objs, BlockedConsumer{
				OrganizerID: organizerID,
				ConsumerID:  consumerID,
			})
		}
	}

	return objs
}

func (r *BlockListRepository) filterQuery(query *gorm.DB, filter *repositories.BlockListFilter) (*gorm.DB, error) {
	if filter == nil {
		return query, nil
	}

	if len(filter.OrganizerIDs) > 0 {
		query = query.Where("organizer_id in ?", filter.OrganizerIDs)
	}

	if len(filter.ConsumerIDs) > 0 {
		consumerQuery := query.Session(&gorm.Session{
			Initialized: true,
		}).
			Model(query.Statement.Model).
			Select("organizer_id").
			Where("consumer_id in ?", filter.ConsumerIDs)

		query = query.Where("organizer_id in ?", consumerQuery)
	}

	return query, nil
}

func (r *BlockListRepository) sliceQuery(query *gorm.DB, order *repositories.BlockListOrder, slice *repositories.ForwardSlice) (*gorm.DB, error) {
	if slice == nil {
		return query, nil
	}

	orderByColumn, exists := blockListFieldToColumn[order.By]
	if !exists {
		return nil, fmt.Errorf("column for field \"%s\" is unknown", order.By)
	}

	if slice.AfterID != "" {
		if orderByColumn == "" {
			orderByColumn = "created_at"
		}

		compareSymbol := ">"
		if order.Desc {
			compareSymbol = "<"
		}

		afterQuery := query.
			Session(&gorm.Session{Initialized: true}).
			Model(query.Statement.Model).
			Where("organizer_id = ?", slice.AfterID).
			Select(orderByColumn)

		whereClause := fmt.Sprintf("%s %s ANY( ? )", orderByColumn, compareSymbol)
		query = query.Where(whereClause, afterQuery)
	}

	if slice.Limit > 0 {
		descClause := ""
		if order.Desc {
			descClause = " DESC"
		}
		selectClause := fmt.Sprintf(
			"*, ROW_NUMBER() OVER(PARTITION BY organizer_id ORDER BY %s%s) as black_list_index",
			orderByColumn,
			descClause,
		)
		blackListsIndexed := query.Select(selectClause)
		query = r.db.Table("( ? ) bli", blackListsIndexed).Where("black_list_index <= ?", slice.Limit)
	}

	return query, nil
}

func (r *BlockListRepository) Find(params *repositories.BlockListFindParams) ([]entities.BlockList, error) {
	query := r.db.Model(&BlockedConsumer{})

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

	var objs []BlockedConsumer
	if err := query.Find(&objs).Error; err != nil {
		return nil, Wrap(err, "db find")
	}

	ents := BlockedConsumersToEntities(objs)

	return ents, nil
}

func (r *BlockListRepository) UpdateOrCreateByOrganizer(organizerID string, updateFn func(ent *entities.BlockList) error) (entities.BlockList, error) {
	var objs []BlockedConsumer
	ent := entities.NewBlockList()

	err := r.db.Transaction(func(tx *gorm.DB) error {
		err := tx.Clauses(clause.Locking{
			Strength: "UPDATE",
			Table:    clause.Table{Name: clause.CurrentTable},
		}).
			Take(&objs, "organizer_id = ?", organizerID).
			Error

		if err != nil {
			return Wrap(err, "db select lock for update")
		}

		ents := BlockedConsumersToEntities(objs)

		if len(ents) == 1 {
			ent = ents[0]
		} else if len(ents) > 1 {
			return fmt.Errorf("BlockedConsumersToEntities returned many entities")
		}

		if err = updateFn(&ent); err != nil {
			return Wrap(err, "updateFn")
		}

		newObjs := BlockedConsumersFromEntities([]entities.BlockList{ent})

		objsCreated := make(map[string]BlockedConsumer, len(newObjs))
		objsDeleted := make(map[string]BlockedConsumer, len(objs))

		for _, obj := range objs {
			objsDeleted[obj.key()] = obj
		}

		for _, obj := range newObjs {
			if _, exists := objsDeleted[obj.key()]; exists {
				delete(objsDeleted, obj.key())
			} else {
				objsCreated[obj.key()] = obj
			}
		}

		if len(objsDeleted) > 0 {
			objs := make([]BlockedConsumer, 0, len(objsDeleted))
			for _, obj := range objsDeleted {
				objs = append(objs, obj)
			}

			if err = tx.Delete(&objs).Error; err != nil {
				return Wrap(err, `tx delete rows`)
			}
		}

		if len(objsCreated) > 0 {
			objs := make([]BlockedConsumer, 0, len(objsCreated))
			for _, obj := range objsCreated {
				objs = append(objs, obj)
			}

			if err = tx.Create(&objs).Error; err != nil {
				return Wrap(err, `tx create rows`)
			}
		}

		objs = newObjs

		return nil
	})

	if err != nil {
		return ent, Wrap(err, "db transaction")
	}

	return ent, nil
}
