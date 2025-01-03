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

func SetupApi(cfg config.Config, logManager *logging.LogManager) *gin.Engine {
	r := gin.New()

	if cfg.Env != "dev" {
		r.Use(gintrace.Middleware("tracking-api"))
	}

	logManager.SetupHttp(r)

	r.NoRoute(func(c *gin.Context) {
		logger := logging.FromContext(c.Request.Context())
		logger.Warn("Route not found", "path", c.Request.URL.Path)
		c.JSON(http.StatusNotFound, gin.H{"message": "Not found"})
	})

	api := r.Group("/tracking-api")
	api.Use(JWTInterceptor(cfg))
	SetupRoutes(api)

	return r
}

func SetupRoutes(api *gin.RouterGroup) {
	// Documentation Swagger
	api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Routes de base
	api.GET("/ping", ping)
	api.GET("/health", healthCheck)
}

func ping(c *gin.Context) {
	logger := logging.FromContext(c.Request.Context())
	logger.Debug("Ping request received")
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}

// @Summary     Health check endpoint
// @Description Get API health status
// @Tags        health
// @Accept      json
// @Produce     json
// @Success     200 {object} string
// @Router      /health [get]
func healthCheck(c *gin.Context) {
	logger := logging.FromContext(c.Request.Context())
	logger.Debug("Health check request received")
	c.JSON(http.StatusOK, gin.H{"status": "UP"})
}
