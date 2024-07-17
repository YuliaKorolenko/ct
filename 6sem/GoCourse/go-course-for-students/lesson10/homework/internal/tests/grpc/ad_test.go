package grpc

import (
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"homework10/internal/ports/grpc/proto"
	"testing"
)

func TestGRPCCreateAd(t *testing.T) {
	client, ctx := createGRCPClient(t)

	u, err := client.CreateUser(ctx, &proto.CreateUserRequest{Name: "Oleg", Email: "OlegEmail@"})
	assert.NoError(t, err, "client.GetUser")

	assert.Equal(t, "Oleg", u.Name)
	assert.Equal(t, "OlegEmail@", u.Email)

	res, err := client.CreateAd(ctx, &proto.CreateAdRequest{UserId: u.Id, Title: "fine", Text: "not fine"})
	assert.NoError(t, err, "client.CreateAd")

	_, err = client.CreateAd(ctx, &proto.CreateAdRequest{UserId: 444, Title: "fine", Text: "not fine"})
	assert.ErrorIs(t, err, status.Error(codes.InvalidArgument, "No such user"))

	assert.Equal(t, "fine", res.Title)
	assert.Equal(t, "not fine", res.Text)
	assert.Equal(t, u.Id, res.AuthorId)
}

func TestGRPCChangeAdStatus(t *testing.T) {
	client, ctx := createGRCPClient(t)

	u, err := client.CreateUser(ctx, &proto.CreateUserRequest{Name: "Oleg", Email: "OlegEmail@"})
	assert.NoError(t, err, "client.GetUser")

	u1, err := client.CreateUser(ctx, &proto.CreateUserRequest{Name: "Yulia", Email: "ojpochemy@"})
	assert.NoError(t, err, "client.GetUser")

	post, err := client.CreateAd(ctx, &proto.CreateAdRequest{UserId: u.Id, Title: "fine", Text: "not fine"})
	assert.NoError(t, err, "client.CreateAd")

	response, err := client.ChangeAdStatus(ctx, &proto.ChangeAdStatusRequest{AdId: post.Id, UserId: post.AuthorId, Published: true})
	assert.NoError(t, err, "client.ChangeAdStatus")
	assert.True(t, response.Published)

	_, err = client.ChangeAdStatus(ctx, &proto.ChangeAdStatusRequest{AdId: post.Id, UserId: 125, Published: true})
	assert.ErrorIs(t, err, status.Error(codes.InvalidArgument, "No such user"))

	_, err = client.ChangeAdStatus(ctx, &proto.ChangeAdStatusRequest{AdId: post.Id, UserId: u1.Id, Published: true})
	assert.ErrorIs(t, err, status.Error(codes.PermissionDenied, "Permission denied. You are not owner"))
}

func TestGRPСUpdateAd(t *testing.T) {
	client, ctx := createGRCPClient(t)

	u, err := client.CreateUser(ctx, &proto.CreateUserRequest{Name: "Oleg", Email: "OlegEmail@"})
	assert.NoError(t, err, "client.GetUser")

	post, err := client.CreateAd(ctx, &proto.CreateAdRequest{UserId: u.Id, Title: "fine", Text: "not fine"})
	assert.NoError(t, err, "client.CreateAd")

	assert.Equal(t, "fine", post.Title)
	assert.Equal(t, "not fine", post.Text)
	assert.Equal(t, u.Id, post.AuthorId)

	response, err := client.UpdateAd(ctx, &proto.UpdateAdRequest{AdId: post.Id, UserId: post.AuthorId, Title: "cool", Text: "tired"})
	assert.NoError(t, err, "client.UpdateAdRequest")

	assert.Equal(t, "tired", response.Text)
	assert.Equal(t, "cool", response.Title)
	assert.Equal(t, u.Id, response.AuthorId)
}

func TestGRPСListAds(t *testing.T) {
	client, ctx := createGRCPClient(t)

	u, err := client.CreateUser(ctx, &proto.CreateUserRequest{Name: "Oleg", Email: "OlegEmail@"})
	assert.NoError(t, err, "client.GetUser")

	response1, err := client.CreateAd(ctx, &proto.CreateAdRequest{UserId: u.Id, Title: "fine", Text: "not fine"})
	assert.NoError(t, err, "client.CreateAd")

	_, err = client.CreateAd(ctx, &proto.CreateAdRequest{UserId: u.Id, Title: "finefine", Text: "not finefine"})
	assert.NoError(t, err, "client.CreateAd")

	_, err = client.ChangeAdStatus(ctx, &proto.ChangeAdStatusRequest{UserId: u.Id, AdId: response1.Id, Published: true})
	assert.NoError(t, err)

	answer, err := client.ListAds(ctx, &emptypb.Empty{})
	assert.NoError(t, err)

	assert.Equal(t, answer.List[0].Id, response1.Id)
	assert.Equal(t, answer.List[0].Title, response1.Title)
	assert.Equal(t, answer.List[0].Text, response1.Text)
	assert.Equal(t, answer.List[0].AuthorId, response1.AuthorId)
}

