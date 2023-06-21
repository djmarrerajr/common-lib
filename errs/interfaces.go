package errs

import (
	"fmt"

	"github.com/pkg/errors"
	"golang.org/x/xerrors"
)

type StackTracer interface {
	StackTrace() errors.StackTrace
}

type Typer interface {
	Type() ErrorType
}

type StacklessErrorWithType interface {
	error
	Typer
	WithStack() ErrorWithType
	From(error) ErrorWithType
}

type ErrorWithType interface {
	error
	Typer
	xerrors.Wrapper
	StackTracer
	fmt.Formatter
}

type innerError interface {
	error
	StackTracer
	fmt.Formatter
}
