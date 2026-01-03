package heap

import (
	"fmt"
	"github.com/Johnny1110/gogo_jvm/runtime/rtcore"
)

// Object runtime Object instance
// fields using rtcore.Slots follow LocalVars, OperandStack
//   - every field's slotId are prepared in stage: ClassLoader.prepare()
//   - support extends: parent's fields at front, this class's fields concat by following
//
// ┌─────────────────────────────────────────┐
// │ Object                                  │
// ├─────────────────────────────────────────┤
// │ class ────────────────────────► Class   │
// ├─────────────────────────────────────────┤
// │ fields (Slots)                          │
// │ ┌──────┬──────┬──────┬──────┬──────┐    │
// │ │ [0]  │ [1]  │ [2]  │ [3]  │ ...  │    │
// │ │父類欄 │父類欄 │子類欄 │子類欄 │      │    │
// │ └──────┴──────┴──────┴──────┴──────┘    │
// └─────────────────────────────────────────┘
type Object struct {
	// Object Head (64 bits)
	// including: hashCode (31 bits) + age (4 bits) + biased (1 bit) + lock state (2 bits)
	markWord uint64

	// class pointing to class's metadata (which is in method_area)
	// (使用 interface{} 避免與 method_area 套件循環依賴)
	// actual type is *method_area.Class
	class interface{}

	// fields object's fields ref
	fields rtcore.Slots

	// extra
	// - array Object: store ([]int32, []int64, []*Object ...)
	// - String Object: possibly store Go string
	// - Class Object: store *Class (usage: reflection)
	extra interface{}
}

// NewObject create new object with specified class
// slotCount: object's fields count (including parent's)
// this func only alloc space, not executing constructor
// the object's constructor is calling by `invokespecial` <init>
func NewObject(class interface{}, slotCount uint) *Object {
	fmt.Println("@@ Debug - [NewObject] class:", class, ", slotCount:", slotCount)
	return &Object{
		markWord: InitialMarkWord, // init state: non-lock, age=0, hashCode=0
		class:    class,
		fields:   rtcore.NewSlots(slotCount),
	}
}

// =============== Getters ===============

// Class getter
// return interface{}, caller need convert type
// ex: class := obj.Class().(*method_area.Class)
func (o *Object) Class() interface{} {
	return o.class
}

func (o *Object) Fields() rtcore.Slots {
	return o.fields
}

func (o *Object) Extra() interface{} {
	return o.extra
}

// =============== Setters ===============

func (o *Object) SetExtra(extra interface{}) {
	o.extra = extra
}

// =============== Field Access by SlotId ===============

func (o *Object) GetIntField(slotId uint) int32 {
	return o.fields.GetInt(slotId)
}

func (o *Object) SetIntField(slotId uint, val int32) {
	o.fields.SetInt(slotId, val)
}

func (o *Object) GetLongField(slotId uint) int64 {
	return o.fields.GetLong(slotId)
}

func (o *Object) SetLongField(slotId uint, val int64) {
	o.fields.SetLong(slotId, val)
}

func (o *Object) GetFloatField(slotId uint) float32 {
	return o.fields.GetFloat(slotId)
}

func (o *Object) SetFloatField(slotId uint, val float32) {
	o.fields.SetFloat(slotId, val)
}

func (o *Object) GetDoubleField(slotId uint) float64 {
	return o.fields.GetDouble(slotId)
}

func (o *Object) SetDoubleField(slotId uint, val float64) {
	o.fields.SetDouble(slotId, val)
}

func (o *Object) GetRefField(slotId uint) interface{} {
	return o.fields.GetRef(slotId)
}

func (o *Object) SetRefField(slotId uint, ref interface{}) {
	o.fields.SetRef(slotId, ref)
}

// =============== Type Checking ===============

// IsInstanceOf java if (A instanceof B) {...}
// targetClass: *method_area.Class
//
// 判斷規則：
// 1. if targetClass is class: check super
// 2. if targetClass is interface: check is implemented
// 3. if array type: 特殊處理
func (o *Object) IsInstanceOf(targetClass interface{}) bool {
	// TODO: 實作完整的類型檢查邏輯
	// 需要在 method_area.Class 上新增輔助方法
	// MVP 階段先簡化處理
	return o.class == targetClass
}

// =============== Null Check Helper ===============

// IsNull check this object is nil
// usage: `getfield`/`putfield`
func IsNull(obj *Object) bool {
	return obj == nil
}

// CheckNotNull check non-null
func CheckNotNull(obj *Object) {
	if obj == nil {
		panic("java.lang.NullPointerException")
	}
}

func (o *Object) String() string {
	if o.IsArray() {
		return fmt.Sprintf("<Array type=%v>", o.ArrayType())
	}
	return fmt.Sprintf("<Object class=%v>", o.class)
}
