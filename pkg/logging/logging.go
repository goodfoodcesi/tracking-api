package logging

import (
	"context"
	"math/rand"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

const (
	RequestIDKey = "request-id"
	charset      = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	idLength     = 6
)

type LogManager struct {
	logger *zap.Logger
	env    string
}

func NewLogManager(env string) *LogManager {
	logger := InitLogger(env)
	return &LogManager{
		logger: logger,
		env:    env,
	}
}

func (lm *LogManager) addTraceInfo(ctx context.Context, fields []zap.Field) []zap.Field {
	if span, exists := tracer.SpanFromContext(ctx); exists {
		traceID := span.Context().TraceID()
		spanID := span.Context().SpanID()
		return append(fields,
			zap.Uint64("dd.trace_id", traceID),
			zap.Uint64("dd.span_id", spanID),
		)
	}
	return fields
}

func generateRequestID() string {
	b := make([]byte, idLength)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func (lm *LogManager) GrpcInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		requestID := generateRequestID()

		fields := []zap.Field{
			zap.String("request_id", requestID),
			zap.String("method", info.FullMethod),
		}

		fields = lm.addTraceInfo(ctx, fields)
		requestLogger := lm.logger.With(fields...)

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			md = metadata.MD{}
		}

		newMD := metadata.Join(md, metadata.Pairs(RequestIDKey, requestID))
		ctx = metadata.NewIncomingContext(ctx, newMD)
		ctx = WithContext(ctx, requestLogger)

		requestLogger.Info("Processing gRPC request")
		resp, err := handler(ctx, req)

		requestLogger.Info("Completed gRPC request",
			zap.Duration("duration", time.Since(start)),
			zap.Error(err))

		return resp, err
	}
}

func (lm *LogManager) SetupHttp(r *gin.Engine) {
	r.Use(func(c *gin.Context) {
		c.Set("logger", lm.logger)
		c.Next()
	})

	r.Use(lm.httpRequestIDMiddleware())

	r.Use(ginzap.GinzapWithConfig(lm.logger, &ginzap.Config{
		UTC:        true,
		TimeFormat: time.RFC3339,
		Context: func(c *gin.Context) []zap.Field {
			fields := []zap.Field{
				zap.String("request_id", c.GetString(RequestIDKey)),
			}
			return lm.addTraceInfo(c.Request.Context(), fields)
		},
	}))

	r.Use(ginzap.RecoveryWithZap(lm.logger, true))
	lm.logger.Info("HTTP logging initialized", zap.String("env", lm.env))
}

func (lm *LogManager) SetupGrpcLogging() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		requestID := generateRequestID()
		fields := []zap.Field{
			zap.String("request_id", requestID),
			zap.String("method", info.FullMethod),
		}
		fields = lm.addTraceInfo(ctx, fields)
		requestLogger := lm.logger.With(fields...)

		requestLogger.Info("Processing gRPC request")
		resp, err := handler(ctx, req)

		requestLogger.Info("Completed gRPC request",
			zap.Duration("duration", time.Since(start)),
			zap.Error(err))

		return resp, err
	}
}

func (lm *LogManager) httpRequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := generateRequestID()
		c.Set(RequestIDKey, requestID)
		c.Header("X-Request-ID", requestID)

		fields := []zap.Field{zap.String("request_id", requestID)}
		fields = lm.addTraceInfo(c.Request.Context(), fields)

		requestLogger := lm.logger.With(fields...)
		ctx := WithContext(c.Request.Context(), requestLogger)
		c.Request = c.Request.WithContext(ctx)
		c.Set("logger", requestLogger)

		c.Next()
	}
}
