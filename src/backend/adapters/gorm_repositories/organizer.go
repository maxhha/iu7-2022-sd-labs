package gorm_repositories

//go:generate go run ../../codegen/gorm_repository/main.go --out organizer_gen.go --entity Organizer --methods Get,orderQuery,sliceQuery,Find,Create,Update

import (
	"iu7-2022-sd-labs/buisness/entities"
	"iu7-2022-sd-labs/buisness/ports/repositories"
	"time"

	"gorm.io/gorm"
)

var organizerFieldToColumn = map[repositories.OrganizerOrderField]string{
	repositories.OrganizerOrderFieldCreationDate: "created_at",
}

type Organizer struct {
	ID        string `gorm:"<-:false;default:generated()"`
	Name      string
	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

func (obj *Organizer) From(e *entities.Organizer) *Organizer {
	if e == nil {
		return nil
	}

	obj.ID = e.ID()
	obj.Name = e.Name()

	return obj
}

func (obj *Organizer) Into(e *entities.Organizer) *entities.Organizer {
	if e == nil {
		return nil
	}

	e.SetID(obj.ID)
	e.SetName(obj.Name)

	return e
}

func (r *OrganizerRepository) filterQuery(query *gorm.DB, filter *repositories.OrganizerFilter) (*gorm.DB, error) {
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
