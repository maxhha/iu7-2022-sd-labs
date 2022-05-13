package interactors

import (
	"iu7-2022-sd-labs/buisness/entities"
	"iu7-2022-sd-labs/buisness/errors"
	"iu7-2022-sd-labs/buisness/ports/interactors"
	"iu7-2022-sd-labs/buisness/ports/repositories"
)

type BidStepTableInteractor struct {
	repo repositories.Repository
}

func NewBidStepTableInteractor(
	repo repositories.Repository,
) BidStepTableInteractor {
	return BidStepTableInteractor{repo}
}

func (interactor *BidStepTableInteractor) bidStepRowsToEntities(
	paramRows []interactors.BidStepRow,
) []entities.BidStepRow {
	rows := make([]entities.BidStepRow, 0, len(paramRows))

	for _, r := range paramRows {
		row := entities.NewBidStepRow()
		row.SetFromAmount(r.FromAmount).SetStep(r.Step)
		rows = append(rows, row)
	}

	return rows
}

func (interactor *BidStepTableInteractor) Create(
	params *interactors.BidStepTableCreateParams,
) (entities.BidStepTable, error) {
	org, err := interactor.repo.Organizer().Get(params.OrganizerID)
	if err != nil {
		return entities.BidStepTable{}, errors.Wrap(err, "organizer repo get")
	}
	rows := interactor.bidStepRowsToEntities(params.Rows)
	table := entities.NewBidStepTable()
	table.
		SetOrganizerID(org.ID()).
		SetName(params.Name).
		SetRows(rows)

	if err = table.Validate(); err != nil {
		return table, errors.Wrap(err, "validate")
	}

	err = interactor.repo.BidStepTable().Create(&table)
	return table, errors.Wrap(err, "table repo create")
}

func (interactor *BidStepTableInteractor) Find(
	params *repositories.BidStepTableFindParams,
) ([]entities.BidStepTable, error) {
	tables, err := interactor.repo.BidStepTable().Find(params)
	return tables, errors.Wrap(err, "table repo find")
}

func (interactor *BidStepTableInteractor) Update(
	params *interactors.BidStepTableUpdateParams,
) (entities.BidStepTable, error) {
	table, err := interactor.repo.BidStepTable().Get(params.ID)
	if err != nil {
		return table, errors.Wrap(err, "table repo get")
	}

	rows := interactor.bidStepRowsToEntities(params.Rows)
	table.
		SetName(params.Name).
		SetRows(rows)

	if err = table.Validate(); err != nil {
		return table, errors.Wrap(err, "validate")
	}

	err = interactor.repo.BidStepTable().Update(&table)

	return table, errors.Wrap(err, "table repo update")
}
