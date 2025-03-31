package v1

import (
	"music-lib/internal/dto/request"
	"music-lib/internal/dto/response"
	"music-lib/internal/middleware"
	"music-lib/internal/model"
	"music-lib/pkg/er"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *Handler) initSongRoutes(api *gin.RouterGroup) {
	songInfo := api.Group("/info")
	{
		songInfo.GET("", h.Get())
		songInfo.GET("/group", h.GetByGroup())
	}

	song := api.Group("/song")
	song.Use(middleware.AuthMiddleware(h.config))
	{
		song.POST("", h.Add())
		song.PATCH("/:id", h.Update())
		song.DELETE("/:id", h.Delete())
	}
}

// Add @Summary Добавить новую песню
// @Description Создает новую запись о песне. Требует авторизации.
// @Tags songs
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param input body request.AddSongRequest true "Данные песни"
// @Success 201 {object} model.Song "Созданная песня"
// @Failure 400 {object} map[string]string "Неверные данные/ошибка авторизации"
// @Failure 500 {object} map[string]string "Ошибка сервера"
// @Router /song [post]
func (h *Handler) Add() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var body request.AddSongRequest

		if err := ctx.ShouldBindJSON(&body); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		user, ok := middleware.GetUserData(ctx)
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": er.ErrWrongUserCredentials.Message})
			return
		}

		song, err := h.services.Song.Add(&model.Song{
			Group:       body.Group,
			Song:        body.Song,
			Text:        body.Text,
			ReleaseDate: body.ReleaseDate,
			Link:        body.Link,
			Owner:       user.Email,
		})

		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusCreated, song)
	}
}

// Get @Summary Получить информацию о песне
// @Description Возвращает данные песни по названию и группе
// @Tags songs
// @Accept json
// @Produce json
// @Param input body request.GetSongRequest true "Данные для поиска"
// @Success 200 {object} response.GetSongResponse
// @Failure 400 {object} map[string]string "Неверный формат запроса"
// @Failure 404 {object} map[string]string "Песня не найдена"
// @Router /info [get]
func (h *Handler) Get() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var body request.GetSongRequest

		if err := ctx.ShouldBindJSON(&body); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		song, err := h.services.Song.Get(body.Song, body.Group)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		}

		ctx.JSON(http.StatusOK, response.GetSongResponse{
			ReleaseDate: song.ReleaseDate,
			Text:        song.Text,
		})
	}
}

// GetByGroup @Summary Получить песни по группе
// @Description Возвращает список песен указанной группы с пагинацией
// @Tags songs
// @Accept json
// @Produce json
// @Param input body request.GetSongsRequest true "Название группы"
// @Param limit query int false "Лимит (по умолчанию 10)" minimum(1) maximum(100)
// @Param offset query int false "Смещение (по умолчанию 0)" minimum(0)
// @Success 200 {array} model.Song
// @Failure 400 {object} map[string]string "Неверный формат параметров"
// @Failure 404 {object} map[string]string "Группа не найдена"
// @Router /info/group [get]
func (h *Handler) GetByGroup() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var body request.GetSongsRequest

		if err := ctx.ShouldBindJSON(&body); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		limit := ctx.DefaultQuery("limit", "10")     // ?limit=20
		offset := ctx.DefaultQuery("offset", "0")    // ?offset=100


		limitInt, err := strconv.Atoi(limit)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit"})
			return
		}

		offsetInt, err := strconv.Atoi(offset)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid offset"})
			return
		}

		songs := h.services.Song.GetByGroup(body.Group, limitInt, offsetInt)

		ctx.JSON(http.StatusOK, songs)
	}
}

// Update @Summary Обновить информацию о песне
// @Description Обновляет данные песни по её ID. Требует авторизации.
// @Tags songs
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "ID песни"
// @Param input body request.UpdateSongRequest true "Данные для обновления"
// @Success 200 {object} model.Song "Успешно обновленная песня"
// @Failure 400 {object} map[string]string "Неверный формат ID/данных, ошибка авторизации"
// @Failure 404 {object} map[string]string "Песня не найдена"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /song/{id} [patch]
func (h *Handler) Update() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var body request.UpdateSongRequest

		if err := ctx.ShouldBindJSON(&body); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		user, ok := middleware.GetUserData(ctx)
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": er.ErrWrongUserCredentials.Message})
			return
		}

		song, err := h.services.Song.GetAndCompare(uint(id), user.Email)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		// TODO вынести логику в service
		song.Group = body.Group
		song.Song = body.Song
		song.Text = body.Text
		song.ReleaseDate = body.ReleaseDate
		song.Link = body.Link

		updatedSong, err := h.services.Song.Update(song)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, updatedSong)
	}
}

// Delete @Summary Удалить песню
// @Description Удаляет песню по ID. Требует прав владельца.
// @Tags songs
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "ID песни"
// @Success 200 {object} model.Song "Удаленная песня"
// @Failure 400 {object} map[string]string "Неверный формат ID"
// @Failure 401 {object} map[string]string "Ошибка авторизации"
// @Failure 404 {object} map[string]string "Песня не найдена"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /song/{id} [delete]
func (h *Handler) Delete() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}

		user, ok := middleware.GetUserData(ctx)
		if !ok {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": er.ErrWrongUserCredentials.Message})
			return
		}

		song, err := h.services.Song.GetAndCompare(uint(id), user.Email)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		err = h.services.Song.Delete(song.ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, song)
	}
}
