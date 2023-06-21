package errs

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

// Container struct used to represent an error with a 'type' that can
// be constructed/compared, etc.
type errorWithType struct {
	inner     innerError
	errorType ErrorType
}

func (e *errorWithType) Error() string {
	return e.inner.Error()
}

func (e *errorWithType) Type() ErrorType {
	return e.errorType
}

func (e *errorWithType) Unwrap() error {
	return e.inner
}

func (e *errorWithType) Format(f fmt.State, verb rune) {
	e.inner.Format(f, verb)
}

func (e *errorWithType) StackTrace() errors.StackTrace {
	return e.inner.StackTrace()
}

// Is supports the ability to compare an instance of 'errorWithType' to other errors using
// errors.Is() - this is important because the child errors within 'errorWithType' have a
// stacktrace which will never be considered equal to another error's stacktrace.
//
// Implementing this method allows sentinal errors to be compared with errors that were
// created when a sentinal's WithStack() or From() is invoked.
func (e *errorWithType) Is(err error) bool {
	if typer, isTyper := err.(Typer); isTyper {
		isSameType := typer.Type() == e.Type()
		hasSameError := e.Error() == err.Error()
		hasSameErrorBeforeWrap := strings.HasPrefix(e.Error(), fmt.Sprintf("%s:", err.Error()))

		return isSameType && (hasSameError || hasSameErrorBeforeWrap)
	}

	return false
}
