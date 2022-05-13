package gorm_repositories

import (
	"fmt"
	"iu7-2022-sd-labs/buisness/ports/repositories"

	"gorm.io/gorm"
)

func sliceQuery(query *gorm.DB, orderByColumn string, orderDesc bool, slice *repositories.ForwardSlice) (*gorm.DB, error) {
	if slice == nil {
		return query, nil
	}

	if slice.AfterID != "" {
		if orderByColumn == "" {
			orderByColumn = "id"
		}

		compareSymbol := ">"
		if orderDesc {
			compareSymbol = "<"
		}

		afterQuery := query.
			Session(&gorm.Session{Initialized: true}).
			Model(query.Statement.Model).
			Where("id = ?", slice.AfterID).
			Select(orderByColumn)

		whereClause := fmt.Sprintf("%s %s ANY( ? )", orderByColumn, compareSymbol)
		query = query.Where(whereClause, afterQuery)
	}

	if slice.Limit > 0 {
		query = query.Limit(slice.Limit)
	}

	return query, nil
}
