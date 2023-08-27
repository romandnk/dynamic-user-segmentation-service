package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/romandnk/dynamic-user-segmentation-service/internal/logger"
	"github.com/romandnk/dynamic-user-segmentation-service/internal/service"
)

type Handler struct {
	engine   *gin.Engine
	services service.Services
	logger   logger.Logger
}

func NewHandler(services service.Services, logger logger.Logger) *Handler {
	return &Handler{
		services: services,
		logger:   logger,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.Use(h.loggerMiddleware())
	gin.SetMode(gin.ReleaseMode)
	h.engine = router

	api := router.Group("/api")
	{
		version := api.Group("/v1")
		{
			segments := version.Group("/segments")
			{
				segments.POST("/", h.CreateSegment)
				segments.DELETE("/", h.DeleteSegment)
			}
		}
	}

	return router
}
