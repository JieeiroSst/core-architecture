package v1

import (
	"errors"

	"github.com/JIeeiroSst/core-backend/internal/usecase"
	"github.com/JIeeiroSst/core-backend/pkg/auth"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Handler struct {
	usecase      *usecase.Usecase
	tokenManager auth.TokenManager
}

func NewHandler(usecase *usecase.Usecase, tokenManager auth.TokenManager) *Handler {
	return &Handler{
		usecase:      usecase,
		tokenManager: tokenManager,
	}
}

func (h *Handler) Init(api *gin.RouterGroup) {
	v1 := api.Group("/v1")
	{
		h.initUsersRoutes(v1)
		h.initCoursesRoutes(v1)
		h.initStudentsRoutes(v1)
		h.initCallbackRoutes(v1)
		h.initAdminRoutes(v1)

		v1.GET("/settings", h.setSchoolFromRequest, h.getSchoolSettings)
		v1.GET("/promocodes/:code", h.setSchoolFromRequest, h.getPromo)
		v1.GET("/offers/:id", h.setSchoolFromRequest, h.getOffer)
	}
}

func parseIdFromPath(c *gin.Context, param string) (primitive.ObjectID, error) {
	idParam := c.Param(param)
	if idParam == "" {
		return primitive.ObjectID{}, errors.New("empty id param")
	}

	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return primitive.ObjectID{}, errors.New("invalid id param")
	}

	return id, nil
}
