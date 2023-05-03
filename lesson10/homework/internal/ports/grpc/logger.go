package grpc

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

const (
	green   = "\033[97;42m"
	red     = "\033[97;41m"
	magenta = "\033[97;45m"
	cyan    = "\033[97;46m"
	reset   = "\033[0m"
)

func UnaryLogInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	start := time.Now()

	h, err := handler(ctx, req)

	st, _ := status.FromError(err)
	log.SetPrefix("[ADAPP] (req) - ")
	log.Printf("- <"+errCodeColor(err)+"%d"+reset+"> - meth:"+cyan+"%s"+reset+" - dur:"+magenta+"%s"+reset,
		st.Code(),
		info.FullMethod,
		time.Since(start),
	)

	return h, err
}

func errCodeColor(err error) string {
	if err != nil {
		return red
	}
	return green
}
