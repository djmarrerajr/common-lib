package errs

import (
	stderr "errors"
	"reflect"
	"unicode"

	pkgerr "github.com/pkg/errors"
	"github.com/samber/lo"
)

var basicErrorPackages = []string{
	"fmt",
	"errors",
	"github.com/pkg/errors",
	"golang.org/x/xerrors",
}

// New returns a new ErrorWithType complete with StackTrace and the specified ErrorType.
func New(errType ErrorType, message string) ErrorWithType {
	inner := pkgerr.New(message).(innerError)

	return newWithInner(inner, errType)
}

// Errorf operates like New except it allows you to format the message with parameterized args.
func Errorf(errType ErrorType, format string, args ...any) ErrorWithType {
	err := pkgerr.Errorf(format, args...).(innerError)

	return newWithInner(err, errType)
}

// Sentinel returns a StacklessErrorWithType that can be used as a sentinel error and compared
// to other errors.  When using a sentinel be sure to invoke WithStack() on the sentinel to get
// a copy of it with a complete stack trace.
func Sentinel(errType ErrorType, message string) StacklessErrorWithType {
	return &sentinalErrorWithType{message, errType}
}

// WithType wraps the provided error with a StackTrace and the ErrorType without needing to
// alter the original error message.
func WithType(err error, errType ErrorType) ErrorWithType {
	if err == nil {
		return nil
	}

	wrapped := pkgerr.WithStack(err).(innerError)

	return newWithInner(wrapped, errType)
}

// WithTypeFallback operates similarly to WithType except it will only used the provided
// ErrorType if the provided error has a type that cannot be determined.
func WithTypeFallback(err error, errType ErrorType) ErrorWithType {
	if err == nil {
		return nil
	}

	if errWithType, isErrWithType := err.(ErrorWithType); isErrWithType {
		return errWithType
	}

	currType := GetType(err)
	if currType != ErrTypeUnknown {
		return WithType(err, currType)
	}

	return WithType(err, errType)
}

// Wrap will wrap the provided error with the specified message/type.
//
// This is useful when bubbling errors up the call stack...
// i.e. a 'validation error' can be wrapped with a 'configuration error'
func Wrap(err error, errType ErrorType, message string) ErrorWithType {
	if err == nil {
		return nil
	}

	wrapped := pkgerr.Wrap(err, message).(innerError)

	return newWithInner(wrapped, errType)
}

// Wrapf operates like Wrap except it allows the message string to be formatted.
func Wrapf(err error, errType ErrorType, format string, args ...any) ErrorWithType {
	if err == nil {
		return nil
	}

	wrapped := pkgerr.Wrapf(err, format, args...).(innerError)

	return newWithInner(wrapped, errType)
}

// GetType will return the ErrorType of the provided error.
func GetType(err error) ErrorType {
	if err == nil {
		return ""
	}

	// if it is a Typer, use its type
	if e, OK := err.(Typer); OK {
		return e.Type()
	}

	// otherwise fallback to reflection...
	reflection := getErrorType(err)
	isBasic := isErrorBasic(reflection)
	unwrapped := stderr.Unwrap(err)

	switch {
	case isPublic(reflection) && !isBasic:
		return ErrorType(reflection.Name())
	case unwrapped != nil:
		return GetType(unwrapped)
	default:
		return ErrTypeUnknown
	}
}

// FindOriginalStackTrace will recurse the error 'tree' looking for and
// returning the original (underlying) StackTrace.
func FindOriginalStackTrace(err error) *pkgerr.StackTrace {
	var lastStackTracer StackTracer

	for {
		tracer, isTracer := err.(StackTracer)
		if isTracer {
			lastStackTracer = tracer
		}

		unwrapped := stderr.Unwrap(err)
		if unwrapped != nil {
			err = unwrapped
		} else {
			break
		}
	}

	if lastStackTracer != nil {
		trace := lastStackTracer.StackTrace()
		return &trace
	}

	return nil
}

// isErrorBasic determines if the error in question stems from one of the
// known error packages (errors, fmt.Errorf, etc...)
func isErrorBasic(reflection reflect.Type) bool {
	return lo.Contains(basicErrorPackages, reflection.PkgPath())
}

// getErrorType uses relfection to return the errors GO type.
func getErrorType(err error) reflect.Type {
	reflection := reflect.TypeOf(err)
	if reflection.Kind() == reflect.Pointer {
		reflection = reflection.Elem()
	}

	return reflection
}

// isPublic uses reflection to check the all important first letter of
// the error name to determine its visibility.
func isPublic(reflection reflect.Type) bool {
	if name := reflection.Name(); name != "" {
		return unicode.IsUpper(rune(name[0]))
	}

	return false
}

// newWithInner returns a new ErrorWithType
func newWithInner(wrapped innerError, errType ErrorType) ErrorWithType {
	if inner, OK := wrapped.(ErrorWithType); OK && inner.Type() == errType {
		return inner
	}

	return &errorWithType{
		inner:     wrapped,
		errorType: errType,
	}
}
