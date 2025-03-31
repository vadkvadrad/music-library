package middleware

import (
	"music-lib/internal/config"
	"music-lib/pkg/er"
	"music-lib/pkg/jwt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type UserData struct {
	Id    uint
	Email string
	Role  string
}

type contextKey string

const (
	ContextUserDataKey contextKey = "userData"
)

// New вернет middleware с инжектированной конфигурацией
func AuthMiddleware(config *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			abortWithUnauthorized(c)
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		isValid, data := jwt.NewJwt(config.Auth.Secret).Parse(token)
		if !isValid {
			abortWithUnauthorized(c)
			return
		}

		// Записываем данные пользователя в контекст Gin
		c.Set(string(ContextUserDataKey), UserData{
			Id:    data.Id,
			Email: data.Email,
			Role:  data.Role,
		})

		c.Next()
	}
}

// Вспомогательная функция для обработки неавторизованных запросов
func abortWithUnauthorized(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"error": er.ErrNotAuthorized.Message,
	})
}

// Вспомогательная функция для получения данных пользователя из контекста
func GetUserData(c *gin.Context) (UserData, bool) {
	val, exists := c.Get(string(ContextUserDataKey))
	if !exists {
		return UserData{}, false
	}
	
	userData, ok := val.(UserData)
	return userData, ok
}