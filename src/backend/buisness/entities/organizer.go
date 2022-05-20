package entities

type Organizer struct {
	id   string
	name string
}

func NewOrganizer() Organizer {
	return Organizer{}
}

func NewOrganizerPtr() *Organizer {
	obj := NewOrganizer()
	return &obj
}

func (obj *Organizer) ID() string {
	return obj.id
}

func (obj *Organizer) SetID(id string) *Organizer {
	obj.id = id
	return obj
}

func (obj *Organizer) Name() string {
	return obj.name
}

func (obj *Organizer) SetName(name string) *Organizer {
	obj.name = name
	return obj
}
