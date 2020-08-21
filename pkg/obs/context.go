package obs

import (
	"context"
)

type telemetryContextKey struct{}

// TelemetryFromContext get from context
func TelemetryFromContext(ctx context.Context) *Telemetry {
	return ctx.Value(telemetryContextKey{}).(*Telemetry)
}

// WithTelemetry inject telemetry
func WithTelemetry(ctx context.Context, telemetry *Telemetry) context.Context {
	return context.WithValue(ctx, telemetryContextKey{}, telemetry)
}
