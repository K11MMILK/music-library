package handler

import (
	"time-tracker/pkg/service"

	_ "time-tracker/docs"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	logrus.Info("Initializing routes")

	group := router.Group("/api/group")
	{
		group.POST("/", h.createGroup)
		group.GET("/", h.getAllGroups)
		group.DELETE("/:id", h.deleteGroup)
		group.PUT("/:id", h.updateGroup)
		group.GET("/filter", h.getGroupsWithFilter)
	}

	song := router.Group("/api/song")
	{
		song.POST("/", h.createSong)
		song.GET("/", h.getAllSongs)
		song.DELETE("/:id", h.deleteSong)
		song.PUT("/:id", h.updateSong)
		song.GET("/filter", h.getSongsWithFilter)
	}

	songDetails := router.Group("/api/songDetails")
	{
		songDetails.GET("/:id", h.getSongDetailsById)
		songDetails.PUT("/:id", h.updateSongDetails)
	}

	songText := router.Group("/api/songText")
	{
		songText.GET("/:id/filter", h.getSongText)
	}
	logrus.Info("Routes initialized successfully")
	return router
}

type errorResponse struct {
	Message string `json:"message"`
}
type statusResponse struct {
	Status string `json:"status"`
}
