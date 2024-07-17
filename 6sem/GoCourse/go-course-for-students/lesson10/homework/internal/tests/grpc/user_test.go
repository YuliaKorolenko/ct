package grpc

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"homework10/internal/adapters/usrepo"
	"homework10/internal/ports/grpc/proto"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"homework10/internal/adapters/adrepo"
	"homework10/internal/app"
	grpcPort "homework10/internal/ports/grpc"
)

func createGRCPClient(t *testing.T) (proto.AdServiceClient, context.Context) {
	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() {
		lis.Close()
	})
	srv := grpcPort.NewServerWithLogger()

	t.Cleanup(func() {
		srv.Stop()
	})

	svc := grpcPort.NewAdService(app.NewApp(adrepo.New(), usrepo.NewUserRepo()))
	proto.RegisterAdServiceServer(srv, svc)

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

	return proto.NewAdServiceClient(conn), ctx
}

func TestGRPCCreateUser(t *testing.T) {
	client, ctx := createGRCPClient(t)

	res, err := client.CreateUser(ctx, &proto.CreateUserRequest{Name: "Oleg", Email: "OlegEmail@"})
	assert.NoError(t, err, "client.GetUser")

	assert.Equal(t, "Oleg", res.Name)
	assert.Equal(t, "OlegEmail@", res.Email)
}

func TestGRPСGetUser(t *testing.T) {
	client, ctx := createGRCPClient(t)

	u, err := client.CreateUser(ctx, &proto.CreateUserRequest{Name: "Oleg", Email: "OlegEmail@"})
	assert.NoError(t, err, "client.CreateUser")

	res, err := client.GetUser(ctx, &proto.GetUserRequest{Id: u.Id})
	assert.NoError(t, err, "client.GetUser")

	assert.Equal(t, u.Name, res.Name)
	assert.Equal(t, u.Email, res.Email)
}

func TestGRPСDeleteUser(t *testing.T) {
	client, ctx := createGRCPClient(t)

	u, err := client.CreateUser(ctx, &proto.CreateUserRequest{Name: "Oleg", Email: "OlegEmail@"})
	assert.NoError(t, err, "client.CreateUser")

	_, err = client.GetUser(ctx, &proto.GetUserRequest{Id: u.Id})
	assert.NoError(t, err, "client.GetUser")

	_, err = client.DeleteUser(ctx, &proto.DeleteUserRequest{Id: u.Id})
	assert.NoError(t, err, "client.DeleteUser")

	_, err = client.GetUser(ctx, &proto.GetUserRequest{Id: u.Id})
	assert.Equal(t, err, status.Error(codes.InvalidArgument, "No such user"))

}

func TestGRPСGetAd(t *testing.T) {
	client, ctx := createGRCPClient(t)

	u, err := client.CreateUser(ctx, &proto.CreateUserRequest{Name: "Oleg", Email: "OlegEmail@"})
	assert.NoError(t, err, "client.CreateUser")

	post, err := client.CreateAd(ctx, &proto.CreateAdRequest{UserId: u.Id, Title: "finefine", Text: "not finefine"})
	assert.NoError(t, err, "client.CreateAd")

	response, err := client.GetAd(ctx, &proto.GetAdRequest{AdId: post.Id})
	assert.NoError(t, err, "client.GetAd")

	assert.Equal(t, post.Id, response.Id)
	assert.Equal(t, post.Title, response.Title)
	assert.Equal(t, post.Text, response.Text)
	assert.Equal(t, post.AuthorId, response.AuthorId)
}

func TestDeleteNonExistUser(t *testing.T) {
	client, ctx := createGRCPClient(t)

	_, err := client.DeleteUser(ctx, &proto.DeleteUserRequest{Id: 189})
	assert.ErrorIs(t, err, status.Error(codes.InvalidArgument, "No such user"))

}
