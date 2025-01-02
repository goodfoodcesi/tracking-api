package logging

import (
	"context"

	"go.uber.org/zap"
)

// Logger est une interface qui définit les méthodes de logging communes
type Logger interface {
	Info(msg string, fields ...any)
	Error(msg string, fields ...any)
	Debug(msg string, fields ...any)
	Warn(msg string, fields ...any)
}

// loggerKey est la clé utilisée pour stocker le logger dans le context
type loggerKey struct{}

// FromContext récupère le logger du context
func FromContext(ctx context.Context) Logger {
	if logger, ok := ctx.Value(loggerKey{}).(*zap.Logger); ok {
		return &zapLogger{logger: logger}
	}
	return &zapLogger{logger: zap.L()} // Logger par défaut si aucun n'est trouvé dans le context
}

// WithContext ajoute un logger au context
func WithContext(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerKey{}, logger)
}

// zapLogger est l'implémentation de Logger utilisant zap
type zapLogger struct {
	logger *zap.Logger
}

// convertToZapFields convertit les arguments variadic en zap.Field
func convertToZapFields(args ...any) []zap.Field {
	fields := make([]zap.Field, 0, len(args)/2)
	for i := 0; i < len(args); i += 2 {
		if i+1 < len(args) {
			key, ok := args[i].(string)
			if !ok {
				continue
			}
			fields = append(fields, zap.Any(key, args[i+1]))
		}
	}
	return fields
}

func (l *zapLogger) Info(msg string, fields ...any) {
	l.logger.Info(msg, convertToZapFields(fields...)...)
}

func (l *zapLogger) Error(msg string, fields ...any) {
	l.logger.Error(msg, convertToZapFields(fields...)...)
}

func (l *zapLogger) Debug(msg string, fields ...any) {
	l.logger.Debug(msg, convertToZapFields(fields...)...)
}

func (l *zapLogger) Warn(msg string, fields ...any) {
	l.logger.Warn(msg, convertToZapFields(fields...)...)
}

func InitLogger(env string) *zap.Logger {
	var config zap.Config

	config = zap.NewProductionConfig()

	logger, err := config.Build()
	if err != nil {
		panic(err)
	}

	// Remplace le logger global
	zap.ReplaceGlobals(logger)
	return logger
}
