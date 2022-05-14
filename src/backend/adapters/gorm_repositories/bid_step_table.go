package gorm_repositories

//go:generate go run ../../codegen/gorm_repository/main.go --out bid_step_table_gen.go --entity BidStepTable --methods orderQuery,sliceQuery

import (
	"iu7-2022-sd-labs/buisness/entities"
	"iu7-2022-sd-labs/buisness/errors"
	"iu7-2022-sd-labs/buisness/ports/repositories"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var bidStepTableFieldToColumn = map[repositories.BidStepTableOrderField]string{
	repositories.BidStepTableOrderFieldCreationDate: "created_at",
}

type BidStepTable struct {
	ID          string `gorm:"<-:false;default:generated()"`
	Name        string
	OrganizerID string
	CreatedAt   time.Time `gorm:"<-:create"`
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt
}

type BidStepRow struct {
	TableID    string          `gorm:"<-:create;primaryKey"`
	FromAmount decimal.Decimal `gorm:"<-:create;primaryKey"`
	Step       decimal.Decimal
}

func (obj *BidStepTable) From(e *entities.BidStepTable) *BidStepTable {
	if e == nil {
		return nil
	}

	obj.ID = e.ID()
	obj.Name = e.Name()
	obj.OrganizerID = e.OrganizerID()

	return obj
}

func (obj *BidStepTable) Into(e *entities.BidStepTable) *entities.BidStepTable {
	if e == nil {
		return nil
	}

	e.SetID(obj.ID)
	e.SetName(obj.Name)
	e.SetOrganizerID(obj.OrganizerID)

	return e
}

func (obj *BidStepRow) From(e *entities.BidStepRow) *BidStepRow {
	if e == nil {
		return nil
	}

	obj.FromAmount = e.FromAmount()
	obj.Step = e.Step()

	return obj
}

func (obj *BidStepRow) Into(e *entities.BidStepRow) *entities.BidStepRow {
	if e == nil {
		return nil
	}

	e.SetFromAmount(obj.FromAmount)
	e.SetStep(obj.Step)

	return e
}

func (obj *BidStepRow) Equal(other *BidStepRow) bool {
	return obj.TableID == other.TableID &&
		obj.FromAmount.Equal(other.FromAmount) &&
		obj.Step.Equal(other.Step)
}

func (obj *BidStepRow) key() string {
	return obj.FromAmount.StringFixed(2)
}

func bidStepRowsFromEntity(table entities.BidStepTable) []BidStepRow {
	objs := make([]BidStepRow, 0, len(table.Rows()))
	for _, ent := range table.Rows() {
		obj := BidStepRow{TableID: table.ID()}
		obj.From(&ent)
		objs = append(objs, obj)
	}

	return objs
}

func bidStepRowsIntoEntity(objs []BidStepRow) []entities.BidStepRow {
	rowEnts := make([]entities.BidStepRow, 0, len(objs))
	for _, obj := range objs {
		ent := entities.NewBidStepRow()
		obj.Into(&ent)
		rowEnts = append(rowEnts, ent)
	}

	return rowEnts
}

func (r *BidStepTableRepository) setRows(db *gorm.DB, ent *entities.BidStepTable) ([]BidStepRow, error) {
	var rowObjs []BidStepRow
	if err := db.Find(&rowObjs, "table_id = ?", ent.ID()).Error; err != nil {
		return nil, Wrap(err, "db take rows")
	}

	ent.SetRows(bidStepRowsIntoEntity(rowObjs))
	return rowObjs, nil
}

func (r *BidStepTableRepository) Get(id string) (entities.BidStepTable, error) {
	ent := entities.NewBidStepTable()
	obj := BidStepTable{}
	if err := r.db.Take(&obj, "id = ?", id).Error; err != nil {
		return ent, Wrap(err, "db take")
	}

	obj.Into(&ent)
	_, err := r.setRows(r.db, &ent)

	return ent, Wrap(err, "set rows")
}

func (r *BidStepTableRepository) filterQuery(query *gorm.DB, filter *repositories.BidStepTableFilter) (*gorm.DB, error) {
	if filter == nil {
		return query, nil
	}

	if len(filter.IDs) > 0 {
		query = query.Where("id in ?", filter.IDs)
	}

	if len(filter.NameQuery) > 0 {
		query = query.Where("name ~* ?", filter.NameQuery)
	}

	return query, nil
}

func (r *BidStepTableRepository) Find(params *repositories.BidStepTableFindParams) ([]entities.BidStepTable, error) {
	query := r.db.Model(&BidStepTable{})

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

	var objs []BidStepTable
	if err := query.Find(&objs).Error; err != nil {
		return nil, Wrap(err, "db find")
	}

	ents := make([]entities.BidStepTable, 0, len(objs))

	for _, obj := range objs {
		ent := entities.NewBidStepTable()
		obj.Into(&ent)

		if _, err := r.setRows(r.db, &ent); err != nil {
			return nil, Wrapf(err, `set row for table_id="%s"`, ent.ID())
		}

		ents = append(ents, ent)
	}

	return ents, nil
}

func (r *BidStepTableRepository) Create(ent *entities.BidStepTable) error {
	if err := ent.Validate(); err != nil {
		return Wrap(err, "validate")
	}

	obj := BidStepTable{}
	obj.From(ent)

	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&obj).Error; err != nil {
			return Wrap(err, "db create")
		}

		tableID := obj.ID

		for _, entRow := range ent.Rows() {
			obj := BidStepRow{TableID: tableID}

			if err := tx.Create(obj.From(&entRow)).Error; err != nil {
				return Wrapf(
					err,
					"db create row from_amount=%s",
					obj.FromAmount.StringFixed(2),
				)
			}
		}

		return nil
	})

	if err != nil {
		return Wrap(err, "db transaction")
	}

	obj.Into(ent)
	return nil
}

