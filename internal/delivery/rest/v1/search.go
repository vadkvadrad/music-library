package v1

import (
	"music-lib/pkg/er"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func (h *Handler) initSearchRoutes(api *gin.RouterGroup) {
	search := api.Group("/search")
	{
		search.GET("", h.Search())
	}
}

func (h *Handler) Search() gin.HandlerFunc {
	return func(c *gin.Context) {
		query := c.Query("q")
		types := strings.Split(c.DefaultQuery("type", "artist,album,song"), ",")
		limit, offset := validatePagination(c)

		validTypes := map[string]bool{"artist": true, "album": true, "song": true}

		for _, t := range types {
			if !validTypes[t] {
				c.Error(er.ValidationError{Message: "invalid type parameter"})
				return
			}
		}

		result := h.services.Search.Search(c, types, query, limit, offset)

		c.JSON(http.StatusOK, result)
	}
}


func validatePagination(c *gin.Context) (limit, offset int) {
    limitStr := c.DefaultQuery("limit", "10")
    offsetStr := c.DefaultQuery("offset", "0")

    // Валидация параметров
    limit, err := strconv.Atoi(limitStr)
    if err != nil || limit < 1 || limit > 100 {
		c.Error(er.ValidationError{Message: "invalid limit value (1-100)"})
        return
    }

    offset, err = strconv.Atoi(offsetStr)
    if err != nil || offset < 0 {
		c.Error(er.ValidationError{Message: "invalid offset value"})
        return
    }

	return limit, offset
}
