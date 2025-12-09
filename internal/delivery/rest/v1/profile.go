package v1

import (
	"music-lib/internal/dto/request"
	"music-lib/internal/middleware"
	"music-lib/internal/model"
	"music-lib/pkg/er"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *Handler) initProfileRoutes(api *gin.RouterGroup) {
	profile := api.Group("/profile")
	profile.Use(middleware.AuthMiddleware(h.config))
	{
		profile.POST("", h.NewProfile())
		profile.GET("", h.GetProfile())
	}
}

func (h *Handler) NewProfile() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		body := request.NewProfileRequest{}

		if err := ctx.ShouldBindJSON(&body); err != nil {
			ctx.Error(er.ValidationError{Message: err.Error()})
			return
		}

		user, ok := middleware.GetUserData(ctx)
		if !ok {
			ctx.Error(er.ErrNotAuthorized)
			return
		}

		err := h.services.Profile.NewProfile(ctx, body, user.Id)
		if err != nil {
			ctx.Error(err)
			return
		}

		ctx.JSON(http.StatusCreated, nil)
	}
}

func (h *Handler) GetProfile() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, ok := middleware.GetUserData(ctx)
		if !ok {
			ctx.Error(er.ErrNotAuthorized)
			return
		}

		profile, err := h.services.Profile.GetProfile(ctx, user.Id)
		if err != nil {
			ctx.Error(err)
			return
		}

		// Получаем данные пользователя
		userData, err := h.services.Profile.GetUserByID(user.Id)
		if err != nil {
			ctx.Error(err)
			return
		}

		// Проверяем, есть ли у пользователя артист
		artist, _ := h.services.Profile.GetArtistByUserID(ctx, user.Id)

		response := map[string]interface{}{
			"user_id":    profile.UserID,
			"user_name":  userData.Name,
			"user_email": userData.Email,
			"bio":        profile.Bio,
			"avatar_url": profile.AvatarURL,
			"created_at": profile.CreatedAt,
			"updated_at": profile.UpdatedAt,
		}

		if artist != nil {
			// Загружаем альбомы артиста
			albums, err := h.services.Artist.GetAlbums(ctx, artist.ID)
			if err != nil {
				albums = []model.Album{}
			}

			// Преобразуем альбомы в map для JSON
			albumDTOs := make([]map[string]interface{}, len(albums))
			for i, album := range albums {
				albumDTOs[i] = map[string]interface{}{
					"id":            album.ID,
					"title":         album.Title,
					"release_date":  album.ReleaseDate,
					"cover_art_url": album.CoverArtURL,
					"artist_id":     album.ArtistID,
				}
			}

			response["artist"] = map[string]interface{}{
				"id":             artist.ID,
				"name":           artist.Name,
				"description":    artist.Description,
				"formation_year": artist.FormationYear,
				"albums":         albumDTOs,
			}
		}

		ctx.JSON(http.StatusOK, response)
	}
}
