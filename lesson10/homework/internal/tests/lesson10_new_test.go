package tests

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"homework10/internal/adapters/adrepo"
	"homework10/internal/adapters/userrepo"
	"homework10/internal/app"
	grpcPort "homework10/internal/ports/grpc"
	"log"
	"testing"
)

func FuzzCreateUserHTTP(f *testing.F) {
	client := getTestHTTPClient()
	f.Add("Jenny", "jenny@mail.ru")
	f.Fuzz(func(t *testing.T, name, email string) {
		_, err := client.createUser(name, email)
		if err != nil && !errors.Is(err, ErrBadRequest) {
			t.Errorf("Err\texpect: %s, got: %s", ErrBadRequest.Error(), err.Error())
		}
	})
}

func FuzzCreateUserGRPC(f *testing.F) {
	client, closer := setupGRPC(50054)
	defer closer()

	f.Add("Jenny", "jenny@gmail.com")
	f.Fuzz(func(t *testing.T, name, email string) {
		_, err := client.CreateUser(context.Background(), &grpcPort.CreateUserRequest{Nickname: name, Email: email})
		s, _ := status.FromError(err)
		if err != nil && s.Code() != codes.InvalidArgument {
			t.Logf("Err expect: %s, got: %s", codes.InvalidArgument, s.Code())
		}
	})
}

func BenchmarkCreateUserHTTP(b *testing.B) {
	client := getTestHTTPClient()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = client.createUser("jenny", "jenny@gmail.com")
	}
}

func BenchmarkCreateUserGRPC(b *testing.B) {
	client, closer := setupGRPC(50054)
	defer closer()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = client.CreateUser(context.Background(), &grpcPort.CreateUserRequest{Nickname: "Jenny", Email: "jenny@gmail.com"})
	}
}

func setupGRPC(port uint) (grpcPort.AdServiceClient, func()) {
	lis := bufconn.Listen(1024 * 1024)
	server := grpc.NewServer()
	service := grpcPort.NewService(app.NewApp(adrepo.New(), userrepo.New()))
	grpcPort.RegisterAdServiceServer(server, service)

	go func() {
		if err := server.Serve(lis); err != nil {
			log.Printf("error serving server: %v", err)
		}
	}()

	conn, err := grpc.Dial("localhost"+fmt.Sprintf(":%d", port), grpc.WithTransportCredentials(insecure.NewCredentials()))
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
