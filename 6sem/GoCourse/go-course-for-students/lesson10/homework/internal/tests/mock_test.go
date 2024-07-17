package tests

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"homework10/internal/app"
	mocks "homework10/internal/app/mocks"
	"homework10/internal/entities/ads"
	"homework10/internal/entities/user"
	"strings"
	"testing"
	"time"
)

func TestUpdateUser(t *testing.T) {

	nick := strings.Repeat("oleg", 500)

	repo := &mocks.Repository{}

	urepo := &mocks.URepository{}
	urepo.On("CreateUser", mock.Anything).
		Return(user.User{ID: 0, Nickname: "normal nick", Email: "email"}, nil).
		Once()
	urepo.On("GetCurrentUser", mock.Anything).
		Return(user.User{ID: 0, Nickname: "normal nick", Email: "email"}, nil).
		Once()
	urepo.On("UpdateUser", mock.Anything).
		Return(user.User{ID: 0, Nickname: nick, Email: "email"}, nil).
		Once()

	service := app.NewApp(repo, urepo)

	user, err := service.CreateUser("normal nick", "email")
	assert.NoError(t, err)
	assert.Equal(t, "normal nick", user.Nickname)
	assert.Equal(t, "email", user.Email)

	user2, err := service.UpdateUser(0, nick, "email")
	assert.Equal(t, nick, user2.Nickname)
	assert.Equal(t, "email", user2.Email)
	assert.Equal(t, app.ValidationApError{Err: app.ErrorValidation}, err)
}

func TestCreateUser(t *testing.T) {

	nick := strings.Repeat("oleg", 500)

	repo := &mocks.Repository{}

	urepo := &mocks.URepository{}
	urepo.On("CreateUser", mock.Anything).
		Return(user.User{ID: 0, Nickname: nick, Email: "email"}, nil).
		Once()

	service := app.NewApp(repo, urepo)

	user, err := service.CreateUser(nick, "email")
	assert.Equal(t, app.ValidationApError{Err: app.ErrorValidation}, err)
	assert.Equal(t, nick, user.Nickname)
	assert.Equal(t, "email", user.Email)
}

func TestChangeAdStatus(t *testing.T) {
	repo := &mocks.Repository{}
	repo.On("GetCurrentAd", mock.Anything).
		Return(ads.Ad{ID: 0, Title: "", Text: "text", AuthorID: 0, Published: true, Time: time.Now().UTC()}, nil).
		Once()
	repo.On("ChangeAdStatus", mock.Anything).
		Return(ads.Ad{ID: 0, Title: "", Text: "text", AuthorID: 0, Published: true, Time: time.Now().UTC()}, nil).
		Once()

	urepo := &mocks.URepository{}
	urepo.On("GetCurrentUser", mock.Anything).
		Return(user.User{ID: 0, Nickname: "normal nick", Email: "email"}, nil).
		Once()

	service := app.NewApp(repo, urepo)

	post, err := service.ChangeAdStatus(0, 0, false)
	assert.EqualValues(t, 0, post.ID)
	assert.EqualValues(t, 0, post.AuthorID)
	assert.Equal(t, app.ValidationApError{Err: app.ErrorValidation}, err)
}
