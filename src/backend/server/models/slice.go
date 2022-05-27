package models

import "iu7-2022-sd-labs/buisness/ports/repositories"

func FillForwardSlice(first *int, after *string, ent *repositories.ForwardSlice) *repositories.ForwardSlice {
	if first == nil && after == nil {
		return nil
	}
	if first != nil {
		ent.Limit = *first + 1
	}
	if after != nil {
		ent.AfterID = *after
	}
	return ent
}
