package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/ozontech/allure-go/pkg/framework/provider"
	"github.com/ozontech/allure-go/pkg/framework/suite"
	"github.com/stretchr/testify/assert"
	"homework10/internal/adapters/adrepo"
	"homework10/internal/adapters/usrepo"
	"homework10/internal/app"
	"homework10/internal/ports/httpgin"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type SuiteTests struct {
	suite.Suite
	client *testSuitClient
}

var (
	ErrBadRequest = fmt.Errorf("bad request")
)

func TestSuiteRunner(t *testing.T) {
	suite.RunSuite(t, new(SuiteTests))
}

type testSuitClient struct {
	client  *http.Client
	baseURL string
}

type userResponse struct {
	Data userData `json:"data"`
}

type adData struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Text      string    `json:"text"`
	AuthorID  int64     `json:"author_id"`
	Published bool      `json:"published"`
	Data      time.Time `json:"data"`
}

type adResponse struct {
	Data adData `json:"data"`
}

type userData struct {
	ID       int64  `json:"id"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
}

func (s *SuiteTests) BeforeAll(t provider.T) {
	t.Log("Start test parsing errors")
	s.client = s.Client()
}

func (s *SuiteTests) AfterAll(t provider.T) {
	t.Log("End test parsing errors")
}

func (s *SuiteTests) Client() *testSuitClient {
	server := httpgin.NewHTTPServer(":18080", app.NewApp(adrepo.New(), usrepo.NewUserRepo()))
	testServer := httptest.NewServer(server.Handler)

	return &testSuitClient{
		client:  testServer.Client(),
		baseURL: testServer.URL,
	}
}

func (tc *testSuitClient) getErrorResponse(req *http.Request) error {
	resp, err := tc.client.Do(req)
	if err != nil {
		return fmt.Errorf("unexpected error: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusBadRequest {
			return ErrBadRequest
		}
		return fmt.Errorf("unexpected status code: %s", resp.Status)
	}

	return nil
}

func (tc *testSuitClient) getResponse(req *http.Request, out any) error {
	resp, err := tc.client.Do(req)
	if err != nil {
		return fmt.Errorf("unexpected error: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusBadRequest {
			return ErrBadRequest
		}
		return fmt.Errorf("unexpected status code: %s", resp.Status)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("unable to read response: %w", err)
	}

	err = json.Unmarshal(respBody, out)
	if err != nil {
		return fmt.Errorf("unable to unmarshal: %w", err)
	}

	return nil
}

func (tc *testSuitClient) deleteUserCheckParse(userId string) error {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf(tc.baseURL+"/api/v1/users/%s", userId), nil)
	if err != nil {
		return fmt.Errorf("unable to create request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")

	err = tc.getErrorResponse(req)
	if err != nil {
		return err
	}
	return nil
}

func (tc *testSuitClient) getUser(userId string) (userResponse, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(tc.baseURL+"/api/v1/users/%s", userId), nil)
	if err != nil {
		return userResponse{}, fmt.Errorf("unable to create request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")

	var response userResponse
	err = tc.getResponse(req, &response)
	if err != nil {
		return userResponse{}, err
	}

	return response, nil
}

func (tc *testSuitClient) deleteAd(adId string, userID int64) error {
	body := map[string]any{
		"user_id": userID,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("unable to marshal: %w", err)
	}

	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf(tc.baseURL+"/api/v1/ads/%s", adId), bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("unable to create request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")

	err = tc.getErrorResponse(req)
	if err != nil {
		return err
	}
	return nil
}

func (tc *testSuitClient) updateAd(userID int64, adID string, title string, text string) (adResponse, error) {
	body := map[string]any{
		"user_id": userID,
		"title":   title,
		"text":    text,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return adResponse{}, fmt.Errorf("unable to marshal: %w", err)
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf(tc.baseURL+"/api/v1/ads/%s", adID), bytes.NewReader(data))
	if err != nil {
		return adResponse{}, fmt.Errorf("unable to create request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")

	var response adResponse
	err = tc.getResponse(req, &response)
	if err != nil {
		return adResponse{}, err
	}

	return response, nil
}

func (tc *testSuitClient) updateUser(userId string, nickname string, email string) (userResponse, error) {
	body := map[string]any{
		"nickname": nickname,
		"email":    email,
	}

	data, err := json.Marshal(body)
	if err != nil {
		return userResponse{}, fmt.Errorf("unable to marshal: %w", err)
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf(tc.baseURL+"/api/v1/users/%s", userId), bytes.NewReader(data))
	if err != nil {
		return userResponse{}, fmt.Errorf("unable to create request: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")

	var response userResponse
	err = tc.getResponse(req, &response)
	if err != nil {
		return userResponse{}, err
	}

	return response, nil
}

func (s *SuiteTests) TestDeleteUserCheckParseError(t provider.T) {
	err := s.client.deleteUserCheckParse("life is")
	assert.ErrorIs(t, err, ErrBadRequest)
}

func (s *SuiteTests) TestGetUserCheckParseError(t provider.T) {
	_, err := s.client.getUser("life is")
	assert.ErrorIs(t, err, ErrBadRequest)
}

func (s *SuiteTests) TestDeleteAdCheckParseError(t provider.T) {
	err := s.client.deleteAd("life is", 0)
	assert.ErrorIs(t, err, ErrBadRequest)
}

func (s *SuiteTests) TestUpdateAdCheckParseError(t provider.T) {
	_, err := s.client.updateAd(int64(0), "1-8", "title", "text")
	assert.ErrorIs(t, err, ErrBadRequest)
}

func (s *SuiteTests) TestUpdateUserCheckParseError(t provider.T) {
	_, err := s.client.updateUser("1-8", "name", "email")
	assert.ErrorIs(t, err, ErrBadRequest)
}
