package cerrors

import "fmt"

type NotFoundError struct {
	Subject string
}

func (err NotFoundError) Error() string {
	return fmt.Sprintf("%s not found", err.Subject)
}

func NewNotFoundError(subject string) NotFoundError {
	return NotFoundError{Subject: subject}
}

type AlreadyExistsError struct {
	Subject string
}

func (err AlreadyExistsError) Error() string {
	return fmt.Sprintf("Entity %s alredy exists", err.Subject)
}

func NewAlreadyExistsError(subject string) AlreadyExistsError {
	return AlreadyExistsError{Subject: subject}
}

type CriticalInternalError struct {
	Place   string
	Subject error
}

func (err CriticalInternalError) Error() string {
	return fmt.Sprintf("Critical error in %s: %s", err.Place, err.Subject.Error())
}

func NewCriticalInternalError(place string, subject error) CriticalInternalError {
	return CriticalInternalError{Place: place, Subject: subject}
}

type InvalidCredentialsError struct{}

func (err InvalidCredentialsError) Error() string {
	return fmt.Sprintf("Invalid credentials")
}

func NewInvalidCredentialsError() InvalidCredentialsError {
	return InvalidCredentialsError{}
}

const (
	TokenExpired = "expired"
	TokenBadFormat = "bad format"
)

type InvalidTokenError struct{
	subject string
}

func (err InvalidTokenError) Error() string {
	return fmt.Sprintf("Invalid token: %s", err.subject)
}

func (err InvalidTokenError) Subject() string {
	return err.subject
}

func NewInvalidTokenError(subject string) InvalidTokenError {
	return InvalidTokenError{subject: subject}
}