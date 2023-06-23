package errs_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/suite"

	"github.com/djmarrerajr/common-lib/errs"

	stderr "errors"

	pkgerr "github.com/pkg/errors"
)

// Utility types used in the tests that follow...
type CustomError struct{}

func (c CustomError) Error() string { return "I am a custom error!" }

type TypedError struct{}

func (t TypedError) Error() string        { return "I am a custom typed error" }
func (t TypedError) Type() errs.ErrorType { return "VerySpecificType" }

type HelperTestSuite struct {
	suite.Suite
}

func (h *HelperTestSuite) TestNew_ReturnsErrorWithProvidedType() {
	errType := errs.ErrorType("type")

	err := errs.New(errType, "Testing")
	h.Equal(errType, err.Type())
}

func (h *HelperTestSuite) TestNew_ErrorStringMatchesUnwrappedErrorString() {
	message := "bruno"
	err := errs.New(errs.ErrTypeUnknown, message)

	h.Equal(err.Unwrap().Error(), err.Error())
}

func (h *HelperTestSuite) TestNew_ErrorStackTraceMatchesUnwrappedStackTrace() {
	err := errs.New(errs.ErrTypeUnknown, "bruno")
	h.Equal(err.Unwrap().(errs.StackTracer).StackTrace(), err.StackTrace())
}

func (h *HelperTestSuite) TestErrorf_ErrorIsFormattedString() {
	format, arg := "hello %s", "bruno"
	err := errs.Errorf(errs.ErrTypeUnknown, format, arg)

	h.Equal(fmt.Sprintf(format, arg), err.Error())
}

func (h *HelperTestSuite) TestErrorf_ErrorHasSpecifiedType() {
	errType := errs.ErrTypeValidation
	err := errs.Errorf(errs.ErrTypeValidation, "bruno")

	h.Equal(errType, err.Type())
}

func (h *HelperTestSuite) TestWrap_NilError_ReturnsNil() {
	result := errs.Wrap(nil, errs.ErrTypeUnknown, "bruno")

	h.Nil(result)
}

func (h *HelperTestSuite) TestWrap_ReturnsErrorWrappingSpecifiedError() {
	toWrap := fmt.Errorf("this should be wrapped")
	result := errs.Wrap(toWrap, errs.ErrTypeUnknown, "bruno")

	h.Contains(result.Error(), toWrap.Error())
}

func (h *HelperTestSuite) TestWrap_MessagePassed_ErrorStringContainsMessge() {
	toWrap := fmt.Errorf("this should be wrapped")
	errMsg := "this is the message"
	result := errs.Wrap(toWrap, errs.ErrTypeUnknown, errMsg)

	h.ErrorContains(result, errMsg)
}

func (h *HelperTestSuite) TestWrap_ResultingErrorTypeIsSpecifiedErrorType() {
	toWrap := fmt.Errorf("this is the inner error")
	errTyp := errs.ErrorType("fancy-error")
	result := errs.Wrap(toWrap, errTyp, "bruno")

	h.Equal(errTyp, result.Type())
}

func (h *HelperTestSuite) TestWithType_ErrorStringIsWrappedErrorString() {
	toWrap := fmt.Errorf("this should be wrapped")
	errTyp := errs.ErrorType("fancy-error")
	result := errs.Wrap(toWrap, errTyp, "bruno")

	h.Contains(result.Error(), toWrap.Error())
}

func (h *HelperTestSuite) TestWrapf_NilError_ReturnsNil() {
	result := errs.Wrapf(nil, errs.ErrTypeUnknown, "bruno")

	h.Nil(result)
}

func (h *HelperTestSuite) TestWrapf_ReturnsErrorWrappingSpecifiedError() {
	toWrap := fmt.Errorf("this should be wrapped")
	result := errs.Wrapf(toWrap, errs.ErrTypeUnknown, "hello %s", "bruno")

	h.ErrorContains(result, toWrap.Error())
}

func (h *HelperTestSuite) TestWrapf_MessagePassed_ErrorStringContainsFormattedArgs() {
	toWrap := fmt.Errorf("this should be wrapped")
	errMsg := "hello: %s"
	result := errs.Wrapf(toWrap, errs.ErrTypeUnknown, errMsg, "bruno")

	h.ErrorContains(result, fmt.Sprintf(errMsg, "bruno"))
}

