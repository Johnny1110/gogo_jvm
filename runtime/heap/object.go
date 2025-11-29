package heap

import (
	"github.com/Johnny1110/gogo_jvm/runtime/method_area"
)

// Object runtime Object
// includes:
// - class: from which class
// - fields: object's fields
type Object struct {
	class  *method_area.Class
	fields []*interface{}
}

// NewObject create new object
// TODO: completed in Phase-2.4
func NewObject(class *method_area.Class) *Object {
	return &Object{
		class:  class,
		fields: make([]*interface{}, class.InstanceSlotCount()),
	}
}

// Class get class
func (o *Object) Class() *method_area.Class {
	return o.class
}

// Fields get fields
func (o *Object) Fields() []*interface{} {
	return o.fields
}
