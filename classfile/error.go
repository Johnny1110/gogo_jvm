package classfile

import "fmt"

type ClassFormatError struct {
	reason string
}

func (e *ClassFormatError) Error() string {
	return fmt.Sprintf("class format error: %s", e.reason)
}

func newClassFormatError(reason string) *ClassFormatError {
	return &ClassFormatError{reason: reason}
}
