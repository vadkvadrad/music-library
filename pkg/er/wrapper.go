package er

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Кастомные типы ошибок
type (
	ValidationError struct {
		Message string
	}
	NotFoundError struct {
		Message string
	}
	UnauthorizedError struct {
		Message string
	}
	InternalError struct {
		Message string
	}
	ConflictError struct {
		ResourceType string
	}
)

func (e ValidationError) Error() string   { return e.Message }
func (e NotFoundError) Error() string     { return e.Message }
func (e UnauthorizedError) Error() string { return e.Message }
func (e InternalError) Error() string     { return e.Message }
func (e *ConflictError) Error() string    { return e.ResourceType }

// ErrorResponse - унифицированный формат ответа об ошибке
type ErrorResponse struct {
	Error     string `json:"error"`
	Tip       string `json:"tip"`
	Reference string `json:"reference,omitempty"` // Уникальный ID ошибки для поиска в логах
}

// ErrorHandlerConfig - конфиг обработчика ошибок
type ErrorHandlerConfig struct {
	LogErrors    bool
	ShowInternal bool // Показывать ли внутренние ошибки клиенту
	AppName      string
}

type ErrorHandler struct {
	cfg *ErrorHandlerConfig
}

func NewErrorHandler(cfg *ErrorHandlerConfig) *ErrorHandler {
	return &ErrorHandler{cfg: cfg}
}

// Предопределенные ошибки
var (
	ErrNotAuthorized = &UnauthorizedError{
		Message: "User is not authorized",
	}

	ErrWrongUserCredentials = &ValidationError{
		Message: "Wrong user credentials",
	}

	ErrUserExists = &ConflictError{
		ResourceType: "User already exists",
	}

	ErrArtistExists = &ConflictError{
		ResourceType: "Artist already exists",
	}

	ErrProfileExists = &ConflictError{
		ResourceType: "Profile already exists",
	}

	ErrAlbumExists = &ConflictError{
		ResourceType: "Album already exists",
	}

	ErrSongExists = &ConflictError{
		ResourceType: "Song already exist in this album",
	}

	ErrArtistLinked = &ConflictError{
		ResourceType: "Album already linked to your account",
	}

	ErrUserNotExists = &NotFoundError{
		Message: "User does not exist",
	}

	ErrArtistNotExists = &NotFoundError{
		Message: "Artist does not exists",
	}

	ErrAlbumNotExists = &NotFoundError{
		Message: "Album does not exist",
	}


	ErrUserNotVerified = &ValidationError{
		Message: "User is not verified",
	}

	ErrDateFormat = &ValidationError{
		Message: "Invalid date format: expected YYYY-MM-DD",
	}
)

// Handle - основной метод обработки ошибок
func (h *ErrorHandler) Handle(ctx context.Context, err error) (int, ErrorResponse) {
	select {
	case <-ctx.Done():
		return http.StatusRequestTimeout, ErrorResponse{
			Error: "Request timed out",
			Tip:   "Please try again later",
		}
	default:
		// Генерируем уникальный идентификатор ошибки
		errorID := generateErrorID()

		// Логирование ошибки
		if h.cfg.LogErrors {
			log.Printf("[%s] ERROR [%s]: %v", h.cfg.AppName, errorID, sanitizeError(err))
		}

		// Обработка кастомных ошибок
		switch e := err.(type) {
		case *ValidationError:
			return http.StatusBadRequest, ErrorResponse{
				Error:     e.Error(),
				Tip:       "Check your input data and try again",
				Reference: errorID,
			}
		case *NotFoundError:
			return http.StatusNotFound, ErrorResponse{
				Error:     e.Error(),
				Tip:       "The requested resource was not found",
				Reference: errorID,
			}
		case *UnauthorizedError:
			return http.StatusUnauthorized, ErrorResponse{
				Error:     "Authentication required",
				Tip:       "Please login and try again",
				Reference: errorID,
			}
		case *ConflictError:
			return http.StatusConflict, ErrorResponse{
				Error:     e.Error(),
				Tip:       "Please check if the " + e.ResourceType,
				Reference: errorID,
			}
		default:
			status := http.StatusInternalServerError
			msg := "Internal server error"
			if h.cfg.ShowInternal {
				msg = e.Error()
			}
			return status, ErrorResponse{
				Error:     msg,
				Tip:       "Please try again later",
				Reference: errorID,
			}
		}
	}
}

func (h *ErrorHandler) GinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		for _, ginErr := range c.Errors {
			status, response := h.Handle(c.Request.Context(), ginErr.Err)
			c.JSON(status, response)
			return
		}
	}
}

// Генерация уникального ID ошибки
func generateErrorID() string {
	return uuid.New().String()
}

// Очистка сообщения об ошибке от чувствительных данных
func sanitizeError(err error) string {
	message := err.Error()
	// логика для очистки чувствительных данных
	// Например, замена паролей или токенов на "[REDACTED]"
	return message
}
