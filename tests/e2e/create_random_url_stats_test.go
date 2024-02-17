package e2e

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
	t.Cleanup(tests.TestUtil_FlushEventsDb)

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
			resp := lo.Must(tests.MainApp.Test(req))

			assert.Equal(t, 301, resp.StatusCode)
			assert.Equal(t, "https://google.com", resp.Header.Get("Location"))
			return resp.Header.Get("Location")
		})
	})
	ch := lo.FanIn(5, chans...)
	applogger.Info("Redirected to: ", lo.Times(3, func(i int) string { return <-ch }))
	// give time for analytics to flush
	time.Sleep(200 * time.Millisecond)

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
