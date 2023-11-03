package errorst

import (
	"fmt"
	"runtime"
	"strings"
)

// stackedError is an error chain with stacktrace information.
type stackedError struct {
	Message string
	Code    ErrorCode
	Inner   error
	Loc     trace
}

// trace contains information about the location of the error.
type trace struct {
	File     string
	Line     int
	Function string
}

func create(inner error, code ErrorCode, msgOrFmt string, vals ...interface{}) error {
	if code == NoCode {
		code = GetCode(inner)
	}

	return &stackedError{
		Message: fmt.Sprintf(msgOrFmt, vals...),
		Code:    code,
		Inner:   inner,
		Loc:     getTrace(),
	}
}

func getTrace() trace {
	var tr trace
	// Caller -> NewError/Wrap -> create -> getTrace
	pc, file, line, ok := runtime.Caller(3)
	if !ok {
		return tr
	}

	tr.File, tr.Line = file, line
	if fn := runtime.FuncForPC(pc); fn != nil {
		tr.Function = shortFuncName(fn)
	}
	return tr
}

func shortFuncName(f *runtime.Func) string {
	// f.Name() is like one of these:
	// - "github.com/xxx/shield/package.FuncName"
	// - "github.com/xxx/shield/package.Receiver.MethodName"
	// - "github.com/xxx/shield/package.(*PtrReceiver).MethodName"

	longName := f.Name()

	withoutPath := longName[strings.LastIndex(longName, "/")+1:]
	withoutPackage := withoutPath[strings.Index(withoutPath, ".")+1:]

	shortName := withoutPackage
	shortName = strings.Replace(shortName, "(", "", 1)
	shortName = strings.Replace(shortName, "*", "", 1)
	shortName = strings.Replace(shortName, ")", "", 1)

	return shortName
}

// --------------------------------------
// Following are methods to support the Error interface and pkg/errors package.

func (st *stackedError) Error() string {
	return fmt.Sprintf("%v", st)
}

func (st *stackedError) Unwrap() error {
	return st.Inner
}

func (st *stackedError) Cause() error {
	return st.Inner
}

func (st *stackedError) Is(target error) bool {
	if st.Code != NoCode {
		return st.Code == GetCode(target)
	}
	return st == target
}
