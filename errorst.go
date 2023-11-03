package errorst

import (
	"github.com/pkg/errors"
)

type ErrorCode int

const NoCode ErrorCode = 0

// --------------------------------------

// GetCode returns code of the error.If err is nil or if there is
// no error code attached to err, it returns NoCode.
func GetCode(err error) ErrorCode {
	var st *stackedError
	if errors.As(err, &st) {
		return st.Code
	}
	return NoCode
}

// -----------------------------------------

// NewError creates a new error with the given message, and records the trace.
func NewError(msgOrFmt string, vals ...interface{}) error {
	return create(nil, NoCode, msgOrFmt, vals...)
}

// NewErrorWithCode creates a new error with the given message and code, and records the trace.
func NewErrorWithCode(code ErrorCode, msgOrFmt string, vals ...interface{}) error {
	return create(nil, code, msgOrFmt, vals...)
}

// Wrap wraps the given error with the given message, and records the trace.
func Wrap(inner error, msgOrFmt string, vals ...interface{}) error {
	return create(inner, NoCode, msgOrFmt, vals...)
}

// WrapWithCode wraps the given error with the given message and code, and records the trace.
func WrapWithCode(inner error, code ErrorCode, msgOrFmt string, vals ...interface{}) error {
	return create(inner, code, msgOrFmt, vals...)
}

// -----------------------------------------

// RootCause unwraps the original error that caused the current one.
func RootCause(err error) error {
	for {
		var st *stackedError
		if ok := errors.As(err, &st); !ok || st.Inner == nil {
			return err
		}
	}
}

// Current returns the top level error.
func Current(err error) error {
	var st *stackedError
	if ok := errors.As(err, &st); !ok {
		return err
	}
	return &stackedError{
		Message: st.Message,
		Code:    st.Code,
		Inner:   nil,
		Loc:     st.Loc,
	}
}

// -----------------------------------------

// StackTrace create an error with full stack trace.
// Because all trace info has been recorded in the error, so it's not recommended to use Wrap on
// this error.
func StackTrace(msgOrFmt string, vals ...interface{}) error {
	// TODO: not impl
	return errors.New("not impl")
}
