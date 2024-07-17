package tests

import (
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
	"time"
)

func TestFilterPublished(t *testing.T) {
	client := getTestClient()

	user1, err := client.createUser("life is", "good")
	assert.NoError(t, err)

	response, err := client.createAd(user1.Data.ID, "hello", "world")
	assert.NoError(t, err)

	_, err = client.changeAdStatus(user1.Data.ID, response.Data.ID, true)
	assert.NoError(t, err)

	notPublished, err := client.createAd(user1.Data.ID, "best cat", "not for sale")
	assert.NoError(t, err)

	params := make(map[string]string, 0)

	params["published"] = "false"

	ads, err := client.listAds(params)
	assert.NoError(t, err)
	assert.Len(t, ads.Data, 1)
	assert.Equal(t, ads.Data[0].ID, notPublished.Data.ID)
	assert.Equal(t, ads.Data[0].Title, notPublished.Data.Title)
	assert.Equal(t, ads.Data[0].Text, notPublished.Data.Text)
	assert.Equal(t, ads.Data[0].AuthorID, notPublished.Data.AuthorID)
	assert.False(t, ads.Data[0].Published)
}

func TestFilterAuthor(t *testing.T) {
	client := getTestClient()

	user1, err := client.createUser("life is", "good")
	assert.NoError(t, err)

	user2, err := client.createUser("sleep", "kill")
	assert.NoError(t, err)

	author1, err := client.createAd(user1.Data.ID, "hello", "world")
	assert.NoError(t, err)

	_, err = client.changeAdStatus(author1.Data.AuthorID, author1.Data.ID, true)
	assert.NoError(t, err)

	author3, err := client.createAd(user2.Data.ID, "hello", "world")
	assert.NoError(t, err)

	_, err = client.changeAdStatus(author3.Data.AuthorID, author3.Data.ID, true)
	assert.NoError(t, err)

	author2, err := client.createAd(user1.Data.ID, "best cat", "not for sale")
	assert.NoError(t, err)

	_, err = client.changeAdStatus(author2.Data.AuthorID, author2.Data.ID, true)
	assert.NoError(t, err)

	params := make(map[string]string, 0)

	params["author"] = "0"

	ads, err := client.listAds(params)
	sort.Slice(
		ads.Data,
		func(i, j int) bool {
			return ads.Data[i].ID < ads.Data[j].ID
		},
	)
	assert.NoError(t, err)
	assert.Len(t, ads.Data, 2)
	assert.Equal(t, ads.Data[0].ID, author1.Data.ID)
	assert.Equal(t, ads.Data[0].AuthorID, int64(0))
	assert.Equal(t, ads.Data[1].ID, author2.Data.ID)
	assert.Equal(t, ads.Data[1].AuthorID, int64(0))

}

func TestFilterAuthorAndPublished(t *testing.T) {
	client := getTestClient()

	user1, err := client.createUser("life is", "good")
	assert.NoError(t, err)

	user2, err := client.createUser("sleep", "kill")
	assert.NoError(t, err)

	author1, err := client.createAd(user1.Data.ID, "hello", "world")
	assert.NoError(t, err)

	author3, err := client.createAd(user2.Data.ID, "fine", "I like sleep and eat")
	assert.NoError(t, err)

	_, err = client.changeAdStatus(author3.Data.AuthorID, author3.Data.ID, true)
	assert.NoError(t, err)

	author2, err := client.createAd(user1.Data.ID, "best cat", "not for sale")
	assert.NoError(t, err)

	_, err = client.changeAdStatus(author2.Data.AuthorID, author2.Data.ID, true)
	assert.NoError(t, err)

	params := make(map[string]string, 0)

	params["author"] = "0"
	params["published"] = "false"

	ads, err := client.listAds(params)
	assert.NoError(t, err)
	assert.Len(t, ads.Data, 1)
	assert.Equal(t, ads.Data[0].ID, author1.Data.ID)
	assert.Equal(t, ads.Data[0].Title, author1.Data.Title)
	assert.Equal(t, ads.Data[0].Text, author1.Data.Text)
	assert.Equal(t, ads.Data[0].AuthorID, author1.Data.AuthorID)
	assert.Equal(t, ads.Data[0].Published, author1.Data.Published)

}

