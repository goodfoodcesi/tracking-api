package main

import (
	"github.com/gin-gonic/gin"
	"github.com/goodfoodcesi/tracking-api/pkg/config"
	"github.com/goodfoodcesi/tracking-api/pkg/logging"
	"go.uber.org/zap"
	gintrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gin-gonic/gin"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"io"
	"net/http"
)

func main() {
	loadConfig := config.LoadConfig()

	if loadConfig.Env != "dev" {
		tracer.Start(
			tracer.WithService("tracking-api"),
			tracer.WithEnv(loadConfig.Env),
			tracer.WithServiceVersion("0.0.5"),
		)
		defer tracer.Stop()
		gin.DefaultWriter = io.Discard
	}

	r := gin.New()
	r.Use(gintrace.Middleware("tracking-api"))
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"message": "Not found"})
	})
	api := r.Group("/tracking-api")

	logger := logging.SetupLogging(r, loadConfig.Env)
	defer logger.Sync()

	api.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	if err := r.Run(); err != nil {
		logger.Fatal("Cannot run API", zap.Error(err))
	}
}
