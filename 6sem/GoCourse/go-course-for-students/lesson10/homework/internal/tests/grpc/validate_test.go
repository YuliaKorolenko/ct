package grpc

import (
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"homework10/internal/ports/grpc/proto"
	"strings"
	"testing"
)

func TestValidateCreateAd(t *testing.T) {
	client, ctx := createGRCPClient(t)

	u, err := client.CreateUser(ctx, &proto.CreateUserRequest{Name: "Oleg", Email: "OlegEmail@"})
	assert.NoError(t, err, "client.CreateUser")

	text := strings.Repeat("hi", 200)

	_, err = client.CreateAd(ctx, &proto.CreateAdRequest{UserId: u.Id, Title: "hi", Text: text})
	assert.ErrorIs(t, err, status.Error(codes.InvalidArgument, "Invalid Argument"))
}

func TestValidateUpdateAd(t *testing.T) {
	client, ctx := createGRCPClient(t)

	u, err := client.CreateUser(ctx, &proto.CreateUserRequest{Name: "Oleg", Email: "OlegEmail@"})
	assert.NoError(t, err, "client.CreateUser")

	resp, err := client.CreateAd(ctx, &proto.CreateAdRequest{UserId: u.Id, Title: "hi", Text: "hello"})
	assert.NoError(t, err)

	_, err = client.UpdateAd(ctx, &proto.UpdateAdRequest{AdId: resp.Id, UserId: resp.AuthorId, Title: "cool", Text: ""})
	assert.ErrorIs(t, err, status.Error(codes.InvalidArgument, "Invalid Argument"))
}
