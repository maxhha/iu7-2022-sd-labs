// This file was generated by codegen/gorm_repository. DO NOT EDIT.
package gorm_repositories

import (
	"iu7-2022-sd-labs/buisness/ports/repositories"

	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"iu7-2022-sd-labs/buisness/entities"
)

type ProductRepository struct{ db *gorm.DB }

func (r *GORMRepository) Product() repositories.ProductRepository {
	return &ProductRepository{r.db}
}

func (r *ProductRepository) Get(id string) (entities.Product, error) {
	obj := Product{}
	ent := entities.NewProduct()

	if err := r.db.Take(&obj, "id = ?", id).Error; err != nil {
		return ent, Wrap(err, "db take")
	}

	obj.Into(&ent)
	return ent, nil
}

func (r *ProductRepository) ShareLock(id string) (entities.Product, error) {
	obj := Product{}
	ent := entities.NewProduct()

	err := r.db.Clauses(clause.Locking{
		Strength: "SHARE",
		Table:    clause.Table{Name: clause.CurrentTable},
	}).
		Take(&obj, "id = ?", id).
		Error
	if err != nil {
		return ent, Wrap(err, "db select for update")
	}

	obj.Into(&ent)
	return ent, nil
}

func (r *ProductRepository) orderQuery(query *gorm.DB, order *repositories.ProductOrder) (*gorm.DB, error) {
	if order == nil {
		return query, nil
	}

	column, exists := productFieldToColumn[order.By]
	if !exists {
		return nil, fmt.Errorf("column for field \"%s\" is unknown", order.By)
	}

	query = query.Order(clause.OrderByColumn{
		Column: clause.Column{Name: column},
		Desc:   order.Desc,
	})

	return query, nil
}

func (r *ProductRepository) sliceQuery(query *gorm.DB, order *repositories.ProductOrder, slice *repositories.ForwardSlice) (*gorm.DB, error) {
	if slice == nil {
		return query, nil
	}

	column, exists := productFieldToColumn[order.By]
	if !exists {
		return nil, fmt.Errorf("column for field \"%s\" is unknown", order.By)
	}

	var err error
	query, err = sliceQuery(query, column, order.Desc, slice)
	return query, Wrap(err, "sliceQuery")
}

func (r *ProductRepository) Find(params *repositories.ProductFindParams) ([]entities.Product, error) {
	query := r.db.Model(&Product{})

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

	var objs []Product
	if err := query.Find(&objs).Error; err != nil {
		return nil, Wrap(err, "db find")
	}

	ents := make([]entities.Product, 0, len(objs))

	for _, obj := range objs {
		ent := entities.NewProduct()
		obj.Into(&ent)
		ents = append(ents, ent)
	}

	return ents, nil
}

func (r *ProductRepository) Create(ent *entities.Product) error {
	obj := Product{}
	obj.From(ent)

	if err := r.db.Create(&obj).Error; err != nil {
		return Wrap(err, "db create")
	}

	obj.Into(ent)
	return nil
}

func (r *ProductRepository) Update(id string, updateFn func(ent *entities.Product) error) (entities.Product, error) {
	obj := Product{}
	ent := entities.NewProduct()

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

		if err = updateFn(obj.Into(&ent)); err != nil {
			return Wrap(err, "updateFn")
		}

		obj.From(&ent)

		err = tx.Save(obj).Error
		return Wrap(err, "db save")
	})

	if err != nil {
		return ent, Wrap(err, "db transaction")
	}

	obj.Into(&ent)
	return ent, nil
}

func (r *ProductRepository) Delete(id string) (entities.Product, error) {
	obj := Product{}
	ent := entities.NewProduct()

	err := r.db.Transaction(func(tx *gorm.DB) error {
		err := r.db.Clauses(clause.Locking{
			Strength: "UPDATE",
			Table:    clause.Table{Name: clause.CurrentTable},
		}).
			Take(&obj, "id = ?", id).
			Error

		if err != nil {
			return Wrap(err, "db select for update")
		}

		if err := r.db.Delete(&obj).Error; err != nil {
			return Wrap(err, "db delete")
		}

		return nil
	})

	if err != nil {
		return ent, Wrap(err, "db transaction")
	}

	obj.Into(&ent)
	return ent, nil
}
