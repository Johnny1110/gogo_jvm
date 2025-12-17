package heap

// ============================================================
// Exception Object Factory - v0.2.10
// ============================================================

// ExceptionData store in object.extra
type ExceptionData struct {
	Message string // ex: "/ by zero"
}

func NewExceptionObject(exClass interface{}, message string) *Object {
	return &Object{
		class: exClass,
		extra: &ExceptionData{
			Message: message,
		},
	}
}

// IsExceptionObject check object is ex or not
func (o *Object) IsExceptionObject() bool {
	return o.GetExceptionData() != nil
}

// GetExceptionData get ex data from object
func (o *Object) GetExceptionData() *ExceptionData {
	if o.extra == nil {
		return nil
	}

	if data, ok := o.extra.(*ExceptionData); ok {
		return data
	}

	return nil
}
