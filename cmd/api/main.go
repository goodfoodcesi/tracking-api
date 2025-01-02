package main

import (
	"context"
	"crypto/tls"
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
	loadConfig := config.LoadConfig()

	// Initialisation du logger
	logger := logging.InitLogger(loadConfig.Env)
	defer logger.Sync()

	// Création du context avec le logger
	ctx := logging.WithContext(context.Background(), logger)

	if loadConfig.Env != "dev" {
		tracer.Start(
			tracer.WithService("tracking-api"),
			tracer.WithEnv(loadConfig.Env),
			tracer.WithServiceVersion("0.0.5"),
		)
		defer tracer.Stop()
		gin.DefaultWriter = io.Discard
	}

	// Initialisation des serveurs
	ginHandler := api.SetupApi(loadConfig)
	grpcServer := grpc.NewServer()

	// Configuration du multiplexage
	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: server.NewMultiplexHandler(ginHandler, grpcServer),
		TLSConfig: &tls.Config{
			NextProtos: []string{"h2", "http/1.1"},
		},
	}
	// Utilisation du logger depuis le context
	loggerFromCtx := logging.FromContext(ctx)
	loggerFromCtx.Info("Starting server", zap.String("env", loadConfig.Env))

	if err := httpServer.ListenAndServeTLS("server.crt", "server.key"); err != nil {
		loggerFromCtx.Error("Failed to start server", zap.Error(err))
	}
}
