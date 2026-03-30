package e2e

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http/httptest"
	"onepixel_backend/src/config"
	"onepixel_backend/src/db"
	"onepixel_backend/src/db/models"
	"onepixel_backend/src/dtos"
	"onepixel_backend/tests"
	"strconv"
	"testing"
	"time"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestUrlsRoute_GroupedUrls(t *testing.T) {
	t.Cleanup(tests.TestUtil_FlushEventsDb)

	owner := tests.TestUtil_CreateUser(t, "group-owner@test.com", "123456")
	other := tests.TestUtil_CreateUser(t, "group-other@test.com", "123456")

	groupReqBody := []byte(`{"short_path":"grpE2E","creator_id":` + strconv.FormatUint(owner.ID, 10) + `}`)
	groupReq := httptest.NewRequest("POST", "/api/v1/urls/groups", bytes.NewBuffer(groupReqBody))
	groupReq.Header.Set("Content-Type", "application/json; charset=UTF-8")
	groupReq.Header.Set("X-API-Key", config.AdminApiKey)
	groupResp := lo.Must(tests.App.Test(groupReq))

	assert.Equal(t, 201, groupResp.StatusCode)
	var groupResponseBody dtos.UrlGroupResponse
	body, err := io.ReadAll(groupResp.Body)
	if err != nil {
		t.Fatalf("Error reading response body: %v", err)
	}
	if err := json.Unmarshal(body, &groupResponseBody); err != nil {
		t.Fatalf("Error unmarshalling response body: %v", err)
	}
	assert.Equal(t, "grpE2E", groupResponseBody.ShortPath)
	assert.Equal(t, owner.ID, groupResponseBody.CreatorID)

	specificReq := httptest.NewRequest("POST", "/api/v1/urls/groups/grpE2E/shorten/grp123", bytes.NewBuffer([]byte(`{"long_url":"https://example.com/grouped-specific"}`)))
	specificReq.Header.Set("Content-Type", "application/json; charset=UTF-8")
	specificReq.Header.Set("Authorization", *owner.Token)
	specificResp := lo.Must(tests.App.Test(specificReq))

	assert.Equal(t, 201, specificResp.StatusCode)
	var specificUrlResponse dtos.UrlResponse
	body, err = io.ReadAll(specificResp.Body)
	if err != nil {
		t.Fatalf("Error reading response body: %v", err)
	}
	if err := json.Unmarshal(body, &specificUrlResponse); err != nil {
		t.Fatalf("Error unmarshalling response body: %v", err)
	}
	assert.Equal(t, config.RedirUrlBase+"grpE2E/grp123", specificUrlResponse.ShortURL)

	randomReq := httptest.NewRequest("POST", "/api/v1/urls/groups/grpE2E/shorten", bytes.NewBuffer([]byte(`{"long_url":"https://example.com/grouped-random"}`)))
	randomReq.Header.Set("Content-Type", "application/json; charset=UTF-8")
	randomReq.Header.Set("Authorization", *owner.Token)
	randomResp := lo.Must(tests.App.Test(randomReq))

	assert.Equal(t, 201, randomResp.StatusCode)
	var randomUrlResponse dtos.UrlResponse
	body, err = io.ReadAll(randomResp.Body)
	if err != nil {
		t.Fatalf("Error reading response body: %v", err)
	}
	if err := json.Unmarshal(body, &randomUrlResponse); err != nil {
		t.Fatalf("Error unmarshalling response body: %v", err)
	}
	assert.Contains(t, randomUrlResponse.ShortURL, config.RedirUrlBase+"grpE2E/")

	forbiddenReq := httptest.NewRequest("POST", "/api/v1/urls/groups/grpE2E/shorten/forbid1", bytes.NewBuffer([]byte(`{"long_url":"https://example.com/forbidden"}`)))
	forbiddenReq.Header.Set("Content-Type", "application/json; charset=UTF-8")
	forbiddenReq.Header.Set("Authorization", *other.Token)
	forbiddenResp := lo.Must(tests.App.Test(forbiddenReq))

	assert.Equal(t, 403, forbiddenResp.StatusCode)

	redirectReq := httptest.NewRequest("GET", config.RedirUrlBase+"grpE2E/grp123", nil)
	redirectReq.Header.Set("User-Agent", "grouped-test-agent")
	redirectReq.Header.Set("X-Forwarded-For", "42.108.28.82")
	redirectResp := lo.Must(tests.MainApp.Test(redirectReq))

	assert.Equal(t, 301, redirectResp.StatusCode)
	assert.Equal(t, "https://example.com/grouped-specific", redirectResp.Header.Get("Location"))

	time.Sleep(200 * time.Millisecond)

	var redirectCount int64
	err = db.GetEventsDB().Model(&models.EventRedirect{}).Where("short_url = ?", "grpE2E/grp123").Count(&redirectCount).Error
	assert.Nil(t, err)
	assert.GreaterOrEqual(t, redirectCount, int64(1))
}
