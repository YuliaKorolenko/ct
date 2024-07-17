package tests

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"homework9/internal/adapters/usrepo"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"homework9/internal/adapters/adrepo"
	"homework9/internal/app"
	grpcPort "homework9/internal/ports/grpc"
)

func createGRCPClient(t *testing.T) (grpcPort.AdServiceClient, context.Context) {
	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() {
		lis.Close()
	})
	srv := grpcPort.NewServerWithLogger()

	t.Cleanup(func() {
		srv.Stop()
	})

	svc := grpcPort.NewAdService(app.NewApp(adrepo.New(), usrepo.NewUserRepo()))
	grpcPort.RegisterAdServiceServer(srv, svc)

	go func() {
		assert.NoError(t, srv.Serve(lis), "srv.Serve")
	}()

	dialer := func(context.Context, string) (net.Conn, error) {
		return lis.Dial()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	t.Cleanup(func() {
		cancel()
	})

	conn, err := grpc.DialContext(ctx, "", grpc.WithContextDialer(dialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err, "grpc.DialContext")

	t.Cleanup(func() {
		conn.Close()
	})

	return grpcPort.NewAdServiceClient(conn), ctx
}

func TestGRPCCreateUser(t *testing.T) {
	client, ctx := createGRCPClient(t)

	res, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "Oleg", Email: "OlegEmail@"})
	assert.NoError(t, err, "client.GetUser")

	assert.Equal(t, "Oleg", res.Name)
	assert.Equal(t, "OlegEmail@", res.Email)
}

func TestGRPCCreateAd(t *testing.T) {
	client, ctx := createGRCPClient(t)

	u, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "Oleg", Email: "OlegEmail@"})
	assert.NoError(t, err, "client.GetUser")

	assert.Equal(t, "Oleg", u.Name)
	assert.Equal(t, "OlegEmail@", u.Email)

	res, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{UserId: u.Id, Title: "fine", Text: "not fine"})
	assert.NoError(t, err, "client.GetAd")

	assert.Equal(t, "fine", res.Title)
	assert.Equal(t, "not fine", res.Text)
	assert.Equal(t, u.Id, res.AuthorId)
}

func TestGRPCChangeAdStatus(t *testing.T) {
	client, ctx := createGRCPClient(t)

	u, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "Oleg", Email: "OlegEmail@"})
	assert.NoError(t, err, "client.GetUser")

	post, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{UserId: u.Id, Title: "fine", Text: "not fine"})
	assert.NoError(t, err, "client.GetAd")

	response, err := client.ChangeAdStatus(ctx, &grpcPort.ChangeAdStatusRequest{AdId: post.Id, UserId: post.AuthorId, Published: true})
	assert.NoError(t, err, "client.ChangeAdStatus")
	assert.True(t, response.Published)
}

func TestGRPСUpdateAd(t *testing.T) {
	client, ctx := createGRCPClient(t)

	u, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "Oleg", Email: "OlegEmail@"})
	assert.NoError(t, err, "client.GetUser")

	post, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{UserId: u.Id, Title: "fine", Text: "not fine"})
	assert.NoError(t, err, "client.GetAd")

	assert.Equal(t, "fine", post.Title)
	assert.Equal(t, "not fine", post.Text)
	assert.Equal(t, u.Id, post.AuthorId)

	response, err := client.UpdateAd(ctx, &grpcPort.UpdateAdRequest{AdId: post.Id, UserId: post.AuthorId, Title: "cool", Text: "tired"})
	assert.NoError(t, err, "client.UpdateAdRequest")

	assert.Equal(t, "tired", response.Text)
	assert.Equal(t, "cool", response.Title)
	assert.Equal(t, u.Id, response.AuthorId)
}

func TestGRPСListAds(t *testing.T) {
	client, ctx := createGRCPClient(t)

	u, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "Oleg", Email: "OlegEmail@"})
	assert.NoError(t, err, "client.GetUser")

	response1, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{UserId: u.Id, Title: "fine", Text: "not fine"})
	assert.NoError(t, err, "client.GetAd")

	_, err = client.CreateAd(ctx, &grpcPort.CreateAdRequest{UserId: u.Id, Title: "finefine", Text: "not finefine"})
	assert.NoError(t, err, "client.GetAd")

	_, err = client.ChangeAdStatus(ctx, &grpcPort.ChangeAdStatusRequest{UserId: u.Id, AdId: response1.Id, Published: true})
	assert.NoError(t, err)

	answer, err := client.ListAds(ctx, &emptypb.Empty{})
	assert.NoError(t, err)

	assert.Equal(t, answer.List[0].Id, response1.Id)
	assert.Equal(t, answer.List[0].Title, response1.Title)
	assert.Equal(t, answer.List[0].Text, response1.Text)
	assert.Equal(t, answer.List[0].AuthorId, response1.AuthorId)
}

func TestGRPСGetUser(t *testing.T) {
	client, ctx := createGRCPClient(t)

	u, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "Oleg", Email: "OlegEmail@"})
	assert.NoError(t, err, "client.CreateUser")

	res, err := client.GetUser(ctx, &grpcPort.GetUserRequest{Id: u.Id})
	assert.NoError(t, err, "client.GetUser")

	assert.Equal(t, u.Name, res.Name)
	assert.Equal(t, u.Email, res.Email)
}

func TestGRPСDeleteUser(t *testing.T) {
	client, ctx := createGRCPClient(t)

	u, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "Oleg", Email: "OlegEmail@"})
	assert.NoError(t, err, "client.CreateUser")

	_, err = client.GetUser(ctx, &grpcPort.GetUserRequest{Id: u.Id})
	assert.NoError(t, err, "client.GetUser")

	_, err = client.DeleteUser(ctx, &grpcPort.DeleteUserRequest{Id: u.Id})
	assert.NoError(t, err, "client.DeleteUser")

	_, err = client.GetUser(ctx, &grpcPort.GetUserRequest{Id: u.Id})
	assert.Equal(t, err, status.Error(codes.InvalidArgument, "No such user"))

}

func TestGRPСGetAd(t *testing.T) {
	client, ctx := createGRCPClient(t)

	u, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "Oleg", Email: "OlegEmail@"})
	assert.NoError(t, err, "client.CreateUser")

	post, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{UserId: u.Id, Title: "finefine", Text: "not finefine"})
	assert.NoError(t, err, "client.CreateAd")

	response, err := client.GetAd(ctx, &grpcPort.GetAdRequest{AdId: post.Id})
	assert.NoError(t, err, "client.GetAd")

	assert.Equal(t, post.Id, response.Id)
	assert.Equal(t, post.Title, response.Title)
	assert.Equal(t, post.Text, response.Text)
	assert.Equal(t, post.AuthorId, response.AuthorId)
}

func TestGRPСDeleteAd(t *testing.T) {
	client, ctx := createGRCPClient(t)

	u, err := client.CreateUser(ctx, &grpcPort.CreateUserRequest{Name: "Oleg", Email: "OlegEmail@"})
	assert.NoError(t, err, "client.CreateUser")

	post, err := client.CreateAd(ctx, &grpcPort.CreateAdRequest{UserId: u.Id, Title: "finefine", Text: "not finefine"})
	assert.NoError(t, err, "client.GetAd")

	_, err = client.DeleteAd(ctx, &grpcPort.DeleteAdRequest{AdId: post.Id, AuthorId: post.AuthorId})
	assert.NoError(t, err, "client.DeleteUser")

	_, err = client.GetAd(ctx, &grpcPort.GetAdRequest{AdId: post.Id})
	assert.ErrorIs(t, err, status.Error(codes.InvalidArgument, "No such add"))
}
