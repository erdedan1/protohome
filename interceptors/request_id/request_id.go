package requestid

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const header = "x-request-id"

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

		ids := md.Get(header)

		if len(ids) == 0 {
			rid = uuid.New().String()
			md.Set(header, rid)
		}

		ctx = metadata.NewIncomingContext(ctx, md)
		ctx = context.WithValue(ctx, header, rid)

		return handler(ctx, req)
	}
}
