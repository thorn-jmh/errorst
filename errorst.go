package errorst

import (
	"fmt"
	"github.com/pkg/errors"
	"runtime"
	"strings"
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
		err = st.Inner
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
	const depth = 32
	var pcs [depth]uintptr

	// Caller -> StackTrace
	n := runtime.Callers(2, pcs[:])

	str := fmt.Sprintf(msgOrFmt, vals...)
	newline := func() {
		if str != "" && !strings.HasSuffix(str, "\n") {
			str += "\n"
		}
	}
	for i := 0; i < n; i++ {
		fn := runtime.FuncForPC(pcs[i])
		if fn == nil {
			newline()
			str += fmt.Sprintf(" --- pc: %#v ---", pcs[i])
		} else {
			file, line := fn.FileLine(pcs[i])
			trace := trace{
				Function: shortFuncName(fn),
				File:     file,
				Line:     line,
			}
			newline()
			str += formatTrace(trace)
		}
	}

	return errors.New(str)
}
