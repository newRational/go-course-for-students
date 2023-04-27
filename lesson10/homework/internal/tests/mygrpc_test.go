package tests

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"

	grpcPort "homework10/internal/ports/grpc"
)

func TestGRPCCreateAd(t *testing.T) {
	ctx, client := getTestGRCPClient(t)

	_, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Nickname: "Jenny", Email: "jenny@gmail.com"})
	assert.NoError(t, err, "client.CreateUser")

	res, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{Title: "Title", Text: "Text", UserId: 0})
	assert.NoError(t, err, "client.CreateAd")

	assert.Zero(t, res.Id)
	assert.Equal(t, "Title", res.Title)
	assert.Equal(t, "Text", res.Text)
	assert.False(t, res.Published)
	assert.Zero(t, res.UserId)
}

func TestGRPCChangeAdStatus(t *testing.T) {
	ctx, client := getTestGRCPClient(t)

	_, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Nickname: "Jenny", Email: "jenny@gmail.com"})
	assert.NoError(t, err, "client.CreateUser")

	_, err = client.CreateAd(ctx, &grpcPort.CreateAdRequest{Title: "Title", Text: "Text", UserId: 0})
	assert.NoError(t, err, "client.CreateAd")

	res, err := client.ChangeAdStatus(ctx, &grpcPort.ChangeAdStatusRequest{AdId: 0, UserId: 0, Published: true})
	assert.NoError(t, err, "client.ChangeAdStatus")

	assert.Zero(t, res.Id)
	assert.Equal(t, "Title", res.Title)
	assert.Equal(t, "Text", res.Text)
	assert.True(t, res.Published)
	assert.Zero(t, res.UserId)
}

func TestGRPCUpdateAd(t *testing.T) {
	ctx, client := getTestGRCPClient(t)

	_, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Nickname: "Jenny", Email: "jenny@gmail.com"})
	assert.NoError(t, err, "client.CreateUser")

	_, err = client.CreateAd(ctx, &grpcPort.CreateAdRequest{Title: "Title", Text: "Text", UserId: 0})
	assert.NoError(t, err, "client.CreateAd")

	res, err := client.UpdateAd(ctx, &grpcPort.UpdateAdRequest{AdId: 0, Title: "New title", Text: "New text", UserId: 0})
	assert.NoError(t, err, "client.UpdateAd")

	assert.Zero(t, res.Id)
	assert.Equal(t, "New title", res.Title)
	assert.Equal(t, "New text", res.Text)
	assert.False(t, res.Published)
	assert.Zero(t, res.UserId)
}

func TestGRPCGetAd(t *testing.T) {
	ctx, client := getTestGRCPClient(t)

	_, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Nickname: "Jenny", Email: "jenny@gmail.com"})
	assert.NoError(t, err, "client.CreateUser")

	_, err = client.CreateAd(ctx, &grpcPort.CreateAdRequest{Title: "First title", Text: "First text", UserId: 0})
	assert.NoError(t, err, "client.CreateAd")

	_, err = client.CreateAd(ctx, &grpcPort.CreateAdRequest{Title: "Second title", Text: "Second text", UserId: 0})
	assert.NoError(t, err, "client.CreateAd")

	res, err := client.GetAd(ctx, &grpcPort.GetAdRequest{Id: 1})
	assert.NoError(t, err, "client.GetAd")

	assert.Equal(t, int64(1), res.Id)
	assert.Equal(t, "Second title", res.Title)
	assert.Equal(t, "Second text", res.Text)
	assert.False(t, res.Published)
	assert.Zero(t, res.UserId)
}

