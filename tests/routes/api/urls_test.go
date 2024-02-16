package api

import (
	"bytes"
	"encoding/json"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http/httptest"
	"onepixel_backend/src/db/models"
	"onepixel_backend/src/dtos"
	"onepixel_backend/src/utils/applogger"
	"onepixel_backend/tests"
	"testing"
	"time"
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

	// ------ CHECK REDIRECT ------
	chans := lo.Times(3, func(i int) <-chan string {
		return lo.Async(func() string {
			req := httptest.NewRequest("GET", "/"+urlResponseBody.ShortURL, nil)
			req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3")
			req.Header.Set("X-Forwarded-For", "42.108.28.82")
			resp = lo.Must(tests.MainApp.Test(req))

			assert.Equal(t, 301, resp.StatusCode)
			assert.Equal(t, "https://google.com", resp.Header.Get("Location"))
			return resp.Header.Get("Location")
		})
	})
	ch := lo.FanIn(5, chans...)
	applogger.Info("Redirected to: ", lo.Times(3, func(i int) string { return <-ch }))
	// give time for analytics to flush
	time.Sleep(1 * time.Second)

	// ------ CHECK STATS ------

	req = httptest.NewRequest("GET", "/api/v1/stats", nil)
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.Header.Set("Authorization", *responseBody.Token)
	resp = lo.Must(tests.App.Test(req))

	assert.Equal(t, 200, resp.StatusCode)

	var statsResponseBody []models.EventRedirectCountView

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Error reading response body: %v", err)
	}
	if err := json.Unmarshal(body, &statsResponseBody); err != nil {
		t.Fatalf("Error unmarshalling response body: %v", err)
	}
	assert.GreaterOrEqual(t, len(statsResponseBody), 1)

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

	// ------ CHECK REDIRECT ------
	chans := lo.Times(3, func(i int) <-chan string {
		return lo.Async(func() string {
			req = httptest.NewRequest("GET", "/"+urlResponseBody.ShortURL, nil)
			resp = lo.Must(tests.MainApp.Test(req))

			assert.Equal(t, 301, resp.StatusCode)
			assert.Equal(t, "https://example.com", resp.Header.Get("Location"))
			return resp.Header.Get("Location")
		})
	})
	ch := lo.FanIn(5, chans...)
	applogger.Info("Redirected to: ", lo.Times(3, func(i int) string { return <-ch }))

	// give time for analytics to flush
	time.Sleep(1 * time.Second)

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
