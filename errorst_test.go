package errorst

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFormat(t *testing.T) {
	for ti, tcase := range []struct {
		err      error
		format   string
		expected string
	}{
		{
			// testcase: brief
			err:      NewError("plain"),
			format:   "%s",
			expected: "plain",
		},
		{
			// testcase: brief format
			err:      NewError("plain: %s", "test"),
			format:   "%s",
			expected: "plain: test",
		},
		{
			// testcase: full
			err:      NewError("plain"),
			format:   "%v",
			expected: formatLoc(pathWithPrefix("errorst_test.go"), "TestFormat") + "\nplain",
		},
		{
			// testcase: wrapped brief
			err:      Wrap(NewError("plain"), "wrapped"),
			format:   "%s",
			expected: "wrapped",
		},
		{
			// testcase: wrapped full
			err:    Wrap(NewError("plain"), "wrapped"),
			format: "%v",
			expected: formatLoc(pathWithPrefix("errorst_test.go"), "TestFormat") + "\nwrapped\n" +
				formatLoc(pathWithPrefix("errorst_test.go"), "TestFormat") + "\nplain",
		},
		{
			// testcase: wrap without msg
			err:    Wrap(NewError("plain"), ""),
			format: "%v",
			expected: formatLoc(pathWithPrefix("errorst_test.go"), "TestFormat") + "\n" +
				formatLoc(pathWithPrefix("errorst_test.go"), "TestFormat") + "\nplain",
		},
	} {
		actual := fmt.Sprintf(tcase.format, tcase.err)
		actual = hideLineNumbers(actual)
		assert.Equal(t, actual, tcase.expected, "testcase %d", ti)
	}
}

func TestCode(t *testing.T) {
	const ErrT1 ErrorCode = 1
	const ErrT2 ErrorCode = 2

	for ti, tcase := range []struct {
		err      error
		expected ErrorCode
	}{
		{
			// testcase: NoCode
			err:      NewError("plain"),
			expected: NoCode,
		},
		{
			// testcase: new with code
			err:      NewErrorWithCode(ErrT1, "plain"),
			expected: ErrT1,
		},
		{
			// testcase: wrap with coded
			err:      Wrap(NewErrorWithCode(ErrT2, "plain"), "wrapped"),
			expected: ErrT2,
		},
		{
			// testcase: wrap with code
			err:      WrapWithCode(NewErrorWithCode(ErrT1, "plain"), ErrT2, "wrapped"),
			expected: ErrT2,
		},
	} {
		actual := GetCode(tcase.err)
		assert.Equal(t, actual, tcase.expected, "testcase %d", ti)
	}
}

func TestErrorIs(t *testing.T) {

	// constants
	const ErrT1 ErrorCode = 1
	const ErrT2 ErrorCode = 2
	var baseErr = errors.New("base")

	for ti, tcase := range []struct {
		err      error
		equal    error
		expected bool
	}{
		{
			err:      WrapWithCode(baseErr, ErrT1, "wrapped"),
			equal:    baseErr,
			expected: true,
		},
		{
			err:      WrapWithCode(baseErr, ErrT1, "wrapped"),
			equal:    NewErrorWithCode(ErrT1, "new"),
			expected: true,
		},
		{
			err:      Current(Wrap(WrapWithCode(baseErr, ErrT1, "wrapped"), "")),
			equal:    NewErrorWithCode(ErrT1, "new"),
			expected: true,
		},
		{
			err:      RootCause(WrapWithCode(baseErr, ErrT1, "wrapped")),
			equal:    NewErrorWithCode(ErrT1, "new"),
			expected: false,
		},
		{
			err:      WrapWithCode(baseErr, ErrT1, "wrapped"),
			equal:    WrapWithCode(baseErr, ErrT2, "new"),
			expected: false,
		},
	} {
		assert.Equal(t, tcase.expected, errors.Is(tcase.err, tcase.equal), "testcase %d", ti)
	}
}