func (r *BidStepTableRepository) Update(id string, updateFn func(ent *entities.BidStepTable) error) (entities.BidStepTable, error) {
	obj := BidStepTable{}
	ent := entities.NewBidStepTable()

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
		objRows, err := r.setRows(tx, &ent)
		if err != nil {
			return Wrap(err, "set rows")
		}

		if err = updateFn(&ent); err != nil {
			return Wrap(err, "updateFn")
		}

		if err := ent.Validate(); err != nil {
			return errors.Wrap(err, "validate")
		}

		obj.From(&ent)
		if err = tx.Save(obj).Error; err != nil {
			return Wrap(err, "db save")
		}

		objRowsCreated := make(map[string]BidStepRow, len(ent.Rows()))
		objRowsDeleted := make(map[string]BidStepRow, len(objRows))
		objRowsUpdated := make(map[string]BidStepRow, len(objRows))

		for _, obj := range objRows {
			objRowsDeleted[obj.key()] = obj
		}

		for _, ent := range ent.Rows() {
			obj := BidStepRow{TableID: id}
			obj.From(&ent)

			if oldObj, exists := objRowsDeleted[obj.key()]; exists {
				delete(objRowsDeleted, obj.key())

				if !oldObj.Equal(&obj) {
					objRowsUpdated[obj.key()] = obj
				}
			} else {
				objRowsCreated[obj.key()] = obj
			}
		}

		objs := make([]BidStepRow, 0, len(objRowsDeleted))
		for _, obj := range objRowsDeleted {
			objs = append(objs, obj)
		}

		if err = tx.Delete(&objs).Error; err != nil {
			return Wrap(err, `tx delete rows`)
		}

		objs = make([]BidStepRow, 0, len(objRowsCreated))
		for _, obj := range objRowsCreated {
			objs = append(objs, obj)
		}

		if err = tx.Create(&objs).Error; err != nil {
			return Wrap(err, `tx create rows`)
		}

		objs = make([]BidStepRow, 0, len(objRowsUpdated))
		for _, obj := range objRowsUpdated {
			objs = append(objs, obj)
		}

		if err = tx.Save(&objs).Error; err != nil {
			return Wrap(err, `tx save rows`)
		}

		return nil
	})

	if err != nil {
		return ent, Wrap(err, "db transaction")
	}

	obj.Into(&ent)
	return ent, nil
}
