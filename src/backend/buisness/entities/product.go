package entities

type Product struct {
	id          string
	organizerID string
	name        string
}

func NewProduct() Product {
	return Product{}
}

func (obj *Product) ID() string {
	return obj.id
}

func (obj *Product) SetID(id string) *Product {
	obj.id = id
	return obj
}

func (obj *Product) OrganizerID() string {
	return obj.organizerID
}

func (obj *Product) SetOrganizerID(id string) *Product {
	obj.organizerID = id
	return obj
}

func (obj *Product) Name() string {
	return obj.name
}

func (obj *Product) SetName(name string) *Product {
	obj.name = name
	return obj
}