func TestFilterTitle(t *testing.T) {
	client := getTestClient()

	user1, err := client.createUser("life is", "good")
	assert.NoError(t, err)

	user2, err := client.createUser("sleep", "kill")
	assert.NoError(t, err)

	user3, err := client.createUser("pure", "try")
	assert.NoError(t, err)

	_, err = client.createAd(user1.Data.ID, "hello", "world")
	assert.NoError(t, err)

	title2, err := client.createAd(user1.Data.ID, "hello", "not for sale")
	assert.NoError(t, err)

	_, err = client.changeAdStatus(title2.Data.AuthorID, title2.Data.ID, true)
	assert.NoError(t, err)

	title3, err := client.createAd(user2.Data.ID, "hello", "I like sleep and eat")
	assert.NoError(t, err)

	_, err = client.changeAdStatus(title3.Data.AuthorID, title3.Data.ID, true)
	assert.NoError(t, err)

	title4, err := client.createAd(user3.Data.ID, "fine", "I like sleep and eat")
	assert.NoError(t, err)

	_, err = client.changeAdStatus(title4.Data.AuthorID, title4.Data.ID, true)
	assert.NoError(t, err)

	params := make(map[string]string, 0)

	params["title"] = "hello"

	ads, err := client.listAds(params)
	sort.Slice(
		ads.Data,
		func(i, j int) bool {
			return ads.Data[i].ID < ads.Data[j].ID
		},
	)
	assert.NoError(t, err)
	assert.Len(t, ads.Data, 2)
	assert.Equal(t, ads.Data[0].ID, title2.Data.ID)
	assert.Equal(t, ads.Data[0].AuthorID, title2.Data.AuthorID)
	assert.Equal(t, ads.Data[1].ID, title3.Data.ID)
	assert.Equal(t, ads.Data[1].AuthorID, title3.Data.AuthorID)
	assert.Equal(t, ads.Data[0].Title, ads.Data[0].Title)

}

func TestFilterTitleAuthorPublished(t *testing.T) {
	client := getTestClient()

	user1, err := client.createUser("life is", "good")
	assert.NoError(t, err)
	user2, err := client.createUser("sleep", "kill")
	assert.NoError(t, err)
	user3, err := client.createUser("pure", "try")
	assert.NoError(t, err)
	user4, err := client.createUser("ginger", "upload")
	assert.NoError(t, err)

	_, err = client.createAd(user1.Data.ID, "hello", "world")
	assert.NoError(t, err)

	test2, err := client.createAd(user1.Data.ID, "hello", "not for sale")
	assert.NoError(t, err)

	_, err = client.changeAdStatus(test2.Data.AuthorID, test2.Data.ID, true)
	assert.NoError(t, err)

	test3, err := client.createAd(user2.Data.ID, "hello", "I like sleep and eat")
	assert.NoError(t, err)

	_, err = client.changeAdStatus(test3.Data.AuthorID, test3.Data.ID, true)
	assert.NoError(t, err)

	test4, err := client.createAd(user3.Data.ID, "fine", "I like sleep and eat")
	assert.NoError(t, err)

	_, err = client.changeAdStatus(test4.Data.AuthorID, test4.Data.ID, true)
	assert.NoError(t, err)

	test5, err := client.createAd(user4.Data.ID, "Hey", "Cool story")
	assert.NoError(t, err)

	params := make(map[string]string, 0)

	params["title"] = "Hey"
	params["published"] = "false"
	params["author"] = "3"

	ads, err := client.listAds(params)
	assert.NoError(t, err)
	assert.Len(t, ads.Data, 1)
	assert.Equal(t, ads.Data[0].ID, test5.Data.ID)
	assert.Equal(t, ads.Data[0].Title, test5.Data.Title)
	assert.Equal(t, ads.Data[0].Text, test5.Data.Text)
	assert.Equal(t, ads.Data[0].AuthorID, test5.Data.AuthorID)
	assert.False(t, ads.Data[0].Published, test5.Data.Published)
}

func TestFilterTime(t *testing.T) {
	client := getTestClient()

	user1, err := client.createUser("life is", "good")
	assert.NoError(t, err)

	response1, err := client.createAd(user1.Data.ID, "hello", "world")
	assert.NoError(t, err)

	response, err := client.changeAdStatus(user1.Data.ID, response1.Data.ID, true)
	assert.NoError(t, err)

	params := make(map[string]string, 0)

	params["time"] = time.Now().UTC().Format(time.DateOnly)

	ads, err := client.listAds(params)
	assert.NoError(t, err)
	assert.Len(t, ads.Data, 1)
	assert.Equal(t, ads.Data[0].ID, response.Data.ID)
	assert.Equal(t, ads.Data[0].Title, response.Data.Title)
	assert.Equal(t, ads.Data[0].Text, response.Data.Text)
	assert.Equal(t, ads.Data[0].AuthorID, response.Data.AuthorID)
	assert.Equal(t, ads.Data[0].Published, response.Data.Published)
}

func TestGetAdById(t *testing.T) {
	client := getTestClient()

	user1, err := client.createUser("life is", "good")
	assert.NoError(t, err)

	response, err := client.createAd(user1.Data.ID, "hello", "world")
	assert.NoError(t, err)

	_, err = client.createAd(user1.Data.ID, "best cat", "not for sale")
	assert.NoError(t, err)

	ad, err := client.getAd(response.Data.ID)
	assert.NoError(t, err)
	assert.Equal(t, ad.Data.ID, response.Data.ID)
	assert.Equal(t, ad.Data.Title, response.Data.Title)
	assert.Equal(t, ad.Data.Text, response.Data.Text)
	assert.Equal(t, ad.Data.AuthorID, response.Data.AuthorID)
	assert.Equal(t, ad.Data.Published, response.Data.Published)
}
