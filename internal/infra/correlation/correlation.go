package correlation

import "context"

type contextKey string

const CorrelationIDKey contextKey = "correlation_id"

func WithCorrelationID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, CorrelationIDKey, id)
}

func FromContext(ctx context.Context) string {
	id, ok := ctx.Value(CorrelationIDKey).(string)
	if !ok {
		return "unknown"
	}
	return id
}
