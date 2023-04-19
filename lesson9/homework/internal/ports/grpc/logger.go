package grpc

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
)

func UnaryLogInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		start := time.Now()

		h, err := handler(ctx, req)

		log.SetPrefix("[GRPC] - ")
		log.Printf("Request - Method:%s\tDuration:%s\tError:%v\n",
			info.FullMethod,
			time.Since(start),
			err,
		)

		return h, err
	}
}
