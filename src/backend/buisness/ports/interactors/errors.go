package interactors

import "errors"

var ErrOfferedAmountIsLessThanMinAmount = errors.New("offered amount is less than min amount")
var ErrOfferIsNotMax = errors.New("offer is not max")
