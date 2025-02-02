package cerrors

import "fmt"

type NotFoundError struct {
	Subject string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s not found", e.Subject)
}

func NewNotFoundError(subject string) *NotFoundError {
	return &NotFoundError{Subject: subject}
}

type AlreadyExistsError struct {
	Subject string
}

func (e *AlreadyExistsError) Error() string {
	return fmt.Sprintf("Entity %s alredy exists", e.Subject)
}

func NewAlreadyExistsError(subject string) *AlreadyExistsError {
	return &AlreadyExistsError{Subject: subject}
}

type CriticalInternalError struct {
	Place   string
	Subject error
}

func (e *CriticalInternalError) Error() string {
	return fmt.Sprintf("Critical error in %s: %s", e.Place, e.Subject.Error())
}

func NewCriticalInternalError(place string, subject error) *CriticalInternalError {
	return &CriticalInternalError{Place: place, Subject: subject}
}

type InvalidCredentialsError struct{}

func (e *InvalidCredentialsError) Error() string {
	return fmt.Sprintf("Invalid credentials")
}

func NewInvalidCredentialsError() *InvalidCredentialsError {
	return &InvalidCredentialsError{}
}
