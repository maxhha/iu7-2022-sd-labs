package interactors

import (
	"iu7-2022-sd-labs/buisness/entities"
	"iu7-2022-sd-labs/buisness/ports/repositories"

	"github.com/shopspring/decimal"
)

type BidStepRow struct {
	FromAmount decimal.Decimal
	Step       decimal.Decimal
}

type BidStepTableCreateParams struct {
	OrganizerID string
	Name        string
	Rows        []BidStepRow
}

type BidStepTableUpdateParams struct {
	ID   string
	Name string
	Rows []BidStepRow
}

type BidStepTableInteractor interface {
	Create(params *BidStepTableCreateParams) (entities.BidStepTable, error)
	Find(params *repositories.BidStepTableFindParams) ([]entities.BidStepTable, error)
	Update(params *BidStepTableUpdateParams) (entities.BidStepTable, error)
}
