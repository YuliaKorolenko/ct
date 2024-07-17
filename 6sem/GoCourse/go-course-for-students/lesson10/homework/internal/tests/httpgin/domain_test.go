package httpgin

import (
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"testing"

	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/stretchr/testify/assert"
)

type SuiteAdTests struct {
	suite.Suite
	client *testClient
	user   userResponse
}

func (s *SuiteAdTests) BeforeEach(t provider.T) {
	t.Log("Start test parsing errors")
	s.client = getTestClient()
	s.user = createAdUser(t, s.client)
}

func createAdUser(t provider.T, tc *testClient) userResponse {
	user, err := tc.createUser("life is", "good")
	assert.NoError(t, err)

	return user
}

func TestSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(SuiteAdTests))
}

func (s *SuiteAdTests) TestChangeStatusAdOfAnotherUser(t provider.T) {

	user2, err := s.client.createUser("sleep", "kill")
	assert.NoError(t, err)

	resp, err := s.client.createAd(s.user.Data.ID, "hello", "world")
	assert.NoError(t, err)

	_, err = s.client.changeAdStatus(user2.Data.ID, resp.Data.ID, true)
	assert.ErrorIs(t, err, ErrForbidden)
}

func (s *SuiteAdTests) TestUpdateAdOfAnotherUser(t provider.T) {

	user2, err := s.client.createUser("sleep", "kill")
	assert.NoError(t, err)

	resp, err := s.client.createAd(s.user.Data.ID, "hello", "world")
	assert.NoError(t, err)

	_, err = s.client.updateAd(user2.Data.ID, resp.Data.ID, "title", "text")
	assert.ErrorIs(t, err, ErrForbidden)
}

func (s *SuiteAdTests) TestCreateAd_ID(t provider.T) {

	resp, err := s.client.createAd(s.user.Data.ID, "hello", "world")
	assert.NoError(t, err)
	assert.Equal(t, resp.Data.ID, int64(0))

	resp, err = s.client.createAd(s.user.Data.ID, "hello", "world")
	assert.NoError(t, err)
	assert.Equal(t, resp.Data.ID, int64(1))

	resp, err = s.client.createAd(s.user.Data.ID, "hello", "world")
	assert.NoError(t, err)
	assert.Equal(t, resp.Data.ID, int64(2))
}
