package tests

import (
	"context"
	"log"
	"net"
	"testing"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	"homework10/internal/adapters/adrepo"
	"homework10/internal/adapters/userrepo"
	"homework10/internal/app"
	grpcPort "homework10/internal/ports/grpc"
)

func FuzzCreateUserHTTP(f *testing.F) {
	client := getTestHTTPClient()
	f.Add("Jenny", "jenny@mail.ru")
	f.Fuzz(func(t *testing.T, name, email string) {
		_, err := client.createUser(name, email)
		if err != nil && !errors.Is(err, ErrBadRequest) {
			t.Logf("Err\texpect: %s, got: %s", ErrBadRequest.Error(), err.Error())
		}
	})
}

func FuzzCreateAdHTTP(f *testing.F) {
	client := getTestHTTPClient()
	f.Add(int64(0), "Title", "Text")
	f.Fuzz(func(t *testing.T, userID int64, name, email string) {
		_, err := client.createAd(userID, name, email)
		if err != nil && !errors.Is(err, ErrBadRequest) {
			t.Logf("Err\texpect: %s, got: %s", ErrBadRequest.Error(), err.Error())
		}
	})
}

func BenchmarkCreateUserHTTP(b *testing.B) {
	client := getTestHTTPClient()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = client.createUser("Jenny", "jenny@gmail.com")
	}
}

func BenchmarkCreateUserGRPC(b *testing.B) {
	client, closer := setupGRPC()
	defer closer()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = client.CreateUser(context.Background(), &grpcPort.CreateUserRequest{Nickname: "Jenny", Email: "jenny@gmail.com"})
	}
}

func setupGRPC() (grpcPort.AdServiceClient, func()) {
	lis := bufconn.Listen(1024 * 1024)
	server := grpc.NewServer()
	service := grpcPort.NewService(app.NewApp(adrepo.New(), userrepo.New()))
	grpcPort.RegisterAdServiceServer(server, service)

	go func() {
		if err := server.Serve(lis); err != nil {
			log.Printf("error serving server: %v", err)
		}
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}
	conn, err := grpc.Dial("", grpc.WithContextDialer(dialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("error connecting to server: %v", err)
	}

	closer := func() {
		if err = lis.Close(); err != nil {
			log.Printf("error closing listener: %v", err)
		}
		server.Stop()
	}

	client := grpcPort.NewAdServiceClient(conn)

	return client, closer
}
