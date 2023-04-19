package tests

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"

	"homework9/internal/adapters/adrepo"
	"homework9/internal/adapters/userrepo"
	"homework9/internal/app"
	grpcPort "homework9/internal/ports/grpc"
)

func TestGRRPCCreateUser(t *testing.T) {
	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() {
		lis.Close()
	})

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpcPort.UnaryLogInterceptor(),
			recovery.UnaryServerInterceptor(),
		),
	)
	t.Cleanup(func() {
		s.Stop()
	})

	srv := grpcPort.NewService(app.NewApp(adrepo.New(), userrepo.New()))
	grpcPort.RegisterAdServiceServer(s, srv)

	go func() {
		assert.NoError(t, s.Serve(lis), "s.Serve")
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(func() {
		cancel()
	})

	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(dialer), grpc.WithInsecure())
	assert.NoError(t, err, "grpc.DialContext")

	t.Cleanup(func() {
		conn.Close()
	})

	client := grpcPort.NewAdServiceClient(conn)
	res, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Nickname: "Oleg", Email: "oleg@gmail.com"})
	assert.NoError(t, err, "client.GetUser")

	assert.Equal(t, "Oleg", res.Nickname)
	assert.Equal(t, "oleg@gmail.com", res.Email)
}
