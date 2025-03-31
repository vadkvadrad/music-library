package v1

import (
	"music-lib/internal/dto/request"
	"music-lib/internal/dto/response"
	"music-lib/pkg/jwt"
	"net/http"

	"github.com/gin-gonic/gin"
)


// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Введите "Bearer [ваш_токен]"
func (h *Handler) initAuthRoutes(api *gin.RouterGroup) {
	auth := api.Group("/auth")
	{
		auth.POST("/login", h.Login())
		auth.POST("/register", h.Register())
		auth.POST("/verify", h.Verify())
	}
}


// Login @Summary Аутентификация пользователя
// @Description Вход в систему с email и паролем
// @Tags auth
// @Accept json
// @Produce json
// @Param input body request.LoginRequest true "Данные для входа"
// @Success 200 {object} response.LoginResponse
// @Failure 400 {object} map[string]string "Неверный формат данных"
// @Failure 401 {object} map[string]string "Неверные учетные данные"
// @Failure 500 {object} map[string]string "Ошибка генерации токена"
// @Router /auth/login [post]
func (handler *Handler) Login() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var body request.LoginRequest

		if err := ctx.ShouldBindJSON(&body); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, err := handler.services.Auth.Login(body.Email, body.Password)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		token, err := jwt.NewJwt(handler.config.Auth.Secret).Create(jwt.JWTData{
			Id:    user.ID,
			Email: user.Email,
			Role:  string(user.Role),
		})

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// отправить ответ
		data := response.LoginResponse{
			Token: token,
		}
		ctx.JSON(http.StatusOK, data)
	}
}


// Register @Summary Регистрация нового пользователя
// @Description Создание нового аккаунта с подтверждением по email
// @Tags auth
// @Accept json
// @Produce json
// @Param input body request.RegisterRequest true "Данные для регистрации"
// @Success 201 {object} response.RegisterResponse
// @Failure 400 {object} map[string]string "Неверные данные/пользователь существует"
// @Router /auth/register [post]
func (handler Handler) Register() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var body request.RegisterRequest

		if err := ctx.ShouldBindJSON(&body); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		sessionId, err := handler.services.Auth.Register(body.Email, body.Password, body.Name)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// отправить ответ
		data := response.RegisterResponse{
			SessionId: sessionId,
		}
		ctx.JSON(http.StatusCreated, data)
	}
}


// Verify @Summary Подтверждение регистрации
// @Description Верификация email с кодом подтверждения
// @Tags auth
// @Accept json
// @Produce json
// @Param input body request.VerifyRequest true "Данные для верификации"
// @Success 202 {object} response.VerifyResponse
// @Failure 400 {object} map[string]string "Неверный код/сессия"
// @Failure 500 {object} map[string]string "Ошибка генерации токена"
// @Router /auth/verify [post]
func (handler *Handler) Verify() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var body request.VerifyRequest

		if err := ctx.ShouldBindJSON(&body); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		user, err := handler.services.Auth.Verify(body.SessionId, body.Code)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		token, err := jwt.NewJwt(handler.config.Auth.Secret).Create(jwt.JWTData{
			Id:    user.ID,
			Email: user.Email,
			Role:  string(user.Role),
		})
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		data := response.VerifyResponse{
			Token: token,
		}
		ctx.JSON(http.StatusAccepted, data)
	}
}
