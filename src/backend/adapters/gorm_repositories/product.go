package gorm_repositories

//go:generate go run ../../codegen/gorm_repository/main.go --out product_gen.go --entity Product --methods Get,ShareLock,orderQuery,sliceQuery,Find,Create,Update(updateFn),Delete

import (
	"iu7-2022-sd-labs/buisness/entities"
	"iu7-2022-sd-labs/buisness/ports/repositories"
	"time"

	"gorm.io/gorm"
)

var productFieldToColumn = map[repositories.ProductOrderField]string{
	repositories.ProductOrderFieldCreationDate: "created_at",
}

type Product struct {
	ID          string `gorm:"<-:false;default:generated()"`
	OrganizerID string
	Name        string
	CreatedAt   time.Time `gorm:"<-:create"`
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt
}

func (obj *Product) From(e *entities.Product) *Product {
	if e == nil {
		return nil
	}

	obj.ID = e.ID()
	obj.OrganizerID = e.OrganizerID()
	obj.Name = e.Name()

	return obj
}

func (obj *Product) Into(e *entities.Product) *entities.Product {
	if e == nil {
		return nil
	}

	e.SetID(obj.ID)
	e.SetOrganizerID(obj.OrganizerID)
	e.SetName(obj.Name)

	return e
}

func (r *ProductRepository) filterQuery(query *gorm.DB, filter *repositories.ProductFilter) (*gorm.DB, error) {
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

	return query, nil
}
