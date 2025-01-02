package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/goodfoodcesi/tracking-api/pkg/config"
	"github.com/goodfoodcesi/tracking-api/pkg/logging"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	gintrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gin-gonic/gin"
)

func SetupApi(loadConfig config.Config) *gin.Engine {
	r := gin.New()
	r.Use(gintrace.Middleware("tracking-api"))
	logging.SetupHttpLogging(r, loadConfig.Env)

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"message": "Not found"})
	})
	api := r.Group("/tracking-api")
	SetupRoutes(api)

	return r
}

func SetupRoutes(api *gin.RouterGroup) {
	api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})
	api.GET("/health", healthCheck)
}

// @Summary     Health check endpoint
// @Description Get API health status
// @Tags        health
// @Accept      json
// @Produce     json
// @Success     200 {object} string
// @Router      /health [get]
func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "UP"})
}
