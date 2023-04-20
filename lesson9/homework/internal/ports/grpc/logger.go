package grpc

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
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

	log.SetPrefix("[ADAPP] (req) - ")
	log.Printf("- meth:"+cyan+"%s"+reset+" - dur:"+magenta+"%s"+reset+" - err:"+errColor(err)+"%v"+reset+"\n",
		info.FullMethod,
		time.Since(start),
		err,
	)

	return h, err
}

func errColor(err error) string {
	if err != nil {
		return red
	}
	return green
}
