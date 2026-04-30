// Package flattener provides functionality to flatten YAML structures into key-value pairs.
package flattener

import (
	"fmt"
)

// ErrorType represents the category of error that occurred
type ErrorType string

const (
	// ErrTypeValidation indicates input validation failures
	ErrTypeValidation ErrorType = "validation"
	// ErrTypeParsing indicates YAML parsing failures
	ErrTypeParsing ErrorType = "parsing"
	// ErrTypeDepthLimit indicates maximum nesting depth exceeded
	ErrTypeDepthLimit ErrorType = "depth_limit"
	// ErrTypeSizeLimit indicates size limit exceeded
	ErrTypeSizeLimit ErrorType = "size_limit"
	// ErrTypeTimeout indicates operation timed out
	ErrTypeTimeout ErrorType = "timeout"
	// ErrTypePathSecurity indicates a file path failed security checks
	ErrTypePathSecurity ErrorType = "path_security"
	// ErrTypeFileAccess indicates a file could not be accessed
	ErrTypeFileAccess ErrorType = "file_access"
)

// Error represents a structured error from the flattener
type Error struct {
	Type    ErrorType
	Message string
	Err     error
}

// Error implements the error interface
func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s error: %s: %v", e.Type, e.Message, e.Err)
	}
	return fmt.Sprintf("%s error: %s", e.Type, e.Message)
}

// Unwrap returns the underlying error
func (e *Error) Unwrap() error {
	return e.Err
}

// Is implements error comparison for errors.Is
func (e *Error) Is(target error) bool {
	t, ok := target.(*Error)
	if !ok {
		return false
	}
	return e.Type == t.Type
}

// ValidationError creates a validation error
func ValidationError(message string, err error) *Error {
	return &Error{
		Type:    ErrTypeValidation,
		Message: message,
		Err:     err,
	}
}

// ParsingError creates a parsing error
func ParsingError(message string, err error) *Error {
	return &Error{
		Type:    ErrTypeParsing,
		Message: message,
		Err:     err,
	}
}

// DepthLimitError creates a depth limit error
func DepthLimitError(depth int) *Error {
	return &Error{
		Type:    ErrTypeDepthLimit,
		Message: fmt.Sprintf("maximum nesting depth of %d exceeded", depth),
		Err:     nil,
	}
}

// SizeLimitError creates a size limit error
func SizeLimitError(size int, limitType string) *Error {
	return &Error{
		Type:    ErrTypeSizeLimit,
		Message: fmt.Sprintf("maximum %s size of %d exceeded", limitType, size),
		Err:     nil,
	}
}

// TimeoutError creates a timeout error
func TimeoutError(operation string) *Error {
	return &Error{
		Type:    ErrTypeTimeout,
		Message: fmt.Sprintf("%s timed out, content may be too complex", operation),
		Err:     nil,
	}
}

// PathSecurityError creates a path security error
func PathSecurityError(message string) *Error {
	return &Error{
		Type:    ErrTypePathSecurity,
		Message: message,
	}
}

// FileAccessError creates a file access error
func FileAccessError(message string, err error) *Error {
	return &Error{
		Type:    ErrTypeFileAccess,
		Message: message,
		Err:     err,
	}
}
