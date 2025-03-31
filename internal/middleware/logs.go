package middleware

import (
	"time"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)


func Logger() gin.HandlerFunc {
	logrus.SetFormatter(&logrus.JSONFormatter{}) 

	return func(c *gin.Context) {
		// Засекаем время начала обработки запроса
		startTime := time.Now()

		// Обрабатываем запрос
		c.Next()

		// Рассчитываем время выполнения
		duration := time.Since(startTime)

		// Собираем информацию для лога
		entry := logrus.WithFields(logrus.Fields{
			"client_ip":  c.ClientIP(),
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"proto":      c.Request.Proto,
			"status":     c.Writer.Status(),
			"duration":   duration,
			"user_agent": c.Request.UserAgent(),
			"bytes_in":   c.Request.ContentLength,
			"bytes_out":  c.Writer.Size(),
		})

		// Логируем с соответствующим уровнем
		if len(c.Errors) > 0 {
			entry.Error(c.Errors.String())
		} else {
			if c.Writer.Status() >= 500 {
				entry.Error("Server error")
			} else if c.Writer.Status() >= 400 {
				entry.Warn("Client error")
			} else {
				entry.Info("Request processed")
			}
		}
	}
}