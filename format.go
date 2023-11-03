package errorst

import (
	"fmt"
	"github.com/pkg/errors"
	"io"
	"strings"
)

// Format defines the behavior of Format on a stackedError.
// %v produces a full stacktrace including line number information.
// %s,%q only produce the error message.
func (st *stackedError) Format(s fmt.State, verb rune) {
	var text string
	switch verb {
	case 'v':
		text = formatWithTrace(st)
	case 's', 'q':
		text = st.Message
	}
	io.WriteString(s, text)
}

// formatWithTrace formats the error as a full stacktrace including Line number information.
// format:
// --- at <File>:<Line> (<Function>) ---
// <msg>
// --- at <File>:<Line> (<Function>) ---
// --- at <File>:<Line> (<Function>) ---
// ...
// Caused by: <Inner error>
func formatWithTrace(st *stackedError) string {
	var str string
	newline := func() {
		if str != "" && !strings.HasSuffix(str, "\n") {
			str += "\n"
		}
	}

	// iterate through the stacktrace
	var cur = st
	for {
		str += formatTrace(cur.Loc)
		newline()
		str += cur.Message

		var cause *stackedError
		if !errors.As(cur.Inner, &cause) {
			break
		}
		cur = cause
		newline()
	}

	if cur.Inner != nil {
		newline()
		str += "Caused by: "
		str += cur.Inner.Error()
	}
	return str
}

// formatTrace formats the trace information.
// format:
// --- at <File>:<Line> (<Function>) ---
func formatTrace(tr trace) string {
	if tr.File == "" && tr.Line == 0 && tr.Function == "" {
		return ""
	} else if tr.Function == "" {
		return fmt.Sprintf(" --- at %v:%v ---", tr.File, tr.Line)
	}
	return fmt.Sprintf(" --- at %v:%v (%v) ---", tr.File, tr.Line, tr.Function)
}
