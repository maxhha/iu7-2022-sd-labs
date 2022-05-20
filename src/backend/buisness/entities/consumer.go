package entities

type Consumer struct {
	id       string
	nickname string
	form     map[string]interface{}
}

func NewConsumer() Consumer {
	return Consumer{}
}

func NewConsumerPtr() *Consumer {
	obj := NewConsumer()
	return &obj
}

func (obj *Consumer) ID() string {
	return obj.id
}

func (obj *Consumer) SetID(id string) *Consumer {
	obj.id = id
	return obj
}

func (obj *Consumer) Nickname() string {
	return obj.nickname
}

func (obj *Consumer) SetNickname(nickname string) *Consumer {
	obj.nickname = nickname
	return obj
}

func (obj *Consumer) Form() map[string]interface{} {
	return obj.form
}

func (obj *Consumer) SetForm(form map[string]interface{}) *Consumer {
	obj.form = form
	return obj
}
