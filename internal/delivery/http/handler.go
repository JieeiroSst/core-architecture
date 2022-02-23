package http

import (
	"net/http"

	"github.com/JIeeiroSst/core-backend/internal/config"
	v1 "github.com/JIeeiroSst/core-backend/internal/delivery/http/v1"
	"github.com/JIeeiroSst/core-backend/internal/usecase"
	"github.com/JIeeiroSst/core-backend/pkg/auth"
	"github.com/JIeeiroSst/core-backend/pkg/limiter"
	"github.com/gin-gonic/gin"
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

func (h *Handler) Init(cfg *config.Config) *gin.Engine {
	// Init gin handler
	router := gin.Default()

	router.Use(
		gin.Recovery(),
		gin.Logger(),
		limiter.Limit(cfg.Limiter.RPS, cfg.Limiter.Burst, cfg.Limiter.TTL),
		corsMiddleware,
	)

	// docs.SwaggerInfo.Host = fmt.Sprintf("%s:%s", cfg.HTTP.Host, cfg.HTTP.Port)
	// if cfg.Environment != config.EnvLocal {
	// 	docs.SwaggerInfo.Host = cfg.HTTP.Host
	// }

	// if cfg.Environment != config.Prod {
	// 	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// }

	// Init router
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	h.initAPI(router)

	return router
}

func (h *Handler) initAPI(router *gin.Engine) {
	handlerV1 := v1.NewHandler(h.usecase, h.tokenManager)
	api := router.Group("/api")
	{
		handlerV1.Init(api)
	}
}
