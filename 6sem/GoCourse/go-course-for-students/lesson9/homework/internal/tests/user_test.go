package tests

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateUser(t *testing.T) {
	client := getTestClient()

	response1, err := client.createUser("hello", "king")
	assert.NoError(t, err)

	response2, err := client.createUser("like", "sub")
	assert.NoError(t, err)

	response3, err := client.createUser("apple", "tree")
	assert.NoError(t, err)

	assert.Equal(t, response1.Data.ID, int64(0))
	assert.Equal(t, response2.Data.ID, int64(1))
	assert.Equal(t, response3.Data.ID, int64(2))
}

func TestUpdate(t *testing.T) {
	client := getTestClient()

	response1, err := client.createUser("like", "sub")
	assert.NoError(t, err)

	_, err = client.createUser("aaa", "hahaha")
	assert.NoError(t, err)

	response3, err := client.updateUser(response1.Data.ID, "hey", "king")
	assert.NoError(t, err)

	assert.Equal(t, response3.Data.ID, response1.Data.ID)
	assert.NoError(t, err)

}

func TestUpdateNoSuchUser(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser("like", "sub")
	assert.NoError(t, err)

	_, err = client.createUser("aaa", "hahaha")
	assert.NoError(t, err)

	_, err = client.updateUser(56, "hey", "king")
	assert.ErrorIs(t, err, ErrForbidden)

}

func TestCreateAdNoSuchUser(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser("like", "sub")
	assert.NoError(t, err)

	_, err = client.createAd(299, "hello", "world")
	assert.ErrorIs(t, err, ErrForbidden)
}

func TestGetUser(t *testing.T) {
	client := getTestClient()

	response, err := client.createUser("like", "sub")
	assert.NoError(t, err)

	answer, err := client.getUser(response.Data.ID)
	assert.NoError(t, err)
	assert.Equal(t, response.Data.ID, answer.Data.ID)
	assert.Equal(t, response.Data.Email, answer.Data.Email)
	assert.Equal(t, response.Data.Nickname, answer.Data.Nickname)
}

func TestGetNonExistentUser(t *testing.T) {
	client := getTestClient()

	_, err := client.createUser("like", "sub")
	assert.NoError(t, err)

	_, err = client.createUser("like", "sub")
	assert.NoError(t, err)

	_, err = client.createUser("like", "sub")
	assert.NoError(t, err)

	_, err = client.getUser(8)
	assert.ErrorIs(t, err, ErrForbidden)
}