func (h *HelperTestSuite) TestSentinelWithStack_CanBeComparedWithIs() {
	sentinel := errs.Sentinel(errs.ErrTypeUnmarshal, "testing")

	s1 := sentinel.WithStack()
	s2 := sentinel.WithStack()

	h.NotEqual(s1, s2)
	h.ErrorIs(s1, sentinel)
	h.ErrorIs(s2, sentinel)
	h.ErrorIs(s1, s2)
}

func (h *HelperTestSuite) TestSentinelFrom_CanBeComparedWithIs() {
	sentinel := errs.Sentinel(errs.ErrTypeUnmarshal, "testing")
	otherErr := pkgerr.New("testing also")
	wrapped := sentinel.From(otherErr)

	h.ErrorIs(wrapped, sentinel)
}

func (h *HelperTestSuite) TestFindStackTrace_NotNestedStdlibError_ReturnsNil() {
	err := stderr.New("test error")
	trace := errs.FindOriginalStackTrace(err)

	h.Nil(trace)
}

func (h *HelperTestSuite) TestFindStackTrace_NotNestedPkgError_GetStacKTraceFromError() {
	err := pkgerr.New("test error")
	trace := errs.FindOriginalStackTrace(err)

	h.Equal(err.(errs.StackTracer).StackTrace(), *trace)
}

func (h *HelperTestSuite) TestFindStackTrace_NotNestedErrsError_GetStacKTraceFromError() {
	err := errs.New(errs.ErrTypeValidation, "test error")
	trace := errs.FindOriginalStackTrace(err)

	h.Equal(err.(errs.StackTracer).StackTrace(), *trace)
}

func (h *HelperTestSuite) TestFindStackTrace_OnceWrapped_StdlibError_GetStacKTraceFromDeepestError() {
	original := pkgerr.New("test error")
	wrapped := fmt.Errorf("wrapped %w", original)
	trace := errs.FindOriginalStackTrace(wrapped)

	h.Equal(original.(errs.StackTracer).StackTrace(), *trace)
}

func (h *HelperTestSuite) TestFindStackTrace_OnceWrapped_PkgError_GetStacKTraceFromDeepestError() {
	original := pkgerr.New("test error")
	wrapped := pkgerr.Wrap(original, "wrapped")
	trace := errs.FindOriginalStackTrace(wrapped)

	h.Equal(original.(errs.StackTracer).StackTrace(), *trace)
}

func (h *HelperTestSuite) TestFindStackTrace_ThriceWrapped_GetStacKTraceFromDeepestError() {
	original := pkgerr.New("original")
	wrapped := pkgerr.Wrap(original, "wrapped")
	wrapped = errs.Wrap(wrapped, errs.ErrTypeValidation, "wrapped again")
	wrapped = fmt.Errorf("wrapped yet again: %w", wrapped)
	trace := errs.FindOriginalStackTrace(wrapped)

	h.Equal(original.(errs.StackTracer).StackTrace(), *trace)
}

func (h *HelperTestSuite) TestFindStackTrace_ThriceWrapped_OriginalHasNoStackTrace_GetStacKTraceFromDeepestError() {
	original := stderr.New("original")
	second := pkgerr.Wrap(original, "wrapped")
	wrapped := errs.Wrap(second, errs.ErrTypeValidation, "wrapped again").(error)
	wrapped = fmt.Errorf("wrapped yet again: %w", wrapped)
	trace := errs.FindOriginalStackTrace(wrapped)

	h.Equal(second.(errs.StackTracer).StackTrace(), *trace)
}

