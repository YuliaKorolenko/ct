package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChangeStatusAdOfAnotherUser(t *testing.T) {
	client := getTestClient()

	user1, err := client.createUser("life is", "good")
	assert.NoError(t, err)
	user2, err := client.createUser("sleep", "kill")
	assert.NoError(t, err)

	resp, err := client.createAd(user1.Data.ID, "hello", "world")
	assert.NoError(t, err)

	_, err = client.changeAdStatus(user2.Data.ID, resp.Data.ID, true)
	assert.ErrorIs(t, err, ErrForbidden)
}

func TestUpdateAdOfAnotherUser(t *testing.T) {
	client := getTestClient()

	user1, err := client.createUser("life is", "good")
	assert.NoError(t, err)
	user2, err := client.createUser("sleep", "kill")
	assert.NoError(t, err)

	resp, err := client.createAd(user1.Data.ID, "hello", "world")
	assert.NoError(t, err)

	_, err = client.updateAd(user2.Data.ID, resp.Data.ID, "title", "text")
	assert.ErrorIs(t, err, ErrForbidden)
}

func TestCreateAd_ID(t *testing.T) {
	client := getTestClient()

	user1, err := client.createUser("life is", "good")
	assert.NoError(t, err)

	resp, err := client.createAd(user1.Data.ID, "hello", "world")
	assert.NoError(t, err)
	assert.Equal(t, resp.Data.ID, int64(0))

	resp, err = client.createAd(user1.Data.ID, "hello", "world")
	assert.NoError(t, err)
	assert.Equal(t, resp.Data.ID, int64(1))

	resp, err = client.createAd(user1.Data.ID, "hello", "world")
	assert.NoError(t, err)
	assert.Equal(t, resp.Data.ID, int64(2))
}
