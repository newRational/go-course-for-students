package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"

	"homework10/internal/adapters/adrepo"
	"homework10/internal/adapters/userrepo"
	"homework10/internal/app"
	grpcPort "homework10/internal/ports/grpc"
)

const port = ":50054"

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer(grpc.ChainUnaryInterceptor(grpcPort.UnaryLogInterceptor, recovery.UnaryServerInterceptor()))
	service := grpcPort.NewService(app.NewApp(adrepo.New(), userrepo.New()))
	grpcPort.RegisterAdServiceServer(server, service)

	eg, ctx := errgroup.WithContext(context.Background())

	sigQuit := make(chan os.Signal, 1)
	signal.Ignore(syscall.SIGHUP, syscall.SIGPIPE)
	signal.Notify(sigQuit, syscall.SIGINT, syscall.SIGTERM)

	eg.Go(func() error {
		select {
		case s := <-sigQuit:
			log.Printf("captured signal: %v\n", s)
			return fmt.Errorf("captured signal: %v", s)
		case <-ctx.Done():
			return nil
		}
	})

	eg.Go(func() error {
		log.Printf("starting grpc server, listening on %s\n", port)
		defer log.Printf("close grpc server listening on %s\n", port)

		errCh := make(chan error)

		defer func() {
			server.GracefulStop()
			_ = lis.Close()

			close(errCh)
		}()

		go func() {
			if err = server.Serve(lis); err != nil {
				errCh <- err
			}
		}()

		select {
		case <-ctx.Done():
			return ctx.Err()
		case err = <-errCh:
			return fmt.Errorf("grpc server can't listen and serve requests: %w", err)
		}
	})

	if err = eg.Wait(); err != nil {
		log.Printf("gracefully shutting down the servers: %s\n", err.Error())
	}

	log.Println("service was successfully shutdown")
}
