package errors

import "fmt"

type DomainError struct {
	code    string
	message string
}

func NewDomainError(code, message string) *DomainError {
	return &DomainError{
		code:    code,
		message: message,
	}
}

func (e *DomainError) Error() string {
	return fmt.Sprintf("[%s] %s", e.code, e.message)
}

func (e *DomainError) Code() string {
	return e.code
}

func (e *DomainError) Message() string {
	return e.message
}

func (e *DomainError) IsDomainError() bool {
	return true
}

type AuthenticationError struct {
	*DomainError
}

func NewAuthenticationError(code, message string) *AuthenticationError {
	return &AuthenticationError{
		DomainError: NewDomainError(code, message),
	}
}

type ValidationError struct {
	*DomainError
	field string
}

func NewValidationError(field, code, message string) *ValidationError {
	return &ValidationError{
		DomainError: NewDomainError(code, message),
		field:       field,
	}
}

func (e *ValidationError) Field() string {
	return e.field
}

type NotFoundError struct {
	*DomainError
	resource string
}

func NewNotFoundError(resource, message string) *NotFoundError {
	return &NotFoundError{
		DomainError: NewDomainError("not_found", message),
		resource:    resource,
	}
}

func (e *NotFoundError) Resource() string {
	return e.resource
}