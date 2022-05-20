// This file was generated by codegen/gorm_repository. DO NOT EDIT.
package gorm_repositories

import (
	"iu7-2022-sd-labs/buisness/ports/repositories"

	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"iu7-2022-sd-labs/buisness/entities"
)

type OrganizerRepository struct{ db *gorm.DB }

func (r *GORMRepository) Organizer() repositories.OrganizerRepository {
	return &OrganizerRepository{r.db}
}

func (r *OrganizerRepository) Get(id string) (entities.Organizer, error) {
	obj := Organizer{}
	ent := entities.NewOrganizer()

	if err := r.db.Take(&obj, "id = ?", id).Error; err != nil {
		return ent, Wrap(err, "db take")
	}

	obj.Into(&ent)
	return ent, nil
}

func (r *OrganizerRepository) orderQuery(query *gorm.DB, order *repositories.OrganizerOrder) (*gorm.DB, error) {
	if order == nil {
		return query, nil
	}

	column, exists := organizerFieldToColumn[order.By]
	if !exists {
		return nil, fmt.Errorf("column for field \"%s\" is unknown", order.By)
	}

	query = query.Order(clause.OrderByColumn{
		Column: clause.Column{Name: column},
		Desc:   order.Desc,
	})

	return query, nil
}

func (r *OrganizerRepository) sliceQuery(query *gorm.DB, order *repositories.OrganizerOrder, slice *repositories.ForwardSlice) (*gorm.DB, error) {
	if slice == nil {
		return query, nil
	}

	column, exists := organizerFieldToColumn[order.By]
	if !exists {
		return nil, fmt.Errorf("column for field \"%s\" is unknown", order.By)
	}

	var err error
	query, err = sliceQuery(query, column, order.Desc, slice)
	return query, Wrap(err, "sliceQuery")
}

func (r *OrganizerRepository) Find(params *repositories.OrganizerFindParams) ([]entities.Organizer, error) {
	query := r.db.Model(&Organizer{})

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

	var objs []Organizer
	if err := query.Find(&objs).Error; err != nil {
		return nil, Wrap(err, "db find")
	}

	ents := make([]entities.Organizer, 0, len(objs))

	for _, obj := range objs {
		ent := entities.NewOrganizer()
		obj.Into(&ent)
		ents = append(ents, ent)
	}

	return ents, nil
}

func (r *OrganizerRepository) Create(ent *entities.Organizer) error {
	obj := Organizer{}
	obj.From(ent)

	if err := r.db.Create(&obj).Error; err != nil {
		return Wrap(err, "db create")
	}

	obj.Into(ent)
	return nil
}

func (r *OrganizerRepository) Update(ent *entities.Organizer) error {
	obj := Organizer{}
	obj.From(ent)

	if err := r.db.Save(&obj).Error; err != nil {
		return Wrap(err, "db save")
	}

	obj.Into(ent)
	return nil
}