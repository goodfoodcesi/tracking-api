package logging

import (
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
	"math/rand"
	"time"
)

const (
	RequestIDKey = "RequestID"
	charset      = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	idLength     = 6
)

func generateRequestID() string {
	b := make([]byte, idLength)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func RequestIDMiddleware() gin.HandlerFunc {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	return func(c *gin.Context) {
		requestID := generateRequestID()
		c.Set(RequestIDKey, requestID)
		c.Header("X-Request-ID", requestID)

		// Récupérer le logger de base
		logger := c.MustGet("logger").(*zap.Logger)

		// Récupérer le span du context
		span, exists := tracer.SpanFromContext(c.Request.Context())
		if exists {
			// Ajouter trace_id et span_id au logger
			traceID := span.Context().TraceID()
			spanID := span.Context().SpanID()

			requestLogger := logger.With(
				zap.String("request_id", requestID),
				zap.Uint64("dd.trace_id", traceID),
				zap.Uint64("dd.span_id", spanID),
			)
			c.Set("logger", requestLogger)
		} else {
			requestLogger := logger.With(zap.String("request_id", requestID))
			c.Set("logger", requestLogger)
		}

		c.Next()
	}
}

func SetupLogging(r *gin.RouterGroup, env string) *zap.Logger {
	var logger *zap.Logger

	gin.SetMode(gin.ReleaseMode)
	logger = zap.Must(zap.NewProduction())

	r.Use(func(c *gin.Context) {
		c.Set("logger", logger)
		c.Next()
	})

	r.Use(RequestIDMiddleware())

	r.Use(ginzap.GinzapWithConfig(logger, &ginzap.Config{
		UTC:        true,
		TimeFormat: time.RFC3339,
		Context: func(c *gin.Context) []zap.Field {
			fields := []zap.Field{
				zap.String("request_id", c.GetString(RequestIDKey)),
			}

			if span, exists := tracer.SpanFromContext(c.Request.Context()); exists {
				traceID := span.Context().TraceID()
				spanID := span.Context().SpanID()
				fields = append(fields,
					zap.Uint64("dd.trace_id", traceID),
					zap.Uint64("dd.span_id", spanID),
				)
			}

			return fields
		},
	}))

	r.Use(ginzap.RecoveryWithZap(logger, true))

	return logger
}
