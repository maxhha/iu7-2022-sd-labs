package models

import (
	"iu7-2022-sd-labs/buisness/entities"
	"iu7-2022-sd-labs/buisness/ports/repositories"
)

type Consumer struct {
	ID       string                 `json:"id"`
	Nickname string                 `json:"nickname"`
	Form     map[string]interface{} `json:"form"`
}

func (Consumer) IsViewer() {}

func (obj *Consumer) From(ent *entities.Consumer) *Consumer {
	if ent == nil {
		return nil
	}

	obj.ID = ent.ID()
	obj.Nickname = ent.Nickname()
	obj.Form = ent.Form()

	return obj
}

func (obj *ConsumerFilter) Into(ent *repositories.ConsumerFilter) *repositories.ConsumerFilter {
	if obj == nil {
		return nil
	}

	ent.IDs = obj.Ids
	if obj.Nickname == nil {
		ent.NickameQuery = ""
	} else {
		ent.NickameQuery = *obj.Nickname
	}

	ent.FormFieldQueries = nil
	for field, query := range obj.Form {
		ent.FormFieldQueries = append(ent.FormFieldQueries, repositories.ConsumerFormFieldQuery{
			Field: field,
			Query: query,
		})
	}

	return ent
}

func ConsumerEdgesArrayFromEntites(ents []entities.Consumer) []ConsumerConnectionEdge {
	objs := make([]ConsumerConnectionEdge, 0, len(ents))
	for _, ent := range ents {
		obj := Consumer{}
		objs = append(objs, ConsumerConnectionEdge{
			Cursor: ent.ID(),
			Node:   obj.From(&ent),
		})
	}
	return objs
}

func ConsumerArrayFromEntites(ents []entities.Consumer) []Consumer {
	objs := make([]Consumer, 0, len(ents))
	for _, ent := range ents {
		obj := Consumer{}
		obj.From(&ent)
		objs = append(objs, obj)
	}
	return objs
}
