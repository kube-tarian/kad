package gerrors

import (
	"fmt"
	"reflect"
	"runtime"
)

const stackSize = 4096

// A Tag represents an gerror identifier of any type.
type Tag interface{}

// A Gerror is a tagged gerror with a stack trace embedded in the Error() string.
type Gerror interface {
	// Returns the tag used to create this gerror.
	Tag() Tag

	// Returns the concrete type of the tag used to create this gerror.
	TagType() reflect.Type

	// Returns the string form of this gerror, which includes the tag value, the tag type, the gerror message, and a stack trace.
	Error() string

	// Test the tag used to create this gerror for equality with a given tag. Returns `true` if and only if the two are equal.
	EqualTag(Tag) bool

	// Message
	Message() string

	// Cause
	Cause() error
}

// New Returns an gerror containing the given tag and message and the current stack trace.
func New(tag Tag, message string) Gerror {
	var stack [stackSize]byte
	n := runtime.Stack(stack[:], false)
	return &err{tag, reflect.TypeOf(tag), message, stack[:n], nil}
}

// Newf Returns an gerror containing the given tag and format string and the current stack trace. The given inserts are applied to the format string to produce an gerror message.
func Newf(tag Tag, format string, insert ...interface{}) Gerror {
	return New(tag, fmt.Sprintf(format, insert...))
}

// NewFromError Return an gerror containing the given tag, the cause of the gerror, and the current stack trace.
func NewFromError(tag Tag, cause error) Gerror {
	if cause != nil {
		var stack [stackSize]byte
		n := runtime.Stack(stack[:], false)
		return &err{tag, reflect.TypeOf(tag), "Error caused by: " + cause.Error(), stack[:n], cause}
	}
	return nil
}

type err struct {
	tag        Tag
	typ        reflect.Type
	message    string
	stackTrace []byte
	cause      error
}

func (e *err) Error() string {
	return fmt.Sprintf("%v %v", e.tag, e.typ) + ": " + e.message + "\n" + string(e.stackTrace)
}

func (e *err) Tag() Tag {
	return e.tag
}

func (e *err) TagType() reflect.Type {
	return e.typ
}

func (e *err) EqualTag(tag Tag) bool {
	return e.typ == reflect.TypeOf(tag) && e.tag == tag
}

func (e *err) Message() string {
	return e.message
}

func (e *err) Cause() error {
	return e.cause
}

func (e *err) StackTrace() string {
	return string(e.stackTrace)
}

func (e ErrorCode) String() string {
	return string(e)
}

// GetErrorType ...
func GetErrorType(err error) ErrorCode {
	gerr, ok := err.(Gerror)
	if ok {
		return gerr.Tag().(ErrorCode)
	}
	return InternalError
}

// GetErrorMessage ...
func GetErrorMessage(err error) string {
	if gerr, ok := err.(Gerror); ok {
		if cause := gerr.Cause(); cause != nil {
			return fmt.Sprintf("%s: %s", gerr.Tag(), GetErrorMessage(cause))
		}
		return fmt.Sprintf("%s: %s", gerr.Tag(), gerr.Message())
	}
	return err.Error()
}
