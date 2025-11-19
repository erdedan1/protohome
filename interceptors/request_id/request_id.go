package requestid

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const Header = "x-request-id"

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			md = metadata.New(nil)
		}

		ids := md.Get(Header)

		if len(ids) == 0 {
			rid := uuid.New().String()
			md.Set(Header, rid)
		}

		ctx = metadata.NewIncomingContext(ctx, md)

		return handler(ctx, req)
	}
}
