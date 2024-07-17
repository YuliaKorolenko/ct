package tests

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDeleteAd(t *testing.T) {
	client := getTestClient()

	user1, err := client.createUser("life is", "good")
	assert.NoError(t, err)

	r1, err := client.createAd(user1.Data.ID, "hello", "world")
	assert.NoError(t, err)

	response, err := client.createAd(user1.Data.ID, "i like fruits", "Pepper is a very spicy.")
	assert.NoError(t, err)

	r3, err := client.createAd(user1.Data.ID, "third ad", "Ho-ho-ho. I'm santa")
	assert.NoError(t, err)

	err = client.deleteAd(response.Data.ID, user1.Data.ID)
	assert.NoError(t, err)

	_, err = client.getAd(r1.Data.ID)
	assert.NoError(t, err)

	_, err = client.getAd(r3.Data.ID)
	assert.NoError(t, err)

	_, err = client.getAd(response.Data.ID)
	assert.ErrorIs(t, err, ErrForbidden)
}

func TestDoubleDelete(t *testing.T) {
	client := getTestClient()

	user1, err := client.createUser("life is", "good")
	assert.NoError(t, err)

	r1, err := client.createAd(user1.Data.ID, "hello", "world")
	assert.NoError(t, err)

	_, err = client.getAd(r1.Data.ID)
	assert.NoError(t, err)

	err = client.deleteAd(r1.Data.ID, user1.Data.ID)
	assert.NoError(t, err)

	_, err = client.getAd(r1.Data.ID)
	assert.ErrorIs(t, err, ErrForbidden)

	err = client.deleteAd(r1.Data.ID, user1.Data.ID)
	assert.ErrorIs(t, err, ErrForbidden)
}

func TestUserDelete(t *testing.T) {
	client := getTestClient()

	user0, err := client.createUser("life is", "good")
	assert.NoError(t, err)

	user1, err := client.createUser("life is", "good")
	assert.NoError(t, err)

	user2, err := client.createUser("life is", "good")
	assert.NoError(t, err)

	err = client.deleteUser(user1.Data.ID)
	assert.NoError(t, err)

	answer0, err := client.getUser(user0.Data.ID)
	assert.NoError(t, err)
	assert.Equal(t, user0.Data.ID, answer0.Data.ID)
	assert.Equal(t, user0.Data.Email, answer0.Data.Email)
	assert.Equal(t, user0.Data.Nickname, answer0.Data.Nickname)

	_, err = client.getUser(user1.Data.ID)
	assert.ErrorIs(t, err, ErrForbidden)

	answer2, err := client.getUser(user2.Data.ID)
	assert.NoError(t, err)
	assert.Equal(t, user2.Data.ID, answer2.Data.ID)
	assert.Equal(t, user2.Data.Email, answer2.Data.Email)
	assert.Equal(t, user2.Data.Nickname, answer2.Data.Nickname)
}

func TestGetAdByDeleteUser(t *testing.T) {
	client := getTestClient()

	user1, err := client.createUser("life is", "good")
	assert.NoError(t, err)

	r1, err := client.createAd(user1.Data.ID, "hello", "world")
	assert.NoError(t, err)

	err = client.deleteUser(user1.Data.ID)
	assert.NoError(t, err)

	_, err = client.getAd(r1.Data.ID)
	assert.ErrorIs(t, err, ErrForbidden)

}
