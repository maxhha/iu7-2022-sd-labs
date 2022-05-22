package models

import (
	"iu7-2022-sd-labs/buisness/entities"
	"iu7-2022-sd-labs/buisness/ports/repositories"
)

type Product struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	OrganizerID string
}

func (obj *Product) From(ent *entities.Product) *Product {
	if ent == nil {
		return nil
	}

	obj.ID = ent.ID()
	obj.Name = ent.Name()
	obj.OrganizerID = ent.OrganizerID()

	return obj
}

func (obj *ProductFilter) Into(ent *repositories.ProductFilter) *repositories.ProductFilter {
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

func ProductEdgesArrayFromEntites(ents []entities.Product) []ProductConnectionEdge {
	objs := make([]ProductConnectionEdge, 0, len(ents))
	for _, ent := range ents {
		obj := Product{}
		objs = append(objs, ProductConnectionEdge{
			Cursor: ent.ID(),
			Node:   obj.From(&ent),
		})
	}
	return objs
}
