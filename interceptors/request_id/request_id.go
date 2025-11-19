package requestid

import (
	"context"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const requestIDKey string = "x-request-id"

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		var rid string
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		}

		ids := md.Get(requestIDKey)

		if len(ids) == 0 {
			rid = uuid.New().String()
			md.Set(requestIDKey, rid)
		}

		ctx = metadata.NewIncomingContext(ctx, md)
		ctx = context.WithValue(ctx, requestIDKey, rid)
		span := trace.SpanFromContext(ctx)
		if span != nil && span.SpanContext().IsValid() {
			span.SetAttributes(attribute.String(string(requestIDKey), rid))
		}

		return handler(ctx, req)
	}
}
