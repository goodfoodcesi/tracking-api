package main

import (
	"context"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/goodfoodcesi/tracking-api/pkg/api"
	"github.com/goodfoodcesi/tracking-api/pkg/config"
	"github.com/goodfoodcesi/tracking-api/pkg/grpc"
	"github.com/goodfoodcesi/tracking-api/pkg/logging"
	"github.com/goodfoodcesi/tracking-api/pkg/server"
	"go.uber.org/zap"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func main() {
	cfg := config.LoadConfig()
	logManager := logging.NewLogManager(cfg.Env)

	if cfg.Env != "dev" {
		tracer.Start(
			tracer.WithService("tracking-api"),
			tracer.WithEnv(cfg.Env),
			tracer.WithServiceVersion("0.0.5"),
		)
		defer tracer.Stop()
		gin.DefaultWriter = io.Discard
	}

	server := setupServer(cfg, logManager)
	loggerFromCtx := logging.FromContext(context.Background())
	loggerFromCtx.Info("Starting server", zap.String("env", cfg.Env))

	if err := server.ListenAndServe(); err != nil {
		loggerFromCtx.Error("Failed to start server", zap.Error(err))
	}
}

func setupServer(cfg config.Config, logManager *logging.LogManager) *http.Server {
	return &http.Server{
		Addr: ":8080",
		Handler: server.NewMultiplexHandler(
			api.SetupApi(cfg, logManager),
			grpc.NewServer(logManager, cfg),
		),
	}
}
