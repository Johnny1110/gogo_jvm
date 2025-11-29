package heap

// Object runtime Object
// includes:
// - class: from which class
// - fields: object's fields
type Object struct {
	class  *Class
	fields Slots
}

// NewObject create new object
// TODO: completed in Phase-2.4
func NewObject(class *Class) *Object {
	return &Object{
		class:  class,
		fields: NewSlots(class.instanceSlotCount),
	}
}

// Class get class
func (o *Object) Class() *Class {
	return o.class
}

// Fields get fields
func (o *Object) Fields() Slots {
	return o.fields
}
