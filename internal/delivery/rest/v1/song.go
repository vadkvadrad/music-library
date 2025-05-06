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
	song := api.Group("/song")
	song.Use(middleware.AuthMiddleware(h.config))
	{
		song.POST("/:album_id", h.AddSong())
	}
}


func (h *Handler) AddSong() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		albumID, err := strconv.Atoi(ctx.Param("album_id"))
		if err != nil {
			h.logger.Debug("Invalid album ID format",
                "received", ctx.Param("album_id"),
                "error", err,
            )
			ctx.Error(&er.ValidationError{Message: err.Error()})
			return
		}

		var body request.NewSongRequest
		if err := ctx.ShouldBindJSON(&body); err != nil {
			ctx.Error(err)
			return
		}

		user, ok := middleware.GetUserData(ctx)
		if !ok {
			h.logger.Debug("User credentials not found",
                "error", er.ErrWrongUserCredentials.Error(),
            )
			ctx.Error(er.ErrWrongUserCredentials)
			return
		}

		album, err := h.services.Album.GetArtistAlbum(ctx, user.Id, uint(albumID))
		if err != nil {
			ctx.Error(err)
			return
		}

		h.logger.Infow("Adding new song",
            "song_title", body.Title,
            "album_id", album.ID,
            "artist_id", album.ArtistID,
        )

		song, err := h.services.Song.AddSong(ctx, album, body)
		if err != nil {
			ctx.Error(err)
			return
		}

		// Добавить разрешение для песен
		err = h.services.Permission.AddPermission(ctx, user.Id, song.ID, model.SongResource, model.EditPermission)
		if err != nil {
			ctx.Error(err)
			return
		}

		ctx.JSON(http.StatusCreated, response.AddSongResponse{
			AlbumID: album.ID,
			AlbumName: album.Title,
			ArtistID: album.ArtistID,
		})
	}
}


