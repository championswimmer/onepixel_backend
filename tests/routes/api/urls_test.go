package api

import (
	"bytes"
	"encoding/json"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http/httptest"
	"onepixel_backend/src/dtos"
	"onepixel_backend/src/utils/applogger"
	"onepixel_backend/tests"
	"testing"
)

func TestUrlsRoute_CreateRandomUrl(t *testing.T) {

	// ------ REGISTER USER ------

	responseBody := tests.TestUtil_CreateUser(t, "user3689@test.com", "123456")

	// ------ CREATE URL ------
	reqBody := []byte(`{"long_url": "https://google.com"}`)
	req := httptest.NewRequest("POST", "/api/v1/urls", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("Authorization", *responseBody.Token)
	resp := lo.Must(tests.App.Test(req))

	assert.Equal(t, 201, resp.StatusCode)
	var urlResponseBody dtos.UrlResponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Error reading response body: %v", err)
	}
	if err := json.Unmarshal(body, &urlResponseBody); err != nil {
		t.Fatalf("Error unmarshalling response body: %v", err)
	}
	assert.Equal(t, responseBody.ID, urlResponseBody.CreatorID)

	applogger.Info("Short URL Created", urlResponseBody.ShortURL)

}

func TestUrlsRoute_CreateSpecificUrl(t *testing.T) {

	// ------ REGISTER USER ------
	responseBody := tests.TestUtil_CreateUser(t, "user2584@test.com", "123456")

	// ------ CREATE URL ------
	reqBody := []byte(`{"long_url": "https://example.com"}`)
	req := httptest.NewRequest("PUT", "/api/v1/urls/my_code", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("Authorization", *responseBody.Token)
	resp := lo.Must(tests.App.Test(req))

	assert.Equal(t, 201, resp.StatusCode)
	var urlResponseBody dtos.UrlResponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Error reading response body: %v", err)
	}
	if err := json.Unmarshal(body, &urlResponseBody); err != nil {
		t.Fatalf("Error unmarshalling response body: %v", err)
	}
	assert.Equal(t, responseBody.ID, urlResponseBody.CreatorID)
	assert.Equal(t, "my_code", urlResponseBody.ShortURL)

	// ------ CREATE URL WITH SAME CODE ------
	reqBody = []byte(`{"long_url": "https://example2.com"}`)
	req = httptest.NewRequest("PUT", "/api/v1/urls/my_code", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("Authorization", *responseBody.Token)
	resp = lo.Must(tests.App.Test(req))

	assert.Equal(t, 409, resp.StatusCode)
	var errorResponseBody dtos.ErrorResponse
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Error reading response body: %v", err)
	}
	if err := json.Unmarshal(body, &errorResponseBody); err != nil {
		t.Fatalf("Error unmarshalling response body: %v", err)
	}

	assert.Equal(t, "URL already exists", errorResponseBody.Message)

}
