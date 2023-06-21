package errs

type ErrorType string

// Collection of high-level ErrorTypes that can be used
//   - when constructing new errors
//   - checking an errors type
const (
	ErrTypeUnknown        ErrorType = "Unknown"
	ErrTypeConfiguration  ErrorType = "Configuration"
	ErrTypeUnmarshal      ErrorType = "Unmarshalling"
	ErrTypeMarshal        ErrorType = "Marshalling"
	ErrTypeValidation     ErrorType = "Validation"
	ErrTypeInvalidNumber  ErrorType = "InvalidNumber"
	ErrTypeInvalidBoolean ErrorType = "InvalidBoolean"
)
