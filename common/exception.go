package common

type JavaException struct {
	ClassName string
	Message   string
}

func NewJavaException(className string, message string) *JavaException {
	return &JavaException{
		ClassName: className,
		Message:   message,
	}
}
