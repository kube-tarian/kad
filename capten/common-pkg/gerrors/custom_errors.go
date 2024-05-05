// Package gerrors ...
package gerrors

// ErrorCode service gerror constants
type ErrorCode string

const (
	InternalError          ErrorCode = "Internal Error"
	ServiceSetup           ErrorCode = "ServiceSetup"
	ValidationFailed       ErrorCode = "Validations Failed"
	BadRequest             ErrorCode = "Bad Request"
	NotFound               ErrorCode = "Not Found"
	TokenNotFound          ErrorCode = "Token NotFound"
	AuthenticationFailed   ErrorCode = "Authentication Failed"
	MailConfigurationError ErrorCode = "Mail Configuration Error"
	RecordAlreadyExists    ErrorCode = "RecordAlreadyExists"
)
