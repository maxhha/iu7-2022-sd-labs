package entities

import "errors"

var ErrIsNil = errors.New("is nil")
var ErrIsEmpty = errors.New("is empty")

var ErrMustStartFromZeroAmount = errors.New("must start from zero amount")
var ErrRowsCollision = errors.New("rows collision")
var ErrNewBidIsLessOrEqualPrevious = errors.New("new bid is less or equal previous bid")
var ErrBidStepIsLessThenTable = errors.New("bid step is less then table")

var ErrAlreadyExists = errors.New("already exists")
var ErrNotFound = errors.New("not found")
