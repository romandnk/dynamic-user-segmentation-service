package v1

import (
	"github.com/gin-gonic/gin"
	_ "github.com/romandnk/dynamic-user-segmentation-service/docs"
	"github.com/romandnk/dynamic-user-segmentation-service/internal/logger"
	"github.com/romandnk/dynamic-user-segmentation-service/internal/service"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	engine        *gin.Engine
	services      service.Services
	logger        logger.Logger
	pathToReports string
}

func NewHandler(services service.Services, logger logger.Logger, pathToReports string) *Handler {
	return &Handler{
		services:      services,
		logger:        logger,
		pathToReports: pathToReports,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.Use(h.loggerMiddleware())
	gin.SetMode(gin.ReleaseMode)
	h.engine = router

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api")
	{
		version := api.Group("/v1")
		{
			segments := version.Group("/segments")
			{
				segments.POST("/", h.CreateSegment)
				segments.DELETE("/", h.DeleteSegment)
			}

			users := version.Group("/users")
			{
				users.POST("/", h.UpdateUserSegments)
				users.POST("/active_segments", h.GetActiveUserSegments)

				report := users.Group("/report")
				{
					report.POST("/", h.CreateCSVReportAndURL)
					report.GET("/:id", h.GetReportByID)
				}
			}
		}
	}

	return router
}
