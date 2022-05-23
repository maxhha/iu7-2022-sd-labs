package models

import (
	"iu7-2022-sd-labs/buisness/entities"
	"iu7-2022-sd-labs/buisness/ports/interactors"
	"iu7-2022-sd-labs/buisness/ports/repositories"

	"github.com/shopspring/decimal"
)

type BidStepRow struct {
	FromAmount decimal.Decimal `json:"fromAmount"`
	Step       decimal.Decimal `json:"step"`
}

type BidStepTable struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	OrganizerID string
	Rows        []BidStepRow `json:"rows"`
}

func (obj *BidStepTable) From(ent *entities.BidStepTable) *BidStepTable {
	if ent == nil {
		return nil
	}

	obj.ID = ent.ID()
	obj.Name = ent.Name()
	obj.OrganizerID = ent.OrganizerID()
	obj.Rows = BidStepRowsArrayFromEntites(ent.Rows())

	return obj
}

func (obj *BidStepRow) From(ent *entities.BidStepRow) *BidStepRow {
	if ent == nil {
		return nil
	}

	obj.FromAmount = ent.FromAmount()
	obj.Step = ent.Step()

	return obj
}

func (obj *BidStepTableFilter) Into(ent *repositories.BidStepTableFilter) *repositories.BidStepTableFilter {
	if obj == nil {
		return nil
	}

	ent.IDs = obj.Ids
	if obj.Name == nil {
		ent.NameQuery = ""
	} else {
		ent.NameQuery = *obj.Name
	}
	ent.OrganizerIDs = obj.Organizers

	return ent
}

func BidStepTableEdgesArrayFromEntites(ents []entities.BidStepTable) []BidStepTableConnectionEdge {
	objs := make([]BidStepTableConnectionEdge, 0, len(ents))
	for _, ent := range ents {
		obj := BidStepTable{}
		objs = append(objs, BidStepTableConnectionEdge{
			Cursor: ent.ID(),
			Node:   obj.From(&ent),
		})
	}
	return objs
}

func BidStepRowsArrayFromEntites(ents []entities.BidStepRow) []BidStepRow {
	objs := make([]BidStepRow, 0, len(ents))
	for _, ent := range ents {
		obj := BidStepRow{}
		obj.From(&ent)
		objs = append(objs, obj)
	}
	return objs
}

func BidStepRowInputsArrayIntoInteractorRows(objs []BidStepRowInput) []interactors.BidStepRow {
	rows := make([]interactors.BidStepRow, 0, len(objs))
	for _, obj := range objs {
		rows = append(rows, interactors.BidStepRow{
			FromAmount: obj.FromAmount,
			Step:       obj.Step,
		})
	}
	return rows
}
