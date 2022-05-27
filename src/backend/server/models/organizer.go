package models

import (
	"iu7-2022-sd-labs/buisness/entities"
	"iu7-2022-sd-labs/buisness/ports/repositories"
)

type Organizer struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (Organizer) IsViewer() {}

type TokenResult struct {
	Token string `json:"token"`
}

func (obj *Organizer) From(ent *entities.Organizer) *Organizer {
	if ent == nil {
		return nil
	}

	obj.ID = ent.ID()
	obj.Name = ent.Name()

	return obj
}

func (obj *OrganizerFilter) Into(ent *repositories.OrganizerFilter) *repositories.OrganizerFilter {
	if obj == nil {
		return nil
	}

	ent.IDs = obj.Ids

	if obj.Name == nil {
		ent.NameQuery = ""
	} else {
		ent.NameQuery = *obj.Name
	}

	return ent
}

func OrganizerEdgesArrayFromEntites(ents []entities.Organizer) []OrganizerConnectionEdge {
	objs := make([]OrganizerConnectionEdge, 0, len(ents))
	for _, ent := range ents {
		obj := Organizer{}
		objs = append(objs, OrganizerConnectionEdge{
			Cursor: ent.ID(),
			Node:   obj.From(&ent),
		})
	}
	return objs
}