func TestGRPCListAds(t *testing.T) {
	ctx, client := getTestGRCPClient(t)

	_, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Nickname: "Jenny", Email: "jenny@gmail.com"})
	assert.NoError(t, err, "client.CreateUser")

	_, err = client.CreateUser(ctx, &grpcPort.CreateUserRequest{Nickname: "Polly", Email: "polly@gmail.com"})
	assert.NoError(t, err, "client.CreateUser")

	tc := time.Now().UTC()

	ad0, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{Title: "hello", Text: "world", UserId: 0})
	assert.NoError(t, err, "client.CreateAd")

	_, err = client.CreateAd(ctx, &grpcPort.CreateAdRequest{Title: "hello", Text: "friend", UserId: 1})
	assert.NoError(t, err, "client.CreateAd")

	_, err = client.CreateAd(ctx, &grpcPort.CreateAdRequest{Title: "hey", Text: "привет", UserId: 0})
	assert.NoError(t, err, "client.CreateAd")

	res, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{Title: "goodbye", Text: "friend", UserId: 1})
	assert.NoError(t, err, "client.CreateAd")

	_, err = client.ChangeAdStatus(ctx, &grpcPort.ChangeAdStatusRequest{AdId: res.Id, UserId: res.UserId, Published: true})
	assert.NoError(t, err, "client.ChangeAdStatus")

	published := false
	title := "hello"
	userId := int64(0)
	created := timestamppb.Timestamp{Seconds: tc.Unix()}
	ads, err := client.ListAds(ctx, &grpcPort.ListAdsRequest{Published: &published, Title: &title, UserId: &userId, Created: &created})
	assert.NoError(t, err)

	assert.Len(t, ads.List, 1)
	assert.Equal(t, ads.List[0].Id, ad0.Id)
	assert.Equal(t, ads.List[0].Title, ad0.Title)
	assert.Equal(t, ads.List[0].Text, ad0.Text)
	assert.Equal(t, ads.List[0].UserId, ad0.UserId)
	assert.Equal(t, ads.List[0].Published, ad0.Published)
}

func TestGRPCDeleteAd(t *testing.T) {
	ctx, client := getTestGRCPClient(t)

	_, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Nickname: "Jenny", Email: "jenny@gmail.com"})
	assert.NoError(t, err, "client.CreateUser")

	ad, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{Title: "Title", Text: "Text", UserId: 0})
	assert.NoError(t, err, "client.CreateAd")

	res, err := client.DeleteAd(ctx, &grpcPort.DeleteAdRequest{AdId: ad.Id, UserId: ad.UserId})
	assert.NoError(t, err, "client.DeleteAd")

	assert.Equal(t, ad.Id, res.Id)
	assert.Equal(t, "Title", res.Title)
	assert.Equal(t, "Text", res.Text)
	assert.False(t, res.Published)
	assert.Zero(t, res.UserId)
}

func TestGRPCGetUser(t *testing.T) {
	ctx, client := getTestGRCPClient(t)

	_, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Nickname: "Jenny", Email: "jenny@gmail.com"})
	assert.NoError(t, err, "client.CreateUser")

	_, err = client.CreateUser(ctx, &grpcPort.CreateUserRequest{Nickname: "Polly", Email: "polly@gmail.com"})
	assert.NoError(t, err, "client.CreateUser")

	res, err := client.GetUser(ctx, &grpcPort.GetUserRequest{Id: 1})
	assert.NoError(t, err, "client.GetUser")

	assert.Equal(t, int64(1), res.Id)
	assert.Equal(t, "Polly", res.Nickname)
	assert.Equal(t, "polly@gmail.com", res.Email)
}

func TestGRPCUpdateUser(t *testing.T) {
	ctx, client := getTestGRCPClient(t)

	u, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Nickname: "Jenny", Email: "jenny@gmail.com"})
	assert.NoError(t, err, "client.CreateUser")

	res, err := client.UpdateUser(ctx, &grpcPort.UpdateUserRequest{Id: u.Id, Nickname: "Polly", Email: "polly@gmail.com"})
	assert.NoError(t, err, "client.DeleteUser")

	assert.Equal(t, u.Id, res.Id)
	assert.Equal(t, "Polly", res.Nickname)
	assert.Equal(t, "polly@gmail.com", res.Email)
}

func TestGRPCDeleteUser(t *testing.T) {
	ctx, client := getTestGRCPClient(t)

	u, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Nickname: "Jenny", Email: "jenny@gmail.com"})
	assert.NoError(t, err, "client.CreateUser")

	res, err := client.DeleteUser(ctx, &grpcPort.DeleteUserRequest{Id: u.Id})
	assert.NoError(t, err, "client.DeleteUser")

	assert.Equal(t, u.Id, res.Id)
	assert.Equal(t, "Jenny", res.Nickname)
	assert.Equal(t, "jenny@gmail.com", res.Email)
}
