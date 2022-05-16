package models

import "iu7-2022-sd-labs/buisness/entities"

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
