package heap

// ============================================================
// Exception Object Factory - v0.2.10
// ============================================================

// ExceptionData store in object.extra
type ExceptionData struct {
	ClassName string // ex: "java/lang/ArithmeticException"
	Message   string // ex: "/ by zero"
}

// NewExceptionObject create ex Object
// TODO: this is MVP simplify
func NewExceptionObject(className, message string) *Object {
	return &Object{
		class: nil, // TODO: real JVM will load ex class and create instance
		extra: &ExceptionData{
			ClassName: className,
			Message:   message,
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
