package errs

import "github.com/pkg/errors"

// Container struct used to represent a sentinel error with a 'type' that can
// be constructed/compared, etc.
type sentinalErrorWithType struct {
	message   string
	errorType ErrorType
}

func (s *sentinalErrorWithType) Error() string {
	return s.message
}

func (s *sentinalErrorWithType) Type() ErrorType {
	return s.errorType
}

func (s *sentinalErrorWithType) WithStack() ErrorWithType {
	innerWithStack := errors.New(s.message).(innerError)

	return &errorWithType{
		inner:     innerWithStack,
		errorType: s.errorType,
	}
}

func (s *sentinalErrorWithType) From(err error) ErrorWithType {
	return Wrap(err, s.errorType, s.message)
}
