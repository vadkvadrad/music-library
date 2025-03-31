package er

import (
	"fmt"
	"net/http"
)

// Тип для структурированных ошибок
type Error struct {
	Code    string // Уникальный код ошибки
	Message string // Человекочитаемое сообщение
	Status  int    // HTTP-статус код
	Cause   error  // Исходная ошибка (при наличии)
}

// Реализуем интерфейс error
func (e *Error) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

// Unwrap для совместимости с errors.Is/As
func (e *Error) Unwrap() error {
	return e.Cause
}

// Константы ошибок
const (
	CodeUnauthorized        = "UNAUTHORIZED"
	CodeBadRequest          = "BAD_REQUEST"
	CodeNotFound            = "NOT_FOUND"
	CodeInternalServerError = "INTERNAL_ERROR"
	CodeConflict            = "CONFLICT"
)

// Предопределенные ошибки
var (
	ErrNotAuthorized = &Error{
		Code:    CodeUnauthorized,
		Message: "User is not authorized",
		Status:  http.StatusUnauthorized,
	}

	ErrWrongUserCredentials = &Error{
		Code:    CodeBadRequest,
		Message: "Wrong user credentials",
		Status:  http.StatusBadRequest,
	}

	ErrUserExists = &Error{
		Code:    CodeConflict,
		Message: "User already exists",
		Status:  http.StatusConflict,
	}

	ErrUserNotExists = &Error{
		Code:    CodeNotFound,
		Message: "User does not exist",
		Status:  http.StatusNotFound,
	}

	ErrUserNotVerified = &Error{
		Code:    CodeBadRequest,
		Message: "User is not verified",
		Status:  http.StatusBadRequest,
	}

	ErrNegativeIncome = &Error{
		Code:    CodeBadRequest,
		Message: "Income can't be negative or zero",
		Status:  http.StatusBadRequest,
	}
)

// New создает новую структурированную ошибку
func New(code, message string, status int) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Status:  status,
	}
}

// Wrap оборачивает ошибку с сохранением структуры
func Wrap(err *Error, cause error) *Error {
	return &Error{
		Code:    err.Code,
		Message: err.Message,
		Status:  err.Status,
		Cause:   cause,
	}
}

// WrapMessage оборачивает ошибку с новым сообщением
func WrapMessage(err *Error, message string, args ...interface{}) *Error {
	return &Error{
		Code:    err.Code,
		Message: fmt.Sprintf(message, args...),
		Status:  err.Status,
		Cause:   err,
	}
}

// Is проверяет тип ошибки
func Is(err error, target *Error) bool {
	if e, ok := err.(*Error); ok {
		return e.Code == target.Code
	}
	return false
}

// As возвращает структурированную ошибку
func As(err error) (*Error, bool) {
	e, ok := err.(*Error)
	return e, ok
}

// HTTPStatus возвращает HTTP статус для ошибки
func HTTPStatus(err error) int {
	if e, ok := err.(*Error); ok {
		return e.Status
	}
	return http.StatusInternalServerError
}