package grpc

import (
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"homework10/internal/ports/grpc/proto"
	"strings"
	"testing"
)

type Test struct {
	Title string
	Text  string
}

func TestValidate(t *testing.T) {
	tests := [...]Test{
		{"", "cool text"},
		{"cool title", ""},
		{"", ""},
		{strings.Repeat("cool", 500), ""},
		{"", strings.Repeat("hi", 100)},
		{strings.Repeat("o", 500), "normal text"},
		{"normal title", strings.Repeat("table text", 500)},
		{strings.Repeat("title", 500), strings.Repeat("text", 500)},
	}

	client, ctx := createGRCPClient(t)

	u, err := client.CreateUser(ctx, &proto.CreateUserRequest{Name: "Oleg", Email: "OlegEmail@"})
	assert.NoError(t, err, "client.CreateUser")

	for _, test := range tests {
		test := test
		t.Run(test.Title, func(t *testing.T) {
			_, err = client.CreateAd(ctx, &proto.CreateAdRequest{UserId: u.Id, Title: test.Title, Text: test.Text})
			assert.ErrorIs(t, err, status.Error(codes.InvalidArgument, "Invalid Argument"))
		})
	}
}
