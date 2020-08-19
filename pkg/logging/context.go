package logging

import (
	"context"

	"github.com/sirupsen/logrus"
)

type logContextKey struct{}

// LoggerFromContext get logger from context
func LoggerFromContext(ctx context.Context) *logrus.Entry {
	return ctx.Value(logContextKey{}).(*logrus.Entry)
}

// WithLogger inject logger
func WithLogger(ctx context.Context, logger *logrus.Entry) context.Context {
	return context.WithValue(ctx, logContextKey{}, logger)
}