func (h *HelperTestSuite) TestGetType_ReturnsExpectedErrorType() {
	stdLibError := stderr.New("error!")

	tests := []lo.Tuple3[string, error, errs.ErrorType]{
		{A: "not an error", B: nil, C: ""},
		{A: "stdlib error", B: stdLibError, C: errs.ErrTypeUnknown},
		{A: "pkgerr error", B: pkgerr.New("error!"), C: errs.ErrTypeUnknown},
		{A: "stdlib wrapped error", B: fmt.Errorf("wrapped: %w", stdLibError), C: errs.ErrTypeUnknown},
		{A: "pkgerr wrapped error", B: pkgerr.Wrap(stdLibError, "wrapped"), C: errs.ErrTypeUnknown},
		{A: "public error type", B: CustomError{}, C: "CustomError"},
		{A: "pkgerr wrapping public error type", B: pkgerr.Wrap(CustomError{}, "wrapped"), C: "CustomError"},
		{A: "double wrapped custom error type", B: pkgerr.Wrap(fmt.Errorf("wrapping: %w", CustomError{}), "wrapped"), C: "CustomError"},
		{A: "error with type", B: errs.New("SpecialType", "some message"), C: "SpecialType"},
		{A: "wrapped error with type", B: pkgerr.Wrap(errs.New("SpecialType", "some message"), "wrapped"), C: "SpecialType"},
		{A: "error with custom typer", B: TypedError{}, C: "VerySpecificType"},
		{A: "wrapped error with custom typer", B: pkgerr.Wrap(TypedError{}, "wrapped"), C: "VerySpecificType"},
		{A: "imported error", B: &json.SyntaxError{}, C: "SyntaxError"},
		{A: "sentinel error", B: errs.Sentinel(errs.ErrTypeUnknown, "danger!"), C: errs.ErrTypeUnknown},
		{A: "wrapped sentinel error", B: pkgerr.Wrap(errs.Sentinel(errs.ErrTypeUnknown, "danger!"), "wrapped"), C: errs.ErrTypeUnknown},
	}

	for _, test := range tests {
		h.Run(test.A, func() {
			result := errs.GetType(test.B)
			h.Equal(test.C, result)
		})
	}
}

func (h *HelperTestSuite) TestWithTypeFallback_ReturnsExpectedErrorType() {
	stdLibError := stderr.New("error!")
	fallBackErr := errs.ErrorType("fallback!")

	tests := []lo.Tuple3[string, error, errs.ErrorType]{
		{A: "not an error", B: nil, C: ""},
		{A: "stdlib error", B: stdLibError, C: fallBackErr},
		{A: "pkgerr error", B: pkgerr.New("error!"), C: fallBackErr},
		{A: "stdlib wrapped error", B: fmt.Errorf("wrapped: %w", stdLibError), C: fallBackErr},
		{A: "pkgerr wrapped error", B: pkgerr.Wrap(stdLibError, "wrapped"), C: fallBackErr},
		{A: "public error type", B: CustomError{}, C: "CustomError"},
		{A: "pkgerr wrapping public error type", B: pkgerr.Wrap(CustomError{}, "wrapped"), C: "CustomError"},
		{A: "double wrapped custom error type", B: pkgerr.Wrap(fmt.Errorf("wrapping: %w", CustomError{}), "wrapped"), C: "CustomError"},
		{A: "error with type", B: errs.New("SpecialType", "some message"), C: "SpecialType"},
		{A: "wrapped error with type", B: pkgerr.Wrap(errs.New("SpecialType", "some message"), "wrapped"), C: "SpecialType"},
		{A: "error with custom typer", B: TypedError{}, C: "VerySpecificType"},
		{A: "wrapped error with custom typer", B: pkgerr.Wrap(TypedError{}, "wrapped"), C: "VerySpecificType"},
		{A: "imported error", B: &json.SyntaxError{}, C: "SyntaxError"},
		{A: "sentinel error", B: errs.Sentinel(errs.ErrTypeUnmarshal, "danger!"), C: errs.ErrTypeUnmarshal},
		{A: "wrapped sentinel error", B: pkgerr.Wrap(errs.Sentinel(errs.ErrTypeUnmarshal, "danger!"), "wrapped"), C: errs.ErrTypeUnmarshal},
	}

	for _, test := range tests {
		h.Run(test.A, func() {
			err := errs.WithTypeFallback(test.B, fallBackErr)
			result := errs.GetType(err)
			h.Equal(test.C, result)
		})
	}
}

func TestHelpers(t *testing.T) {
	suite.Run(t, new(HelperTestSuite))
}
