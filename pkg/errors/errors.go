package errors

import (
	"errors"
	"fmt"
)

var (
	ErrFieldsIsEmpty         = errors.New("fields is empty")
	ErrFieldIsEmpty          = errors.New("field is empty")
	ErrUnbalancedParenthesis = errors.New("unbalanced parentheses")
	ErrUnrecognizedFormat    = errors.New("unrecognized format")
)

var _ error = new(FieldIsNotObjectError)

func NotAnObjectError(field string) FieldIsNotObjectError {
	return FieldIsNotObjectError{field: field}
}

type FieldIsNotObjectError struct {
	field string
}

func (f FieldIsNotObjectError) Error() string {
	return fmt.Sprintf("Field '%s' is not object", f.field)
}
