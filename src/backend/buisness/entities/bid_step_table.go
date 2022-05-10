package entities

import (
	"fmt"
	"iu7-2022-sd-labs/buisness/errors"
	"sort"

	"github.com/hashicorp/go-multierror"
	"github.com/shopspring/decimal"
)

type BidStepTable struct {
	id          string
	name        string
	organizerID string
	rows        []BidStepRow
}

type BidStepRow struct {
	fromAmount decimal.Decimal
	step       decimal.Decimal
}

func NewBidStepTable() BidStepTable {
	return BidStepTable{}
}

func (obj *BidStepTable) ID() string {
	return obj.id
}

func (obj *BidStepTable) SetID(id string) *BidStepTable {
	obj.id = id
	return obj
}

func (obj *BidStepTable) Name() string {
	return obj.name
}

func (obj *BidStepTable) SetName(name string) *BidStepTable {
	obj.name = name
	return obj
}

func (obj *BidStepTable) OrganizerID() string {
	return obj.organizerID
}

func (obj *BidStepTable) SetOrganizerID(id string) *BidStepTable {
	obj.organizerID = id
	return obj
}

func (obj *BidStepTable) Rows() []BidStepRow {
	return obj.rows
}

func (obj *BidStepTable) SetRows(rows []BidStepRow) *BidStepTable {
	obj.rows = rows
	return obj
}

func (obj *BidStepTable) Validate() error {
	var errs error

	if obj.organizerID == "" {
		errs = multierror.Append(errs, errors.Wrap(ErrIsEmpty, "organizer id"))
	}

	if obj.name == "" {
		errs = multierror.Append(errs, errors.Wrap(ErrIsEmpty, "name id"))
	}

	if len(obj.rows) == 0 {
		errs = multierror.Append(errs, errors.Wrap(ErrIsEmpty, "rows"))
	} else {
		sort.SliceStable(obj.rows, func(i, j int) bool {
			return obj.rows[i].fromAmount.LessThan(obj.rows[j].fromAmount)
		})

		prev := obj.rows[0]
		if !prev.fromAmount.Equal(decimal.Zero) {
			errs = multierror.Append(errs, ErrMustStartFromZeroAmount)
		}

		for _, r := range obj.rows[1:] {
			if r.fromAmount.Equal(prev.fromAmount) {
				errs = multierror.Append(
					errs,
					fmt.Errorf(
						"%w: %s and %s",
						ErrRowsCollision,
						prev.String(),
						r.String(),
					),
				)
			}

			prev = r
		}
	}

	return errs
}

func (obj *BidStepTable) IsAllowedBid(prevBid decimal.Decimal, newBid decimal.Decimal) error {
	if newBid.LessThanOrEqual(prevBid) {
		return errors.Wrapf(
			ErrNewBidIsLessOrEqualPrevious,
			"prev=%s new=%s",
			prevBid.String(),
			newBid.String(),
		)
	}

	stepRow := obj.rows[0]
	for _, row := range obj.rows {
		if row.fromAmount.LessThan(prevBid) {
			break
		}
		stepRow = row
	}

	if newBid.Sub(prevBid).LessThan(stepRow.step) {
		return errors.Wrapf(
			ErrBidStepIsLessThenTable,
			"actual=%s table=%s",
			newBid.Sub(prevBid).String(),
			stepRow.step.String(),
		)
	}

	return nil
}

func NewBidStepRow() BidStepRow {
	return BidStepRow{}
}

func NewBidStepRowPtr() *BidStepRow {
	org := NewBidStepRow()
	return &org
}

func (obj *BidStepRow) String() string {
	return fmt.Sprintf(
		"BidStepRow[from=%s; step=%s]",
		obj.fromAmount.String(),
		obj.step.String(),
	)
}

func (obj *BidStepRow) FromAmount() decimal.Decimal {
	return obj.fromAmount
}

func (obj *BidStepRow) SetFromAmount(minAmount decimal.Decimal) *BidStepRow {
	obj.fromAmount = minAmount
	return obj
}

func (obj *BidStepRow) Step() decimal.Decimal {
	return obj.step
}

func (obj *BidStepRow) SetStep(step decimal.Decimal) *BidStepRow {
	obj.step = step
	return obj
}
