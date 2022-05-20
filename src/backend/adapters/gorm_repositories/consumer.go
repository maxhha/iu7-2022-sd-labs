package gorm_repositories

//go:generate go run ../../codegen/gorm_repository/main.go --out consumer_gen.go --entity Consumer --methods Get,orderQuery,sliceQuery,Find,Create,Update

import (
	"iu7-2022-sd-labs/buisness/entities"
	"iu7-2022-sd-labs/buisness/ports/repositories"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

var consumerFieldToColumn = map[repositories.ConsumerOrderField]string{
	repositories.ConsumerOrderFieldCreationDate: "created_at",
}

type Consumer struct {
	ID        string `gorm:"<-:false;default:generated()"`
	Nickname  string
	Form      datatypes.JSONMap
	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}

func (obj *Consumer) From(e *entities.Consumer) *Consumer {
	if e == nil {
		return nil
	}

	obj.ID = e.ID()
	obj.Nickname = e.Nickname()
	obj.Form = e.Form()

	return obj
}

func (obj *Consumer) Into(e *entities.Consumer) *entities.Consumer {
	if e == nil {
		return nil
	}

	e.SetID(obj.ID)
	e.SetNickname(obj.Nickname)
	e.SetForm(obj.Form)

	return e
}

func (r *ConsumerRepository) filterQuery(query *gorm.DB, filter *repositories.ConsumerFilter) (*gorm.DB, error) {
	if filter == nil {
		return query, nil
	}

	if len(filter.IDs) > 0 {
		query = query.Where("id in ?", filter.IDs)
	}

	if len(filter.NickameQuery) > 0 {
		query = query.Where("nickname ~* ?", filter.NickameQuery)
	}

	for _, fieldQuery := range filter.FormFieldQueries {
		query = query.Where(
			"form->>? ~* ?",
			fieldQuery.Field,
			fieldQuery.Query,
		)
	}

	return query, nil
}