func TestGRPСDeleteAd(t *testing.T) {
	client, ctx := createGRCPClient(t)

	u, err := client.CreateUser(ctx, &proto.CreateUserRequest{Name: "Oleg", Email: "OlegEmail@"})
	assert.NoError(t, err, "client.CreateUser")

	post, err := client.CreateAd(ctx, &proto.CreateAdRequest{UserId: u.Id, Title: "finefine", Text: "not finefine"})
	assert.NoError(t, err, "client.CreateAd")

	_, err = client.DeleteAd(ctx, &proto.DeleteAdRequest{AdId: post.Id, AuthorId: post.AuthorId})
	assert.NoError(t, err, "client.DeleteAd")

	_, err = client.GetAd(ctx, &proto.GetAdRequest{AdId: post.Id})
	assert.ErrorIs(t, err, status.Error(codes.InvalidArgument, "No such add"))
}

func TestUpdateAdByNotOwner(t *testing.T) {
	client, ctx := createGRCPClient(t)

	u, err := client.CreateUser(ctx, &proto.CreateUserRequest{Name: "Oleg", Email: "OlegEmail@"})
	assert.NoError(t, err, "client.CreateUser")

	user2, err := client.CreateUser(ctx, &proto.CreateUserRequest{Name: "olga", Email: "OlgaEmail@"})
	assert.NoError(t, err, "client.CreateUser")

	resp, err := client.CreateAd(ctx, &proto.CreateAdRequest{UserId: u.Id, Title: "hi", Text: "hello"})
	assert.NoError(t, err)

	_, err = client.UpdateAd(ctx, &proto.UpdateAdRequest{AdId: resp.Id, UserId: 197, Title: "cool", Text: "bye"})
	assert.ErrorIs(t, err, status.Error(codes.InvalidArgument, "No such user"))

	_, err = client.UpdateAd(ctx, &proto.UpdateAdRequest{AdId: resp.Id, UserId: user2.Id, Title: "cool", Text: "bye"})
	assert.ErrorIs(t, err, status.Error(codes.PermissionDenied, "Permission denied. You are not owner"))

}

func TestDeleteNotExistedAd(t *testing.T) {
	client, ctx := createGRCPClient(t)

	u, err := client.CreateUser(ctx, &proto.CreateUserRequest{Name: "Oleg", Email: "OlegEmail@"})
	assert.NoError(t, err, "client.CreateUser")

	post, err := client.CreateAd(ctx, &proto.CreateAdRequest{UserId: u.Id, Title: "finefine", Text: "not finefine"})
	assert.NoError(t, err, "client.CreateAd")

	user2, err := client.CreateUser(ctx, &proto.CreateUserRequest{Name: "olga", Email: "OlgaEmail@"})
	assert.NoError(t, err, "client.CreateUser")

	_, err = client.DeleteAd(ctx, &proto.DeleteAdRequest{AdId: 12, AuthorId: post.AuthorId})
	assert.ErrorIs(t, err, status.Error(codes.InvalidArgument, "No such add"))

	_, err = client.DeleteAd(ctx, &proto.DeleteAdRequest{AdId: post.Id, AuthorId: 124})
	assert.ErrorIs(t, err, status.Error(codes.InvalidArgument, "No such user"))

	_, err = client.DeleteAd(ctx, &proto.DeleteAdRequest{AdId: post.Id, AuthorId: user2.Id})
	assert.ErrorIs(t, err, status.Error(codes.PermissionDenied, "Permission denied. You are not owner"))
}

func TestGetAd(t *testing.T) {
	client, ctx := createGRCPClient(t)

	u, err := client.CreateUser(ctx, &proto.CreateUserRequest{Name: "Oleg", Email: "OlegEmail@"})
	assert.NoError(t, err, "client.CreateUser")

	_, err = client.CreateAd(ctx, &proto.CreateAdRequest{UserId: u.Id, Title: "finefine", Text: "not finefine"})
	assert.NoError(t, err, "client.CreateAd")

	_, err = client.GetAd(ctx, &proto.GetAdRequest{AdId: 192734})
	assert.ErrorIs(t, err, status.Error(codes.InvalidArgument, "No such add"))
}
