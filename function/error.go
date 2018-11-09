package function

import "fmt"

// Error is a function error which can contain an argument number
type Error struct {
	Err      error
	ArgIndex int
}

// NewError creates a new function error
func NewError(argIndex int, err error) *Error {
	return &Error{
		ArgIndex: argIndex,
		Err:      err,
	}
}

// NewErrorf creates a new function error
func NewErrorf(argIndex int, format string, args ...interface{}) *Error {
	return &Error{
		ArgIndex: argIndex,
		Err:      fmt.Errorf(format, args...),
	}
}

// Error returns with the error message
func (e *Error) Error() string {
	return e.Err.Error()
}
